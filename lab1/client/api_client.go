package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const apiURL = "http://localhost:8080"

// APIClient предоставляет методы для работы с REST API
type APIClient struct {
	baseURL string
}

// NewAPIClient создаёт новый клиент API
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{baseURL: baseURL}
}

// GetOrder получает заказ по ID
func (c *APIClient) GetOrder(orderID int) (*Order, error) {
	url := fmt.Sprintf("%s/orders/%d", c.baseURL, orderID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var order Order
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return &order, nil
}

// CreateOrder создаёт новый заказ
func (c *APIClient) CreateOrder(orderName, startDate string) (int, error) {
	orderData := map[string]string{
		"order_name": orderName,
		"start_date": startDate,
	}

	orderJSON, err := json.Marshal(orderData)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга: %w", err)
	}

	url := fmt.Sprintf("%s/orders", c.baseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(orderJSON))
	if err != nil {
		return 0, fmt.Errorf("ошибка создания заказа: %w", err)
	}
	defer resp.Body.Close()

	var order Order
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if err := json.Unmarshal(body, &order); err != nil {
		return 0, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return order.ID, nil
}

// AddTask добавляет работу к заказу
func (c *APIClient) AddTask(orderID int, taskData map[string]interface{}) error {
	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга: %w", err)
	}

	url := fmt.Sprintf("%s/orders/%d/tasks", c.baseURL, orderID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(taskJSON))
	if err != nil {
		return fmt.Errorf("ошибка добавления задачи: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// CreateTestOrder создаёт тестовый заказ с работами
func (c *APIClient) CreateTestOrder() (int, error) {
	orderID, err := c.CreateOrder("изготовить изделие", "01-01-2022")
	if err != nil {
		return 0, err
	}

	tasks := []map[string]interface{}{
		{"task": "задача 1", "duration": 5, "resource": 3, "pred": []int{}},
		{"task": "задача 2", "duration": 3, "resource": 5, "pred": []int{}},
		{"task": "задача 3", "duration": 4, "resource": 2, "pred": []int{}},
		{"task": "задача 4", "duration": 6, "resource": 4, "pred": []int{}},
		{"task": "задача 5", "duration": 2, "resource": 3, "pred": []int{}},
	}

	for _, taskData := range tasks {
		if err := c.AddTask(orderID, taskData); err != nil {
			log.Printf("Предупреждение: %v", err)
		}
	}

	return orderID, nil
}
