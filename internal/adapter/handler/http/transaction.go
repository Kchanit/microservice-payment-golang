package handler

import (
	"fmt"

	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService ports.TransactionService
}

func NewTransactionHandler(transactionService ports.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	transactions, err := h.transactionService.GetAllTransactions()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(transactions)
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	transactionID := c.Params("id")
	transaction, err := h.transactionService.GetTransaction(transactionID)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(transaction)
}
