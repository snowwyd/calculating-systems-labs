package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Product struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
}

type ProductItem struct {
	ProductCode string  `json:"product_code"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Cost        float64 `json:"cost"`
}

type Order struct {
	OrderCode    string        `json:"order_code"`
	OrderDate    string        `json:"order_date"`
	ProductItems []ProductItem `json:"product_items"`
}

type Inventory struct {
	ProductCode  string  `json:"product_code"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	CurrentPrice string  `json:"current_price"`
}

type Notification struct {
	MessageType string `json:"message_type"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type TestResult struct {
	Service string
	Test    string
	Status  string
	Error   string
}

var results []TestResult

func addResult(service, test, status, err string) {
	results = append(results, TestResult{
		Service: service,
		Test:    test,
		Status:  status,
		Error:   err,
	})
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func httpPost(url string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func httpDelete(url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func testProductService() {
	fmt.Println("\n=== Тестирование Product Service (8081) ===")

	_, err := httpGet("http://localhost:8081/products")
	if err != nil {
		addResult("Product", "Доступность", "FAIL", err.Error())
		return
	}
	addResult("Product", "Доступность", "OK", "")

	product := Product{
		Code:        "TEST001",
		Name:        "Тестовый продукт",
		Weight:      1.5,
		Description: "Описание тестового продукта",
	}

	_, err = httpPost("http://localhost:8081/products", product)
	if err != nil {
		addResult("Product", "Создание", "FAIL", err.Error())
		return
	}
	addResult("Product", "Создание", "OK", "")

	body, err := httpGet("http://localhost:8081/products/TEST001")
	if err != nil {
		addResult("Product", "Чтение", "FAIL", err.Error())
		return
	}

	var retrieved Product
	if err := json.Unmarshal(body, &retrieved); err != nil {
		addResult("Product", "Чтение", "FAIL", err.Error())
		return
	}

	if retrieved.Code != product.Code {
		addResult("Product", "Чтение", "FAIL", "Данные не совпадают")
		return
	}
	addResult("Product", "Чтение", "OK", "")

	err = httpDelete("http://localhost:8081/products/TEST001")
	if err != nil {
		addResult("Product", "Удаление", "FAIL", err.Error())
		return
	}
	addResult("Product", "Удаление", "OK", "")
}

func testOrderService() {
	fmt.Println("\n=== Тестирование Order Service (8082) ===")

	_, err := httpGet("http://localhost:8082/orders")
	if err != nil {
		addResult("Order", "Доступность", "FAIL", err.Error())
		return
	}
	addResult("Order", "Доступность", "OK", "")

	order := Order{
		OrderCode: "TEST_ORDER001",
		OrderDate: "2024-01-01",
		ProductItems: []ProductItem{
			{
				ProductCode: "P001",
				Name:        "Продукт 1",
				Quantity:    5,
				Cost:        100.50,
			},
		},
	}

	_, err = httpPost("http://localhost:8082/orders", order)
	if err != nil {
		addResult("Order", "Создание", "FAIL", err.Error())
		return
	}
	addResult("Order", "Создание", "OK", "")

	body, err := httpGet("http://localhost:8082/orders/TEST_ORDER001")
	if err != nil {
		addResult("Order", "Чтение", "FAIL", err.Error())
		return
	}

	var retrieved Order
	if err := json.Unmarshal(body, &retrieved); err != nil {
		addResult("Order", "Чтение", "FAIL", err.Error())
		return
	}

	if retrieved.OrderCode != order.OrderCode {
		addResult("Order", "Чтение", "FAIL", "Данные не совпадают")
		return
	}
	addResult("Order", "Чтение", "OK", "")

	err = httpDelete("http://localhost:8082/orders/TEST_ORDER001")
	if err != nil {
		addResult("Order", "Удаление", "FAIL", err.Error())
		return
	}
	addResult("Order", "Удаление", "OK", "")
}

func testInventoryService() {
	fmt.Println("\n=== Тестирование Inventory Service (8083) ===")

	_, err := httpGet("http://localhost:8083/inventory")
	if err != nil {
		addResult("Inventory", "Доступность", "FAIL", err.Error())
		return
	}
	addResult("Inventory", "Доступность", "OK", "")

	inventory := Inventory{
		ProductCode:  "TEST_INV001",
		Name:         "Тестовый товар",
		Quantity:     100,
		CurrentPrice: "500.00",
	}

	_, err = httpPost("http://localhost:8083/inventory", inventory)
	if err != nil {
		addResult("Inventory", "Создание", "FAIL", err.Error())
		return
	}
	addResult("Inventory", "Создание", "OK", "")

	body, err := httpGet("http://localhost:8083/inventory/TEST_INV001")
	if err != nil {
		addResult("Inventory", "Чтение", "FAIL", err.Error())
		return
	}

	var retrieved Inventory
	if err := json.Unmarshal(body, &retrieved); err != nil {
		addResult("Inventory", "Чтение", "FAIL", err.Error())
		return
	}

	if retrieved.ProductCode != inventory.ProductCode {
		addResult("Inventory", "Чтение", "FAIL", "Данные не совпадают")
		return
	}
	addResult("Inventory", "Чтение", "OK", "")
}

func testNotificationService() {
	fmt.Println("\n=== Тестирование Notification Service (8084) ===")

	_, err := httpGet("http://localhost:8084/notifications")
	if err != nil {
		addResult("Notification", "Доступность", "FAIL", err.Error())
		return
	}
	addResult("Notification", "Доступность", "OK", "")

	notification := Notification{
		MessageType: "TEST",
		Description: "Тестовое уведомление",
		Date:        time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err = httpPost("http://localhost:8084/notifications", notification)
	if err != nil {
		addResult("Notification", "Создание", "FAIL", err.Error())
		return
	}
	addResult("Notification", "Создание", "OK", "")

	body, err := httpGet("http://localhost:8084/notifications")
	if err != nil {
		addResult("Notification", "Чтение", "FAIL", err.Error())
		return
	}

	var notifications []Notification
	if err := json.Unmarshal(body, &notifications); err != nil {
		addResult("Notification", "Чтение", "FAIL", err.Error())
		return
	}

	if len(notifications) == 0 {
		addResult("Notification", "Чтение", "FAIL", "Нет уведомлений")
		return
	}
	addResult("Notification", "Чтение", "OK", "")
}

func testIntegrationScenario() {
	fmt.Println("\n=== Интеграционный сценарий ===")

	// 1. Создать продукт
	product := Product{
		Code:        "SCENARIO001",
		Name:        "Интеграционный товар",
		Weight:      2.0,
		Description: "Товар для интеграционного теста",
	}

	_, err := httpPost("http://localhost:8081/products", product)
	if err != nil {
		addResult("Scenario", "Шаг 1: Создание продукта", "FAIL", err.Error())
		return
	}
	addResult("Scenario", "Шаг 1: Создание продукта", "OK", "")

	// 2. Добавить на склад
	inventory := Inventory{
		ProductCode:  "SCENARIO001",
		Name:         "Интеграционный товар",
		Quantity:     50,
		CurrentPrice: "1000.00",
	}

	_, err = httpPost("http://localhost:8083/inventory", inventory)
	if err != nil {
		addResult("Scenario", "Шаг 2: Добавление на склад", "FAIL", err.Error())
		return
	}
	addResult("Scenario", "Шаг 2: Добавление на склад", "OK", "")

	// 3. Создать заказ
	order := Order{
		OrderCode: "SCENARIO_ORDER001",
		OrderDate: time.Now().Format("2006-01-02"),
		ProductItems: []ProductItem{
			{
				ProductCode: "SCENARIO001",
				Name:        "Интеграционный товар",
				Quantity:    3,
				Cost:        1000.00,
			},
		},
	}

	_, err = httpPost("http://localhost:8082/orders", order)
	if err != nil {
		addResult("Scenario", "Шаг 3: Создание заказа", "FAIL", err.Error())
		return
	}
	addResult("Scenario", "Шаг 3: Создание заказа", "OK", "")

	// 4. Создать уведомление
	notification := Notification{
		MessageType: "ORDER_CREATED",
		Description: fmt.Sprintf("Заказ %s успешно создан", order.OrderCode),
		Date:        time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err = httpPost("http://localhost:8084/notifications", notification)
	if err != nil {
		addResult("Scenario", "Шаг 4: Создание уведомления", "FAIL", err.Error())
		return
	}
	addResult("Scenario", "Шаг 4: Создание уведомления", "OK", "")

	addResult("Scenario", "Полный сценарий", "OK", "Все шаги выполнены успешно")
}

func printResults() {
	separator := strings.Repeat("=", 80)
	fmt.Println("\n" + separator)
	fmt.Println("РЕЗУЛЬТАТЫ ИНТЕГРАЦИОННОГО ТЕСТИРОВАНИЯ")
	fmt.Println(separator + "\n")

	services := make(map[string][]TestResult)
	for _, result := range results {
		services[result.Service] = append(services[result.Service], result)
	}

	totalTests := 0
	passedTests := 0

	for service, tests := range services {
		fmt.Printf("%s \n", service)
		for _, test := range tests {
			totalTests++
			status := "✗ FAIL"
			if test.Status == "OK" {
				status = "✓ OK"
				passedTests++
			}

			fmt.Printf("│  %s  %s", status, test.Test)
			if test.Error != "" {
				fmt.Printf(" (%s)", test.Error)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	fmt.Println(separator)
	fmt.Printf("Всего тестов: %d | Пройдено: %d | Провалено: %d\n",
		totalTests, passedTests, totalTests-passedTests)

	successRate := float64(passedTests) / float64(totalTests) * 100
	fmt.Printf("Процент успеха: %.1f%%\n", successRate)
	fmt.Println(separator + "\n")

	if passedTests == totalTests {
		fmt.Println("🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!")
	} else {
		fmt.Printf("⚠️  НЕКОТОРЫЕ ТЕСТЫ НЕ ПРОЙДЕНЫ (%d из %d)\n", totalTests-passedTests, totalTests)
	}
}

func main() {
	fmt.Println("ИНТЕГРАЦИОННОЕ ТЕСТИРОВАНИЕ МИКРОСЕРВИСОВ LAB2")

	fmt.Println("Проверка доступности сервисов...")
	time.Sleep(1 * time.Second)

	// Запуск тестов
	testProductService()
	testOrderService()
	testInventoryService()
	testNotificationService()
	testIntegrationScenario()

	// Вывод результатов
	printResults()
}
