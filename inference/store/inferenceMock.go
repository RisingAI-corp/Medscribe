package inferencestore

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockInferenceStore struct {
	mock.Mock
}
func (m *MockInferenceStore) Query(ctx context.Context, request string, tokens int) (InferenceResponse, error) {
	args := m.Called(ctx, request, tokens)
	return args.Get(0).(InferenceResponse), args.Error(1)
}

