package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	productServiceURL = "http://localhost:8081"
	orderServiceURL   = "http://localhost:8082"
	graphqlURL        = "http://localhost:8080/graphql"
)

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

// Проверка доступности сервисов
func checkServiceAvailability(url, name string) bool {
	resp, err := http.Get(url)
	if err != nil {
		addResult(name, "Доступность", "FAIL", err.Error())
		return false
	}
	resp.Body.Close()
	addResult(name, "Доступность", "OK", "")
	return true
}

// Тестирование Product Service
func testProductService() {
	log.Println("\n=== Тестирование Product Service ===")

	// Получить все продукты
	resp, err := http.Get(productServiceURL + "/products")
	if err != nil {
		addResult("Product Service", "GET /products", "FAIL", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		addResult("Product Service", "GET /products", "FAIL", fmt.Sprintf("Status: %d", resp.StatusCode))
		return
	}

	var products []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		addResult("Product Service", "GET /products", "FAIL", err.Error())
		return
	}

	if len(products) == 0 {
		addResult("Product Service", "GET /products", "FAIL", "Нет продуктов")
		return
	}

	addResult("Product Service", "GET /products", "OK", fmt.Sprintf("Получено %d продуктов", len(products)))

	// Получить конкретный продукт
	resp, err = http.Get(productServiceURL + "/products/1")
	if err != nil {
		addResult("Product Service", "GET /products/{id}", "FAIL", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		addResult("Product Service", "GET /products/{id}", "OK", "")
	} else {
		addResult("Product Service", "GET /products/{id}", "FAIL", fmt.Sprintf("Status: %d", resp.StatusCode))
	}
}

// Тестирование Order Service
func testOrderService() string {
	log.Println("\n=== Тестирование Order Service ===")

	// Создать заказ
	orderData := map[string]interface{}{
		"customer_id": "TEST_CUSTOMER",
		"product_id":  "1",
		"quantity":    2,
		"total_price": 100000.0,
	}

	jsonData, _ := json.Marshal(orderData)
	resp, err := http.Post(orderServiceURL+"/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		addResult("Order Service", "POST /orders", "FAIL", err.Error())
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		addResult("Order Service", "POST /orders", "FAIL", fmt.Sprintf("Status: %d", resp.StatusCode))
		return ""
	}

	var createdOrder map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &createdOrder); err != nil {
		addResult("Order Service", "POST /orders", "FAIL", err.Error())
		return ""
	}

	orderID := createdOrder["id"].(string)
	addResult("Order Service", "POST /orders", "OK", fmt.Sprintf("Создан заказ %s", orderID))

	// Получить все заказы
	resp, err = http.Get(orderServiceURL + "/orders")
	if err != nil {
		addResult("Order Service", "GET /orders", "FAIL", err.Error())
		return orderID
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		addResult("Order Service", "GET /orders", "OK", "")
	} else {
		addResult("Order Service", "GET /orders", "FAIL", fmt.Sprintf("Status: %d", resp.StatusCode))
	}

	// Получить конкретный заказ
	resp, err = http.Get(orderServiceURL + "/orders/" + orderID)
	if err != nil {
		addResult("Order Service", "GET /orders/{id}", "FAIL", err.Error())
		return orderID
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		addResult("Order Service", "GET /orders/{id}", "OK", "")
	} else {
		addResult("Order Service", "GET /orders/{id}", "FAIL", fmt.Sprintf("Status: %d", resp.StatusCode))
	}

	return orderID
}

