package utils

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrDuplicateEmail     = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrInvalidInput       = errors.New("invalid input")
)
