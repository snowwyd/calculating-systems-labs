package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Order struct {
	OrderID    string      `json:"order_id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	CreatedAt  string      `json:"created_at"`
}

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type Notification struct {
	NotificationID string  `json:"notification_id"`
	Type           string  `json:"type"`
	OrderID        string  `json:"order_id"`
	CustomerID     string  `json:"customer_id"`
	Message        string  `json:"message"`
	ReceivedAt     string  `json:"received_at"`
}

const (
	rabbitMQURL = "amqp://admin:password@localhost:5672/"
	queueName   = "order_notifications"
)

var notificationCount = 0

func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	var conn *amqp.Connection
	var err error

	// Повторные попытки подключения
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Попытка подключения к RabbitMQ (%d/10)...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func processOrder(order *Order) {
	notificationCount++
	notification := Notification{
		NotificationID: time.Now().Format("20060102150405") + "-" + order.OrderID,
		Type:           "ORDER_CREATED",
		OrderID:        order.OrderID,
		CustomerID:     order.CustomerID,
		Message:        formatNotificationMessage(order),
		ReceivedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	log.Println("================================================================")
	log.Printf("НОВОЕ УВЕДОМЛЕНИЕ #%d", notificationCount)
	log.Println("================================================================")
	log.Printf("Тип: %s", notification.Type)
	log.Printf("Заказ: %s", notification.OrderID)
	log.Printf("Клиент: %s", notification.CustomerID)
	log.Printf("Сообщение: %s", notification.Message)
	log.Printf("Время: %s", notification.ReceivedAt)
	log.Println("================================================================")
}

func formatNotificationMessage(order *Order) string {
	itemCount := 0
	for _, item := range order.Items {
		itemCount += item.Quantity
	}
	return fmt.Sprintf("Создан заказ на сумму %.2f руб. (%d товаров)", order.TotalPrice, itemCount)
}

func main() {
	log.Println("Notification Service запускается...")

	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	log.Println("Подключено к RabbitMQ")
	log.Printf("Подписка на очередь: %s", queueName)

	msgs, err := ch.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Не удалось подписаться на очередь: %v", err)
	}

	log.Println("Ожидание сообщений... (Нажмите Ctrl+C для выхода)")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var order Order
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				log.Printf("Ошибка парсинга сообщения: %v", err)
				continue
			}

			processOrder(&order)
		}
	}()

	<-forever
}

