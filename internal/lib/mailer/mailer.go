package mailer

import (
	"bytes"
	"html/template"
	"log"

	"github.com/stanislavCasciuc/atom-fit/internal/lib/config"
	"gopkg.in/gomail.v2"
)

const userVerificationTemplPath = "./internal/lib/mailer/templates/verify-email.html"

func send(to []string, subject string, body string, emailCfg config.MailCfg) error {
	m := gomail.NewMessage()
	m.SetHeader("From", emailCfg.Addr)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(emailCfg.Host, emailCfg.Port, emailCfg.Addr, emailCfg.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendVerifyUser(username, email, code string, emailCfg config.MailCfg) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(userVerificationTemplPath)
	if err != nil {
		log.Fatal("cannot to parse email template")
	}
	t.Execute(
		&body, struct {
			Name string
			Code string
		}{Name: username, Code: code},
	)
	bodyStr := body.String()
	err = send([]string{email}, "User Verification", bodyStr, emailCfg)
	if err != nil {
		return err
	}

	return nil
}
