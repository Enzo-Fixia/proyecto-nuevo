package order

import "gorm.io/gorm"

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusShipped   = "shipped"
	StatusDelivered = "delivered"
	StatusCancelled = "cancelled"
)

type Order struct {
	gorm.Model
	UserID uint        `gorm:"not null" json:"user_id"`
	Total  float64     `gorm:"not null" json:"total"`
	Status string      `gorm:"default:'pending'" json:"status"`
	Items  []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"not null" json:"unit_price"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity"   binding:"required,min=1"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}
