package main

import (
	"github.com/go-mail/mail"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
)

func (s *Server) SendMail(message *mail.Message) error {
	log := logging.Log("sendMail")
	recipient := message.GetHeader("To")[0]
	log.Tracef("Sending mail to %s", recipient)
	if *env == domain.LocalEnv {
		log.Tracef("Current environment is local, Skipping sending mail to %s", recipient)
		return nil
	}
	return s.dialer.DialAndSend(message)
}
