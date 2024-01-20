package domain

import "time"

type Transaction struct {
	ID       string    `json:"id"`
	Amount   int64     `json:"amount"`
	Currency string    `json:"currency"`
	Created  time.Time `json:"created"`
	UserID   uint      `json:"user_id"`
}
