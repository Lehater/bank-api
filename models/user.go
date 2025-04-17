package models

import (
	"time"
)

// User представляет пользователя системы.
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email" validate:"required,email"`
	Username     string    `json:"username" validate:"required,min=3,max=30"`
	PasswordHash string    `json:"-"` // Пароль хранится в виде хеша; не выводится в JSON
	CreatedAt    time.Time `json:"created_at"`
}
