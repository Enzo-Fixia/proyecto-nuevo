package order

import "gorm.io/gorm"

type Repository interface {
	DB() *gorm.DB
	FindByID(id uint) (*Order, error)
	FindByUser(userID uint) ([]Order, error)
	UpdateStatus(id uint, status string) (*Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) DB() *gorm.DB { return r.db }

func (r *repository) FindByID(id uint) (*Order, error) {
	var o Order
	if err := r.db.Preload("Items").First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repository) FindByUser(userID uint) ([]Order, error) {
	var orders []Order
	if err := r.db.Preload("Items").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repository) UpdateStatus(id uint, status string) (*Order, error) {
	var o Order
	if err := r.db.First(&o, id).Error; err != nil {
		return nil, err
	}
	o.Status = status
	if err := r.db.Save(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}
