package main

import "fmt"

// ErrOrderNotFound возвращается когда заказ не найден
type ErrOrderNotFound struct {
	OrderID int
}

func (e ErrOrderNotFound) Error() string {
	return fmt.Sprintf("заказ с ID %d не найден", e.OrderID)
}

// ErrTaskNotFound возвращается когда задача не найдена
type ErrTaskNotFound struct {
	OrderID int
	TaskID  int
}

func (e ErrTaskNotFound) Error() string {
	return fmt.Sprintf("задача с ID %d не найдена в заказе %d", e.TaskID, e.OrderID)
}

// ErrInvalidInput возвращается при невалидных входных данных
type ErrInvalidInput struct {
	Field   string
	Message string
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("невалидное поле %s: %s", e.Field, e.Message)
}

