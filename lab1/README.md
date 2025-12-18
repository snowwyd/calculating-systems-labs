# Lab1

## API

### Задачи (Orders)

- `POST /orders` - Создать задачу
- `GET /orders` - Получить все задачи
- `GET /orders/{id}` - Получить задачу по ID
- `PUT /orders/{id}` - Обновить задачу
- `DELETE /orders/{id}` - Удалить задачу

### Работы (Tasks)

- `POST /orders/{id}/tasks` - Добавить работу к задаче
- `DELETE /orders/{orderId}/tasks/{taskId}` - Удалить работу
- `PUT /orders/{orderId}/tasks/{taskId}/predecessors` - Обновить предшествующие работы

## Запуск

### 1. Установка зависимостей

```bash
go mod tidy
```

### 2. Запуск API сервера

В первом терминале:

```bash
go run .
```

Сервер запустится на `http://localhost:8080`

### 3. Запуск вычислительного сервиса

Во втором терминале (после запуска API сервера):

```bash
cd client
go run .
```