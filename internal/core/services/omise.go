package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OmiseService struct {
	userRepo           ports.UserRepository
	transactionRepo    ports.TransactionRepository
	transactionService ports.TransactionService
}

func NewOmiseService(userRepo ports.UserRepository, transactionRepo ports.TransactionRepository, transactionService ports.TransactionService) *OmiseService {
	return &OmiseService{
		userRepo:           userRepo,
		transactionRepo:    transactionRepo,
		transactionService: transactionService,
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

	_, err = s.userRepo.UpdateUser(userID, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
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

func (h *OmiseService) processSuccessfulTransaction(userID string, amount int64, payload map[string]interface{}) error {
	// Convert userID to uint
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create a new transaction object
	newTransaction := &domain.Transaction{
		ID:       payload["transaction"].(string),
		UserID:   uint(userIDUint),
		Amount:   int64(amount),
		Currency: payload["currency"].(string),
		Status:   payload["status"].(string),
	}
	fmt.Println("TransactionID: ", newTransaction.ID)
	transaction, err := h.transactionRepo.CreateTransaction(newTransaction)
	if err != nil {
		fmt.Println(err)
	}

	customer := domain.User{
		Name: "John Doe",
		Addresses: []domain.Address{
			{
				Address:    "89 somewhere",
				PostalCode: "12345",
				City:       "Phuket",
				Country:    "Thailand",
			},
		},
	}

	products := []domain.Product{
		{
			Name:        "T-shirt",
			Description: "White t-shirt, Size L",
			Price:       9300,
			Quantity:    1,
		},
		{
			Name:        "Test2",
			Description: "Temp Description",
			Price:       4800,
			Quantity:    5,
		},
	}

	err = h.generateAndUploadInvoice(customer, products, transaction, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *OmiseService) HandleWebhook(payload map[string]interface{}) error {
	// Retrieve the status and user ID from the payload
	status, ok := payload["status"].(string)
	if !ok {
		fmt.Println("Error retrieving status from payload")
	}
	userID, ok := payload["metadata"].(map[string]interface{})["user_id"].(string)
	if !ok {
		fmt.Println("Error retrieving user ID from payload")
	}

	// Covert amount to int64
	amount, ok := payload["amount"].(float64)
	if !ok {
		fmt.Println("Error retrieving amount from payload")
	}

	// If status is successful, create the transaction and add it to the user
	if status == "successful" {
		err := s.processSuccessfulTransaction(userID, int64(amount), payload)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (h *OmiseService) generateAndUploadInvoice(customer domain.User, products []domain.Product, transaction *domain.Transaction, userID string) error {
	// Generate Invoice
	outputName, err := utils.GenerateInvoice(customer, products, transaction.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	bucketName := "pixelmanstorage"
	objectName := "invoices/" + transaction.ID + time.Now().Format("2006-01-02") + ".pdf"

	// Get Minio client instance
	minioClientInstance, err := utils.GetMinioClient()
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Upload PDF file to Minio
	err = minioClientInstance.UploadImage(bucketName, objectName, outputName)
	if err != nil {
		log.Fatal(err)
		return err
	}

	os.Remove(outputName)
	log.Println("PDF file uploaded successfully.")

	// Add the transaction to the user
	_, err = h.transactionService.AddTransactionToUser(userID, *transaction)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
