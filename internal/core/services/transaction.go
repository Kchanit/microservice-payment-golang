package services

import (
	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
)

type TransactionService struct {
	transactionRepo ports.TransactionRepository
}

func NewTransactionService(transactionRepo ports.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) GetAllTransactions() ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetAllTransactions()
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
