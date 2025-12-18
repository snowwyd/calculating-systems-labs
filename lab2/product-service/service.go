package main

// ProductService предоставляет бизнес-логику для работы с продуктами
type ProductService struct {
	repo ProductRepository
}

// NewProductService создаёт новый сервис продуктов
func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct создаёт новый продукт
func (s *ProductService) CreateProduct(req CreateProductRequest) (*Product, error) {
	product := &Product{
		Code:        req.Code,
		Name:        req.Name,
		Weight:      req.Weight,
		Description: req.Description,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct получает продукт по коду
func (s *ProductService) GetProduct(code string) (*Product, error) {
	return s.repo.FindByCode(code)
}

// GetAllProducts получает все продукты
func (s *ProductService) GetAllProducts() ([]*Product, error) {
	return s.repo.FindAll()
}

// DeleteProduct удаляет продукт
func (s *ProductService) DeleteProduct(code string) error {
	return s.repo.Delete(code)
}

