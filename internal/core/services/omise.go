package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OmiseService struct{}

func NewOmiseService() *OmiseService {
	return &OmiseService{}
}

func NewOmiseClient() (*omise.Client, error) {
	OmisePublicKey := os.Getenv("OMISE_PUBLIC_KEY")
	OmiseSecretKey := os.Getenv("OMISE_SECRET_KEY")
	client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if e != nil {
		fmt.Println(e)
	}
	return client, e
}

// ChargeCreditCard charges a credit card with the given amount and token.
func (s *OmiseService) ChargeCreditCard(amount string, token string) (*omise.Charge, error) {
	client, e := NewOmiseClient()
	if e != nil {
		return nil, e
	}

	// Parse amount to int
	amountInt, _ := strconv.Atoi(amount)

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   int64(amountInt),
		Currency: "thb",
		Card:     token,
	}

	if e := client.Do(charge, createCharge); e != nil {
		return nil, e
	}

	return charge, nil
}

// ChargeBanking charges a specified amount from a banking source.
func (s *OmiseService) ChargeBanking(amount string, source string) (*omise.Charge, error) {
	client, e := NewOmiseClient()
	if e != nil {
		return nil, e
	}

	// Parse amount to int
	amountInt, _ := strconv.Atoi(amount)

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:    int64(amountInt),
		Currency:  "thb",
		Source:    source,
		ReturnURI: "https://example.com/orders/345678/complete",
	}

	if e := client.Do(charge, createCharge); e != nil {
		return nil, e
	}

	return charge, nil
}

// CreateToken creates a token
func (s *OmiseService) CreateToken(name string, number string, expirationMonth time.Month, expirationYear int) (*omise.Card, error) {
	client, e := NewOmiseClient()
	if e != nil {
		return nil, e
	}

	result := &omise.Card{}
	err := client.Do(result, &operations.CreateToken{
		Name:            name,
		Number:          number,
		ExpirationMonth: expirationMonth,
		ExpirationYear:  expirationYear,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListCustomers lists all customers
func (s *OmiseService) ListCustomers() (*omise.CustomerList, error) {
	client, e := NewOmiseClient()
	if e != nil {
		return nil, e
	}

	result := &omise.CustomerList{}

	err := client.Do(result, &operations.ListCustomers{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AttachCardToCustomer attaches a card to a customer
func (s *OmiseService) AttachCardToCustomer(customerID string, card string) (*omise.Customer, error) {
	client, e := NewOmiseClient()
	if e != nil {
		return nil, e
	}

	//handle token was already attached

	result := &omise.Customer{}

	err := client.Do(result, &operations.UpdateCustomer{
		CustomerID: customerID,
		Card:       card,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RetrieveCharge retrieves a charge
func (s *OmiseService) RetrieveCharge(chargeID string) (*omise.Charge, error) {
	client, e := NewOmiseClient()
	if e != nil {
		return nil, e
	}

	result := &omise.Charge{}
	err := client.Do(result, &operations.RetrieveCharge{
		ChargeID: chargeID,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCharges get all charges
func (s *OmiseService) GetCharges() (*omise.ChargeList, error) {
	client, e := NewOmiseClient()
	if e != nil {
		return nil, e
	}

	result := &omise.ChargeList{}

	err := client.Do(result, &operations.ListCharges{})
	if err != nil {
		return nil, err
	}
	return result, nil

}
