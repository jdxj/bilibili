package email

import (
	"fmt"
	"net/smtp"

	"github.com/jdxj/bilibili/config"
	"github.com/jordan-wright/email"
)

const SMTPHost = "smtp.qq.com"
const SMTPAddr = "smtp.qq.com:587"

var log *Email

func init() {
	log = NewEmail(config.Cfg.Email)
}

func NewEmail(config *config.Email) *Email {
	jEmail := email.NewEmail()
	jEmail.From = fmt.Sprintf("bilibili <%s>", config.User)
	jEmail.To = []string{config.User}

	auth := smtp.PlainAuth("", config.User, config.Password, SMTPHost)

	e := &Email{
		config: config,
		auth:   auth,
		jEmail: jEmail,
	}
	return e
}

type Email struct {
	config *config.Email
	auth   smtp.Auth
	jEmail *email.Email
}

func (e *Email) Log(format string, a ...interface{}) {
	log := fmt.Sprintf(format, a...)
	e.send("bilibili sign log", log)
}

func (e *Email) send(subject, text string) {
	je := e.jEmail
	je.Subject = subject
	je.Text = []byte(text)

	je.Send(SMTPAddr, e.auth)
}

func Log(format string, a ...interface{}) {
	log.Log(format, a...)
}
