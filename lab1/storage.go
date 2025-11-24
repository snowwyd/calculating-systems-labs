package main

import "sync"

// Storage управляет in-memory хранилищем данных
type Storage struct {
	orders  map[int]*Order
	orderID int
	taskID  int
	mu      sync.RWMutex
}

// NewStorage создаёт новое хранилище
func NewStorage() *Storage {
	return &Storage{
		orders: make(map[int]*Order),
	}
}

// CreateOrder создаёт новый заказ
func (s *Storage) CreateOrder(orderName, startDate string) *Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orderID++
	order := &Order{
		ID:        s.orderID,
		OrderName: orderName,
		StartDate: startDate,
		Tasks:     []Task{},
	}
	s.orders[order.ID] = order

	return order
}

// GetOrder возвращает заказ по ID
func (s *Storage) GetOrder(id int) (*Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	return order, exists
}

// GetAllOrders возвращает все заказы
func (s *Storage) GetAllOrders() []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}

// UpdateOrder обновляет информацию о заказе
func (s *Storage) UpdateOrder(id int, orderName, startDate string) (*Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, false
	}

	order.OrderName = orderName
	order.StartDate = startDate
	return order, true
}

// DeleteOrder удаляет заказ
func (s *Storage) DeleteOrder(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[id]; !exists {
		return false
	}
	delete(s.orders, id)
	return true
}

// AddTask добавляет работу к заказу
func (s *Storage) AddTask(orderID int, taskName string, duration, resource int, pred []int) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderID]
	if !exists {
		return nil, false
	}

	s.taskID++
	task := Task{
		ID:       s.taskID,
		TaskName: taskName,
		Duration: duration,
		Resource: resource,
		Pred:     pred,
	}

	if task.Pred == nil {
		task.Pred = []int{}
	}

	order.Tasks = append(order.Tasks, task)
	return &task, true
}

// DeleteTask удаляет работу из заказа
func (s *Storage) DeleteTask(orderID, taskID int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderID]
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

// UpdateTaskPredecessors обновляет предшествующие работы
func (s *Storage) UpdateTaskPredecessors(orderID, taskID int, pred []int) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderID]
	if !exists {
		return nil, false
	}

	for i := range order.Tasks {
		if order.Tasks[i].ID == taskID {
			order.Tasks[i].Pred = pred
			return &order.Tasks[i], true
		}
	}
	return nil, false
}
