package ports

import "github.com/Kchanit/microservice-payment-golang/internal/core/domain"

type TransactionRepository interface {
	CreateTransaction(transaction *domain.Transaction) (*domain.Transaction, error)
	GetTransactionByID(id string) (*domain.Transaction, error)
	GetAllTransactions() ([]*domain.Transaction, error)
}

type TransactionService interface {
	GetAllTransactions() ([]*domain.Transaction, error)
	AddTransactionToUser(userID string, transaction domain.Transaction) (*domain.User, error)
	CreateTransaction(transaction *domain.Transaction) (*domain.Transaction, error)
}
