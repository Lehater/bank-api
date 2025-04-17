package repositories

import (
	"database/sql"
	"fmt"
	"bank-api/models"
)

// CardRepository определяет методы для работы с картами.
type CardRepository interface {
	Create(card *models.Card) error
	GetByID(id int) (*models.Card, error)
}

type cardRepository struct {
	db *sql.DB
}

// NewCardRepository возвращает новую реализацию CardRepository.
func NewCardRepository(db *sql.DB) CardRepository {
	return &cardRepository{db: db}
}

// Create вставляет новую карту в базу данных.
func (r *cardRepository) Create(card *models.Card) error {
	query := `
		INSERT INTO cards (user_id, account_id, card_number, expiration_date, cvv_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	err := r.db.QueryRow(query, card.UserID, card.AccountID, card.CardNumber, card.ExpirationDate, card.CVVHash, card.CreatedAt).
		Scan(&card.ID)
	if err != nil {
		return fmt.Errorf("error inserting card: %w", err)
	}
	return nil
}

// GetByID возвращает карту по ID.
func (r *cardRepository) GetByID(id int) (*models.Card, error) {
	var card models.Card
	query := `
		SELECT id, user_id, account_id, card_number, expiration_date, cvv_hash, created_at
		FROM cards WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&card.ID, &card.UserID, &card.AccountID, &card.CardNumber, &card.ExpirationDate, &card.CVVHash, &card.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card not found")
		}
		return nil, fmt.Errorf("error fetching card: %w", err)
	}
	return &card, nil
}
