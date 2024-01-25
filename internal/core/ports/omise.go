package ports

import (
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/omise/omise-go"
)

type OmiseService interface {
	ChargeBanking(amount int64, source string, userID string) (*omise.Charge, error)
	ChargeCreditCard(amount int64, token string, userID string) (*omise.Charge, error)
	CreateToken(name string, number string, expirationMonth time.Month, expirationYear int) (*omise.Card, error)
	ListCustomers() (*omise.CustomerList, error)
	GetCustomer(customerID string) (*omise.Customer, error)
	AttachCardToCustomer(customerID string, card string) (*omise.Customer, error)
	RetrieveCharge(chargeID string) (*omise.Charge, error)
	GetCharges() (*omise.ChargeList, error)
	GetTransaction(transactionID string) (*omise.Transaction, error)
	AddTransactionToUser(userID string, transaction domain.Transaction) error
}
