package transcriber

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTranscription struct {
	mock.Mock
}

func (m *MockTranscription) Transcribe(ctx context.Context, audio []byte) (string, error) {
	args := m.Called(ctx, audio)
	return args.String(0), args.Error(1)
}

func (m *MockTranscription) TranscribeWithDiarization(ctx context.Context, audio []byte) (TranscriptTurn, error) {
	args := m.Called(ctx, audio)
	return args.Get(0).(TranscriptTurn), args.Error(1)
}

func (m *MockTranscription) TranscribeToDiarizedString(transcript []TranscriptTurn) (string, error) {
	args := m.Called(transcript)
	return args.String(0), args.Error(1)
}


