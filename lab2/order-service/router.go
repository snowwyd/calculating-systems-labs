package main

import "github.com/gorilla/mux"

// SetupRoutes настраивает маршруты для API
func SetupRoutes(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	router.HandleFunc("/orders", handler.GetOrders).Methods("GET")
	router.HandleFunc("/orders/{code}", handler.GetOrder).Methods("GET")
	router.HandleFunc("/orders/{code}", handler.DeleteOrder).Methods("DELETE")

	return router
}

