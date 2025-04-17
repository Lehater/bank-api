package services

import (
	"errors"
	"strconv"
	"time"

	"bank-api/models"
	"bank-api/repositories"
	"bank-api/utils"

	"github.com/golang-jwt/jwt/v5"
)

// UserService описывает бизнес-логику, связанную с пользователями.
type UserService interface {
	Register(input models.RegistrationInput) (*models.User, error)
	Authenticate(email, password string) (string, error)
}

type userService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(userRepo repositories.UserRepository, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register проводит регистрацию нового пользователя.
func (s *userService) Register(input models.RegistrationInput) (*models.User, error) {
	// Проверка входных данных.
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Хеширование пароля.
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Формирование структуры пользователя.
	user := &models.User{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
	}

	// Сохранение пользователя в базе.
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Authenticate проверяет учетные данные и генерирует JWT-токен в случае успеха.
func (s *userService) Authenticate(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("неверный пароль")
	}

	// Генерация JWT-токена.
	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateJWTToken создает JWT-токен с id пользователя в качестве Subject.
func (s *userService) generateJWTToken(userID int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
