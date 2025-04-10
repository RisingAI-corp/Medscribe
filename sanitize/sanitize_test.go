package sanitize

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func loadEnv(t *testing.T) {
	t.Helper()
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func TestDeidentifyWithReplace(t *testing.T) {
	loadEnv(t)
	// Load environment variables for testing (optional, depending on how you want to test project ID)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "test-project" // Use a default for testing if env var is not set
	}


	textToDeidentify := `
This email address is email@example.com and phone number is 12223334444.
`
	expectedDeidentifiedText := `
This email address is REDACTED and phone number is REDACTED.
`
	ctx := context.Background()
	sanitizedResult, err := SanitizeTranscript(ctx, projectID, textToDeidentify); 
	fmt.Println("result",sanitizedResult)
	assert.Equal(t,expectedDeidentifiedText,sanitizedResult)
	assert.NoError(t, err)
}