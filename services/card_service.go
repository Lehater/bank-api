package services

import (
	"fmt"
	"time"

	"bank-api/models"
	"bank-api/repositories"
	"bank-api/utils"
)

// CardService описывает методы работы с картами.
type CardService interface {
	CreateCard(userID, accountID int) (*models.Card, error)
	GetCardByID(id int) (*models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
}

// NewCardService возвращает CardService.
func NewCardService(repo repositories.CardRepository) CardService {
	return &cardService{cardRepo: repo}
}

// CreateCard генерирует виртуальную карту и сохраняет в БД.
func (s *cardService) CreateCard(userID, accountID int) (*models.Card, error) {
	// 1. генерация данных
	number := utils.GenerateCardNumber()
	exp := utils.GenerateExpirationDate(5)
	cvvPlain, err := utils.GenerateCVV()
	if err != nil {
		return nil, err
	}

	// 2. шифрование (PGP + HMAC)
	encNum, macNum, err := utils.EncryptPGP(number)
	if err != nil {
		return nil, err
	}
	encExp, macExp, err := utils.EncryptPGP(exp)
	if err != nil {
		return nil, err
	}

	// 3. хеш CVV
	cvvHash, err := utils.HashCVV(cvvPlain)
	if err != nil {
		return nil, err
	}

	card := &models.Card{
		UserID:         userID,
		AccountID:      accountID,
		CardNumber:     encNum,
		CardNumberMAC:  macNum,
		ExpirationDate: encExp,
		ExpirationMAC:  macExp,
		CVVHash:        cvvHash,
		CreatedAt:      time.Now(),
	}

	if err := s.cardRepo.Create(card); err != nil {
		return nil, fmt.Errorf("save card: %w", err)
	}
	return card, nil
}

// GetCardByID возвращает карту с расшифрованными полями.
func (s *cardService) GetCardByID(id int) (*models.Card, error) {
	card, err := s.cardRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	num, err := utils.DecryptPGP(card.CardNumber, card.CardNumberMAC)
	if err != nil {
		return nil, err
	}
	exp, err := utils.DecryptPGP(card.ExpirationDate, card.ExpirationMAC)
	if err != nil {
		return nil, err
	}

	card.CardNumber = num
	card.ExpirationDate = exp
	return card, nil
}
