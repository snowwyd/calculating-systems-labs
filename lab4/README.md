# Lab 4: GraphQL Gateway для микросервисов

Единый GraphQL API для взаимодействия с Product Service и Order Service.

## Архитектура

```
Клиент
  ↓
  ↓ GraphQL Query
  ↓
GraphQL Gateway (порт 8080)
  ↓
  ├─→ REST API → Product Service (8081)
  └─→ REST API → Order Service (8082)
```

### Сервисы

1. **Product Service** (порт 8081)
   - REST API для продуктов
   - GET /products
   - GET /products/{id}

2. **Order Service** (порт 8082)
   - REST API для заказов
   - GET /orders
   - GET /orders/{id}
   - POST /orders

3. **GraphQL Gateway** (порт 8080)
   - Объединяет оба REST API в единую GraphQL схему
   - GraphiQL UI для тестирования запросов

## Быстрый запуск

### 1. Установка зависимостей

```bash
# Product Service
cd product-service
go mod tidy

# Order Service
cd ../order-service
go mod tidy

# GraphQL Gateway
cd ../graphql-gateway
go mod tidy
```

### 2. Запуск сервисов

```cmd
start-services.bat
```

### 3. Запуск интеграционных тестов

```cmd
run-tests.bat
```

Тест автоматически проверит:
- Доступность всех сервисов
- REST API Product Service
- REST API Order Service
- GraphQL запросы всех типов
- Связанные данные (order с вложенным product)
- Комбинированные запросы

### 4. Остановка

```cmd
stop-services.bat
```

## Использование GraphQL

### GraphiQL UI (рекомендуется)

Откройте в браузере: **http://localhost:8080**

Интерактивный интерфейс для тестирования GraphQL запросов с автодополнением.

### Примеры GraphQL запросов

#### 1. Получить все продукты

```graphql
{
  products {
    id
    name
    description
    price
    stock
  }
}
```

#### 2. Получить конкретный продукт

```graphql
{
  product(id: "1") {
    id
    name
    description
    price
    stock
  }
}
```

#### 3. Получить все заказы

```graphql
{
  orders {
    id
    customerId
    productId
    quantity
    totalPrice
    status
    createdAt
  }
}
```

#### 4. Получить заказ с информацией о продукте

```graphql
{
  order(id: "ORDER0001") {
    id
    customerId
    quantity
    totalPrice
    status
    createdAt
    product {
      id
      name
      price
    }
  }
}
```

#### 5. Комбинированный запрос (продукты + заказы)

```graphql
{
  products {
    id
    name
    price
    stock
  }
  orders {
    id
    customerId
    quantity
    totalPrice
    product {
      name
    }
  }
}
```

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