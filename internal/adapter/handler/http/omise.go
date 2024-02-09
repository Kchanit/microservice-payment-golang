package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
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
	City            string     `json:"city"`
	PostalCode      string     `json:"postal_code"`
	SecurityCode    string     `json:"security_code"`
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

	result, err := h.omiseService.CreateToken(token.Name, token.Number, token.ExpirationMonth, token.ExpirationYear, token.City, token.PostalCode, token.SecurityCode)
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

	err := h.omiseService.HandleWebhook(payload)
	if err != nil {
		fmt.Println("Error handling webhook:", err)
		return c.SendStatus(500)
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
	log.Println(transactionID)
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
