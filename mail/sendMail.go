package mail

import (
	"flight-ticket-aggregator/domain"

	"github.com/go-mail/mail"
)

func sendMail(message *mail.Message) error {
	// Init Mail connection
	dialer := mail.NewDialer(domain.SMTP_HOST, domain.SMTP_PORT, domain.SMTP_USER, domain.SMTP_PASSWORD)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	// send the mail
	return dialer.DialAndSend(message)
}
