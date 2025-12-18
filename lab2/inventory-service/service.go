package main

// InventoryService предоставляет бизнес-логику для работы со складом
type InventoryService struct {
	repo InventoryRepository
}

// NewInventoryService создаёт новый сервис склада
func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// CreateInventory создаёт новый товар на складе
func (s *InventoryService) CreateInventory(req CreateInventoryRequest) (*Inventory, error) {
	inventory := &Inventory{
		ProductCode:  req.ProductCode,
		Name:         req.Name,
		Quantity:     req.Quantity,
		CurrentPrice: req.CurrentPrice,
	}

	if err := s.repo.Create(inventory); err != nil {
		return nil, err
	}

	return inventory, nil
}

// GetInventory получает товар по коду
func (s *InventoryService) GetInventory(code string) (*Inventory, error) {
	return s.repo.FindByCode(code)
}

// GetAllInventory получает все товары на складе
func (s *InventoryService) GetAllInventory() ([]*Inventory, error) {
	return s.repo.FindAll()
}

// UpdateInventory обновляет товар на складе
func (s *InventoryService) UpdateInventory(code string, req UpdateInventoryRequest) (*Inventory, error) {
	// Получаем существующий товар для сохранения имени
	existing, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	inventory := &Inventory{
		ProductCode:  code,
		Name:         existing.Name,
		Quantity:     req.Quantity,
		CurrentPrice: req.CurrentPrice,
	}

	if err := s.repo.Update(code, inventory); err != nil {
		return nil, err
	}

	return inventory, nil
}

