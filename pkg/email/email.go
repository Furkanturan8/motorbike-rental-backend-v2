package email

import (
	"net/smtp"
)

type Email struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
}

func NewEmail(from, password, smtpHost, smtpPort string) *Email {
	return &Email{
		From:     from,
		Password: password,
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
	}
}

func (e *Email) Send(to string, subject string, body string) error {
	msg := "From: " + e.From + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", e.From, e.Password, e.SMTPHost)

	return smtp.SendMail(e.SMTPHost+":"+e.SMTPPort, auth, e.From, []string{to}, []byte(msg))
}