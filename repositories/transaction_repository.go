package repositories

import (
	"database/sql"
	"fmt"
	"time"
	"bank-api/models"
)

// TransactionRepository определяет методы для работы с транзакциями.
type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByAccountID(accountID int) ([]models.Transaction, error)
	SumByType(userID int, txType string, since time.Time) (float64, error)

}

type transactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository возвращает новую реализацию TransactionRepository.
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create вставляет новую транзакцию и возвращает сгенерированный ID.
func (r *transactionRepository) Create(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (account_id, amount, type, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	err := r.db.QueryRow(query, transaction.AccountID, transaction.Amount, transaction.Type, transaction.CreatedAt).
		Scan(&transaction.ID)
	if err != nil {
		return fmt.Errorf("error inserting transaction: %w", err)
	}
	return nil
}

// GetByAccountID возвращает все транзакции по ID счета.
func (r *transactionRepository) GetByAccountID(accountID int) ([]models.Transaction, error) {
	query := `
		SELECT id, account_id, amount, type, created_at
		FROM transactions WHERE account_id = $1
	`
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("error fetching transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.Amount, &t.Type, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning transaction: %w", err)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *transactionRepository) SumByType(userID int, txType string, since time.Time) (float64, error) {
	var sum sql.NullFloat64
	err := r.db.QueryRow(
		`SELECT COALESCE(SUM(amount),0) FROM transactions
		 WHERE user_id = $1 AND type = $2 AND created_at >= $3`,
		userID, txType, since,
	).Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum.Float64, nil
}