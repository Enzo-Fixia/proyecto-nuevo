package product

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Stock       int     `gorm:"default:0" json:"stock"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	UserID      uint    `json:"user_id"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"        binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"       binding:"required,gt=0"`
	Stock       int     `json:"stock"       binding:"min=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"       binding:"omitempty,gt=0"`
	Stock       *int     `json:"stock"       binding:"omitempty,min=0"`
	Category    *string  `json:"category"`
	ImageURL    *string  `json:"image_url"`
}
