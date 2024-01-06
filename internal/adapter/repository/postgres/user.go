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
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {

	if err := r.db.Create(user).Error; err != nil {
		fmt.Println("Error while creating user", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]*domain.User, error) {

	users := []*domain.User{}
	if err := r.db.Find(&users).Error; err != nil {
		fmt.Println("Error while getting all users", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(existingUser *domain.User, user *domain.User) (*domain.User, error) {

	if err := r.db.Model(existingUser).Updates(user).Error; err != nil {
		fmt.Println("Error while updating user", err)
		return nil, err
	}

	return existingUser, nil
}

func (r *UserRepository) DeleteUser(id string) error {

	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		fmt.Println("Error while deleting user", err)
		return err
	}

	return nil
}
