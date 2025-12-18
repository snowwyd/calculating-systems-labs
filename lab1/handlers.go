package main

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func extractIDFromPath(vars map[string]string, key string) (int, error) {
	return strconv.Atoi(vars[key])
}

// HTTP обработчики для Orders

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(req)
	if err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, order, http.StatusCreated)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.service.GetAllOrders()
	writeJSONResponse(w, orders, http.StatusOK)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(mux.Vars(r), "id")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "id", Message: "невалидный ID"}, http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(id)
	if err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, order, http.StatusOK)
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(mux.Vars(r), "id")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "id", Message: "невалидный ID"}, http.StatusBadRequest)
		return
	}

	var req UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	order, err := h.service.UpdateOrder(id, req)
	if err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, order, http.StatusOK)
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(mux.Vars(r), "id")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "id", Message: "невалидный ID"}, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteOrder(id); err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HTTP обработчики для Tasks

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	orderID, err := extractIDFromPath(mux.Vars(r), "id")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "id", Message: "невалидный ID заказа"}, http.StatusBadRequest)
		return
	}

	var req AddTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	task, err := h.service.AddTask(orderID, req)
	if err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, task, http.StatusCreated)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := extractIDFromPath(vars, "orderId")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "orderId", Message: "невалидный ID заказа"}, http.StatusBadRequest)
		return
	}

	taskID, err := extractIDFromPath(vars, "taskId")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "taskId", Message: "невалидный ID задачи"}, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(orderID, taskID); err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else if _, ok := err.(ErrTaskNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTaskPredecessors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := extractIDFromPath(vars, "orderId")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "orderId", Message: "невалидный ID заказа"}, http.StatusBadRequest)
		return
	}

	taskID, err := extractIDFromPath(vars, "taskId")
	if err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "taskId", Message: "невалидный ID задачи"}, http.StatusBadRequest)
		return
	}

	var req UpdatePredecessorsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	task, err := h.service.UpdateTaskPredecessors(orderID, taskID, req)
	if err != nil {
		if _, ok := err.(ErrOrderNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else if _, ok := err.(ErrTaskNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, task, http.StatusOK)
}
