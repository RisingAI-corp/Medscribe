package transcriber

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type TranscriptTurn struct {
	Speaker   string  `json:"speaker"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	Text      string  `json:"text"`
}

type Transcription interface {
	Transcribe(ctx context.Context, audio []byte) (string, error)
	TranscribeWithDiarization(ctx context.Context, audio []byte) ([]TranscriptTurn, error)
}

func DiarizedTranscriptToString(transcript []TranscriptTurn) (string, error) {
	jsonData, err := json.Marshal(transcript)
	if err != nil {
		return "", fmt.Errorf("failed to marshal transcript turns to JSON: %w", err)
	}
	return string(jsonData), nil
}

func StringToDiarizedTranscript(transcript string) ([]TranscriptTurn, error) {
	var turns []TranscriptTurn
	err := json.Unmarshal([]byte(transcript), &turns)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to TranscriptTurn slice: %w", err)
	}
	return turns, nil
}

func CompressDiarizedText(transcript string) (string, error) {
	if len(transcript) == 0 {
		return "", fmt.Errorf("empty transcript provided")
	}

	fmt.Println("compressing",transcript)
	transcriptTurns, err := UnmarshalTranscript([]byte(transcript))
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal transcript: %w", err)
	}

	var diarizedText strings.Builder
	for _, turn := range transcriptTurns {
		diarizedText.WriteString(fmt.Sprintf("[%s]: %s\n", turn.Speaker, turn.Text))
	}
	return diarizedText.String(), nil
}

// UnmarshalTranscript turns the JSON byte slice into a slice of transcriber.TranscriptTurn.
func UnmarshalTranscript(jsonData []byte) ([]TranscriptTurn, error) {
	var transcript []TranscriptTurn
	err := json.Unmarshal(jsonData, &transcript)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to TranscriptTurn slice: %w", err)
	}
	return transcript, nil
}

