# Lab2

## Сервисы

- **Product Service** (порт 8081) - PostgreSQL
- **Order Service** (порт 8082) - MongoDB
- **Inventory Service** (порт 8083) - PostgreSQL
- **Notification Service** (порт 8084) - MongoDB

## Запуск

```cmd
start-services.bat
```

Остановка:
```cmd
stop-services.bat
```

## API

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

### Интеграционные тесты

```cmd
run-tests.bat
```

### Быстрые тесты API

```powershell
.\test-api.ps1
```