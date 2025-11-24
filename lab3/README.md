# Lab 3: Микросервисы с синхронной и асинхронной коммуникацией

## Архитектура

### Сервисы

1. **Inventory Service** (порт 8083)
   - REST API для управления складом
   - Проверка наличия товаров
   - Резервирование товаров

2. **Order Service** (порт 8082)
   - REST API для создания заказов
   - **Синхронная коммуникация**: вызывает Inventory Service по REST API для проверки товаров
   - **Асинхронная коммуникация**: отправляет уведомления в RabbitMQ после создания заказа

3. **Notification Service**
   - Подписан на RabbitMQ
   - Получает и обрабатывает уведомления о заказах асинхронно

### Инфраструктура

- **RabbitMQ** (порты 5672, 15672) - брокер сообщений для асинхронной коммуникации

## Схема взаимодействия

```
Клиент
  ↓
  ↓ POST /orders
  ↓
Order Service
  ↓
  ├─→ (Синхронно) POST /inventory/check → Inventory Service
  ├─→ (Синхронно) POST /inventory/reserve → Inventory Service
  └─→ (Асинхронно) RabbitMQ → Notification Service
```

## Быстрый запуск

### 1. Запуск всех сервисов

```cmd
start-services.bat
```

Это запустит:
- RabbitMQ в Docker
- Inventory Service
- Order Service
- Notification Service

### 2. Тестирование сценария

```cmd
test-scenario.bat
```

Скрипт автоматически проверит доступность сервисов перед запуском тестов.

### 3. Отладка (если возникают проблемы)

```cmd
test-debug.bat
```

Этот скрипт покажет детальные логи запросов.

### 4. Остановка

```cmd
stop-services.bat
```

## Ручное тестирование

### Проверка склада

```bash
curl http://localhost:8083/inventory
```

### Создание заказа

```bash
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "CUST001",
    "items": [
      {
        "product_code": "PROD001",
        "product_name": "Ноутбук",
        "quantity": 2,
        "price": 50000
      }
    ]
  }'
```

**Что происходит:**
1. Order Service **синхронно** проверяет наличие товара в Inventory Service
2. Если товар доступен, резервирует его
3. Создаёт заказ
4. **Асинхронно** отправляет уведомление в RabbitMQ
5. Notification Service получает уведомление и выводит информацию в консоль

### Просмотр заказов

```bash
curl http://localhost:8082/orders
```

### Просмотр остатков

```bash
curl http://localhost:8083/inventory
```

## RabbitMQ Management UI

Откройте http://localhost:15672

- Логин: `admin`
- Пароль: `password`

Здесь можно посмотреть очереди, сообщения и статистику.

## Особенности реализации

### Синхронная коммуникация (REST API)

Order Service делает HTTP запросы к Inventory Service:
- `POST /inventory/check` - проверка наличия товара
- `POST /inventory/reserve` - резервирование товара

**Клиент ждёт** ответа от обоих сервисов перед получением результата.

### Асинхронная коммуникация (RabbitMQ)

Order Service отправляет сообщение в очередь RabbitMQ после создания заказа.
Notification Service **независимо** получает и обрабатывает эти сообщения.

**Клиент не ждёт** обработки уведомления.

## Устранение проблем

### Зависает при создании заказа?

1. Убедитесь, что **все сервисы запущены**:
   - Inventory Service (8083)
   - Order Service (8082)
   - Notification Service
   - RabbitMQ

2. Запустите отладочный скрипт:
   ```cmd
   test-debug.bat
   ```

3. Проверьте логи в окнах сервисов

4. Убедитесь, что RabbitMQ запущен:
   ```cmd
   docker ps
   ```

### Порты заняты?

```cmd
netstat -ano | findstr "8082"
netstat -ano | findstr "8083"
```

Остановите процессы:
```cmd
stop-services.bat
```

## Установка зависимостей

```bash
# Inventory Service
cd inventory-service
go mod tidy

# Order Service
cd ../order-service
go mod tidy

# Notification Service
cd ../notification-service
go mod tidy
```
