package repository

import (
	"fmt"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (r *UserRepository) GetUserByID(id string) (*domain.User, error) {

	user := &domain.User{}
	if err := r.db.First(user, id).Error; err != nil {
		fmt.Println("Error while getting user", err)
		return nil, err
	}

	return user, nil
}
