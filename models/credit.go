package models

import (
	"time"
)

// Credit представляет кредит, оформленный пользователем.
type Credit struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id" validate:"required"`
	AccountID    int       `json:"account_id" validate:"required"`
	Amount       float64   `json:"amount" validate:"required,gt=0"`
	InterestRate float64   `json:"interest_rate" validate:"required"` // Процентная ставка
	CreatedAt    time.Time `json:"created_at"`
}
