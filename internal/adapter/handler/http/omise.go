package handler

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OmiseHandler struct{}

type Token struct {
	Name            string     `json:"name"`
	Number          string     `json:"number"`
	ExpirationMonth time.Month `json:"expiration_month"`
	ExpirationYear  int        `json:"expiration_year"`
}

func NewOmiseHandler() *OmiseHandler {
	return &OmiseHandler{}
}

func (h *OmiseHandler) Charge(c *fiber.Ctx) error {
	OmisePublicKey := os.Getenv("OMISE_PUBLIC_KEY")
	OmiseSecretKey := os.Getenv("OMISE_SECRET_KEY")

	client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if e != nil {
		fmt.Println(e)
	}

	token := c.Params("token")
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   99500,
		Currency: "thb",
		Card:     token,
	}
	if e := client.Do(charge, createCharge); e != nil {
		return c.Status(500).SendString(e.Error())
	}
	fmt.Printf("[%s]charge: %s  amount: %s %d\n", charge.Status, charge.ID, charge.Currency, charge.Amount)
	return c.JSON(fiber.Map{"Charge ID": charge.ID, "Amount": charge.Amount, "Currency": charge.Currency, "Charge": charge})
}

func (h *OmiseHandler) CreateToken(c *fiber.Ctx) error {
	OmisePublicKey := os.Getenv("OMISE_PUBLIC_KEY")
	OmiseSecretKey := os.Getenv("OMISE_SECRET_KEY")
	client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if e != nil {
		fmt.Println(e)
	}

	token := &Token{}
	if err := c.BodyParser(token); err != nil {
		fmt.Println(err)
	}

	result := &omise.Card{}
	err := client.Do(result, &operations.CreateToken{
		Name:            token.Name,
		Number:          token.Number,
		ExpirationMonth: token.ExpirationMonth,
		ExpirationYear:  token.ExpirationYear,
	})
	if err != nil {
		fmt.Println(err)
	}

	log.Println(result)
	return c.JSON(result)
}

func (h *OmiseHandler) ListCustomers(c *fiber.Ctx) error {
	OmisePublicKey := os.Getenv("OMISE_PUBLIC_KEY")
	OmiseSecretKey := os.Getenv("OMISE_SECRET_KEY")
	client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if e != nil {
		fmt.Println(e)
	}

	result := &omise.CustomerList{}

	err := client.Do(result, &operations.ListCustomers{})
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(result)
}
