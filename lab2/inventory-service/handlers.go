package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler предоставляет HTTP обработчики для API
type Handler struct {
	service *InventoryService
}

// NewHandler создаёт новый обработчик
func NewHandler(service *InventoryService) *Handler {
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

// HTTP обработчики

func (h *Handler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var req CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	inventory, err := h.service.CreateInventory(req)
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, inventory, http.StatusCreated)
}

func (h *Handler) GetInventory(w http.ResponseWriter, r *http.Request) {
	inventory, err := h.service.GetAllInventory()
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, inventory, http.StatusOK)
}

func (h *Handler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	inventory, err := h.service.GetInventory(code)
	if err != nil {
		if _, ok := err.(ErrInventoryNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, inventory, http.StatusOK)
}

func (h *Handler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	var req UpdateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	inventory, err := h.service.UpdateInventory(code, req)
	if err != nil {
		if _, ok := err.(ErrInventoryNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, inventory, http.StatusOK)
}

