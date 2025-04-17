package models

import (
	"time"
)

// Card представляет виртуальную банковскую карту.
type Card struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	AccountID       int       `json:"account_id"`
	// Зашифрованный номер карты (PGP)
	CardNumber      string    `json:"card_number"`
	// HMAC от зашифрованного номера для проверки целостности
	CardNumberMAC   string    `json:"card_number_mac"`
	// Зашифрованная дата истечения срока (PGP)
	ExpirationDate  string    `json:"expiration_date"`
	// HMAC от зашифрованной даты для проверки целостности
	ExpirationMAC   string    `json:"expiration_mac"`
	// Хеш CVV (bcrypt). Не выводится в JSON.
	CVVHash         string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
}
