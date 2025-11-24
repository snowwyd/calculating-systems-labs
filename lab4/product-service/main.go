package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

var (
	products = map[string]*Product{
		"1": {ID: "1", Name: "Ноутбук", Description: "Игровой ноутбук", Price: 50000, Stock: 10},
		"2": {ID: "2", Name: "Мышь", Description: "Беспроводная мышь", Price: 1000, Stock: 50},
		"3": {ID: "3", Name: "Клавиатура", Description: "Механическая клавиатура", Price: 3000, Stock: 30},
	}
	mu sync.RWMutex
)

func getProducts(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	productList := make([]*Product, 0, len(products))
	for _, p := range products {
		productList = append(productList, p)
	}
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productList)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	mu.RLock()
	product, exists := products[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", getProduct).Methods("GET")

	log.Println("Product Service запущен на :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

