package models

import (
	"time"
)

// Transaction представляет операцию по счету.
type Transaction struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
	// Type может принимать значения: deposit, withdrawal, transfer и т.д.
	Type      string    `json:"type" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
