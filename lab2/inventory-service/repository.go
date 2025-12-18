package main

import (
	"database/sql"
	"fmt"
)

// InventoryRepository определяет интерфейс для работы со складом
type InventoryRepository interface {
	Create(inventory *Inventory) error
	FindByCode(code string) (*Inventory, error)
	FindAll() ([]*Inventory, error)
	Update(code string, inventory *Inventory) error
}

// PostgreSQLInventoryRepository реализует репозиторий для PostgreSQL
type PostgreSQLInventoryRepository struct {
	db *sql.DB
}

// NewPostgreSQLInventoryRepository создаёт новый репозиторий
func NewPostgreSQLInventoryRepository(db *sql.DB) *PostgreSQLInventoryRepository {
	return &PostgreSQLInventoryRepository{db: db}
}

// Create создаёт новый товар на складе
func (r *PostgreSQLInventoryRepository) Create(inventory *Inventory) error {
	query := "INSERT INTO inventory (product_code, name, quantity, current_price) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, inventory.ProductCode, inventory.Name, inventory.Quantity, inventory.CurrentPrice)
	if err != nil {
		return ErrDatabaseError{Operation: "создание товара на складе", Err: err}
	}
	return nil
}

// FindByCode находит товар по коду
func (r *PostgreSQLInventoryRepository) FindByCode(code string) (*Inventory, error) {
	var inventory Inventory
	query := "SELECT product_code, name, quantity, current_price FROM inventory WHERE product_code = $1"
	err := r.db.QueryRow(query, code).Scan(&inventory.ProductCode, &inventory.Name, &inventory.Quantity, &inventory.CurrentPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInventoryNotFound{Code: code}
		}
		return nil, ErrDatabaseError{Operation: "поиск товара на складе", Err: err}
	}
	return &inventory, nil
}

// FindAll возвращает все товары на складе
func (r *PostgreSQLInventoryRepository) FindAll() ([]*Inventory, error) {
	query := "SELECT product_code, name, quantity, current_price FROM inventory"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, ErrDatabaseError{Operation: "получение всех товаров на складе", Err: err}
	}
	defer rows.Close()

	inventory := []*Inventory{}
	for rows.Next() {
		var inv Inventory
		if err := rows.Scan(&inv.ProductCode, &inv.Name, &inv.Quantity, &inv.CurrentPrice); err != nil {
			return nil, ErrDatabaseError{Operation: "сканирование товара", Err: err}
		}
		inventory = append(inventory, &inv)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrDatabaseError{Operation: "итерация по товарам", Err: err}
	}

	return inventory, nil
}

// Update обновляет товар на складе
func (r *PostgreSQLInventoryRepository) Update(code string, inventory *Inventory) error {
	query := "UPDATE inventory SET quantity = $1, current_price = $2 WHERE product_code = $3"
	result, err := r.db.Exec(query, inventory.Quantity, inventory.CurrentPrice, code)
	if err != nil {
		return ErrDatabaseError{Operation: "обновление товара на складе", Err: err}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrDatabaseError{Operation: "проверка обновления", Err: err}
	}

	if rowsAffected == 0 {
		return ErrInventoryNotFound{Code: code}
	}

	return nil
}

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

