package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/Kchanit/microservice-payment-golang/internal/core/utils"
	"github.com/gofiber/fiber/v2"
)

type OmiseHandler struct {
	omiseService       ports.OmiseService
	userService        ports.UserService
	transactionService ports.TransactionService
}

func NewOmiseHandler(omiseService ports.OmiseService, userService ports.UserService, transactionService ports.TransactionService) *OmiseHandler {
	return &OmiseHandler{
		omiseService:       omiseService,
		userService:        userService,
		transactionService: transactionService,
	}
}

type Token struct {
	Name            string     `json:"name"`
	Number          string     `json:"number"`
	ExpirationMonth time.Month `json:"expiration_month"`
	ExpirationYear  int        `json:"expiration_year"`
}

type ChargeCreditCardInput struct {
	Token  string `json:"token"`
	Amount int64  `json:"amount"`
}

func (h *OmiseHandler) ChargeCreditCard(c *fiber.Ctx) error {
	chargeInput := &ChargeCreditCardInput{}
	if err := c.BodyParser(chargeInput); err != nil {
		fmt.Println(err)
	}
	amount := chargeInput.Amount
	token := chargeInput.Token
	userID := c.Params("userID")

	charge, e := h.omiseService.ChargeCreditCard(amount, token, userID)
	if e != nil {
		return c.Status(500).SendString(e.Error())
	}
	fmt.Printf("[%s]charge: %s  amount: %s %d\n", charge.Status, charge.ID, charge.Currency, charge.Amount)
	if charge.Status == "failed" {
		fmt.Println(charge.FailureCode, charge.FailureMessage)
		return c.Status(500).JSON(fiber.Map{"failure_code": charge.FailureCode, "message": charge.FailureMessage})
	}
	return c.JSON(fiber.Map{"Charge ID": charge.ID, "Amount": charge.Amount, "Status": charge.Status, "Charge": charge})
}

type ChargeBankingInput struct {
	Source string `json:"source"`
	Amount int64  `json:"amount"`
}

