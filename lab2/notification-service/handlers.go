package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Handler предоставляет HTTP обработчики для API
type Handler struct {
	service *NotificationService
}

// NewHandler создаёт новый обработчик
func NewHandler(service *NotificationService) *Handler {
	return &Handler{service: service}
}

// Вспомогательные функции для HTTP ответов

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: err.Error(),
	}
	writeJSONResponse(w, response, statusCode)
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// HTTP обработчики

func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContext()
	defer cancel()

	var req CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	notification, err := h.service.CreateNotification(ctx, req)
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, notification, http.StatusCreated)
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContext()
	defer cancel()

	notifications, err := h.service.GetAllNotifications(ctx)
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, notifications, http.StatusOK)
}

