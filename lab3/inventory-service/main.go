package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type InventoryItem struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type CheckRequest struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

type CheckResponse struct {
	Available bool   `json:"available"`
	Message   string `json:"message"`
}

type ReserveRequest struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

// In-memory хранилище
var (
	inventory = map[string]*InventoryItem{
		"PROD001": {ProductCode: "PROD001", ProductName: "Ноутбук", Quantity: 10, Price: 50000},
		"PROD002": {ProductCode: "PROD002", ProductName: "Мышь", Quantity: 50, Price: 1000},
		"PROD003": {ProductCode: "PROD003", ProductName: "Клавиатура", Quantity: 30, Price: 3000},
	}
	mu sync.RWMutex
)

func checkAvailability(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.RLock()
	item, exists := inventory[req.ProductCode]
	mu.RUnlock()

	response := CheckResponse{}
	if !exists {
		response.Available = false
		response.Message = "Продукт не найден"
	} else if item.Quantity < req.Quantity {
		response.Available = false
		response.Message = "Недостаточно товара на складе"
	} else {
		response.Available = true
		response.Message = "Товар доступен"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func reserveInventory(w http.ResponseWriter, r *http.Request) {
	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	item, exists := inventory[req.ProductCode]
	if !exists {
		http.Error(w, "Продукт не найден", http.StatusNotFound)
		return
	}

	if item.Quantity < req.Quantity {
		http.Error(w, "Недостаточно товара на складе", http.StatusBadRequest)
		return
	}

	item.Quantity -= req.Quantity
	log.Printf("Зарезервировано: %s x %d, осталось: %d", req.ProductCode, req.Quantity, item.Quantity)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "reserved"})
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	items := make([]*InventoryItem, 0, len(inventory))
	for _, item := range inventory {
		items = append(items, item)
	}
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/inventory/check", checkAvailability).Methods("POST")
	r.HandleFunc("/inventory/reserve", reserveInventory).Methods("POST")
	r.HandleFunc("/inventory", getInventory).Methods("GET")

	log.Println("Inventory Service запущен на :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}

