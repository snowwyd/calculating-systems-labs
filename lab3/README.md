# Lab3
## Архитектура

### Сервисы

1. **Inventory Service** (порт 8083)

2. **Order Service** (порт 8082)

3. **Notification Service**

### Инфраструктура

- **RabbitMQ** (порты 5672, 15672)

## Быстрый запуск

### 1. Запуск всех сервисов

```cmd
start-services.bat
```

### 2. Тестирование сценария

```cmd
test-scenario.bat
```
### 3. Остановка

```cmd
stop-services.bat
```


## RabbitMQ

Откройте http://localhost:15672

- Логин: `admin`
- Пароль: `password`