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

// TestQuery_RoleValidation tests the Query function for both system and user roles
// func TestQuery(t *testing.T) {
// 	setupEnv(t)

// 	store := NewGPTInferenceStore(os.Getenv("OPENAI_API_CHAT_URL"), os.Getenv("OPENAI_API_KEY"))

// 	testCases := []struct {
// 		name                  string
// 		request               string
// 		tokens                int
// 		expectedError         error
// 		expectedResultIsEmpty bool
// 		description           string
// 	}{
// 		{
// 			name:                  "successful response",
// 			request:               "test request",
// 			expectedError:         nil,
// 			tokens:                10,
// 			expectedResultIsEmpty: false,
// 			description:           "should return a valid response when request input isn't empty",
// 		},
// 		{
// 			name:                  "request body is empty",
// 			request:               "",
// 			expectedError:         errors.New("request cannot be empty"),
// 			expectedResultIsEmpty: true,
// 			description:           "should return an error when request is empty",
// 		},
// 		{
// 			name:                  "SystemRole",
// 			request:               "test request",
// 			expectedError:         fmt.Errorf("number of tokens has to be greater and 0 not %d", -10),
// 			tokens:                -10,
// 			expectedResultIsEmpty: true,
// 			description:           "should return an error when there are too many maxTokens suggestions",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.description, func(t *testing.T) {
// 			ctx := context.Background()
// 			response, err := store.Query(ctx, tc.request, tc.tokens)
// 			fmt.Println("this is response",response)
// 			assert.Equal(t, tc.expectedError, err)
// 			assert.Equal(t, tc.expectedResultIsEmpty, len(response.Content) == 0)
// 		})
// 	}
// }


func TestNewGeminiVertexInferenceStore(t *testing.T) {
	setupEnv(t)
	testCases := []struct {
		name                   string
		apiKey                 string
		expectedError          error
		expectedResultIsEmpty   bool
		description            string
	}{
		{
			name:                   "successful response",
			apiKey:                 os.Getenv("GEMINI_API_KEY"),
			expectedError:          nil,
			expectedResultIsEmpty:   false,
			description:            "should return a valid response when request input isn't empty",
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

	store,err := NewGeminiInferenceStore(
		client,
	)
	assert.NoError(t, err)
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()
			
			result, err := store.Query(ctx, "test request", 10)
			fmt.Println("this is result",result,err)
			// assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResultIsEmpty, result.Content == "")
		})
	}
}





