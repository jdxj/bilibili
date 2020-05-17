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
	emailConfig := config.GetEmail()
	log = NewEmail()
	log.AddRecipients(emailConfig.User)
}

func NewEmail() *Email {
	emailConfig := config.GetEmail()

	jEmail := email.NewEmail()
	jEmail.From = fmt.Sprintf("bilibili <%s>", emailConfig.User)

	auth := smtp.PlainAuth("", emailConfig.User, emailConfig.Password, SMTPHost)

	e := &Email{
		auth:   auth,
		jEmail: jEmail,
	}
	return e
}

type Email struct {
	auth   smtp.Auth
	jEmail *email.Email

	recipients map[string]struct{}
}

func (e *Email) AddRecipients(recipients ...string) {
	if e.recipients == nil {
		e.recipients = make(map[string]struct{})
	}

	for _, v := range recipients {
		if v == "" {
			continue
		}

		e.recipients[v] = struct{}{}
	}

	je := e.jEmail
	je.To = nil
	for addr := range e.recipients {
		je.To = append(je.To, addr)
	}
}

func (e *Email) To() []string {
	return e.jEmail.To
}

func (e *Email) ResetRecipients() {
	e.recipients = nil
	e.jEmail.To = nil
}

func (e *Email) RunLog(format string, a ...interface{}) {
	log := fmt.Sprintf(format, a...)
	e.send("bilibili run log", log)
}

func (e *Email) SignLog(format string, a ...interface{}) {
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
	log.RunLog(format, a...)
}
