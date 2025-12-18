package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

const (
	productServiceBaseURL = "http://localhost:8081"
	orderServiceBaseURL   = "http://localhost:8082"
	serverPort            = ":8080"
	requestTimeout        = 10 * time.Second
)

// Модели данных
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type Order struct {
	ID         string  `json:"id"`
	CustomerID string  `json:"customer_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
}

// HTTP клиент с таймаутом
var restClient = &http.Client{
	Timeout: requestTimeout,
}

// Функции для вызова REST API

func getProductsFromService() ([]Product, error) {
	resp, err := restClient.Get(productServiceBaseURL + "/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var products []Product
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func getProductFromService(id string) (*Product, error) {
	resp, err := restClient.Get(productServiceBaseURL + "/products/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var product Product
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func getOrdersFromService() ([]Order, error) {
	resp, err := restClient.Get(orderServiceBaseURL + "/orders")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(body, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func getOrderFromService(id string) (*Order, error) {
	resp, err := restClient.Get(orderServiceBaseURL + "/orders/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("order not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// GraphQL типы

var productGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"stock": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var orderGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"customerId": &graphql.Field{
			Type: graphql.String,
		},
		"productId": &graphql.Field{
			Type: graphql.String,
		},
		"quantity": &graphql.Field{
			Type: graphql.Int,
		},
		"totalPrice": &graphql.Field{
			Type: graphql.Float,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"product": &graphql.Field{
			Type: productGraphQLType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Проверяем разные типы источников данных
				switch order := p.Source.(type) {
				case *Order:
					return getProductFromService(order.ProductID)
				case map[string]interface{}:
					if productID, ok := order["product_id"].(string); ok && productID != "" {
						return getProductFromService(productID)
					}
					if productID, ok := order["productId"].(string); ok && productID != "" {
						return getProductFromService(productID)
					}
				case Order:
					return getProductFromService(order.ProductID)
				}
				return nil, nil
			},
		},
	},
})

// Query
var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"products": &graphql.Field{
			Type: graphql.NewList(productGraphQLType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return getProductsFromService()
			},
		},
		"product": &graphql.Field{
			Type: productGraphQLType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				return getProductFromService(id)
			},
		},
		"orders": &graphql.Field{
			Type: graphql.NewList(orderGraphQLType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return getOrdersFromService()
			},
		},
		"order": &graphql.Field{
			Type: orderGraphQLType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				return getOrderFromService(id)
			},
		},
	},
})

// Schema
var graphQLSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

func processGraphQLQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		log.Printf("GraphQL errors: %v", result.Errors)
	}
	return result
}

func handleGraphQLRequest(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result := processGraphQLQuery(requestBody.Query, graphQLSchema)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/graphql", handleGraphQLRequest)

	log.Printf("GraphQL Gateway запущен на %s\n", serverPort)
	log.Printf("GraphQL endpoint: http://localhost%s/graphql\n", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
