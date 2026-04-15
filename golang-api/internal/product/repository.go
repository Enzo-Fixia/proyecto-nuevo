package product

import "gorm.io/gorm"

type Repository interface {
	Create(p *Product) error
	FindByID(id uint) (*Product, error)
	FindAll() ([]Product, error)
	Update(p *Product) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(p *Product) error {
	return r.db.Create(p).Error
}

func (r *repository) FindByID(id uint) (*Product, error) {
	var p Product
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindAll() ([]Product, error) {
	var ps []Product
	if err := r.db.Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (r *repository) Update(p *Product) error {
	return r.db.Save(p).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}
