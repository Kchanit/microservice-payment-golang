package repository

import (
	"fmt"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{
		db,
	}
}

func (r *InvoiceRepository) CreateInvoice(invoice *domain.Invoice) (*domain.Invoice, error) {
	if err := r.db.Create(invoice).Error; err != nil {
		fmt.Println("Error while creating invoice", err)
		return nil, err
	}

	return invoice, nil
}

func (r *InvoiceRepository) GetInvoiceByID(id string) (*domain.Invoice, error) {
	invoice := &domain.Invoice{}
	if err := r.db.First(invoice, id).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}
