package main

// CreateInventoryRequest представляет запрос на создание товара на складе
type CreateInventoryRequest struct {
	ProductCode  string  `json:"product_code"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	CurrentPrice string  `json:"current_price"`
}

// UpdateInventoryRequest представляет запрос на обновление товара на складе
type UpdateInventoryRequest struct {
	Quantity     float64 `json:"quantity"`
	CurrentPrice string  `json:"current_price"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

