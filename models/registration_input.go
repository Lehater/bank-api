package models

import (
	"github.com/go-playground/validator/v10"
)

// RegistrationInput описывает входные данные при регистрации пользователя.
type RegistrationInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8"`
}

// Validate выполняет валидацию входных данных.
func (input *RegistrationInput) Validate() error {
	validate := validator.New()
	return validate.Struct(input)
}
