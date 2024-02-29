package email

import (
	"bytes"
	"html/template"
	"net/smtp"

	"github.com/Foedie/foedie-server-v2/email/pkg/config"
)

type Email struct {
	FromEmail    string
	FromPassword string
	// This is list of email address that email is to be sent.
	ToEmail []string
	// This is to create subject of the email
	Subject string
	// This is the message to send in the mail
	Message string
	// This is a template of sending an email
	EmailTemplate string
}

func NewEmail(toEmail []string, subject string, message string, emailTemplate string, config config.Config) *Email {
	return &Email{
		ToEmail:       toEmail,
		Subject:       subject,
		Message:       message,
		EmailTemplate: emailTemplate,
		FromEmail:     config.FromEmail,
		FromPassword:  config.FromPassword,
	}
}

func (email Email) SendEmail() error {

	auth := smtp.PlainAuth("", email.FromEmail, email.FromPassword, "smtp.gmail.com")

	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Dhanush",
		URL:  email.Message,
	}

	buf, err := ParseTemplate(email, email.EmailTemplate, templateData)

	if err != nil {
		return err
	}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + email.Subject + "\n"
	msg := []byte(subject + mime + "\n" + buf.String())
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, email.FromEmail, email.ToEmail, msg); err != nil {
		return err
	}
	return nil
}

func ParseTemplate(email Email, templateFileName string, data interface{}) (*bytes.Buffer, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf, nil
}
