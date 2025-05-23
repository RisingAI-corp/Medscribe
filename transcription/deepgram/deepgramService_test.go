package deepgram

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

func TestTranscribe(t *testing.T) {
	setupEnv(t)

	apiUrl := os.Getenv("DEEPGRAM_API_URL")
	apiKey := os.Getenv("DEEPGRAM_API_KEY")

	assert.NotEmpty(t, apiUrl, "DEEPGRAM_API_URL should not be empty")
	assert.NotEmpty(t, apiKey, "DEEPGRAM_API_KEY should not be empty")

	// Define test cases
	testCases := []struct {
		name        string
		audioData   []byte
		expectErr   error
		expectEmpty bool
	}{
		{
			name:        "should return populated transcript when supplied with audio data",
			audioData:   loadAudioFile(t, "../../testdata/sample1.wav"),
			expectErr:   nil,
			expectEmpty: false,
		},
		{
			name:        "should return empty string when supplied with not audio data",
			audioData:   []byte{},
			expectErr:   nil,
			expectEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			txn := NewDeepgramTranscriber(apiUrl, apiKey)

			transcript, err := txn.Transcribe(ctx, tc.audioData)
			assert.Equal(t, err, tc.expectErr)
			assert.Equal(t, tc.expectEmpty, len(transcript) == 0)
		})
	}
}
