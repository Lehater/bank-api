package services

import (
	"bank-api/repositories"
	"time"
)

// Добавили CreditLoad
type AnalyticsData struct {
	TotalDeposits    float64 `json:"total_deposits"`
	TotalWithdrawals float64 `json:"total_withdrawals"`
	NetChange        float64 `json:"net_change"`
	CreditLoad       float64 `json:"credit_load"` // доля платежей от доходов
}

type AnalyticsService interface {
	GetAnalytics(userID int) (*AnalyticsData, error)
	PredictBalance(accountID int, days int) (float64, error)
}

type analyticsService struct {
	transactionRepo    repositories.TransactionRepository
	accountRepo        repositories.AccountRepository
	creditRepo         repositories.CreditRepository
	paymentScheduleRepo repositories.PaymentScheduleRepository
}

func NewAnalyticsService(
	trRepo repositories.TransactionRepository,
	accRepo repositories.AccountRepository,
	crRepo repositories.CreditRepository,
	psRepo repositories.PaymentScheduleRepository,
) AnalyticsService {
	return &analyticsService{
		transactionRepo:    trRepo,
		accountRepo:        accRepo,
		creditRepo:         crRepo,
		paymentScheduleRepo: psRepo,
	}
}

func (s *analyticsService) GetAnalytics(userID int) (*AnalyticsData, error) {
	// Получаем доходы/расходы за текущий месяц из транзакций
	start := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	deposits, err := s.transactionRepo.SumByType(userID, "deposit", start)
	if err != nil {
		return nil, err
	}
	withdrawals, err := s.transactionRepo.SumByType(userID, "withdraw", start)
	if err != nil {
		return nil, err
	}
	net := deposits - withdrawals

	// Считаем суммарные ежемесячные платежи по всем кредитам
	credits, err := s.creditRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	var totalMonthlyPayments float64
	for _, c := range credits {
		// Находим платежи в этом месяце для кредита
		schedules, _ := s.paymentScheduleRepo.GetByCreditID(c.ID)
		for _, p := range schedules {
			if p.DueDate.After(start) && p.DueDate.Before(time.Now().AddDate(0, 1, 0)) {
				totalMonthlyPayments += p.Amount
			}
		}
	}

	creditLoad := 0.0
	if deposits > 0 {
		creditLoad = totalMonthlyPayments / deposits
	}

	return &AnalyticsData{
		TotalDeposits:    deposits,
		TotalWithdrawals: withdrawals,
		NetChange:        net,
		CreditLoad:       creditLoad,
	}, nil
}

func (s *analyticsService) PredictBalance(accountID int, days int) (float64, error) {
	acc, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return 0, err
	}
	return acc.Balance - float64(days)*50.0, nil
}
