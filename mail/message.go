package mail

import (
	"github.com/go-mail/mail"
)

// New Message
func NewMessage() *mail.Message {
	// Init new message
	return mail.NewMessage()
}
