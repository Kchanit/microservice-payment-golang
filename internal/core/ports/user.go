package ports

import "github.com/Kchanit/microservice-payment-golang/internal/core/domain"

type UserRepository interface {
	GetUserByID(id string) (*domain.User, error)
}

type UserService interface {
	GetUser(id string) (*domain.User, error)
}
