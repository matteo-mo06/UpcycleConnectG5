package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"upcycle_connect-api/internal/config"
)

func loginAuth(username, password string) smtp.Auth {
	return &loginAuthMethod{username: username, password: password}
}

type loginAuthMethod struct {
	username, password string
}

func (a *loginAuthMethod) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuthMethod) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}
	switch {
	case strings.Contains(strings.ToLower(string(fromServer)), "username"):
		return []byte(a.username), nil
	case strings.Contains(strings.ToLower(string(fromServer)), "password"):
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("smtp: prompt LOGIN inattendu: %s", fromServer)
	}
}

func SendVerificationEmail(to, link string) error {
	host := config.SMTPHost()
	port := config.SMTPPort()
	username := config.SMTPUsername()
	password := config.SMTPPassword()
	from := config.MailFrom()
	certPEM := config.SMTPTLSCert()

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

	c, err := smtp.Dial(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: host}
		if certPEM != "" {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM([]byte(certPEM)) {
				tlsConfig.RootCAs = pool
			}
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
			return err
		}
	}

	auth := loginAuth(username, password)
	if err := c.Auth(auth); err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}

	if err := c.Mail(from); err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}
	if err := c.Rcpt(to); err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}

	w, err := c.Data()
	if err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}
	if err := w.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "SendVerificationEmail error:", err)
		return err
	}

	return c.Quit()
}

func sendPlainTextEmail(to, subject, body, logPrefix string) error {
	host := config.SMTPHost()
	port := config.SMTPPort()
	username := config.SMTPUsername()
	password := config.SMTPPassword()
	from := config.MailFrom()
	certPEM := config.SMTPTLSCert()

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		err := fmt.Errorf("configuration SMTP incomplète (SMTP_HOST/SMTP_PORT/SMTP_USERNAME/SMTP_PASSWORD/MAIL_FROM)")
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		from, to, subject, body,
	)

	addr := host + ":" + port

	c, err := smtp.Dial(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: host}
		if certPEM != "" {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM([]byte(certPEM)) {
				tlsConfig.RootCAs = pool
			}
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
			return err
		}
	}

	auth := loginAuth(username, password)
	if err := c.Auth(auth); err != nil {
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}

	if err := c.Mail(from); err != nil {
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}
	if err := c.Rcpt(to); err != nil {
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}

	w, err := c.Data()
	if err != nil {
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}
	if err := w.Close(); err != nil {
		fmt.Fprintln(os.Stderr, logPrefix+" error:", err)
		return err
	}

	return c.Quit()
}

func SendProfessionalApprovedEmail(to string) error {
	subject := "Votre compte professionnel a été validé - UpcycleConnect"
	body := "Bonjour,\r\n\r\n" +
		"Votre demande de compte professionnel vient d'être approuvée par un administrateur.\r\n" +
		"Vous avez maintenant accès à l'espace artisan sur UpcycleConnect.\r\n"

	return sendPlainTextEmail(to, subject, body, "SendProfessionalApprovedEmail")
}

func SendAnnouncementSoldEmail(to, announcementTitle string) error {
	subject := "Votre annonce a été achetée - UpcycleConnect"
	body := "Bonjour,\r\n\r\n" +
		"Votre annonce \"" + announcementTitle + "\" vient d'être achetée.\r\n" +
		"Connectez-vous à votre espace UpcycleConnect pour voir les détails.\r\n"

	return sendPlainTextEmail(to, subject, body, "SendAnnouncementSoldEmail")
}
