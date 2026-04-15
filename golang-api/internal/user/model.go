package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	Role      string `gorm:"default:'user'" json:"role"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name"  binding:"required,min=2,max=50"`
	Email     string `json:"email"      binding:"required,email"`
	Password  string `json:"password"   binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name"  binding:"omitempty,min=2,max=50"`
	Role      *string `json:"role"       binding:"omitempty,oneof=user admin"`
	IsActive  *bool   `json:"is_active"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
