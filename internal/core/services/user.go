package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"gorm.io/gorm"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	var user *domain.User
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	var user *domain.User
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(user *domain.User) (*domain.User, *omise.Customer, error) {
	// Check if the user exists
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return nil, nil, fmt.Errorf("User with email %s already exists", user.Email)
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, nil, err
	}

	client, _ := omise.NewClient(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"))
	customer := &omise.Customer{}
	err = client.Do(customer, &operations.CreateCustomer{
		Email:       user.Email,
		Description: fmt.Sprintf("%s (id: %d)", user.Name, user.ID),
	})
	if err != nil {
		return nil, nil, err
	}

	updatedUser := &domain.User{
		ID:            createdUser.ID,
		Name:          createdUser.Name,
		Email:         createdUser.Email,
		CustomerToken: customer.ID,
	}

	updatedUser, err = s.repo.UpdateUser(strconv.Itoa(int(user.ID)), updatedUser)
	return updatedUser, customer, err
}

func (s *UserService) GetAllUsers() ([]*domain.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(id string, user *domain.User) (*domain.User, error) {
	existingUser, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the user exists
	if existingUser == nil {
		return nil, gorm.ErrRecordNotFound
	}

	updatedUser, err := s.repo.UpdateUser(id, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
