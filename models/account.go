package models

import (
	"time"
)

// Account представляет банковский счёт пользователя.
type Account struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" validate:"required"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
