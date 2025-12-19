# Как рассказывать

Скажи, что у тебя тут 4 микросервиса (перечисли все папки с подписью `service`)

Дальше скажи, что в каждом сервисе ты сделал структуру такую же, как в 1 лабе. Открой `inventory-service`, в нем так же есть
main
models
repository

вот тут остановись и открой ему файл `repository.go`
читай коменты 
и расскажи про вот эту штучку:
```go
func (r *PostgreSQLInventoryRepository) Create(inventory *Inventory) error {
	query := "INSERT INTO inventory (product_code, name, quantity, current_price) VALUES ($1, $2, $3, $4)" // эту
	_, err := r.db.Exec(query, inventory.ProductCode, inventory.Name, inventory.Quantity, inventory.CurrentPrice)
	if err != nil {
		return ErrDatabaseError{Operation: "создание товара на складе", Err: err}
	}
	return nil
}
```

эти значки доллара называются плейсхолдерами, они защищают нашу программу от SQL-инъекций

**ОБЯЗАТЕЛЬНО РАССКАЖИ ПРО ЭТО**
```go
// InitDatabase инициализирует базу данных
func InitDatabase(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS inventory (
		product_code VARCHAR PRIMARY KEY,
		name VARCHAR NOT NULL,
		quantity NUMERIC NOT NULL,
		current_price VARCHAR NOT NULL
	)`

	if _, err := db.Exec(createTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	return db, nil
}

```
Вот этот код у тебя подключает твой сервис к Postgres по `connStr` (адрес подключения), а затем SQL-запросом у нас создается сразу же структура нашей БД

Затем покажи ему файл `docker-compose.yml` - это yaml-файл, по которому docker понимает, какие образы нам нужно создать. Зачем это нужно - чтобы мы не скачивали себе на комп БД, а быстренько развернули все, что нам нужно для работы сервиса

# Как показывать

1. Открываешь в терминале директорию, скорее всего `cd ../lab2`
2. Открой ещё 4 терминала через `+`. У тебя в итоге должно быть открыто 5 терминалов в `cd lab2` 
3. В первом терминале пиши команду `docker-compose up -d`. Это у тебя выполнит все то, что прописано в `docker-compose.yaml` - а именно поднимет две базы данных - **mongo** (для `order service`) и **postgres** - для `inventory`. Ты увидишь их в Docker
4. Дождись, как выполнится шаг 3
5. Во втором терминале `cd lab2/inventory-service` и пишешь `go run .` - запустит сервис
6. В третьем терминале `cd lab2/notification-service` и пишешь `go run .`
7. В четвертом терминале `cd lab2/order-service` и пишешь `go run .`
8. В пятом терминале `cd lab2/product-service` и пишешь `go run .`
9. Все, у тебя запустились все сервисы. Для наглядности просто можешь ему запустить интеграционный тест между всеми сервисами. Скажи, что попросил нейронку написать скрипт для этого теста, чтобы наглядно преподу все показать
10. Скрипт запускается в терминале, открытом в `lab2` командой `run-tests.bat`
11. Потом закрой все сервисы: `Ctrl + C` в 4х терминалах с сервисами. И потом на крестик убей эти терминалы
12. Почисти докер через `docker-compose down -v`