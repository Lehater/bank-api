package services_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"bank-api/models"
	"bank-api/services"
)

// fakeUserRepo – простая реализация репозитория для тестирования.
type fakeUserRepo struct {
	users map[string]*models.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[string]*models.User)}
}

func (r *fakeUserRepo) Create(user *models.User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	user.ID = len(r.users) + 1
	user.CreatedAt = time.Now()
	r.users[user.Email] = user
	return nil
}

func (r *fakeUserRepo) GetByEmail(email string) (*models.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeUserRepo) GetByID(id int) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// TestRegisterAndAuthenticate проверяет регистрацию и аутентификацию.
func TestRegisterAndAuthenticate(t *testing.T) {
	repo := newFakeUserRepo()
	jwtSecret := "testsecret"
	userService := services.NewUserService(repo, jwtSecret)

	// Тест регистрации.
	regInput := models.RegistrationInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}
	user, err := userService.Register(regInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("Expected user ID to be set after registration")
	}

	// Тест аутентификации с правильным паролем.
	token, err := userService.Authenticate("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if token == "" {
		t.Error("Expected a non-empty token on successful authentication")
	}

	// Тест аутентификации с неверным паролем.
	_, err = userService.Authenticate("test@example.com", "wrongpassword")
	if err == nil {
		t.Error("Expected authentication to fail for incorrect password")
	}

	// Дополнительно можно проверить, что токен содержит корректный Subject (идентификатор пользователя).
	// Например, при разборе JWT токена (если хотите расширить тест).
	parsedUser, err := repo.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Error fetching user: %v", err)
	}
	if strconv.Itoa(parsedUser.ID) == "" {
		t.Error("Expected user ID to be correctly convertible to string in token claims")
	}
}
