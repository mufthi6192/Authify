package email

import (
	"SMM-PPOB/config"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Mail struct {
	Id          uint
	FromMail    string
	ToMail      string
	ToName      string
	SubjectMail string
	BodyMail    string
	TypeMail    string
}

type Additional struct {
	Url string
}

type dataHtml struct {
	Mail       Mail
	Additional Additional
}

func SendVerificationMail(data Mail) error {
	auth := smtp.PlainAuth("", "verification@diselesain.my.id", "verificatioN123*", "mail.diselesain.my.id")

	mailData := dataHtml{
		Mail: data,
		Additional: Additional{
			Url: fmt.Sprintf("%s/api/client/verification/email?code=%s", config.LOCAL_BASE_URL, data.BodyMail),
		},
	}

	from := data.FromMail
	toMail := []string{data.ToMail}
	subject := data.SubjectMail

	//For Unit Test
	//tmpl := template.Must(template.ParseFiles("verification_email.html"))

	//For Production
	tmpl := template.Must(template.ParseFiles("app/email/verification_email.html"))

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, mailData); err != nil {
		return err
	}

	body := tpl.String()

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html\r\n\r\n%s", toMail[0], subject, body))

	err := smtp.SendMail("mail.diselesain.my.id:587", auth, from, toMail, msg)

	if err != nil {
		return err
	}

	return nil
}

func SendForgetPasswordMail(data Mail) error {
	auth := smtp.PlainAuth("", "verification@diselesain.my.id", "verificatioN123*", "mail.diselesain.my.id")

	mailData := dataHtml{
		Mail: data,
		Additional: Additional{
			Url: fmt.Sprintf("%s/api/auth/reset-password?code=%s", config.LOCAL_BASE_URL, data.BodyMail),
		},
	}

	from := data.FromMail
	toMail := []string{data.ToMail}
	subject := data.SubjectMail

	//For Unit Test
	//tmpl := template.Must(template.ParseFiles("verification_email.html"))

	//For Production
	tmpl := template.Must(template.ParseFiles("app/email/reset_password_email.html"))

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, mailData); err != nil {
		return err
	}

	body := tpl.String()

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html\r\n\r\n%s", toMail[0], subject, body))

	err := smtp.SendMail("mail.diselesain.my.id:587", auth, from, toMail, msg)

	if err != nil {
		return err
	}

	return nil
}
