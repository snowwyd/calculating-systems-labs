package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Handler предоставляет HTTP обработчики для API
type Handler struct {
	service *OrderService
}

// NewHandler создаёт новый обработчик
func NewHandler(service *OrderService) *Handler {
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

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContext()
	defer cancel()

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(ctx, req)
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, order, http.StatusCreated)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContext()
	defer cancel()

	orders, err := h.service.GetAllOrders(ctx)
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, orders, http.StatusOK)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContext()
	defer cancel()

	code := mux.Vars(r)["code"]
	order, err := h.service.GetOrder(ctx, code)
	if err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, order, http.StatusOK)
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContext()
	defer cancel()

	code := mux.Vars(r)["code"]
	if err := h.service.DeleteOrder(ctx, code); err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

