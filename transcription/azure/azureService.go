package azure

import (
	transcriber "Medscribe/Transcription"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type azureTranscriber struct {
	apiKey string
	diarizationURL string
	apiUrl string
}

type azureResponse struct {
	CombinedPhrases []struct {
		ChannelId int    `json:"channelId"`
		SpeakerId string `json:"speakerId"`
		Text      string `json:"text"`
	} `json:"combinedPhrases"`
}

type azureResponseWithDiarization struct {
	DurationMilliseconds int `json:"durationMilliseconds"`
	CombinedPhrases    []struct {
		Text string `json:"text"`
	} `json:"combinedPhrases"`
	Phrases []struct {
		Speaker            int     `json:"speaker"` // This is the key field for diarization
		OffsetMilliseconds int     `json:"offsetMilliseconds"`
		DurationMilliseconds int     `json:"durationMilliseconds"`
		Text               string  `json:"text"`
		Words              []struct {
			Text               string  `json:"text"`
			OffsetMilliseconds int     `json:"offsetMilliseconds"`
			DurationMilliseconds int     `json:"durationMilliseconds"`
		} `json:"words"`
		Locale     string  `json:"locale"`
		Confidence float64 `json:"confidence"`
	} `json:"phrases"`
}

func NewAzureTranscriber(apiUrl,diarizationURL, apiKey string) transcriber.Transcription {
	return &azureTranscriber{apiUrl: apiUrl, apiKey: apiKey, diarizationURL: diarizationURL}
}
func (t *azureTranscriber) doAzureRequest(ctx context.Context, apiURL string, audio []byte, definition map[string]interface{}) (io.ReadCloser, error) {
	if len(audio) == 0 {
		fmt.Println("there was no audio provided for the request")
		return nil, nil // Or return an error if empty input should be an error
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	audioPart, err := writer.CreateFormFile("audio", "audio")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file for audio: %w", err)
	}
	if _, err := audioPart.Write(audio); err != nil {
		return nil, fmt.Errorf("failed to write audio data to form file: %w", err)
	}

	definitionJSON, err := json.Marshal(definition)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal definition: %w", err)
	}
	if err := writer.WriteField("definition", string(definitionJSON)); err != nil {
		return nil, fmt.Errorf("failed to write definition field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close form writer: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", t.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close() // Close the body here to avoid resource leaks
		return nil, fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, string(b))
	}

	return resp.Body, nil
}

func (t *azureTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	fmt.Println("entered azure transcriber")

	definition := map[string]interface{}{
		"locales":             []string{"en-US"},
		"profanityFilterMode": "Masked",
		"channels":            []int{0, 1},
	}

	respBody, err := t.doAzureRequest(ctx, t.apiUrl, audio, definition)
	if err != nil {
		return "", err
	}
	if respBody == nil {
		return "", nil // Handle case where no audio was provided
	}
	defer respBody.Close()

	var response azureResponse
	decoder := json.NewDecoder(respBody)
	if err := decoder.Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("this is transcript here", response.CombinedPhrases[0].Text)
	if len(response.CombinedPhrases) > 0 {
		return response.CombinedPhrases[0].Text, nil
	}

	return "", fmt.Errorf("transcript not found in the response")
}

func (t *azureTranscriber) TranscribeWithDiarization(ctx context.Context, audio []byte) ([]transcriber.TranscriptTurn, error) {
	fmt.Println("entered azure transcriber with diarization")

	// Modify the definition to enable speaker diarization
	definition := map[string]interface{}{
		"locales":             []string{"en-US"},
		"profanityFilterMode": "Masked",
		"diarization": map[string]interface{}{
			"enabled":     true,
			"maxSpeakers": 2,
			//"minSpeakers": 2, // Remove or comment out minSpeakers if it's causing issues
		},
	}

	respBody, err := t.doAzureRequest(ctx, t.diarizationURL, audio, definition)
	if err != nil {
		return nil, err
	}
	if respBody == nil {
		return nil, nil // Handle case where no audio was provided
	}
	defer respBody.Close()

	var response azureResponseWithDiarization
	decoder := json.NewDecoder(respBody)
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response for diarization: %w", err)
	}
	fmt.Println("this is transcript here", response.Phrases)
	transcriptTurns := make([]transcriber.TranscriptTurn, 0)
	for _, phrase := range response.Phrases {
		turn := transcriber.TranscriptTurn{
			Speaker:   "Speaker" + strconv.Itoa(phrase.Speaker), // Use the "speaker" field from the response
			StartTime: float64(phrase.OffsetMilliseconds) / 1000.0,
			EndTime:   float64(phrase.OffsetMilliseconds+phrase.DurationMilliseconds) / 1000.0,
			Text:      phrase.Text,
		}
		transcriptTurns = append(transcriptTurns, turn)
	}
	return transcriptTurns,nil
}