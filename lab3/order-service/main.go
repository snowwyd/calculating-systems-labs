package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type Order struct {
	OrderID    string      `json:"order_id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	CreatedAt  string      `json:"created_at"`
}

type CreateOrderRequest struct {
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
}

type CheckRequest struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

type CheckResponse struct {
	Available bool   `json:"available"`
	Message   string `json:"message"`
}

type ReserveRequest struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

const (
	inventoryServiceURL = "http://localhost:8083"
	rabbitMQURL         = "amqp://admin:password@localhost:5672/"
	queueName           = "order_notifications"
)

var (
	orders     = make(map[string]*Order)
	orderCount = 0
	mu         sync.RWMutex
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
)

func initRabbitMQ() error {
	var err error

	// Повторные попытки подключения
	for i := 0; i < 10; i++ {
		rabbitConn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Попытка подключения к RabbitMQ (%d/10)...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}

	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("не удалось открыть канал: %w", err)
	}

	_, err = rabbitCh.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("не удалось объявить очередь: %w", err)
	}

	log.Println("Подключено к RabbitMQ")
	return nil
}

func checkInventory(productCode string, quantity int) (bool, string, error) {
	checkReq := CheckRequest{
		ProductCode: productCode,
		Quantity:    quantity,
	}

	jsonData, _ := json.Marshal(checkReq)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(
		inventoryServiceURL+"/inventory/check",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return false, "", fmt.Errorf("ошибка запроса к Inventory Service: %w", err)
	}
	defer resp.Body.Close()

	var checkResp CheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		return false, "", fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	return checkResp.Available, checkResp.Message, nil
}

func reserveInventory(productCode string, quantity int) error {
	reserveReq := ReserveRequest{
		ProductCode: productCode,
		Quantity:    quantity,
	}

	jsonData, _ := json.Marshal(reserveReq)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(
		inventoryServiceURL+"/inventory/reserve",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("ошибка резервирования: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка резервирования: %s", string(body))
	}

	return nil
}

func publishOrderNotification(order *Order) error {
	message, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = rabbitCh.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	log.Printf("Отправлено уведомление о заказе %s в RabbitMQ", order.OrderID)
	return nil
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 1. Синхронная проверка наличия всех товаров в Inventory Service
	log.Println("Проверка наличия товаров в Inventory Service...")
	for _, item := range req.Items {
		available, message, err := checkInventory(item.ProductCode, item.Quantity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка проверки товара: %v", err), http.StatusInternalServerError)
			return
		}
		if !available {
			http.Error(w, fmt.Sprintf("Товар %s: %s", item.ProductCode, message), http.StatusBadRequest)
			return
		}
		log.Printf("Товар %s доступен", item.ProductCode)
	}

	// 2. Резервирование товаров
	log.Println("Резервирование товаров...")
	for _, item := range req.Items {
		if err := reserveInventory(item.ProductCode, item.Quantity); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка резервирования: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// 3. Создание заказа
	mu.Lock()
	orderCount++
	order := &Order{
		OrderID:    fmt.Sprintf("ORDER%04d", orderCount),
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Status:     "created",
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	// Расчёт общей стоимости
	totalPrice := 0.0
	for _, item := range order.Items {
		totalPrice += item.Price * float64(item.Quantity)
	}
	order.TotalPrice = totalPrice

	orders[order.OrderID] = order
	mu.Unlock()

	log.Printf("Заказ %s создан", order.OrderID)

	// 4. Асинхронная отправка уведомления в RabbitMQ
	if err := publishOrderNotification(order); err != nil {
		log.Printf("Предупреждение: не удалось отправить уведомление: %v", err)
		// Не возвращаем ошибку клиенту, так как заказ уже создан
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	orderList := make([]*Order, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, order)
	}
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderList)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]

	mu.RLock()
	order, exists := orders[orderID]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func main() {
	// Инициализация RabbitMQ
	if err := initRabbitMQ(); err != nil {
		log.Fatalf("Ошибка инициализации RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	r := mux.NewRouter()

	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")

	log.Println("Order Service запущен на :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
