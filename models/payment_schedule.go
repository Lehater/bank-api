package models

import (
	"time"
)

// PaymentSchedule представляет запись в графике платежей по кредиту.
type PaymentSchedule struct {
	ID        int       `json:"id"`
	CreditID  int       `json:"credit_id" validate:"required"`
	DueDate   time.Time `json:"due_date" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
	IsPaid    bool      `json:"is_paid"`
	CreatedAt time.Time `json:"created_at"`
}
