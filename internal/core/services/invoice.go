package services

import (
	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
)

type InvoiceService struct {
	InvoiceRepository ports.InvoiceRepository
}

func NewInvoiceService(invoiceRepository ports.InvoiceRepository) *InvoiceService {
	return &InvoiceService{
		InvoiceRepository: invoiceRepository,
	}
}

func (s *InvoiceService) CreateInvoice(invoice *domain.Invoice) (*domain.Invoice, error) {
	invoice, err := s.InvoiceRepository.CreateInvoice(invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}
