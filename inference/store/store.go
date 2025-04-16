package inferencestore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/genai"
)

// QueryTunables defines constants for tuning query outputs.
const (
	Temperature = 0.7
	TopP        = 0.95
	MaxTokens   = 1000
	SystemRole  = "system"
	UserRole    = "user"
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
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
}

type InferenceResponse struct {
	Content string
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
}

// InferenceStore defines the interface for querying the chat models, requiring a Query method to send requests and return responses.
type InferenceStore interface {
	Query(ctx context.Context, request string, tokens int) (InferenceResponse, error)
}

type inferenceGPTStore struct {
	apiUrl string
	apiKey string
}

func NewGPTInferenceStore(apiUrl, apiKey string) InferenceStore {
	return &inferenceGPTStore{apiUrl: apiUrl, apiKey: apiKey}
}

type inferenceGeminiStore struct {
	client *genai.Client
}

func NewGeminiInferenceStore(client *genai.Client) (InferenceStore, error) {
	return &inferenceGeminiStore{client: client}, nil
}

// Query sends a request to the OpenAI API with the specified role and message, and returns the assistant's response.
// Query sends a request to the OpenAI API with the specified role and message, and returns the assistant's response.
func (i *inferenceGPTStore) Query(ctx context.Context, request string, tokens int) (InferenceResponse, error) {

	if request == "" {
		return InferenceResponse{}, errors.New("request cannot be empty")
	}

	if tokens < 1 {
		return InferenceResponse{}, fmt.Errorf("number of tokens has to be greater and 0 not %d", tokens)
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
		Temperature: Temperature, // Use constants
		TopP:        TopP,        // Use constants
		MaxTokens:   tokens,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return InferenceResponse{}, fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", i.apiUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return InferenceResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", i.apiKey)

	client := &http.Client{}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return InferenceResponse{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	switch resp.StatusCode {
	case http.StatusOK:
		// Successful response, proceed to decode
		var apiResponse AzureAPIResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&apiResponse); err != nil {
			return InferenceResponse{}, fmt.Errorf("error decoding response: %v", err)
		}
		fmt.Println("this is the response from the API", apiResponse)
		if len(apiResponse.Choices) == 0 {
			return InferenceResponse{}, errors.New("no choices found in the response")
		}
		content := apiResponse.Choices[0].Message.Content
		return InferenceResponse{Content: content, Usage: apiResponse.Usage}, nil

	case http.StatusTooManyRequests:
		// Rate limit encountered
		return InferenceResponse{}, fmt.Errorf("rate limit exceeded: %v", resp.Status)

	case http.StatusInternalServerError:
		// Server error, the server might be overloaded or have issues
		return InferenceResponse{}, fmt.Errorf("OpenAI server error: %v", resp.Status)

	case http.StatusServiceUnavailable:
		// Service unavailable, the server might be temporarily down
		return InferenceResponse{}, fmt.Errorf("OpenAI service unavailable: %v", resp.Status)

	default:
		// Other unexpected status codes
		return InferenceResponse{}, fmt.Errorf("unexpected HTTP status code: %v", resp.Status)
	}
}

func (i *inferenceGeminiStore) Query(ctx context.Context, request string, tokens int) (InferenceResponse, error) {
	if i.client == nil {
		return InferenceResponse{}, fmt.Errorf("gemini client is not initialized")
	}

	resp, err := i.client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text(request), nil)
	if err != nil {
		return InferenceResponse{}, fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()
	return InferenceResponse{Content: respText}, nil
}
