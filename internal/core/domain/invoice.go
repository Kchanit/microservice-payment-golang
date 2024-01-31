package domain

import "time"

type Invoice struct {
	ID            string    `json:"id"`
	Date          time.Time `json:"date"`
	TransactionID string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Products      []Product `json:"products"`
	UserID        uint      `json:"user_id"`
}
