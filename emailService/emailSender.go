package emailsender

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

// emailSenderStore holds the private configuration for sending emails.
type emailSenderStore struct {
	smtpServer   string
	smtpPort     int
	username     string
	password     string
	fromEmail    string
	fromName     string
	tlsConfig    *tls.Config
}

// EmailSender defines the interface for sending emails.
type EmailSender interface {
	SendEmail(to, subject, body, htmlBody string) error
}

// newEmailSenderStore creates a new email sender store with default values
// and returns an EmailSender interface.
func NewEmailSenderStore(server string, port int, username, password, fromEmail, fromName string, skipTLS bool) EmailSender {
	tlsConfig := &tls.Config{}
	if skipTLS {
		tlsConfig.InsecureSkipVerify = true // Allow self-signed or untrusted certificates
	}else{
		tlsConfig.InsecureSkipVerify = false // Enforce TLS verification
	}
	tlsConfig.ServerName = server
	return &emailSenderStore{
		smtpServer: server,
		smtpPort:   port,
		username:   username,
		password:   password,
		fromEmail:  fromEmail,
		fromName:   fromName,
		tlsConfig:  tlsConfig,
	}
}


// SendEmail sends an email using the configuration in the emailSenderStore.
func (s *emailSenderStore) SendEmail(to, subject, body, htmlBody string) error {
	fmt.Println("Sending email to:", to)
	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.fromEmail, s.fromName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	if htmlBody != "" {
		m.AddAlternative("text/html", htmlBody)
	}

	d := gomail.NewDialer(s.smtpServer, s.smtpPort, s.username, s.password)
	d.TLSConfig = s.tlsConfig

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
