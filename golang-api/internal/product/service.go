package product

type Service interface {
	Create(req CreateProductRequest, userID uint) (*Product, error)
	GetByID(id uint) (*Product, error)
	List() ([]Product, error)
	Update(id uint, req UpdateProductRequest) (*Product, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req CreateProductRequest, userID uint) (*Product, error) {
	p := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		UserID:      userID,
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetByID(id uint) (*Product, error) { return s.repo.FindByID(id) }

func (s *service) List() ([]Product, error) { return s.repo.FindAll() }

func (s *service) Update(id uint, req UpdateProductRequest) (*Product, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	if req.Category != nil {
		p.Category = *req.Category
	}
	if req.ImageURL != nil {
		p.ImageURL = *req.ImageURL
	}
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) Delete(id uint) error { return s.repo.Delete(id) }
