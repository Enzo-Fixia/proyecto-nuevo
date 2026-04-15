package user

import (
	"errors"

	"github.com/fixia/golang-api/internal/auth"
	"github.com/fixia/golang-api/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(req RegisterRequest) (*User, error)
	Login(req LoginRequest) (*LoginResponse, error)
	GetByID(id uint) (*User, error)
	ListAll() ([]User, error)
	Update(id uint, req UpdateUserRequest) (*User, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(req RegisterRequest) (*User, error) {
	if existing, err := s.repo.FindByEmail(req.Email); err == nil && existing != nil {
		return nil, utils.ErrDuplicateEmail
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	u := &User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hash),
		Role:      "user",
		IsActive:  true,
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Login(req LoginRequest) (*LoginResponse, error) {
	u, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, err
	}

	if !u.IsActive {
		return nil, utils.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(u.ID, u.Email, u.Role)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{Token: token, User: *u}, nil
}

func (s *service) GetByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) ListAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *service) Update(id uint, req UpdateUserRequest) (*User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	if req.Role != nil {
		u.Role = *req.Role
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}
