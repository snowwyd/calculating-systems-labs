package main

// Order представляет задачу проекта
type Order struct {
	ID        int    `json:"id"`
	OrderName string `json:"order_name"`
	StartDate string `json:"start_date"`
	Tasks     []Task `json:"tasks"`
}

// Task представляет отдельную работу
type Task struct {
	ID       int    `json:"id"`
	TaskName string `json:"task"`
	Duration int    `json:"duration"`
	Resource int    `json:"resource"`
	Pred     []int  `json:"pred"`
}

// Result содержит результат расчёта для одной последовательности
type Result struct {
	Sequence []int
	Duration int
}

// TaskExecution представляет план выполнения задачи
type TaskExecution struct {
	StartTime int
	EndTime   int
	Resource  int
}
