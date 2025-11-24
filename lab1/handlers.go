package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Вспомогательные функции для HTTP ответов

func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}

func parseIDFromRequest(vars map[string]string, key string) (int, error) {
	return strconv.Atoi(vars[key])
}

// HTTP обработчики для Orders

func (s *Storage) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var data struct {
		OrderName string `json:"order_name"`
		StartDate string `json:"start_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order := s.CreateOrder(data.OrderName, data.StartDate)
	respondJSON(w, order, http.StatusCreated)
}

func (s *Storage) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	orders := s.GetAllOrders()
	respondJSON(w, orders, http.StatusOK)
}

func (s *Storage) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(mux.Vars(r), "id")
	if err != nil {
		respondError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, exists := s.GetOrder(id)
	if !exists {
		respondError(w, "Order not found", http.StatusNotFound)
		return
	}

	respondJSON(w, order, http.StatusOK)
}

func (s *Storage) handleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(mux.Vars(r), "id")
	if err != nil {
		respondError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var data struct {
		OrderName string `json:"order_name"`
		StartDate string `json:"start_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, exists := s.UpdateOrder(id, data.OrderName, data.StartDate)
	if !exists {
		respondError(w, "Order not found", http.StatusNotFound)
		return
	}

	respondJSON(w, order, http.StatusOK)
}

func (s *Storage) handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(mux.Vars(r), "id")
	if err != nil {
		respondError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if !s.DeleteOrder(id) {
		respondError(w, "Order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HTTP обработчики для Tasks

func (s *Storage) handleAddTask(w http.ResponseWriter, r *http.Request) {
	orderID, err := parseIDFromRequest(mux.Vars(r), "id")
	if err != nil {
		respondError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Task     string `json:"task"`
		Duration int    `json:"duration"`
		Resource int    `json:"resource"`
		Pred     []int  `json:"pred"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, exists := s.AddTask(orderID, data.Task, data.Duration, data.Resource, data.Pred)
	if !exists {
		respondError(w, "Order not found", http.StatusNotFound)
		return
	}

	respondJSON(w, task, http.StatusCreated)
}

func (s *Storage) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := parseIDFromRequest(vars, "orderId")
	if err != nil {
		respondError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	taskID, err := parseIDFromRequest(vars, "taskId")
	if err != nil {
		respondError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if !s.DeleteTask(orderID, taskID) {
		respondError(w, "Order or task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Storage) handleUpdateTaskPredecessors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := parseIDFromRequest(vars, "orderId")
	if err != nil {
		respondError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	taskID, err := parseIDFromRequest(vars, "taskId")
	if err != nil {
		respondError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Pred []int `json:"pred"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, exists := s.UpdateTaskPredecessors(orderID, taskID, data.Pred)
	if !exists {
		respondError(w, "Order or task not found", http.StatusNotFound)
		return
	}

	respondJSON(w, task, http.StatusOK)
}
