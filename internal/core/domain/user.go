package domain

import (
	"time"
)

type User struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	Name          string        `gorm:"not null,unique" json:"name" validate:"required,min=3"`
	Email         string        `gorm:"not null,unique" json:"email" validate:"required,email"`
	CustomerToken string        `gorm:"unique " json:"customer_token"`
	CardTokens    []CardToken   `gorm:"many2many:user_card_tokens" json:"card_tokens"`
	Transactions  []Transaction `gorm:"foreignKey:UserID" json:"transactions"`
	Addresses     []Address     `gorm:"foreignKey:UserID" json:"addresses"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
