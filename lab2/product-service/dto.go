package main

// CreateProductRequest представляет запрос на создание продукта
type CreateProductRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

