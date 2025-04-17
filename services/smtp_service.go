package services

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-mail/mail/v2"
)

// Константы для SMTP (замени на актуальные для своего сервера)
const (
	smtpHost = "smtp.example.com"       // адрес SMTP-сервера
	smtpPort = 587                      // обычно используется порт 587 для TLS
	smtpUser = "noreply@example.com"    // учетная запись отправителя
	smtpPass = "strong_password"        // пароль или токен
)

// createMessage формирует письмо для отправки.
func createMessage(to, subject, body string) *mail.Message {
	m := mail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return m
}

// createDialer настраивает соединение с SMTP-сервером.
func createDialer() *mail.Dialer {
	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{
		ServerName:         smtpHost,
		InsecureSkipVerify: false, // в продакшене не рекомендуется отключать проверку сертификата
	}
	return d
}

// SendPaymentEmail отправляет уведомление об успешном платеже.
func SendPaymentEmail(userEmail string, amount float64) error {
	content := fmt.Sprintf(`
		<h1>Спасибо за оплату!</h1>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<small>Это автоматическое уведомление</small>
	`, amount)
	
	msg := createMessage(userEmail, "Платеж успешно проведен", content)
	dialer := createDialer()

	if err := dialer.DialAndSend(msg); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("Email sent to %s", userEmail)
	return nil
}
