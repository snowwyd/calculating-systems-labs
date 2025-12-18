package main

import "sync"

// OrderRepository определяет интерфейс для работы с заказами
type OrderRepository interface {
	Create(order *Order) *Order
	FindByID(id int) (*Order, bool)
	FindAll() []*Order
	Update(order *Order) bool
	Delete(id int) bool
}

// TaskRepository определяет интерфейс для работы с задачами
type TaskRepository interface {
	AddToOrder(orderID int, task *Task) bool
	RemoveFromOrder(orderID, taskID int) bool
	UpdatePredecessors(orderID, taskID int, predecessors []int) (*Task, bool)
}

// InMemoryRepository реализует хранилище в памяти
type InMemoryRepository struct {
	orders  map[int]*Order
	orderID int
	taskID  int
	mu      sync.RWMutex
}

// NewInMemoryRepository создаёт новое хранилище в памяти
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		orders: make(map[int]*Order),
	}
}

// Create создаёт новый заказ
func (r *InMemoryRepository) Create(order *Order) *Order {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orderID++
	order.ID = r.orderID
	order.Tasks = make([]Task, 0)
	r.orders[order.ID] = order
	return order
}

// FindByID находит заказ по ID
func (r *InMemoryRepository) FindByID(id int) (*Order, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	return order, exists
}

// FindAll возвращает все заказы
func (r *InMemoryRepository) FindAll() []*Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Order, 0, len(r.orders))
	for _, order := range r.orders {
		result = append(result, order)
	}
	return result
}

// Update обновляет заказ
func (r *InMemoryRepository) Update(order *Order) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return false
	}
	r.orders[order.ID] = order
	return true
}

// Delete удаляет заказ
func (r *InMemoryRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[id]; !exists {
		return false
	}
	delete(r.orders, id)
	return true
}

// AddToOrder добавляет задачу к заказу
func (r *InMemoryRepository) AddToOrder(orderID int, task *Task) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[orderID]
	if !exists {
		return false
	}

	r.taskID++
	task.ID = r.taskID
	if task.Pred == nil {
		task.Pred = make([]int, 0)
	}

	order.Tasks = append(order.Tasks, *task)
	return true
}

// RemoveFromOrder удаляет задачу из заказа
func (r *InMemoryRepository) RemoveFromOrder(orderID, taskID int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[orderID]
	if !exists {
		return false
	}

	for i, task := range order.Tasks {
		if task.ID == taskID {
			order.Tasks = append(order.Tasks[:i], order.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

// UpdatePredecessors обновляет предшественников задачи
func (r *InMemoryRepository) UpdatePredecessors(orderID, taskID int, predecessors []int) (*Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[orderID]
	if !exists {
		return nil, false
	}

	for i := range order.Tasks {
		if order.Tasks[i].ID == taskID {
			order.Tasks[i].Pred = predecessors
			return &order.Tasks[i], true
		}
	}
	return nil, false
}
