package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Product struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
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
	CREATE TABLE IF NOT EXISTS products (
		code VARCHAR PRIMARY KEY,
		name VARCHAR NOT NULL,
		weight NUMERIC NOT NULL,
		description VARCHAR NOT NULL
	)`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	_, err := db.Exec("INSERT INTO products (code, name, weight, description) VALUES ($1, $2, $3, $4)",
		product.Code, product.Name, product.Weight, product.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT code, name, weight, description FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		rows.Scan(&p.Code, &p.Name, &p.Weight, &p.Description)
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	var product Product

	err := db.QueryRow("SELECT code, name, weight, description FROM products WHERE code = $1", code).
		Scan(&product.Code, &product.Name, &product.Weight, &product.Description)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	_, err := db.Exec("DELETE FROM products WHERE code = $1", code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{code}", getProduct).Methods("GET")
	r.HandleFunc("/products/{code}", deleteProduct).Methods("DELETE")

	log.Println("Product Service запущен на :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

