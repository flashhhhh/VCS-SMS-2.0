package mailsending

import (
	"github.com/flashhhhh/pkg/logging"
	"gopkg.in/gomail.v2"
)

type MailSending interface {
	SendEmail(to, subject, body string) error
}

type mailSending struct {
	senderEmail string
	senderPassword string
}

func NewMailSending(senderEmail, senderPassword string) MailSending {
	return &mailSending{
		senderEmail: senderEmail,
		senderPassword: senderPassword,
	}
}

func (ms *mailSending) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", ms.senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, ms.senderEmail, ms.senderPassword)

	if err := d.DialAndSend(m); err != nil {
		logging.LogMessage("mail_service", "Failed to send email. Err: " + err.Error(), "ERROR")
		return err
	}

	logging.LogMessage("mail_service", "Successfully send email!", "INFO")
	return nil
}