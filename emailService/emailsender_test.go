package emailsender

import (
	"crypto/tls"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

// TestSendEmailSuccess tests sending a successful email.
func TestSendEmailSuccess(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("TestSendEmailSuccess: Unable to load environment variables: %v", err)
	}

	// Porkbun SMTP settings (REPLACE WITH YOUR ACTUAL DETAILS FOR TESTING)
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT")) // Or 465
	username := os.Getenv("EMAIL_USERNAME")             // REPLACE WITH YOUR ACTUAL EMAIL
	password := os.Getenv("EMAIL_PASSWORD")    // REPLACE WITH YOUR ACTUAL PASSWORD - USE WITH CAUTION IN TESTS
	from := os.Getenv("EMAIL_FROM")                // REPLACE WITH YOUR ACTUAL EMAIL
	fromName := os.Getenv("FROM_NAME")                   // Optional sender name
	to := "emenikeani3@gmail.com"      // REPLACE WITH A VALID RECIPIENT FOR TESTING!
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
		true,
	)

	// Optionally configure TLS (it's usually fine with the defaults for Porkbun)
	if store, ok := sender.(*emailSenderStore); ok {
		store.tlsConfig = &tls.Config{
			InsecureSkipVerify: true, // Allow self-signed or untrusted for *TESTING ONLY*
			ServerName:         smtpServer, // Helps verify the server's certificate
		}
	}

	// Send the email
	err = sender.SendEmail(to, subject, body, htmlBody)
	if err != nil {
		t.Fatalf("TestSendEmailSuccess: Error sending email: %v", err)
	}

	t.Logf("TestSendEmailSuccess: Email sent successfully!")
}
