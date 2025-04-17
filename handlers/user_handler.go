package handlers

import (
	"encoding/json"
	"net/http"

	"bank-api/models"
	"bank-api/services"
)

// UserHandler содержит зависимости для работы с пользователями.
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler возвращает новый экземпляр UserHandler.
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register обрабатывает POST-запрос на регистрацию пользователя.
// URL: /register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegistrationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Login обрабатывает POST-запрос для входа пользователя.
// URL: /login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
