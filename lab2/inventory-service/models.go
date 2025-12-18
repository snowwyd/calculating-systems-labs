package main

// Inventory представляет товар на складе
type Inventory struct {
	ProductCode  string  `json:"product_code"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	CurrentPrice string  `json:"current_price"`
}

