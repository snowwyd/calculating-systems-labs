# Как рассказывать

Рассказать только про graphql, в этом вся суть работы. Все взято с прошлой лабы, но переделано под graphql.

В graphql у нас есть один запрос, который состоит из частей. То есть его можно построить по этим частям. Это надстройка над обычным Rest API, просто более крутая.

GraphQL у нас принимает запросы по одному эндпоинту и начинает резолвить этот запрос.
Вот например как тут (с.182 в main.go в `graphql`):
```go
"product": &graphql.Field{
    Type: productGraphQLType,
    Resolve: func(p graphql.ResolveParams) (interface{}, error) { //  вот это как раз резолвер 
        // Проверяем разные типы источников данных  
        switch order := p.Source.(type) {
        case *Order:
            return getProductFromService(order.ProductID) // это методд бизнес логики (нажми на Ctrl+ПКМ, чтобы перейти к этой функции)
        case map[string]interface{}:
            if productID, ok := order["product_id"].(string); ok && productID != "" {
                return getProductFromService(productID) // и это
            }
            if productID, ok := order["productId"].(string); ok && productID != "" {
                return getProductFromService(productID) // и это
            }
        case Order:
            return getProductFromService(order.ProductID) // и вот это ничесе 
        }
        return nil, nil
    },
},
```

как можно заметить, уже этот резолвер вызывает бизнес-логику `getProductFromService`:
```go
// ну тут вот уже бизнес логика сама
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
```

а вот эти две функции в конце файла у нас служат оберткой над HTTP запросом, чтобы graphQL работал
```go
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
```

# Как запускать
1. Очисти докер `docker-compose down -v`
2. Запусти докер `docker-compose up -d`
3. 3 терминала с `go run .`
4. в `lab4/` введи в терминал `run-tests.bat`
5. Заверши работу сервисов
6. Очисти докер