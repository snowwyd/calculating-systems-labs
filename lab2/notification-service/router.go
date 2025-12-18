package main

import "github.com/gorilla/mux"

// SetupRoutes настраивает маршруты для API
func SetupRoutes(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/notifications", handler.CreateNotification).Methods("POST")
	router.HandleFunc("/notifications", handler.GetNotifications).Methods("GET")

	return router
}

