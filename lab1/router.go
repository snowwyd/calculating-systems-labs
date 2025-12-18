package main

import "github.com/gorilla/mux"

// SetupRoutes настраивает маршруты для API
func SetupRoutes(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	// Маршруты для заказов
	router.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	router.HandleFunc("/orders", handler.GetOrders).Methods("GET")
	router.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	router.HandleFunc("/orders/{id}", handler.UpdateOrder).Methods("PUT")
	router.HandleFunc("/orders/{id}", handler.DeleteOrder).Methods("DELETE")

	// Маршруты для задач
	router.HandleFunc("/orders/{id}/tasks", handler.AddTask).Methods("POST")
	router.HandleFunc("/orders/{orderId}/tasks/{taskId}", handler.DeleteTask).Methods("DELETE")
	router.HandleFunc("/orders/{orderId}/tasks/{taskId}/predecessors", handler.UpdateTaskPredecessors).Methods("PUT")

	return router
}
