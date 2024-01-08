package ports

import (
	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/omise/omise-go"
)

type UserRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	UpdateUser(existingUser *domain.User, user *domain.User) (*domain.User, error)
	DeleteUser(id string) error
}

type UserService interface {
	CreateUser(user *domain.User) (*domain.User, *omise.Customer, error)
	GetUserByID(id string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	UpdateUser(user *domain.User, id string) (*domain.User, error)
	DeleteUser(id string) error
}
