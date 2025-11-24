package main

import (
	"log"
	"net/http"
)

const serverPort = ":8080"

func main() {
	storage := NewStorage()
	router := SetupRoutes(storage)

	log.Printf("Сервер запущен на %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
