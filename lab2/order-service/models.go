package main

// ProductItem представляет элемент продукта в заказе
type ProductItem struct {
	ProductCode string  `json:"product_code" bson:"product_code"`
	Name        string  `json:"name" bson:"name"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Cost        float64 `json:"cost" bson:"cost"`
}

// Order представляет заказ
type Order struct {
	OrderCode    string        `json:"order_code" bson:"order_code"`
	OrderDate    string        `json:"order_date" bson:"order_date"`
	ProductItems []ProductItem `json:"product_items" bson:"product_items"`
}

