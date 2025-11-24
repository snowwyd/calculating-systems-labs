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
- ✅ Доступность всех сервисов
- ✅ REST API Product Service
- ✅ REST API Order Service
- ✅ GraphQL запросы всех типов
- ✅ Связанные данные (order с вложенным product)
- ✅ Комбинированные запросы

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

## Тестирование через curl

### GraphQL запрос

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d "{\"query\": \"{ products { id name price } }\"}"
```

### Создать заказ через REST API (Order Service)

```bash
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d "{\"customer_id\":\"CUST001\",\"product_id\":\"1\",\"quantity\":2,\"total_price\":100000}"
```

Затем запросить через GraphQL:

```graphql
{
  orders {
    id
    customerId
    product {
      name
      price
    }
  }
}
```

## Структура проекта

```
lab4/
├── start-services.bat          # Запуск всех сервисов
├── stop-services.bat           # Остановка
├── run-tests.bat               # Интеграционные тесты
├── integration-test/
│   ├── main.go                 # Автоматические тесты
│   └── go.mod
├── product-service/
│   ├── main.go                 # REST API продуктов
│   └── go.mod
├── order-service/
│   ├── main.go                 # REST API заказов
│   └── go.mod
└── graphql-gateway/
    ├── main.go                 # GraphQL Gateway + GraphiQL UI
    └── go.mod
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

## Особенности реализации

✅ **Единая точка входа** - один GraphQL endpoint объединяет несколько REST API

✅ **GraphiQL UI** - интерактивный интерфейс для тестирования

✅ **Связанные данные** - можно запрашивать order с вложенным product

✅ **Гибкость запросов** - клиент выбирает какие поля ему нужны

✅ **Минимальный MVP** - простая и понятная реализация

## Преимущества GraphQL

1. **Один запрос** вместо нескольких REST вызовов
2. **Только нужные данные** - клиент указывает какие поля запросить
3. **Строгая типизация** - автодополнение и валидация
4. **Документация** - встроенная в GraphiQL
5. **Версионирование не требуется** - schema эволюционирует

## Требования

- Go 1.21+
- Браузер (для GraphiQL UI)
