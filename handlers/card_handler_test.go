package handlers_test

import (
	"bank-api/handlers"
	"bank-api/models"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// fakeCardService реализует интерфейс CardService для тестирования обработчиков.
type fakeCardService struct{}

func (f *fakeCardService) CreateCard(userID, accountID int) (*models.Card, error) {
	// Возвращаем фиксированный объект карты.
	return &models.Card{
		ID:             1,
		UserID:         userID,
		AccountID:      accountID,
		CardNumber:     "encrypted_card_number",
		ExpirationDate: "encrypted_exp_date",
		CVVHash:        "hashed_cvv",
		CreatedAt:      time.Now(),
	}, nil
}

func (f *fakeCardService) GetCardByID(id int) (*models.Card, error) {
	// Для теста возвращаем карту с userID 42.
	return &models.Card{
		ID:             id,
		UserID:         42,
		AccountID:      101,
		CardNumber:     "encrypted_card_number",
		ExpirationDate: "encrypted_exp_date",
		CVVHash:        "hashed_cvv",
		CreatedAt:      time.Now(),
	}, nil
}

func TestCreateCardHandler(t *testing.T) {
	fakeSvc := &fakeCardService{}
	handler := handlers.NewCardHandler(fakeSvc)

	// Тело запроса содержит account_id.
	body := `{"account_id": 101}`
	req := httptest.NewRequest("POST", "/cards", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// Устанавливаем userID в контекст (например, "42").
	ctx := context.WithValue(req.Context(), "userID", "42")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateCard(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var card models.Card
	if err := json.NewDecoder(rr.Body).Decode(&card); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if card.ID != 1 {
		t.Errorf("expected card ID 1, got %d", card.ID)
	}
	if card.UserID != 42 {
		t.Errorf("expected userID 42, got %d", card.UserID)
	}
}

func TestGetCardHandler(t *testing.T) {
	fakeSvc := &fakeCardService{}
	handler := handlers.NewCardHandler(fakeSvc)

	req := httptest.NewRequest("GET", "/cards/1", nil)
	ctx := context.WithValue(req.Context(), "userID", "42")
	req = req.WithContext(ctx)
	// Используем mux.SetURLVars, чтобы установить переменную "id" в URL.
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.GetCard(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var card models.Card
	if err := json.NewDecoder(rr.Body).Decode(&card); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if card.ID != 1 {
		t.Errorf("expected card ID 1, got %d", card.ID)
	}
	if card.UserID != 42 {
		t.Errorf("expected userID 42, got %d", card.UserID)
	}
}
