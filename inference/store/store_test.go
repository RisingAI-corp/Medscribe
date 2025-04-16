package inferencestore

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

// setupEnv loads environment variables for testing
func setupEnv(t *testing.T) {
	t.Helper()
	err := godotenv.Load("../../.env")
	require.NoError(t, err, "Failed to load .env file")
}

func TestNewGeminiVertexInferenceStore(t *testing.T) {
	setupEnv(t)
	testCases := []struct {
		name                  string
		apiKey                string
		expectedError         error
		expectedResultIsEmpty bool
		description           string
	}{
		{
			name:                  "successful response",
			apiKey:                os.Getenv("GEMINI_API_KEY"),
			expectedError:         nil,
			expectedResultIsEmpty: false,
			description:           "should return a valid response when request input isn't empty",
		},
	}

	ctx := context.Background()

	projectID := os.Getenv("GCP_PROJECT_ID")
	location := os.Getenv("VERTEX_LOCATION")

	if projectID == "" || location == "" {
		log.Fatalf("VERTEX_PROJECT_ID and VERTEX_LOCATION environment variables must be set.")
		return
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  projectID,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	assert.NoError(t, err)

	store, err := NewGeminiInferenceStore(
		client,
	)
	assert.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			result, err := store.Query(ctx, "test request", 10)
			fmt.Println("this is result", result, err)
			// assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResultIsEmpty, result.Content == "")
		})
	}
}
