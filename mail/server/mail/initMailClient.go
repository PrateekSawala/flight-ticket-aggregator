package main

import (
	"os"
	"strconv"

	"github.com/go-mail/mail"
)

func InitMailClient() *mail.Dialer {
	port, _ := strconv.Atoi(os.Getenv("FTA_SMTP_PORT"))
	// Init Mail connection
	d := mail.NewDialer(os.Getenv("FTA_SMTP_HOST"), port, os.Getenv("FTA_SMTP_USER"), os.Getenv("FTA_SMTP_PASSWORD"))
	d.StartTLSPolicy = mail.MandatoryStartTLS
	return d
}
