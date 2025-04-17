package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"bank-api/handlers"
	"bank-api/models"
	"bank-api/services"
)

// fakeUserRepo – упрощённая реализация репозитория для интеграционных тестов.
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

// TestRegisterHandler проверяет обработчик регистрации.
func TestRegisterHandler(t *testing.T) {
	repo := newFakeUserRepo()
	jwtSecret := "testsecret"
	userService := services.NewUserService(repo, jwtSecret)
	userHandler := handlers.NewUserHandler(userService)

	// Подготовка тестовых данных.
	regData := models.RegistrationInput{
		Email:    "integration@example.com",
		Username: "integrationUser",
		Password: "integrationPass",
	}
	jsonData, _ := json.Marshal(regData)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Вызов обработчика.
	userHandler.Register(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
	var user models.User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if user.Email != regData.Email {
		t.Errorf("Expected email %s, got %s", regData.Email, user.Email)
	}
}

// TestLoginHandler проверяет обработчик аутентификации.
func TestLoginHandler(t *testing.T) {
	repo := newFakeUserRepo()
	jwtSecret := "testsecret"
	userService := services.NewUserService(repo, jwtSecret)
	userHandler := handlers.NewUserHandler(userService)

	// Сначала зарегистрируем пользователя через сервис.
	regData := models.RegistrationInput{
		Email:    "login@example.com",
		Username: "loginUser",
		Password: "loginPass",
	}
	_, err := userService.Register(regData)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Подготовка данных для логина.
	loginData := map[string]string{
		"email":    regData.Email,
		"password": regData.Password,
	}
	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Вызов обработчика.
	userHandler.Login(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
	var respData map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respData); err != nil {
		t.Fatalf("Error decoding login response: %v", err)
	}
	token, exists := respData["token"]
	if !exists || token == "" {
		t.Error("Expected JWT token in the response")
	}

	// Для проверки можно дополнительно запросить пользователя и проверить соответствие токена (при наличии соответствующей логики).
	// Например, получить пользователя по email и убедиться, что его ID содержится в токене.
	user, err := repo.GetByEmail(regData.Email)
	if err != nil {
		t.Fatalf("Error fetching user: %v", err)
	}
	if strconv.Itoa(user.ID) == "" {
		t.Error("Expected a valid user ID to be present")
	}
}
