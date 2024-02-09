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
	if err := r.db.Preload("Transactions").First(user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.First(user, "email = ?", email).Error; err != nil {
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
	if err := r.db.Preload("Transactions").Find(&users).Error; err != nil {
		fmt.Println("Error while getting all users", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(id string, user *domain.User) (*domain.User, error) {

	if err := r.db.Model(&domain.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		fmt.Println("Error while updating user", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(id string) error {

	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		fmt.Println("Error while deleting user", err)
		return err
	}

	return nil
}
