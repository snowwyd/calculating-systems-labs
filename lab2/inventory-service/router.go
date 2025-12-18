package main

import "github.com/gorilla/mux"

// SetupRoutes настраивает маршруты для API
func SetupRoutes(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/inventory", handler.CreateInventory).Methods("POST")
	router.HandleFunc("/inventory", handler.GetInventory).Methods("GET")
	router.HandleFunc("/inventory/{code}", handler.GetInventoryItem).Methods("GET")
	router.HandleFunc("/inventory/{code}", handler.UpdateInventory).Methods("PUT")

	return router
}

