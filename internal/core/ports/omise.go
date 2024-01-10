package ports

import (
	"time"

	"github.com/omise/omise-go"
)

type OmiseService interface {
	ChargeBanking(amount string, source string) (*omise.Charge, error)
	ChargeCreditCard(amount string, token string) (*omise.Charge, error)
	CreateToken(name string, number string, expirationMonth time.Month, expirationYear int) (*omise.Card, error)
	ListCustomers() (*omise.CustomerList, error)
	AttachCardToCustomer(customerID string, card string) (*omise.Customer, error)
	RetrieveCharge(chargeID string) (*omise.Charge, error)
	GetCharges() (*omise.ChargeList, error)
	GetTransaction(transactionID string) (*omise.Transaction, error)
}
