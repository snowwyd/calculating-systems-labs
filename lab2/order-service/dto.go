package main

// CreateOrderRequest представляет запрос на создание заказа
type CreateOrderRequest struct {
	OrderCode    string        `json:"order_code"`
	OrderDate    string        `json:"order_date"`
	ProductItems []ProductItem `json:"product_items"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

