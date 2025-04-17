package services_test

import (
	"bank-api/models"
	"bank-api/services"
	"errors"
	"testing"
)

// fakeCardRepo реализует интерфейс CardRepository для тестирования.
type fakeCardRepo struct {
	createdCard *models.Card
}

func (f *fakeCardRepo) Create(card *models.Card) error {
	// Простой мок: присваиваем ID и сохраняем ссылку на карту.
	if card.UserID == 0 || card.AccountID == 0 {
		return errors.New("invalid card data")
	}
	card.ID = 1
	f.createdCard = card
	return nil
}

func (f *fakeCardRepo) GetByID(id int) (*models.Card, error) {
	if f.createdCard != nil && f.createdCard.ID == id {
		return f.createdCard, nil
	}
	return nil, errors.New("card not found")
}

func TestCreateCard(t *testing.T) {
	repo := &fakeCardRepo{}
	cardService := services.NewCardService(repo)
	userID := 42
	accountID := 101

	card, err := cardService.CreateCard(userID, accountID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.ID != 1 {
		t.Errorf("expected card ID 1, got %d", card.ID)
	}
	if card.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, card.UserID)
	}
	if card.AccountID != accountID {
		t.Errorf("expected accountID %d, got %d", accountID, card.AccountID)
	}
	if card.CardNumber == "" {
		t.Error("expected non-empty encrypted card number")
	}
	if card.ExpirationDate == "" {
		t.Error("expected non-empty encrypted expiration date")
	}
	if card.CVVHash == "" {
		t.Error("expected non-empty CVV hash")
	}
	if card.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}
