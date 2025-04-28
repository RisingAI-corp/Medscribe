package emailsender

import (
	"crypto/tls"
	"testing"
)

// TestSendEmailSuccess tests sending a successful email.
func TestSendEmailSuccess(t *testing.T) {
	// Porkbun SMTP settings (REPLACE WITH YOUR ACTUAL DETAILS FOR TESTING)
	smtpServer := "smtp-relay.brevo.com"
	smtpPort := 587 // Or 465
	username := "8b7a41001@smtp-brevo.com"             // REPLACE WITH YOUR ACTUAL EMAIL
	password := "Kx13NYGCRnTQ6mOU"    // REPLACE WITH YOUR ACTUAL PASSWORD - USE WITH CAUTION IN TESTS
	from := "dev@medscribe.pro"                // REPLACE WITH YOUR ACTUAL EMAIL
	fromName := "Medscribe Team"                   // Optional sender name
	to := []string{"emenikeani3@gmail.com","shahsatya25@gmail.com"}      // REPLACE WITH A VALID RECIPIENT FOR TESTING!
	subject := "Test Email How do you do my guy"
	body := "satya emenike medscribe"
	htmlBody := ""

	// Create a new EmailSender
	sender := NewEmailSenderStore(
		smtpServer,
		smtpPort,
		username,
		password,
		from,
		fromName,
	)

	// Optionally configure TLS (it's usually fine with the defaults for Porkbun)
	if store, ok := sender.(*emailSenderStore); ok {
		store.tlsConfig = &tls.Config{
			InsecureSkipVerify: true, // Allow self-signed or untrusted for *TESTING ONLY*
			ServerName:         smtpServer, // Helps verify the server's certificate
		}
	}

	// Send the email
	err := sender.SendEmail(to, subject, body, htmlBody)
	if err != nil {
		t.Fatalf("TestSendEmailSuccess: Error sending email: %v", err)
	}

	t.Logf("TestSendEmailSuccess: Email sent successfully!")
}