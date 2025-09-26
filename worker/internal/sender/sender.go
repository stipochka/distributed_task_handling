package sender

import (
	"context"
	"fmt"
	"net/smtp"
	"worker/internal/models"
)

type Sender interface {
	SendEmail(ctx context.Context, email models.Email) error
}

type GmailSender struct {
	mailAddress string
	smtpHost    string
	smtpPort    string
	auth        smtp.Auth
}

func NewGmailSender(email, password, host, port string) *GmailSender {
	auth := smtp.PlainAuth("", email, password, host)

	return &GmailSender{
		mailAddress: email,
		smtpHost:    host,
		smtpPort:    port,
		auth:        auth,
	}
}

func (g *GmailSender) SendEmail(ctx context.Context, email models.Email) error {
	const op = "sender.SendEmail"

	message := []byte("To: " + email.Receiver + "\r\n" + "Subject: " + email.Topic + "\r\n" + "\r\n" + email.MessageBody + "\r\n")

	err := smtp.SendMail(g.smtpHost+":"+g.smtpPort, g.auth, g.mailAddress, []string{email.Receiver}, message)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
