package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler предоставляет HTTP обработчики для API
type Handler struct {
	service *ProductService
}

// NewHandler создаёт новый обработчик
func NewHandler(service *ProductService) *Handler {
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

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, ErrInvalidInput{Field: "body", Message: "невалидный JSON"}, http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, product, http.StatusCreated)
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, products, http.StatusOK)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	product, err := h.service.GetProduct(code)
	if err != nil {
		if _, ok := err.(ErrProductNotFound); ok {
			writeErrorResponse(w, err, http.StatusNotFound)
		} else if _, ok := err.(ErrDatabaseError); ok {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		} else {
			writeErrorResponse(w, err, http.StatusBadRequest)
		}
		return
	}

	writeJSONResponse(w, product, http.StatusOK)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	if err := h.service.DeleteProduct(code); err != nil {
		if _, ok := err.(ErrProductNotFound); ok {
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

