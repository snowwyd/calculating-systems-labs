# Lab 2

## Архитектура

- **Product Service** (порт 8081) - PostgreSQL
- **Order Service** (порт 8082) - MongoDB
- **Inventory Service** (порт 8083) - PostgreSQL
- **Notification Service** (порт 8084) - MongoDB

## Быстрый запуск

### Batch файлы (простой способ)

```cmd
start-services.bat
```

Остановка:
```cmd
stop-services.bat
```

## Ручной запуск

### 1. Запуск баз данных

```bash
docker-compose up -d
```

### 2. Установка зависимостей (первый раз)

```bash
# Product Service
cd product-service
go mod tidy

# Order Service
cd ../order-service
go mod tidy

# Inventory Service
cd ../inventory-service
go mod tidy

# Notification Service
cd ../notification-service
go mod tidy
```

### 3. Запуск сервисов

Откройте 4 терминала:

**Терминал 1 - Product Service:**
```bash
cd product-service
go run .
```

**Терминал 2 - Order Service:**
```bash
cd order-service
go run .
```

**Терминал 3 - Inventory Service:**
```bash
cd inventory-service
go run .
```

**Терминал 4 - Notification Service:**
```bash
cd notification-service
go run .
```

## API Endpoints

### Product Service (8081)

```bash
# Создать продукт
curl -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{"code":"P001","name":"Laptop","weight":2.5,"description":"Gaming laptop"}'

# Получить все продукты
curl http://localhost:8081/products

# Получить продукт по коду
curl http://localhost:8081/products/P001

# Удалить продукт
curl -X DELETE http://localhost:8081/products/P001
```

### Order Service (8082)

```bash
# Создать заказ
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{
    "order_code":"O001",
    "order_date":"2024-01-01",
    "product_items":[
      {"product_code":"P001","name":"Laptop","quantity":2,"cost":1500.00}
    ]
  }'

# Получить все заказы
curl http://localhost:8082/orders

# Получить заказ по коду
curl http://localhost:8082/orders/O001

# Удалить заказ
curl -X DELETE http://localhost:8082/orders/O001
```

### Inventory Service (8083)

```bash
# Добавить товар на склад
curl -X POST http://localhost:8083/inventory \
  -H "Content-Type: application/json" \
  -d '{"product_code":"P001","name":"Laptop","quantity":50,"current_price":"1500.00"}'

# Получить весь склад
curl http://localhost:8083/inventory

# Получить товар по коду
curl http://localhost:8083/inventory/P001

# Обновить товар на складе
curl -X PUT http://localhost:8083/inventory/P001 \
  -H "Content-Type: application/json" \
  -d '{"quantity":45,"current_price":"1450.00"}'
```

### Notification Service (8084)

```bash
# Создать уведомление
curl -X POST http://localhost:8084/notifications \
  -H "Content-Type: application/json" \
  -d '{"message_type":"INFO","description":"Order created","date":"2024-01-01"}'

# Получить все уведомления
curl http://localhost:8084/notifications
```

## Тестирование

### Интеграционные тесты (рекомендуется)

После запуска всех сервисов запустите интеграционные тесты:

```cmd
run-tests.bat
```

**Результат:** Отчёт с количеством пройденных/проваленных тестов и процентом успеха

### Быстрые тесты API

Для быстрой проверки можно использовать:

```powershell
.\test-api.ps1
```

Этот скрипт создаст тестовые данные во всех сервисах и выведет результаты.