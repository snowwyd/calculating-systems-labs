package main

import "github.com/gorilla/mux"

// SetupRoutes настраивает маршруты для API
func SetupRoutes(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	router.HandleFunc("/products", handler.GetProducts).Methods("GET")
	router.HandleFunc("/products/{code}", handler.GetProduct).Methods("GET")
	router.HandleFunc("/products/{code}", handler.DeleteProduct).Methods("DELETE")

	return router
}

