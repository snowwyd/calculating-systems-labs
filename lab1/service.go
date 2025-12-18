package main

// OrderService предоставляет бизнес-логику для работы с заказами
type OrderService struct {
	orderRepo OrderRepository
	taskRepo  TaskRepository
}

// NewOrderService создаёт новый сервис заказов
func NewOrderService(orderRepo OrderRepository, taskRepo TaskRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		taskRepo:  taskRepo,
	}
}

// CreateOrder создаёт новый заказ
func (s *OrderService) CreateOrder(req CreateOrderRequest) (*Order, error) {
	order := &Order{
		OrderName: req.OrderName,
		StartDate: req.StartDate,
	}
	return s.orderRepo.Create(order), nil
}

// GetOrder получает заказ по ID
func (s *OrderService) GetOrder(id int) (*Order, error) {
	order, exists := s.orderRepo.FindByID(id)
	if !exists {
		return nil, ErrOrderNotFound{OrderID: id}
	}
	return order, nil
}

// GetAllOrders получает все заказы
func (s *OrderService) GetAllOrders() []*Order {
	return s.orderRepo.FindAll()
}

// UpdateOrder обновляет заказ
func (s *OrderService) UpdateOrder(id int, req UpdateOrderRequest) (*Order, error) {
	order, exists := s.orderRepo.FindByID(id)
	if !exists {
		return nil, ErrOrderNotFound{OrderID: id}
	}

	order.OrderName = req.OrderName
	order.StartDate = req.StartDate

	if !s.orderRepo.Update(order) {
		return nil, ErrOrderNotFound{OrderID: id}
	}

	return order, nil
}

// DeleteOrder удаляет заказ
func (s *OrderService) DeleteOrder(id int) error {
	if !s.orderRepo.Delete(id) {
		return ErrOrderNotFound{OrderID: id}
	}
	return nil
}

// AddTask добавляет задачу к заказу
func (s *OrderService) AddTask(orderID int, req AddTaskRequest) (*Task, error) {
	_, exists := s.orderRepo.FindByID(orderID)
	if !exists {
		return nil, ErrOrderNotFound{OrderID: orderID}
	}

	task := &Task{
		TaskName: req.Task,
		Duration: req.Duration,
		Resource: req.Resource,
		Pred:     req.Pred,
	}

	if !s.taskRepo.AddToOrder(orderID, task) {
		return nil, ErrOrderNotFound{OrderID: orderID}
	}

	// Задача была модифицирована в репозитории (получила ID), возвращаем её
	return task, nil
}

// DeleteTask удаляет задачу из заказа
func (s *OrderService) DeleteTask(orderID, taskID int) error {
	_, exists := s.orderRepo.FindByID(orderID)
	if !exists {
		return ErrOrderNotFound{OrderID: orderID}
	}

	if !s.taskRepo.RemoveFromOrder(orderID, taskID) {
		return ErrTaskNotFound{OrderID: orderID, TaskID: taskID}
	}

	return nil
}

// UpdateTaskPredecessors обновляет предшественников задачи
func (s *OrderService) UpdateTaskPredecessors(orderID, taskID int, req UpdatePredecessorsRequest) (*Task, error) {
	_, exists := s.orderRepo.FindByID(orderID)
	if !exists {
		return nil, ErrOrderNotFound{OrderID: orderID}
	}

	task, exists := s.taskRepo.UpdatePredecessors(orderID, taskID, req.Pred)
	if !exists {
		return nil, ErrTaskNotFound{OrderID: orderID, TaskID: taskID}
	}

	return task, nil
}
