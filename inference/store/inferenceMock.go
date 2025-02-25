package inferencestore

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockInferenceStore struct {
	mock.Mock
}

func (m *MockInferenceStore) Query(ctx context.Context, request string, tokens int) (string, error) {
	args := m.Called(ctx, request, tokens)
	return args.String(0), args.Error(1)
}
