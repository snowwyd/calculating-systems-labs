package main

import "fmt"

// ErrProductNotFound возвращается когда продукт не найден
type ErrProductNotFound struct {
	Code string
}

func (e ErrProductNotFound) Error() string {
	return fmt.Sprintf("продукт с кодом %s не найден", e.Code)
}

// ErrInvalidInput возвращается при невалидных входных данных
type ErrInvalidInput struct {
	Field   string
	Message string
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("невалидное поле %s: %s", e.Field, e.Message)
}

// ErrDatabaseError возвращается при ошибках базы данных
type ErrDatabaseError struct {
	Operation string
	Err       error
}

func (e ErrDatabaseError) Error() string {
	return fmt.Sprintf("ошибка базы данных при %s: %v", e.Operation, e.Err)
}

func (e ErrDatabaseError) Unwrap() error {
	return e.Err
}

