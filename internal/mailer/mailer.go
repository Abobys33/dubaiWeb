package mailer

import (
	"dubaiWeb/internal/config"
	"fmt"
	"net/smtp"
)

type Mailer struct {
	cfg config.SMTP
}

func New(cfg config.SMTP) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) SendFeedbackNotification(name, phone string) error {
	subject := "Новая заявка с сайта"
	body := fmt.Sprintf("Имя: %s\r\nТелефон: %s\r\n", name, phone)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.cfg.From, m.cfg.To, subject, body,
	)

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)

	return smtp.SendMail(addr, auth, m.cfg.From, []string{m.cfg.To}, []byte(msg))
}
