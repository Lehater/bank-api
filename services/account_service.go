package services

import (
	"context"
	"database/sql"

	"bank-api/models"
	"bank-api/repositories"
)

// AccountService описывает операции над банковскими счетами.
type AccountService interface {
	CreateAccount(a *models.Account) error
	Deposit(accountID int, amount float64) error
	Withdraw(accountID int, amount float64) error
	Transfer(fromAccountID, toAccountID int, amount float64) error
}

type accountService struct {
	accountRepo repositories.AccountRepository
	db          *sql.DB
}

// NewAccountService создает AccountService.
func NewAccountService(repo repositories.AccountRepository, db *sql.DB) AccountService {
	return &accountService{accountRepo: repo, db: db}
}

func (s *accountService) CreateAccount(a *models.Account) error {
	return s.accountRepo.Create(a)
}

func (s *accountService) Deposit(id int, amt float64) error {
	return s.accountRepo.UpdateBalance(id, amt)
}

func (s *accountService) Withdraw(id int, amt float64) error {
	return s.accountRepo.UpdateBalance(id, -amt)
}

func (s *accountService) Transfer(fromID, toID int, amt float64) error {
	ctx := context.Background()
	return s.accountRepo.TransferTx(ctx, fromID, toID, amt)
}
