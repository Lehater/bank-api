package scheduler_test

import (
	"bank-api/models"
	"bank-api/scheduler"
	"testing"
	"time"
)

// fakeCreditService реализует все методы интерфейса services.CreditService.
type fakeCreditService struct{}

func (f *fakeCreditService) ApplyForCredit(c *models.Credit) error {
	return nil
}

func (f *fakeCreditService) GetCreditByID(id int) (*models.Credit, error) {
	return &models.Credit{
		ID:           id,
		UserID:       0,
		AccountID:    0,
		Amount:       1000.0,
		InterestRate: 10.0,
		CreatedAt:    time.Now(),
	}, nil
}

func (f *fakeCreditService) ProcessOverduePayments() error {
	return nil
}

// fakeAccountService реализует все методы интерфейса services.AccountService.
type fakeAccountService struct{}

func (f *fakeAccountService) CreateAccount(a *models.Account) error {
	return nil
}

func (f *fakeAccountService) Deposit(accountID int, amount float64) error {
	return nil
}

func (f *fakeAccountService) Withdraw(accountID int, amount float64) error {
	return nil
}

func (f *fakeAccountService) Transfer(fromAccountID, toAccountID int, amount float64) error {
	return nil
}

// TestSchedulerDoesNotPanic проверяет, что запуск шедулера не вызывает panic.
func TestSchedulerDoesNotPanic(t *testing.T) {
	cs := &fakeCreditService{}
	as := &fakeAccountService{}
	sch := scheduler.NewPaymentScheduler(cs, as)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("scheduler panicked: %v", r)
		}
	}()
	sch.Start()
}
