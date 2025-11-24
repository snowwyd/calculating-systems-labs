package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Order struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	ProductID  string    `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  string    `json:"created_at"`
}

var (
	orders      = make(map[string]*Order)
	orderCount  = 0
	mu          sync.RWMutex
)

func getOrders(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	orderList := make([]*Order, 0, len(orders))
	for _, o := range orders {
		orderList = append(orderList, o)
	}
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderList)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	mu.RLock()
	order, exists := orders[id]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	orderCount++
	order.ID = fmt.Sprintf("ORDER%04d", orderCount)
	order.Status = "pending"
	order.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	orders[order.ID] = &order
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")

	log.Println("Order Service запущен на :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}

