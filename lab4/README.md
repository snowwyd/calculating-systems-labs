# Lab4

### Сервисы

1. **Product Service** (порт 8081)

2. **Order Service** (порт 8082)

3. **GraphQL Gateway** (порт 8080)

## Быстрый запуск

### 1. Запуск сервисов

```cmd
start-services.bat
```

### 2. Запуск интеграционных тестов

```cmd
run-tests.bat
```

### 3. Остановка

```cmd
stop-services.bat
```

## Использование GraphQL

### GraphiQL UI

Откройте в браузере: **http://localhost:8080**

## GraphQL Schema

```graphql
type Product {
  id: String!
  name: String!
  description: String!
  price: Float!
  stock: Int!
}

type Order {
  id: String!
  customerId: String!
  productId: String!
  quantity: Int!
  totalPrice: Float!
  status: String!
  createdAt: String!
  product: Product
}

type Query {
  products: [Product]
  product(id: String!): Product
  orders: [Order]
  order(id: String!): Order
}
```