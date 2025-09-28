package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
	"time"
)

type MailHogMailer struct {
	addr      string
	fromEmail string
}

func NewMailhog(addr, fromEmail string) *MailHogMailer {
	return &MailHogMailer{
		addr:      addr,
		fromEmail: fromEmail,
	}
}

func (m *MailHogMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	// parse and buid templates
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	// build message
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n%s\r\n",
		m.fromEmail,
		email,
		subject,
		body,
	))
	if isSandbox {
		fmt.Println("Sandbox mode is enabled")
		fmt.Println(string(msg))
		return nil
	}

	// send
	for i := 0; i < maxRetries; i++ {
		err := smtp.SendMail(m.addr, nil, m.fromEmail, []string{email}, msg)
		if err != nil {
			log.Printf("failed to send email to %s, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("error: %v", err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("email sent successfuly")
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)

}
