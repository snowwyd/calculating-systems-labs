package main

import "github.com/gorilla/mux"

// SetupRoutes настраивает маршруты для API
func SetupRoutes(storage *Storage) *mux.Router {
	r := mux.NewRouter()

	// Маршруты для заказов
	r.HandleFunc("/orders", storage.handleCreateOrder).Methods("POST")
	r.HandleFunc("/orders", storage.handleGetOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", storage.handleGetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}", storage.handleUpdateOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", storage.handleDeleteOrder).Methods("DELETE")

	// Маршруты для работ
	r.HandleFunc("/orders/{id}/tasks", storage.handleAddTask).Methods("POST")
	r.HandleFunc("/orders/{orderId}/tasks/{taskId}", storage.handleDeleteTask).Methods("DELETE")
	r.HandleFunc("/orders/{orderId}/tasks/{taskId}/predecessors", storage.handleUpdateTaskPredecessors).Methods("PUT")

	return r
}
