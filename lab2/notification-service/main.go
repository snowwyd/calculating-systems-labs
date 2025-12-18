package main

import (
	"log"
	"net/http"
)

const (
	serverPort = ":8084"
	mongoURI   = "mongodb://admin:password@localhost:27017"
)

func main() {
	// Инициализация базы данных
	collection, err := InitDatabase(mongoURI)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	// Инициализация зависимостей
	repository := NewMongoDBNotificationRepository(collection)
	service := NewNotificationService(repository)
	handler := NewHandler(service)
	router := SetupRoutes(handler)

	log.Printf("Notification Service запущен на %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
