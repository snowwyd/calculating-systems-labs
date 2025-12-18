package main

import (
	"database/sql"
	"fmt"
)

// ProductRepository определяет интерфейс для работы с продуктами
type ProductRepository interface {
	Create(product *Product) error
	FindByCode(code string) (*Product, error)
	FindAll() ([]*Product, error)
	Delete(code string) error
}

// PostgreSQLProductRepository реализует репозиторий для PostgreSQL
type PostgreSQLProductRepository struct {
	db *sql.DB
}

// NewPostgreSQLProductRepository создаёт новый репозиторий
func NewPostgreSQLProductRepository(db *sql.DB) *PostgreSQLProductRepository {
	return &PostgreSQLProductRepository{db: db}
}

// Create создаёт новый продукт
func (r *PostgreSQLProductRepository) Create(product *Product) error {
	query := "INSERT INTO products (code, name, weight, description) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, product.Code, product.Name, product.Weight, product.Description)
	if err != nil {
		return ErrDatabaseError{Operation: "создание продукта", Err: err}
	}
	return nil
}

// FindByCode находит продукт по коду
func (r *PostgreSQLProductRepository) FindByCode(code string) (*Product, error) {
	var product Product
	query := "SELECT code, name, weight, description FROM products WHERE code = $1"
	err := r.db.QueryRow(query, code).Scan(&product.Code, &product.Name, &product.Weight, &product.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound{Code: code}
		}
		return nil, ErrDatabaseError{Operation: "поиск продукта", Err: err}
	}
	return &product, nil
}

// FindAll возвращает все продукты
func (r *PostgreSQLProductRepository) FindAll() ([]*Product, error) {
	query := "SELECT code, name, weight, description FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, ErrDatabaseError{Operation: "получение всех продуктов", Err: err}
	}
	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Code, &product.Name, &product.Weight, &product.Description); err != nil {
			return nil, ErrDatabaseError{Operation: "сканирование продукта", Err: err}
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrDatabaseError{Operation: "итерация по продуктам", Err: err}
	}

	return products, nil
}

// Delete удаляет продукт
func (r *PostgreSQLProductRepository) Delete(code string) error {
	query := "DELETE FROM products WHERE code = $1"
	result, err := r.db.Exec(query, code)
	if err != nil {
		return ErrDatabaseError{Operation: "удаление продукта", Err: err}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrDatabaseError{Operation: "проверка удаления", Err: err}
	}

	if rowsAffected == 0 {
		return ErrProductNotFound{Code: code}
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
	CREATE TABLE IF NOT EXISTS products (
		code VARCHAR PRIMARY KEY,
		name VARCHAR NOT NULL,
		weight NUMERIC NOT NULL,
		description VARCHAR NOT NULL
	)`
	
	if _, err := db.Exec(createTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	return db, nil
}

