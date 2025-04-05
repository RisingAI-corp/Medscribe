package azure

import (
	transcriberType "Medscribe/transcription"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type transcriber struct {
	apiKey string
	apiUrl string
}

type azureResponse struct {
	CombinedPhrases []struct {
		Text string `json:"text"`
	} `json:"combinedPhrases"`
}

func NewAzureTranscriber(apiUrl, apiKey string) transcriberType.Transcription {
	return &transcriber{apiUrl: apiUrl, apiKey: apiKey}
}

func (t *transcriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	if len(audio) == 0 {
		return "", nil // Avoid making API calls for empty input
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	audioPart, err := writer.CreateFormFile("audio", "audio")
	if err != nil {
		return "", fmt.Errorf("failed to create form file for audio: %w", err)
	}
	if _, err := audioPart.Write(audio); err != nil {
		return "", fmt.Errorf("failed to write audio data to form file: %w", err)
	}

	definition := map[string]interface{}{
		"locales":             []string{"en-US"},
		"profanityFilterMode": "Masked",
		"channels":            []int{0, 1},
	}
	definitionJSON, _ := json.Marshal(definition)
	if err := writer.WriteField("definition", string(definitionJSON)); err != nil {
		return "", fmt.Errorf("failed to write definition field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close form writer: %w", err)
	}

	req, err := http.NewRequest("POST", t.apiUrl, body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", t.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response azureResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.CombinedPhrases) > 0 {
		return response.CombinedPhrases[0].Text, nil
	}

	return "", fmt.Errorf("transcript not found in the response")
}
