package ports

import "github.com/Kchanit/microservice-payment-golang/internal/core/domain"

type InvoiceRepository interface {
	CreateInvoice(invoice *domain.Invoice) (*domain.Invoice, error)
	GetInvoiceByID(id string) (*domain.Invoice, error)
}

type InvoiceService interface {
	CreateInvoice(invoice *domain.Invoice) (*domain.Invoice, error)
	GetInvoiceByID(id string) (*domain.Invoice, error)
}
