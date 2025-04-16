package geminiTranscriber

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"
)

func setupEnv(t *testing.T) {
	t.Helper()
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Failed to load .env file")
	}
}

func loadAudioFile(t *testing.T, filePath string) []byte {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	audioData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	return audioData
}

// func TestTranscribe(t *testing.T) {
// 	setupEnv(t)

// 	ctx := context.Background()

// 	projectID := os.Getenv("GCP_PROJECT_ID")
// 	location := os.Getenv("VERTEX_LOCATION")

// 	if projectID == "" || location == "" {
// 		log.Fatalf("VERTEX_PROJECT_ID and VERTEX_LOCATION environment variables must be set.")
// 		return
// 	}

// 	client, err := genai.NewClient(ctx, &genai.ClientConfig{
// 		Project:  projectID,
// 		Location: location,
// 		Backend:  genai.BackendVertexAI,
// 	})
// 	assert.NoError(t, err)

// 	geminiTranscriber := NewGeminiTranscriberStore(client)

// 	// Define test cases
// 	testCases := []struct {
// 		name        string
// 		audioData   []byte
// 		expectErr   error
// 		expectEmpty bool
// 	}{
// 		{
// 			name:        "should return populated transcript when supplied with audio data",
// 			audioData:   loadAudioFile(t, "../../testdata/sample1.wav"),
// 			expectErr:   nil,
// 			expectEmpty: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctx := context.Background()
// 			transcript, err := geminiTranscriber.TranscribeWithDiarization(ctx, tc.audioData)
// 			assert.Equal(t, err, tc.expectErr)
// 			assert.Equal(t, tc.expectEmpty, len(transcript) == 0)
// 		})
// 	}
// }

func TestTranscribeWithDiarization(t *testing.T) {
	setupEnv(t)

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

	geminiTranscriber := NewGeminiTranscriberStore(client)

	// Define test cases
	testCases := []struct {
		name        string
		audioData   []byte
		expectErr   error
		expectEmpty bool
	}{
		{
			name:        "should return populated transcript when supplied with audio data",
			audioData:   loadAudioFile(t, "../../testdata/sample3.m4a"),
			expectErr:   nil,
			expectEmpty: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			transcript, err := geminiTranscriber.TranscribeWithDiarization(ctx, tc.audioData)
			fmt.Println("transcript",transcript)
			assert.Equal(t, err, tc.expectErr)
			assert.Equal(t, tc.expectEmpty, len(transcript) == 0)
		})
	}
}
