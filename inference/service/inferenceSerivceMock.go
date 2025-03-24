package inferenceService

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockInferenceService struct {
	mock.Mock
}

func (m *MockInferenceService) GenerateReportPipeline(ctx context.Context, report *ReportRequest, contentChan chan ContentChanPayload) error {
	args := m.Called(report, contentChan)
	return args.Error(0)
}

func (m *MockInferenceService) RegenerateReport(ctx context.Context, contentChan chan ContentChanPayload, report *ReportRequest) error {
	args := m.Called(ctx, contentChan, report)
	return args.Error(0)
}

func (m *MockInferenceService) LearnStyle(ctx context.Context, reportID, contentSection, previous, content string) error {
	args := m.Called(ctx, reportID, contentSection, content)
	return args.Error(0)
}