// Тестирование GraphQL Gateway
func testGraphQL() {
	log.Println("\n=== Тестирование GraphQL Gateway ===")

	// Тест 1: Получить все продукты
	query1 := `{"query": "{ products { id name price stock } }"}`
	resp, err := http.Post(graphqlURL, "application/json", bytes.NewBufferString(query1))
	if err != nil {
		addResult("GraphQL", "Query: products", "FAIL", err.Error())
	} else {
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		if data, ok := result["data"].(map[string]interface{}); ok {
			if products, ok := data["products"].([]interface{}); ok && len(products) > 0 {
				addResult("GraphQL", "Query: products", "OK", fmt.Sprintf("Получено %d продуктов", len(products)))
			} else {
				addResult("GraphQL", "Query: products", "FAIL", "Нет данных о продуктах")
			}
		} else {
			addResult("GraphQL", "Query: products", "FAIL", "Неверная структура ответа")
		}
	}

	// Тест 2: Получить конкретный продукт
	query2 := `{"query": "{ product(id: \"1\") { id name description price } }"}`
	resp, err = http.Post(graphqlURL, "application/json", bytes.NewBufferString(query2))
	if err != nil {
		addResult("GraphQL", "Query: product(id)", "FAIL", err.Error())
	} else {
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		if data, ok := result["data"].(map[string]interface{}); ok {
			if product, ok := data["product"].(map[string]interface{}); ok && product != nil {
				addResult("GraphQL", "Query: product(id)", "OK", "")
			} else {
				addResult("GraphQL", "Query: product(id)", "FAIL", "Продукт не найден")
			}
		} else {
			addResult("GraphQL", "Query: product(id)", "FAIL", "Неверная структура ответа")
		}
	}

	// Тест 3: Получить все заказы
	query3 := `{"query": "{ orders { id customerId productId quantity } }"}`
	resp, err = http.Post(graphqlURL, "application/json", bytes.NewBufferString(query3))
	if err != nil {
		addResult("GraphQL", "Query: orders", "FAIL", err.Error())
	} else {
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		if data, ok := result["data"].(map[string]interface{}); ok {
			if orders, ok := data["orders"].([]interface{}); ok {
				addResult("GraphQL", "Query: orders", "OK", fmt.Sprintf("Получено %d заказов", len(orders)))
			} else {
				addResult("GraphQL", "Query: orders", "FAIL", "Нет данных о заказах")
			}
		} else {
			addResult("GraphQL", "Query: orders", "FAIL", "Неверная структура ответа")
		}
	}

	// Тест 4: Получить заказ с вложенным продуктом (связанные данные)
	query4 := `{"query": "{ orders { id customerId product { id name price } } }"}`
	resp, err = http.Post(graphqlURL, "application/json", bytes.NewBufferString(query4))
	if err != nil {
		addResult("GraphQL", "Query: orders with product", "FAIL", err.Error())
	} else {
		defer resp.Body.Close()
		var result map[string]interface{}
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &result)

		if data, ok := result["data"].(map[string]interface{}); ok {
			if ordersList, ok := data["orders"].([]interface{}); ok && len(ordersList) > 0 {
				// Проверяем, что есть вложенный продукт
				firstOrder := ordersList[0].(map[string]interface{})
				if product, ok := firstOrder["product"].(map[string]interface{}); ok && product != nil {
					addResult("GraphQL", "Query: orders with product", "OK", "Связанные данные работают")
				} else {
					addResult("GraphQL", "Query: orders with product", "FAIL", "Нет вложенного продукта")
				}
			} else {
				addResult("GraphQL", "Query: orders with product", "WARN", "Нет заказов для проверки")
			}
		} else {
			addResult("GraphQL", "Query: orders with product", "FAIL", "Неверная структура ответа")
		}
	}

	// Тест 5: Комбинированный запрос
	query5 := `{"query": "{ products { id name } orders { id customerId } }"}`
	resp, err = http.Post(graphqlURL, "application/json", bytes.NewBufferString(query5))
	if err != nil {
		addResult("GraphQL", "Query: combined", "FAIL", err.Error())
	} else {
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		if data, ok := result["data"].(map[string]interface{}); ok {
			hasProducts := false
			hasOrders := false

			if products, ok := data["products"].([]interface{}); ok && len(products) > 0 {
				hasProducts = true
			}
			if _, ok := data["orders"].([]interface{}); ok {
				hasOrders = true
			}

			if hasProducts && hasOrders {
				addResult("GraphQL", "Query: combined", "OK", "Комбинированный запрос работает")
			} else {
				addResult("GraphQL", "Query: combined", "FAIL", "Не все данные получены")
			}
		} else {
			addResult("GraphQL", "Query: combined", "FAIL", "Неверная структура ответа")
		}
	}
}

// Вывод результатов
func printResults() {
	separator := strings.Repeat("=", 80)
	fmt.Println("\n" + separator)
	fmt.Println("РЕЗУЛЬТАТЫ ИНТЕГРАЦИОННОГО ТЕСТИРОВАНИЯ LAB4")
	fmt.Println(separator + "\n")

	services := make(map[string][]TestResult)
	for _, result := range results {
		services[result.Service] = append(services[result.Service], result)
	}

	totalTests := 0
	passedTests := 0

	for service, tests := range services {
		fmt.Printf("%s\n", service)
		for _, test := range tests {
			totalTests++
			status := "FAIL"
			if test.Status == "OK" {
				status = "OK"
				passedTests++
			} else if test.Status == "WARN" {
				status = "WARN"
				passedTests++
			}

			fmt.Printf("|  %s  %s", status, test.Test)
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
		fmt.Println("ВСЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!")
	} else {
		fmt.Printf("НЕКОТОРЫЕ ТЕСТЫ НЕ ПРОЙДЕНЫ (%d из %d)\n", totalTests-passedTests, totalTests)
	}
}

func main() {
	fmt.Println("\n================================================================")
	fmt.Println("      ИНТЕГРАЦИОННОЕ ТЕСТИРОВАНИЕ LAB4 - GraphQL Gateway")
	fmt.Println("================================================================\n")

	fmt.Println("Проверка доступности сервисов...")
	time.Sleep(1 * time.Second)

	// Проверка доступности
	productOK := checkServiceAvailability(productServiceURL+"/products", "Product Service")
	orderOK := checkServiceAvailability(orderServiceURL+"/orders", "Order Service")
	graphqlOK := checkServiceAvailability(graphqlURL, "GraphQL Gateway")

	if !productOK || !orderOK || !graphqlOK {
		fmt.Println("\nНе все сервисы доступны! Запустите: start-services.bat")
		printResults()
		return
	}

	fmt.Println("\nВсе сервисы доступны. Начинаем тестирование...\n")

	// Запуск тестов
	testProductService()
	testOrderService()
	testGraphQL()

	// Вывод результатов
	printResults()
}

