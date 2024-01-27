package services

import (
	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
)

type TransactionService struct {
	transactionRepo ports.TransactionRepository
	userRepo        ports.UserRepository
}

func NewTransactionService(transactionRepo ports.TransactionRepository, userRepo ports.UserRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *TransactionService) GetAllTransactions() ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetAllTransactions()
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *TransactionService) CreateTransaction(transaction *domain.Transaction) (*domain.Transaction, error) {
	transaction, err := s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) AddTransactionToUser(userID string, transaction domain.Transaction) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	user.Transactions = append(user.Transactions, transaction)

	user, err = s.userRepo.UpdateUser(userID, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
