package repository

import (
	"fmt"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db,
	}
}

func (r *TransactionRepository) CreateTransaction(transaction *domain.Transaction) (*domain.Transaction, error) {
	if err := r.db.Create(transaction).Error; err != nil {
		fmt.Println("Error while creating transaction", err)
		return nil, err
	}

	return transaction, nil
}

func (r *TransactionRepository) GetTransactionByID(id string) (*domain.Transaction, error) {
	transaction := &domain.Transaction{}
	if err := r.db.First(transaction, id).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *TransactionRepository) GetAllTransactions() ([]*domain.Transaction, error) {

	transactions := []*domain.Transaction{}
	if err := r.db.Find(&transactions).Error; err != nil {
		fmt.Println("Error while getting all transactions", err)
		return nil, err
	}

	return transactions, nil
}
