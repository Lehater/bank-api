package scheduler

import (
	"log"
	"time"

	"bank-api/services"

	"github.com/robfig/cron/v3"
)

// PaymentScheduler отвечает за автоматическую обработку платежей и просрочек.
type PaymentScheduler struct {
	creditService  services.CreditService
	accountService services.AccountService // не используется в данном примере, но может понадобиться для списания
	cronScheduler  *cron.Cron
}

func NewPaymentScheduler(creditSvc services.CreditService, accountSvc services.AccountService) *PaymentScheduler {
	return &PaymentScheduler{
		creditService:  creditSvc,
		accountService: accountSvc,
		cronScheduler:  cron.New(cron.WithSeconds()),
	}
}

// Start запускает шедулер, который каждые 12 часов обрабатывает платежи.
func (ps *PaymentScheduler) Start() {
	// Запланировать задачу в 00:00 и 12:00 каждую сутки.
	_, err := ps.cronScheduler.AddFunc("0 0 0,12 * * *", func() {
		log.Println("Starting scheduled payment processing at", time.Now().Format(time.RFC3339))
		if err := ps.creditService.ProcessOverduePayments(); err != nil {
			log.Printf("Error processing overdue payments: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule payments: %v", err)
	}
	ps.cronScheduler.Start()
	log.Println("Payment scheduler started.")
}
