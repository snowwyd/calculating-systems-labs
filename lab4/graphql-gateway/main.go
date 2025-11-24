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
	productServiceURL = "http://localhost:8081"
	orderServiceURL   = "http://localhost:8082"
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
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Функции для вызова REST API

func fetchProducts() ([]Product, error) {
	resp, err := httpClient.Get(productServiceURL + "/products")
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

func fetchProduct(id string) (*Product, error) {
	resp, err := httpClient.Get(productServiceURL + "/products/" + id)
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

func fetchOrders() ([]Order, error) {
	resp, err := httpClient.Get(orderServiceURL + "/orders")
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

func fetchOrder(id string) (*Order, error) {
	resp, err := httpClient.Get(orderServiceURL + "/orders/" + id)
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

var productType = graphql.NewObject(graphql.ObjectConfig{
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

var orderType = graphql.NewObject(graphql.ObjectConfig{
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
			Type: productType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Проверяем разные типы источников данных
				switch order := p.Source.(type) {
				case *Order:
					return fetchProduct(order.ProductID)
				case map[string]interface{}:
					if productID, ok := order["product_id"].(string); ok && productID != "" {
						return fetchProduct(productID)
					}
					if productID, ok := order["productId"].(string); ok && productID != "" {
						return fetchProduct(productID)
					}
				case Order:
					return fetchProduct(order.ProductID)
				}
				return nil, nil
			},
		},
	},
})

// Query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"products": &graphql.Field{
			Type: graphql.NewList(productType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return fetchProducts()
			},
		},
		"product": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				return fetchProduct(id)
			},
		},
		"orders": &graphql.Field{
			Type: graphql.NewList(orderType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return fetchOrders()
			},
		},
		"order": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				return fetchOrder(id)
			},
		},
	},
})

// Schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		log.Printf("GraphQL errors: %v", result.Errors)
	}
	return result
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result := executeQuery(reqBody.Query, schema)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func graphiqlHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>GraphQL Gateway</title>
		<style>
			body {
				height: 100vh;
				margin: 0;
				width: 100%;
				overflow: hidden;
				font-family: Arial, sans-serif;
			}
			#graphiql {
				height: 100vh;
			}
		</style>
		<link rel="stylesheet" href="https://unpkg.com/graphiql/graphiql.min.css" />
	</head>
	<body>
		<div id="graphiql">Loading...</div>
		<script
			crossorigin
			src="https://unpkg.com/react/umd/react.production.min.js"
		></script>
		<script
			crossorigin
			src="https://unpkg.com/react-dom/umd/react-dom.production.min.js"
		></script>
		<script
			crossorigin
			src="https://unpkg.com/graphiql/graphiql.min.js"
		></script>
		<script>
			const fetcher = GraphiQL.createFetcher({
				url: '/graphql',
			});
			ReactDOM.render(
				React.createElement(GraphiQL, { fetcher: fetcher }),
				document.getElementById('graphiql'),
			);
		</script>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/graphql", graphqlHandler)
	http.HandleFunc("/", graphiqlHandler)

	log.Println("GraphQL Gateway запущен на :8080")
	log.Println("GraphiQL UI: http://localhost:8080")
	log.Println("GraphQL endpoint: http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

