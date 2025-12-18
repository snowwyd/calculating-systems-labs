package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	serverPort = ":8081"
	connStr    = "host=localhost port=5432 user=user password=password dbname=lab2db sslmode=disable"
)

func main() {
	// Инициализация базы данных
	db, err := InitDatabase(connStr)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	// Инициализация зависимостей
	repository := NewPostgreSQLProductRepository(db)
	service := NewProductService(repository)
	handler := NewHandler(service)
	router := SetupRoutes(handler)

	log.Printf("Product Service запущен на %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
