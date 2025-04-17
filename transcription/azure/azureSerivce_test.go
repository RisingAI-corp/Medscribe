package azure

import (
	transcriber "Medscribe/transcription"
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

// func TestTranscribe(t *testing.T) {
// 	setupEnv(t)

// 	apiUrl := os.Getenv("OPENAI_API_SPEECH_URL")
// 	apiKey := os.Getenv("OPENAI_API_KEY")

// 	assert.NotEmpty(t, apiUrl, "OPENAI_API_SPEECH_URL should not be empty")
// 	assert.NotEmpty(t, apiKey, "OPENAI_API_KEY should not be empty")

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
// 		{
// 			name:        "should return empty string when supplied with no audio data",
// 			audioData:   []byte{},
// 			expectErr:   nil,
// 			expectEmpty: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctx := context.Background()

// 			txn := NewAzureTranscriber(apiUrl, apiKey)

// 			transcript, err := txn.Transcribe(ctx, tc.audioData)
// 			assert.Equal(t, err, tc.expectErr)
// 			assert.Equal(t, tc.expectEmpty, len(transcript) == 0)
// 		})
// 	}
// }

func TestTranscribeWithDiarization(t *testing.T) {
	// setupEnv(t)

	// apiURL := os.Getenv("OPENAI_API_SPEECH_URL")
	// apiDiarizationURL := os.Getenv("OPENAI_API_DIARIZATION_SPEECH_URL")
	// apiKey := os.Getenv("OPENAI_API_KEY")

	// assert.NotEmpty(t, apiURL, "OPENAI_API_SPEECH_URL should not be empty")
	// assert.NotEmpty(t, apiDiarizationURL, "OPENAI_API_DIARIZATION_SPEECH_URL should not be empty")
	// assert.NotEmpty(t, apiKey, "OPENAI_API_KEY should not be empty")

	// ctx := context.Background()
	// txn := NewAzureTranscriber(apiURL, apiDiarizationURL, apiKey)
	// audioData := loadAudioFile(t, "../../testdata/sample1.wav")

	// transcript, err := txn.TranscribeWithDiarization(ctx, audioData)
	// assert.NoError(t, err)
	// assert.NotEmpty(t, transcript, "Transcript should not be empty")

	// turnsString, err := transcriber.DiarizedTranscriptToString(transcript)
	// assert.NoError(t, err)
	// assert.NotEmpty(t, turnsString, "Diarized transcript should not be empty")

	// turns, err := transcriber.StringToDiarizedTranscript(turnsString)
	// assert.NoError(t, err)
	// assert.NotEmpty(t, turns, "Diarized transcript should not be empty")

	compressedDiarizedText, err := transcriber.CompressDiarizedText(`[{"speaker":"provider","startTime":0,"endTime":0,"text":"1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19"}]`)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressedDiarizedText, "Compressed diarized transcript should not be empty")
}