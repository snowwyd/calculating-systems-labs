package main

// Order представляет задачу проекта
type Order struct {
	ID        int    `json:"id"`
	OrderName string `json:"order_name"`
	StartDate string `json:"start_date"`
	Tasks     []Task `json:"tasks"`
}

// Task представляет отдельную работу в рамках задачи
type Task struct {
	ID       int    `json:"id"`
	TaskName string `json:"task"`
	Duration int    `json:"duration"`
	Resource int    `json:"resource"`
	Pred     []int  `json:"pred"`
}
