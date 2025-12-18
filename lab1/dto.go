package main

// CreateOrderRequest представляет запрос на создание заказа
type CreateOrderRequest struct {
	OrderName string `json:"order_name"`
	StartDate string `json:"start_date"`
}

// UpdateOrderRequest представляет запрос на обновление заказа
type UpdateOrderRequest struct {
	OrderName string `json:"order_name"`
	StartDate string `json:"start_date"`
}

// AddTaskRequest представляет запрос на добавление задачи
type AddTaskRequest struct {
	Task     string `json:"task"`
	Duration int    `json:"duration"`
	Resource int    `json:"resource"`
	Pred     []int  `json:"pred"`
}

// UpdatePredecessorsRequest представляет запрос на обновление предшественников
type UpdatePredecessorsRequest struct {
	Pred []int `json:"pred"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

