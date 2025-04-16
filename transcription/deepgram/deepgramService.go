package deepgram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	transcriber "Medscribe/transcription"
)

type deepGramTranscriber struct {
	apiKey string
	apiUrl string
}

type deepgramResponse struct {
	Results struct {
		Channels []struct {
			Alternatives []struct {
				Transcript string `json:"transcript"`
			} `json:"alternatives"`
		} `json:"channels"`
	} `json:"results"`
}

func NewDeepgramTranscriber(apiUrl, apiKey string) transcriber.Transcription {
	return &deepGramTranscriber{
		apiKey: apiKey, apiUrl: apiUrl}
}
func (t *deepGramTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	if len(audio) == 0 {
		return "", nil // No need to exhaust API usage
	}

	req, err := http.NewRequest("POST", t.apiUrl, bytes.NewReader(audio))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", t.apiKey))

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

	var response deepgramResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Results.Channels) > 0 && len(response.Results.Channels[0].Alternatives) > 0 {
		return response.Results.Channels[0].Alternatives[0].Transcript, nil
	}

	return "", fmt.Errorf("transcript not found in the response")
}

func (t *deepGramTranscriber) TranscribeWithDiarization(ctx context.Context, audio []byte) (transcriber.TranscriptTurn, error) {
	return transcriber.TranscriptTurn{}, nil
}
