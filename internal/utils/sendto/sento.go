package sendto

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/ntquang/ecommerce/global"
	"go.uber.org/zap"
)

const (
	EMAIL_APP_PASSWORD = "xysz idqs awxa xzbd"
	EMAIL_NAME         = "quangtn0607@gmail.com"
	EMAIL_APP_HOST     = "smtp.gmail.com"
	EMAIL_APP_PORT     = "587"
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func BuildMessageMail(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("Form: %s\r\n", mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func SendTextEmailOtp(to []string, from string, otp string) error {
	contextEmail := Mail{
		From:    EmailAddress{Address: from, Name: "test"},
		To:      to,
		Subject: "Otp verification",
		Body:    fmt.Sprintf("Your Otp is %s. Please enter it to verify account.", otp),
	}

	messageMail := BuildMessageMail(contextEmail)

	//send smtp
	authencation := smtp.PlainAuth("", EMAIL_NAME, EMAIL_APP_PASSWORD, EMAIL_APP_HOST)

	err := smtp.SendMail(EMAIL_APP_HOST+":"+EMAIL_APP_PORT, authencation, EMAIL_NAME, to, []byte(messageMail))
	if err != nil {
		global.Logger.Error("Send email error", zap.Error(err))
		return err
	}
	return nil
}

func SendTemplateEmailOtp(
	to []string, from string, htmlTemplate string,
	dataTemplate map[string]interface{},
) error {
	htmlBody, err := getMailTemplate(htmlTemplate, dataTemplate)
	if err != nil {
		return err
	}

	return send(to, from, htmlBody)
}

func getMailTemplate(nameTemplate string, dataTemplate map[string]interface{}) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(nameTemplate).ParseFiles("templates-email/" + nameTemplate))
	err := t.Execute(htmlTemplate, dataTemplate)
	if err != nil {
		return "", err
	}
	return htmlTemplate.String(), nil
}

func send(to []string, from string, htmlTemplate string) error {
	contextEmail := Mail{
		From:    EmailAddress{Address: from, Name: "test"},
		To:      to,
		Subject: "Otp verification",
		Body:    htmlTemplate,
	}

	messageMail := BuildMessageMail(contextEmail)

	//send smtp
	authencation := smtp.PlainAuth("", EMAIL_NAME, EMAIL_APP_PASSWORD, EMAIL_APP_HOST)

	err := smtp.SendMail(EMAIL_APP_HOST+":"+EMAIL_APP_PORT, authencation, EMAIL_NAME, to, []byte(messageMail))
	if err != nil {
		global.Logger.Error("Send email error", zap.Error(err))
		return err
	}
	return nil
}
