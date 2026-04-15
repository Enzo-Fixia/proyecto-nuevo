package user

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]User, error)
	Update(user *User) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *repository) FindByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindAll() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) Update(u *User) error {
	return r.db.Save(u).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
