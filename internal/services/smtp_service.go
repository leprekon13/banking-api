package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-mail/mail/v2"
)

func SendPaymentEmail(to string, amount float64) error {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("неверный порт SMTP: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Платеж успешно проведен")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Спасибо за оплату!</h2>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<small>Это автоматическое уведомление</small>
	`, amount))

	d := mail.NewDialer(host, port, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // на проде ставим false

	if err := d.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("не удалось отправить email: %v", err)
	}

	log.Printf("Письмо успешно отправлено на %s", to)
	return nil
}
