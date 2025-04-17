package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"bank-api/models"
)

// AccountRepository описывает методы работы с аккаунтами.
type AccountRepository interface {
	Create(a *models.Account) error
	GetByID(id int) (*models.Account, error)
	UpdateBalance(accountID int, delta float64) error
	TransferTx(ctx context.Context, fromID, toID int, amount float64) error
}

type accountRepository struct {
	db *sql.DB
}

// NewAccountRepository возвращает реализацию AccountRepository.
func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(a *models.Account) error {
	_, err := r.db.Exec(
		`INSERT INTO accounts (user_id, balance, currency, created_at)
		 VALUES ($1, $2, $3, NOW())`,
		a.UserID, a.Balance, a.Currency,
	)
	return err
}

func (r *accountRepository) GetByID(id int) (*models.Account, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, balance, currency, created_at
		 FROM accounts WHERE id = $1`, id,
	)
	acc := &models.Account{}
	if err := row.Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Balance,
		&acc.Currency,
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}
	return acc, nil
}

func (r *accountRepository) UpdateBalance(accountID int, delta float64) error {
	_, err := r.db.Exec(
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
		delta, accountID,
	)
	return err
}

func (r *accountRepository) TransferTx(ctx context.Context, fromID, toID int, amount float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance = balance - $1 WHERE id = $2`,
		amount, fromID,
	); err != nil {
		return fmt.Errorf("debit from %d: %w", fromID, err)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
		amount, toID,
	); err != nil {
		return fmt.Errorf("credit to %d: %w", toID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
