package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OmiseService struct {
	userRepo        ports.UserRepository
	transactionRepo ports.TransactionRepository
}

func NewOmiseService(userRepo ports.UserRepository, transactionRepo ports.TransactionRepository) *OmiseService {
	return &OmiseService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

// ChargeCreditCard charges a credit card with the given amount and token.
func (s *OmiseService) ChargeCreditCard(amount int64, token string, userID string) (*omise.Charge, error) {
	facade := utils.FacadeSingleton()
	client := facade.Omise

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   amount,
		Currency: "thb",
		Card:     token,
	}

	if e := client.Do(charge, createCharge); e != nil {
		return nil, e
	}
	existingUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}
	newTransaction := domain.Transaction{
		ID:       charge.Transaction,
		Amount:   charge.Amount,
		Currency: charge.Currency,
		Created:  time.Now(),
		UserID:   uint(id),
	}

	transaction, err := s.transactionRepo.CreateTransaction(&newTransaction)
	if err != nil {
		return nil, err
	}
	fmt.Println(transaction)
	existingUser.Transactions = append(existingUser.Transactions, *transaction)

	_, err = s.userRepo.UpdateUser(userID, existingUser)
	if err != nil {
		return nil, err
	}
	return charge, nil
}

// ChargeBanking charges a specified amount from a banking source.
func (s *OmiseService) ChargeBanking(amount int64, source string, userID string) (*omise.Charge, error) {
	facade := utils.FacadeSingleton()
	client := facade.Omise

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:    amount,
		Currency:  "thb",
		Source:    source,
		ReturnURI: "https://example.com/orders/345678/complete",
		Metadata:  map[string]interface{}{"user_id": userID},
	}

	if e := client.Do(charge, createCharge); e != nil {
		return nil, e
	}
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// user.Transactions = append(user.Transactions, charge.Transaction)

	s.userRepo.UpdateUser(userID, user)
	return charge, nil
}

func (s *OmiseService) AddTransactionToUser(userID string, transaction domain.Transaction) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.Transactions = append(user.Transactions, transaction)
	_, err = s.userRepo.UpdateUser(userID, user)
	if err != nil {
		return err
	}
	return nil
}

// CreateToken creates a token
func (s *OmiseService) CreateToken(name string, number string, expirationMonth time.Month, expirationYear int) (*omise.Card, error) {
	facade := utils.FacadeSingleton()
	client := facade.Omise

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
	facade := utils.FacadeSingleton()
	client := facade.Omise

	result := &omise.CustomerList{}

	err := client.Do(result, &operations.ListCustomers{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AttachCardToCustomer attaches a card to a customer
func (s *OmiseService) AttachCardToCustomer(customerID string, card string) (*omise.Customer, error) {
	facade := utils.FacadeSingleton()
	client := facade.Omise

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
	facade := utils.FacadeSingleton()
	client := facade.Omise

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
	facade := utils.FacadeSingleton()
	client := facade.Omise

	result := &omise.ChargeList{}

	err := client.Do(result, &operations.ListCharges{})
	if err != nil {
		return nil, err
	}
	return result, nil

}

// GetTransaction get a transaction
func (s *OmiseService) GetTransaction(transactionID string) (*omise.Transaction, error) {
	facade := utils.FacadeSingleton()
	client := facade.Omise

	result := &omise.Transaction{}

	err := client.Do(result, &operations.RetrieveTransaction{
		TransactionID: transactionID,
	})

	if err != nil {
		return nil, err
	}
	log.Println(result)
	return result, nil
}

// GetCustomer get a customer
func (s *OmiseService) GetCustomer(customerID string) (*omise.Customer, error) {
	facade := utils.FacadeSingleton()
	client := facade.Omise

	customer := &omise.Customer{}

	err := client.Do(customer, &operations.RetrieveCustomer{
		CustomerID: customerID,
	})

	if err != nil {
		return nil, err
	}

	return customer, nil
}
