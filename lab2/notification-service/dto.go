package main

// CreateNotificationRequest представляет запрос на создание уведомления
type CreateNotificationRequest struct {
	MessageType string `json:"message_type"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