func (h *OmiseHandler) ChargeBanking(c *fiber.Ctx) error {
	chargeInput := &ChargeBankingInput{}
	if err := c.BodyParser(chargeInput); err != nil {
		fmt.Println(err)
	}
	amount := chargeInput.Amount
	source := chargeInput.Source
	userID := c.Params("userID")

	charge, err := h.omiseService.ChargeBanking(amount, source, userID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	fmt.Printf("[%s]charge: %s  amount: %s %d\n", charge.Status, charge.ID, charge.Currency, charge.Amount)
	if charge.Status == "failed" {
		fmt.Println(charge.FailureCode, charge.FailureMessage)
		return c.Status(500).JSON(fiber.Map{"failure_code": charge.FailureCode, "message": charge.FailureMessage})
	}
	return c.JSON(fiber.Map{"Charge ID": charge.ID, "Amount": charge.Amount, "Currency": charge.Currency, "AuthorizeURI": charge.AuthorizeURI, "Status": charge.Status, "Charge": charge})
}

func (h *OmiseHandler) CreateToken(c *fiber.Ctx) error {
	token := &Token{}
	if err := c.BodyParser(token); err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	result, err := h.omiseService.CreateToken(token.Name, token.Number, token.ExpirationMonth, token.ExpirationYear)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(result)
	return c.JSON(result)
}

func (h *OmiseHandler) ListCustomers(c *fiber.Ctx) error {
	result, err := h.omiseService.ListCustomers()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(result)
	return c.JSON(result)
}

type AttachCardRequest struct {
	CustomerID string `json:"customer_id"`
	Card       string `json:"card"`
}

func (h *OmiseHandler) AttachCardToCustomer(c *fiber.Ctx) error {
	req := &AttachCardRequest{}
	if err := c.BodyParser(req); err != nil {
		fmt.Println(err)
	}

	customer, err := h.omiseService.AttachCardToCustomer(req.CustomerID, req.Card)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(customer)
	return c.JSON(customer)
}

func (h *OmiseHandler) RetrieveCharge(c *fiber.Ctx) error {
	chargeID := c.Params("charge_id")

	charge, err := h.omiseService.RetrieveCharge(chargeID)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(charge)
	return c.JSON(charge)
}

func (h *OmiseHandler) HandleWebhook(c *fiber.Ctx) error {
	data := c.Body()
	if data == nil {
		return c.SendStatus(500)
	}

	// Parse the JSON payload into a map
	var payload map[string]interface{}
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return c.SendStatus(500)
	}

	// Retrieve the status and user ID from the payload
	status, ok := payload["status"].(string)
	if !ok {
		fmt.Println("Error retrieving status from payload")
		return c.SendStatus(500)
	}
	userID, ok := payload["metadata"].(map[string]interface{})["user_id"].(string)
	if !ok {
		fmt.Println("Error retrieving user ID from payload")
		return c.SendStatus(500)
	}

	// If status is successful, create the transaction and add it to the user
	if status == "successful" {
		// Convert the user ID to uint
		userIDUint, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		// Covert amount to int64
		amount, ok := payload["amount"].(float64)
		if !ok {
			fmt.Println("Error retrieving amount from payload")
			return c.SendStatus(500)
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
		transaction, err := h.transactionService.CreateTransaction(newTransaction)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		// Map products
		productsInput, _ := payload["products"].([]interface{})
		var products []domain.Product
		for _, p := range productsInput {
			productMap, _ := p.(map[string]interface{})
			product := domain.Product{
				Name:        productMap["name"].(string),
				Description: productMap["description"].(string),
				Price:       float64(productMap["price"].(float64)),
				Quantity:    int(productMap["quantity"].(float64)),
			}

			products = append(products, product)
		}

		// Map customer
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

		// products := []domain.Product{
		// 	{
		// 		Name:        "T-shirt",
		// 		Description: "White t-shirt, Size L",
		// 		Price:       9300,
		// 		Quantity:    1,
		// 	},
		// 	{
		// 		Name:        "Test2",
		// 		Description: "Temp Description",
		// 		Price:       4800,
		// 		Quantity:    5,
		// 	},
		// }

		// customerInput, _ := payload["customer"].(map[string]interface{})
		// customer := domain.User{
		// 	Name: customerInput["name"].(string),
		// 	Addresses: []domain.Address{
		// 		{
		// 			Address:    customerInput["address"].(map[string]interface{})["address"].(string),
		// 			PostalCode: customerInput["address"].(map[string]interface{})["postal_code"].(string),
		// 			City:       customerInput["address"].(map[string]interface{})["city"].(string),
		// 			Country:    customerInput["address"].(map[string]interface{})["country"].(string),
		// 		},
		// 	},
		// }

		// Generate Invoice
		outputName, err := utils.GenerateInvoice(customer, products, transaction.ID)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		bucketName := "pixelmanstorage"
		objectName := "invoices/" + transaction.ID + time.Now().Format("2006-01-02") + ".pdf"

		// Get Minio client instance
		minioClientInstance, err := utils.GetMinioClient()
		if err != nil {
			log.Fatal(err)
		}

		// Upload PDF file to Minio
		err = minioClientInstance.UploadImage(bucketName, objectName, outputName)
		if err != nil {
			log.Fatal(err)
		}

		os.Remove(outputName)
		log.Println("PDF file uploaded successfully.")

		// Add the transaction to the user
		user, err := h.transactionService.AddTransactionToUser(userID, *transaction)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		if err := c.JSON(user); err != nil {
			// Handle the error, log it, or return an appropriate status code
			fmt.Println("Error sending JSON response:", err)
			return c.SendStatus(500)
		}
	}
	return c.SendStatus(200)
}

func (h *OmiseHandler) GetCharges(c *fiber.Ctx) error {
	result, err := h.omiseService.GetCharges()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(result)
}

func (h *OmiseHandler) GetTransaction(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")
	result, err := h.omiseService.GetTransaction(transactionID)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(result)
}

func (h *OmiseHandler) GetCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customerToken")
	result, err := h.omiseService.GetCustomer(customerID)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(result)
}
