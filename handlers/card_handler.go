package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/services"

	"github.com/gorilla/mux"
)

// CardHandler отвечает за HTTP-обработку запросов, связанных с картами.
type CardHandler struct {
	cardService services.CardService
}

// NewCardHandler создаёт новый экземпляр CardHandler.
func NewCardHandler(cardService services.CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

// CreateCard обрабатывает POST-запрос на создание виртуальной карты.
// URL: POST /cards
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	// Получение userID из контекста (установленного, например, AuthMiddleware).
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok || userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Ожидаемый JSON-запрос должен содержать account_id.
	var reqBody struct {
		AccountID int `json:"account_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if reqBody.AccountID == 0 {
		http.Error(w, "account_id is required", http.StatusBadRequest)
		return
	}

	// Вызываем CardService для генерации карты.
	card, err := h.cardService.CreateCard(userID, reqBody.AccountID)
	if err != nil {
		http.Error(w, "Failed to create card: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданную карту в JSON-формате.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

// GetCard обрабатывает GET-запрос на просмотр карты по ID.
// URL: GET /cards/{id}
func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста.
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok || userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Извлекаем ID карты из URL.
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Получаем карту через сервис.
	card, err := h.cardService.GetCardByID(cardID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	// Проверяем, что пользователь является владельцем карты.
	if card.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Если требуется, здесь можно реализовать дешифрование зашифрованных полей,
	// но если возвращать в неизменном виде – делаем JSON-ответ.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}
