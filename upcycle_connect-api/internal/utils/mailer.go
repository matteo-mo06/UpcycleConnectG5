package utils

import (
	"fmt"
	"net/smtp"
	"os"

	"upcycle_connect-api/internal/config"
)

func unencryptedPlainAuth(username, password string) smtp.Auth {
	return &plainAuthNoTLS{username: username, password: password}
}

type plainAuthNoTLS struct {
	username, password string
}

func (a *plainAuthNoTLS) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	resp := []byte("\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *plainAuthNoTLS) Next(_ []byte, more bool) ([]byte, error) {
	if more {
		return nil, fmt.Errorf("smtp: unexpected server challenge")
	}
	return nil, nil
}

func SendVerificationEmail(to, link string) error {
	host := config.SMTPHost()
	port := config.SMTPPort()
	username := config.SMTPUsername()
	password := config.SMTPPassword()
	from := config.MailFrom()

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		err := fmt.Errorf("configuration SMTP incomplète (SMTP_HOST/SMTP_PORT/SMTP_USERNAME/SMTP_PASSWORD/MAIL_FROM)")
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}

	subject := "Confirmez votre adresse email - UpcycleConnect"
	body := "Bonjour,\r\n\r\n" +
		"Merci de vous etre inscrit sur UpcycleConnect.\r\n" +
		"Cliquez sur le lien ci-dessous pour confirmer votre adresse email :\r\n\r\n" +
		link + "\r\n\r\n" +
		"Ce lien expire dans 48 heures.\r\n" +
		"Si vous n'etes pas a l'origine de cette inscription, ignorez cet email.\r\n"

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		from, to, subject, body,
	)

	addr := host + ":" + port
	auth := unencryptedPlainAuth(username, password)

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}
	return nil
}
