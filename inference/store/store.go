package inferencestore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// QueryTunables defines constants for tuning query outputs.
const (
	Temperature = 0.7
	TopP        = 0.95
	MaxTokens   = 1000
)

// Message represents a message in a conversation, including the role (e.g., "user" or "system") and content (text and type).
type Message struct {
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

// Payload defines the structure of the request body sent to the OpenAI API, including messages, temperature, top_p, and max_tokens.
type Payload struct {
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	MaxTokens   int       `json:"max_tokens"`
}

// AzureAPIResponse represents the partial structure of the response returned by the OpenAI API, containing choices with message content and role.
type AzureAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
}

// InferenceStore defines the interface for querying the chat models, requiring a Query method to send requests and return responses.
type InferenceStore interface {
	Query(ctx context.Context, request string, tokens int) (string, error)
}

type inferenceStore struct {
	apiUrl string
	apiKey string
}

const SystemRole = "system"
const UserRole = "user"

// NewInferenceStore creates and returns a new instance of an InferenceStore implementation with the provided API URL and API key.
func NewInferenceStore(apiUrl, apiKey string) InferenceStore {
	return &inferenceStore{apiUrl: apiUrl, apiKey: apiKey}
}

// Query sends a request to the OpenAI API with the specified role and message, and returns the assistant's response.
func (i *inferenceStore) Query(ctx context.Context, request string, tokens int) (string, error) {

	if request == "" {
		return "", errors.New("request cannot be empty")
	}

	if tokens < 1 {
		return "", fmt.Errorf("number of tokens has to be greater and 0 not %d", tokens)
	}

	payload := Payload{
		Messages: []Message{
			{
				Role: UserRole,
				Content: []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				}{
					{Type: "text", Text: request},
				},
			},
		},
		Temperature: 0.7,
		TopP:        0.95,
		MaxTokens:   tokens,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", i.apiUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", i.apiKey)

	client := &http.Client{}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var apiResponse AzureAPIResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(apiResponse.Choices) == 0 {
		return "", errors.New("no choices found in the response")
	}
	content := apiResponse.Choices[0].Message.Content
	return content, nil
}
