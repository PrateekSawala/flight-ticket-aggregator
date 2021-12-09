package main

import (
	"github.com/go-mail/mail"
)

func (s *Server) NewMessage() *mail.Message {
	// Init new message
	return mail.NewMessage()
}

// Server implements the mail service
type Server struct {
	dialer mail.Dialer
}
