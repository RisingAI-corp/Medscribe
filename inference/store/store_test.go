package inferencestore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/api/option"
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
// 			assert.Equal(t, tc.expectedError, err)
// 			assert.Equal(t, tc.expectedResultIsEmpty, len(response.Content) == 0)
// 		})
// 	}
// }

func TestNewGeminiInferenceStore(t *testing.T) {
	setupEnv(t)
	testCases := []struct {
		name                   string
		apiKey                 string
		expectedError          error
		expectedResultIsEmpty   bool
		description            string
	}{
		// {
		// 	name:                   "successful response",
		// 	apiKey:                 os.Getenv("GEMINI_API_KEY"),
		// 	expectedError:          nil,
		// 	expectedResultIsEmpty:   false,
		// 	description:            "should return a valid response when request input isn't empty",
		// },
		{
			name:                   "api key is empty",
			apiKey:                 "",
			expectedError:          errors.New("api-key cannot be empty"),
			expectedResultIsEmpty:   true,
			description:            "should return an error when api key is empty",
		},
	}

	geminiClient, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal("Failed to create genai client", zap.Error(err))
		return
	}
	defer geminiClient.Close()

	store,err := NewGeminiInferenceStore(
		geminiClient,
	)
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			fmt.Println("this is err", err)
			result, _ := store.Query(ctx, "test request", 10)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResultIsEmpty, result.Content == "")
		})
	}
}


// func TestQuery_InvalidRole(t *testing.T) {
// 	setupEnv(t)

// 	store := NewInferenceStore(os.Getenv("OPENAI_API_CHAT_URL"), os.Getenv("OPENAI_API_KEY"))
// 	invalidRole := "invalid_role"
// 	_, err := store.Query("invalid_role", "test request")

// 	assert.Error(t, err, fmt.Errorf("role must be '%s' or '%s' not %s", SystemRole, UserRole, invalidRole))

// 	expectedError := "role must be 'system' or 'user' not invalid_role"
// 	assert.EqualError(t, err, expectedError, "Expected error message mismatch")
// }
