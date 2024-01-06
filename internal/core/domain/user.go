package domain

import "time"

type UserRole string

const (
	Admin    UserRole = "admin"
	Customer UserRole = "customer"
)

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	// Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
