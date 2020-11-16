package utils

import (
	"fmt"
	"net/smtp"

	"github.com/jdxj/bilibili/config"
	"github.com/jordan-wright/email"
)

const (
	SMTPHost = "smtp.qq.com"
	SMTPAddr = "smtp.qq.com:587"
)

func SendMessage(to, subject, content string) error {
	ge := config.GetEmail()

	e := email.NewEmail()
	e.From = fmt.Sprintf("bilibili <%s>", ge.User)
	e.To = []string{to}

	e.Subject = subject
	e.Text = []byte(content)

	auth := smtp.PlainAuth("", ge.User, ge.Token, SMTPHost)
	return e.Send(SMTPAddr, auth)
}
