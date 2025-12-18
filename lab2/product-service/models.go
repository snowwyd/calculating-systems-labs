package main

// Product представляет продукт
type Product struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
}

