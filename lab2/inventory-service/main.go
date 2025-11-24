package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Inventory struct {
	ProductCode  string  `json:"product_code"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	CurrentPrice string  `json:"current_price"`
}

var db *sql.DB

func initDB() {
	var err error
	connStr := "host=localhost port=5432 user=user password=password dbname=lab2db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS inventory (
		product_code VARCHAR PRIMARY KEY,
		name VARCHAR NOT NULL,
		quantity NUMERIC NOT NULL,
		current_price VARCHAR NOT NULL
	)`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func createInventory(w http.ResponseWriter, r *http.Request) {
	var inv Inventory
	json.NewDecoder(r.Body).Decode(&inv)

	_, err := db.Exec("INSERT INTO inventory (product_code, name, quantity, current_price) VALUES ($1, $2, $3, $4)",
		inv.ProductCode, inv.Name, inv.Quantity, inv.CurrentPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT product_code, name, quantity, current_price FROM inventory")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	inventory := []Inventory{}
	for rows.Next() {
		var inv Inventory
		rows.Scan(&inv.ProductCode, &inv.Name, &inv.Quantity, &inv.CurrentPrice)
		inventory = append(inventory, inv)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func getInventoryItem(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	var inv Inventory

	err := db.QueryRow("SELECT product_code, name, quantity, current_price FROM inventory WHERE product_code = $1", code).
		Scan(&inv.ProductCode, &inv.Name, &inv.Quantity, &inv.CurrentPrice)
	if err != nil {
		http.Error(w, "Inventory item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func updateInventory(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	var inv Inventory
	json.NewDecoder(r.Body).Decode(&inv)

	_, err := db.Exec("UPDATE inventory SET quantity = $1, current_price = $2 WHERE product_code = $3",
		inv.Quantity, inv.CurrentPrice, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/inventory", createInventory).Methods("POST")
	r.HandleFunc("/inventory", getInventory).Methods("GET")
	r.HandleFunc("/inventory/{code}", getInventoryItem).Methods("GET")
	r.HandleFunc("/inventory/{code}", updateInventory).Methods("PUT")

	log.Println("Inventory Service запущен на :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}

