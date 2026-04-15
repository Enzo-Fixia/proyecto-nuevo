package order

import (
	"github.com/fixia/golang-api/internal/product"
	"github.com/fixia/golang-api/utils"
	"gorm.io/gorm"
)

type Service interface {
	Create(req CreateOrderRequest, userID uint) (*Order, error)
	GetByID(id uint) (*Order, error)
	ListByUser(userID uint) ([]Order, error)
	UpdateStatus(id uint, status string) (*Order, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req CreateOrderRequest, userID uint) (*Order, error) {
	var created *Order

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var total float64
		items := make([]OrderItem, 0, len(req.Items))

		for _, it := range req.Items {
			var p product.Product
			if err := tx.First(&p, it.ProductID).Error; err != nil {
				return err
			}
			if p.Stock < it.Quantity {
				return utils.ErrInsufficientStock
			}

			p.Stock -= it.Quantity
			if err := tx.Save(&p).Error; err != nil {
				return err
			}

			items = append(items, OrderItem{
				ProductID: p.ID,
				Quantity:  it.Quantity,
				UnitPrice: p.Price,
			})
			total += p.Price * float64(it.Quantity)
		}

		o := &Order{
			UserID: userID,
			Total:  total,
			Status: StatusPending,
			Items:  items,
		}
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		created = o
		return nil
	})

	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *service) GetByID(id uint) (*Order, error) { return s.repo.FindByID(id) }

func (s *service) ListByUser(userID uint) ([]Order, error) { return s.repo.FindByUser(userID) }

func (s *service) UpdateStatus(id uint, status string) (*Order, error) {
	return s.repo.UpdateStatus(id, status)
}
