package main

import "context"

// OrderService предоставляет бизнес-логику для работы с заказами
type OrderService struct {
	repo OrderRepository
}

// NewOrderService создаёт новый сервис заказов
func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// CreateOrder создаёт новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	order := &Order{
		OrderCode:    req.OrderCode,
		OrderDate:    req.OrderDate,
		ProductItems: req.ProductItems,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder получает заказ по коду
func (s *OrderService) GetOrder(ctx context.Context, code string) (*Order, error) {
	return s.repo.FindByCode(ctx, code)
}

// GetAllOrders получает все заказы
func (s *OrderService) GetAllOrders(ctx context.Context) ([]*Order, error) {
	return s.repo.FindAll(ctx)
}

// DeleteOrder удаляет заказ
func (s *OrderService) DeleteOrder(ctx context.Context, code string) error {
	return s.repo.Delete(ctx, code)
}

