package services

import (
	"fmt"
	"math"
	"time"

	"bank-api/models"
	"bank-api/repositories"

	"github.com/sirupsen/logrus"
)

// CreditService — базовый интерфейс для кредитов (только что нужно Scheduler, Service и т.д.)
type CreditService interface {
	ApplyForCredit(credit *models.Credit) error
	GetCreditByID(id int) (*models.Credit, error)
	ProcessOverduePayments() error
}

// creditService — реализация CreditService
type creditService struct {
	creditRepo          repositories.CreditRepository
	paymentScheduleRepo repositories.PaymentScheduleRepository
}

// NewCreditService возвращает CreditService
func NewCreditService(
	creditRepo repositories.CreditRepository,
	paymentScheduleRepo repositories.PaymentScheduleRepository,
) CreditService {
	return &creditService{
		creditRepo:          creditRepo,
		paymentScheduleRepo: paymentScheduleRepo,
	}
}

func (s *creditService) ApplyForCredit(credit *models.Credit) error {
	credit.CreatedAt = time.Now()
	if err := s.creditRepo.Create(credit); err != nil {
		return err
	}
	months := 12
	monthlyRate := credit.InterestRate/100/float64(months)
	payment := (credit.Amount*monthlyRate)/(1-math.Pow(1+monthlyRate, -float64(months)))
	for i := 1; i <= months; i++ {
		schedule := &models.PaymentSchedule{
			CreditID:  credit.ID,
			DueDate:   credit.CreatedAt.AddDate(0, i, 0),
			Amount:    payment,
			IsPaid:    false,
			CreatedAt: time.Now(),
		}
		if err := s.paymentScheduleRepo.Create(schedule); err != nil {
			return err
		}
	}
	return nil
}

func (s *creditService) GetCreditByID(id int) (*models.Credit, error) {
	return s.creditRepo.GetByID(id)
}

func (s *creditService) ProcessOverduePayments() error {
	now := time.Now()
	overdue, err := s.paymentScheduleRepo.GetOverdueUnpaid(now)
	if err != nil {
		return fmt.Errorf("failed to get overdue payments: %w", err)
	}

	for _, p := range overdue {
		credit, err := s.creditRepo.GetByID(p.CreditID)
		if err != nil {
			logrus.WithField("creditID", p.CreditID).Errorf("credit not found: %v", err)
			continue
		}

		penalty := p.Amount * 0.10
		newAmount := credit.Amount + penalty

		if err := s.creditRepo.UpdateAmount(credit.ID, newAmount); err != nil {
			logrus.WithField("creditID", credit.ID).Errorf("failed to update credit: %v", err)
			continue
		}

		logrus.WithFields(logrus.Fields{
			"creditID": credit.ID,
			"penalty":  penalty,
		}).Info("Applied overdue penalty")
	}

	return nil
}

// Этот метод НЕ включаем в CreditService интерфейс, но он существует на struct-е.
// Мы будем использовать его только в Handler-е через свой отдельный интерфейс.
func (s *creditService) GetSchedule(creditID int) ([]*models.PaymentSchedule, error) {
	return s.paymentScheduleRepo.GetByCreditID(creditID)
}
