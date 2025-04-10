package inferenceService

import (
	Chat "Medscribe/inference/store"
	contextLogger "Medscribe/logger"
	"Medscribe/reports"
	reportsTokenUsage "Medscribe/reportsTokenUsageStore"
	Transcription "Medscribe/transcription"
	"Medscribe/user"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var contentSections = []string{reports.Subjective, reports.Objective, reports.AssessmentAndPlan, reports.Summary, reports.PatientInstructions}

// InferenceService defines the methods for interacting with the inference service
type InferenceService interface {
	GenerateReportPipeline(ctx context.Context, report *ReportRequest, contentChan chan ContentChanPayload) error
	RegenerateReport(ctx context.Context, contentChan chan ContentChanPayload, report *ReportRequest) error
	LearnStyle(ctx context.Context, providerID, contentSection, previous, content string) error
}

type inferenceService struct {
	reportsStore         reports.Reports
	transcriptionService Transcription.Transcription
	chat                 Chat.InferenceStore
	userStore            user.UserStore
	reportTokenUsageStore reportsTokenUsage.TokenUsageStore
}

// NewInferenceService creates a new instance of InferenceService with the provided dependencies.
//
// Parameters:
// - reportsStore: An instance of Reports.Reports to handle report storage operations.
// - transcriptionService: An instance of Transcription.Transcription to handle transcription operations.
// - chat: An instance of Chat.InferenceStore to handle chat-related operations.
// - userStore: An instance of user.UserStore to handle user-related operations.
//
// Returns:
// - An instance of InferenceService initialized with the provided dependencies.
func NewInferenceService(reportsStore reports.Reports, transcriptionService Transcription.Transcription, chat Chat.InferenceStore, userStore user.UserStore, reportTokenUsageStore reportsTokenUsage.TokenUsageStore) InferenceService {
	return &inferenceService{
		userStore:            userStore,
		reportsStore:         reportsStore,
		transcriptionService: transcriptionService,
		chat:                 chat,
		reportTokenUsageStore: reportTokenUsageStore,

	}
}

// ContentChanPayload represents a key-value pair for sending content to the frontend.
type ContentChanPayload struct {
	Key   string
	Value interface{}
}

// ReportContentSection represents a section of a report with a specific content type and content.
// ContentType specifies the type of content (e.g., text, image, etc.).
// Content holds the actual content of the section.
type ReportContentSection struct {
	ContentType string
	Content     string
}

// ReportRequest holds the configuration request for generating a medical report.
type ReportRequest struct {
	ID                        string
	PatientName               string
	AudioBytes                []byte
	TranscribedAudio          string
	ProviderID                string
	ProviderName              string
	Timestamp                 time.Time
	Duration                  float64
	Updates                   bson.D
	SubjectiveContent         string
	ObjectiveContent          string
	AssessmentAndPlanContent  string `json:"assessmentAndPlanContent"`
	PatientInstructionContent string
	SummaryContent            string
	SubjectiveStyle           string `bson:"subjectiveStyle"`
	ObjectiveStyle            string `bson:"objectiveStyle"`
	AssessmentAndPlanStyle    string `bson:"assessmentStyle"`
	SummaryStyle              string `bson:"summaryStyle"`
	SessionSummary            string
	CondensedSummary          string
	PatientInstructionsStyle  string `bson:"patientInstructionsStyle"`
	LastVisitID string 
	VisitContext string
}

// GenerateReportPipeline is a pipeline that generates a report based on the given audio bytes and report configuration.
// parameters:
// - ctx: The context of the request.
// - report: The report configuration.
// - contentChan: A channel to send content to the frontend.
// returns:
// - An error if the pipeline fails.
func (s *inferenceService) GenerateReportPipeline(ctx context.Context, report *ReportRequest, contentChan chan ContentChanPayload) error {
	logger := contextLogger.FromCtx(ctx)
	defer close(contentChan)

	// Stage 1: Create pre-configured report
	logger.Info("Pipeline: creating new report entry")
	reportID, err := s.reportsStore.Put(ctx, report.PatientName, report.ProviderID, report.Timestamp, report.Duration, false, reports.THEY, report.LastVisitID)
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error storing report: %w", err)
	}
	contentChan <- ContentChanPayload{Key: "_id", Value: reportID}


	// Stage 2: Transcribe audio
	logger.Info("Pipeline: starting transcription")
	start := time.Now()
	transcribedAudio, err := s.transcriptionService.Transcribe(ctx, report.AudioBytes)
	elapsed := time.Since(start)
	log.Printf("transcription took %s", elapsed) // You can convert this to structured logging if preferred
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error transcribing audio: %w", err)
	}
	report.TranscribedAudio = transcribedAudio

	// Stage 3: Generate report sections (SOAP + summary)
	tokenUsage := make(map[string]int)

	defer func() {
		reportIDtoPrimitive, err := primitive.ObjectIDFromHex(reportID)
		if err != nil {
			return
		}
		tokenEntry := reportsTokenUsage.TokenUsageEntry{
			ReportID:   reportIDtoPrimitive,
			ProviderID: report.ProviderID,
			Timestamp:  primitive.NewDateTimeFromTime(time.Now()),
			TokenUsage: tokenUsage,
		}
		_ = s.reportTokenUsageStore.Insert(ctx, tokenEntry)
		logger.Info("Pipeline: token usage recorded")
	}()
	
	logger.Info("Pipeline: generating report sections")
	combinedUpdates, err := s.generateSoapSections(ctx, report, contentChan, bson.D{}, tokenUsage)
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error generating report sections: %w", err)
	}
	// Stage 4: Notify client report finished generating
	contentChan <- ContentChanPayload{Key: reports.FinishedGenerating, Value: true}

	// Stage 5: Update the report with generated content
	combinedUpdates = append(combinedUpdates, bson.E{Key: reports.FinishedGenerating, Value: true}, bson.E{Key: reports.Transcript, Value: transcribedAudio})
	logger.Info("Pipeline: updating report in store")
	if err := s.reportsStore.UpdateReport(ctx, reportID, combinedUpdates); err != nil {
		return fmt.Errorf("GenerateReportPipeline: error updating report: %w", err)
	}

	return nil
}


// RegenerateReport regenerates the SOAP content based on key-value updates.
// probably will not make reportContents a pointer. it doesn't seem like it will have a high access pattern
func (s *inferenceService) RegenerateReport(
	ctx context.Context,
	contentChan chan ContentChanPayload,
	report *ReportRequest,
) error {
	logger := contextLogger.FromCtx(ctx)
	defer close(contentChan)

	if report.Updates == nil {
		logger.Info("Regeneration aborted: no updates provided")
		return fmt.Errorf("RegenerateReport: no updates provided")
	}

	// Stage 1: Validate update keys
	logger.Info("Validating update keys")
	allowedKeys := map[string]bool{
		reports.Pronouns:        true,
		reports.VisitType:       true,
		reports.PatientOrClient: true,
		reports.IsFollowUp:      true,
		reports.LastVisitID:     true,
	}
	for _, update := range report.Updates {
		if !allowedKeys[update.Key] {
			logger.Info("Invalid update key encountered", zap.String("Key", update.Key))
			return fmt.Errorf("invalid update key: %s", update.Key)
		}
	}

	// Stage 2: Pre-mark report as loading (not finished generating)
	logger.Info("Updating report with pre-generation state")
	preUpdates := append(report.Updates, bson.D{{Key: reports.FinishedGenerating, Value: false}}...)
	if err := s.reportsStore.UpdateReport(ctx, report.ID, preUpdates); err != nil {
		return fmt.Errorf("RegenerateReport: error updating loading status before report regeneration: %w", err)
	}

	// Stage 3: Regenerate SOAP sections
	logger.Info("Generating report sections")
	tokenUsage := make(map[string]int)
	combinedUpdates, err := s.generateSoapSections(ctx, report, contentChan, report.Updates, tokenUsage)
	if err != nil {
		return fmt.Errorf("RegenerateReport: error generating report sections while regenerating report: %w", err)
	}

	// Stage 4: Notify client and finalize
	contentChan <- ContentChanPayload{Key: reports.FinishedGenerating, Value: true}

	combinedUpdates = append(combinedUpdates, bson.D{{Key: reports.FinishedGenerating, Value: true}}...)
	logger.Info("Updating report with regenerated content")
	if err := s.reportsStore.UpdateReport(ctx, report.ID, combinedUpdates); err != nil {
		return fmt.Errorf("RegenerateReport: error updating report after regeneration: %w", err)
	}
	return nil
}

// LearnStyle learns the style from the given report and content section.
// LearnStyle learns the style from the given report and content section.
func (s *inferenceService) LearnStyle(ctx context.Context, providerID, contentSection, previous, current string) error {
	logger := contextLogger.FromCtx(ctx)

	if current == "" {
		return errors.New("cannot learn from empty content")
	}

	styleField, err := styleFieldFromContentSection(contentSection)
	if err != nil {
		return fmt.Errorf("LearnStyle: invalid content section%w", err)
	}

	logger.Info("LearnStyle: generating learning prompt and querying chat model")
	learnStylePrompt := GenerateLearnStylePrompt(contentSection, previous, current)
	response, err := s.chat.Query(ctx, learnStylePrompt, 100)
	if err != nil {
		return fmt.Errorf("LearnStyle: error querying for style: %w", err)
	}

	logger.Info("LearnStyle: updating style in user store")
	if err = s.userStore.UpdateStyle(ctx, providerID, styleField, response.Content); err != nil {
		return fmt.Errorf("LearnStyle: error updating style: %w", err)
	}
	return nil
}


// generateReportSections generates all sections of the report concurrently.
// It serves as a helper function for both generateReportPipeline and regenerateReport.
func (s *inferenceService) generateSoapSections(
	ctx context.Context,
	reportRequest *ReportRequest,
	contentChan chan ContentChanPayload,
	updates bson.D,
	tokenUsage map[string]int,
) (bson.D, error) {
	logger := contextLogger.FromCtx(ctx)

	logger.Info("SOAP: starting concurrent section generation")

	g, ctx := errgroup.WithContext(ctx)
	updatesChan := make(chan bson.E, len(contentSections))

	combinedUpdates := bson.D{}
	var m sync.Mutex

	aggregateUpdates := func(update ...bson.E) {
		m.Lock()
		combinedUpdates = append(combinedUpdates, update...)
		m.Unlock()
	}

	for _, section := range contentSections {
		section := section // capture loop variable
		g.Go(func() error {
			logger.Info("SOAP: generating section", zap.String("Section", section))

			style, err := reportRequest.styleFromContentSection(section)
			if err != nil {
				return fmt.Errorf("invalid content Section: %w", err)
			}
			content, err := reportRequest.contentFromContentSection(section)
			if err != nil {
				return fmt.Errorf("invalid content Section: %w", err)
			}

			var contentPrompt string
			if reportRequest.TranscribedAudio != "" {
				contentPrompt = GenerateReportContentPrompt(
					reportRequest.TranscribedAudio,
					section,
					style,
					reportRequest.ProviderName,
					reportRequest.PatientName,
					reportRequest.VisitContext,
				)
			} else {
				contentPrompt = RegenerateReportContentPrompt(
					content,
					section,
					style,
					reportRequest.Updates,
					reportRequest.VisitContext,
				)
			}

			queryResult, err := s.generateReportSection(ctx, contentPrompt, section, contentChan)
			tokenUsage[section] = queryResult.Usage.TotalTokens
			if err != nil {
				return fmt.Errorf("error generating report section: %w", err)
			}

			logger.Info("SOAP: section generated", zap.String("Section", section), zap.Int("TokensUsed", queryResult.Usage.TotalTokens))

			if section == reports.Summary {
				logger.Info("SOAP: generating summaries from summary section")
				summaries, err := s.generateSummaries(ctx,queryResult.Content, contentChan, tokenUsage)
				if err != nil {
					return fmt.Errorf("GenerateReport: error generating report sections while regenerating report: %w", err)
				}
				aggregateUpdates(summaries...)
				logger.Info("SOAP: summaries generated")
			}

			aggregateUpdates(bson.E{
				Key: section,
				Value: bson.D{
					{Key: reports.ContentData, Value: queryResult.Content},
					{Key: reports.Loading, Value: false},
				},
			})
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return bson.D{}, fmt.Errorf("failed to generate report sections: %w", err)
	}
	close(updatesChan)

	logger.Info("SOAP: all sections generated successfully")
	return combinedUpdates, nil
}


// generateReportSection generates a single section of the report.
func (s *inferenceService) generateReportSection(
	ctx context.Context,
	queryMessage string,
	field string,
	contentChan chan ContentChanPayload,
) (Chat.InferenceResponse, error) {
	logger := contextLogger.FromCtx(ctx)

	logger.Info("generateReportSection: querying chat model", zap.String("Section", field))

	response, err := s.chat.Query(ctx, queryMessage, Chat.MaxTokens)
	if err != nil {
		return Chat.InferenceResponse{}, fmt.Errorf("error generating report section: %w", err)
	}

	logger.Info("generateReportSection: received content", zap.String("Section", field), zap.Int("TokensUsed", response.Usage.TotalTokens))
	contentChan <- ContentChanPayload{Key: field, Value: response.Content}

	return response, nil
}


func (s *inferenceService) generateSummaries(
	ctx context.Context,
	summary string,
	contentChan chan ContentChanPayload,
	tokenUsage map[string]int,
) (bson.D, error) {
	logger := contextLogger.FromCtx(ctx)

	logger.Info("generateSummaries: generating condensed summary")
	condensed, err := s.chat.Query(ctx, fmt.Sprintf(condensedSummary, summary), Chat.MaxTokens)
	if err != nil {
		return bson.D{}, fmt.Errorf("error generating condensed summary: %w", err)
	}
	logger.Info("generateSummaries: condensed summary complete", zap.Int("TokensUsed", condensed.Usage.TotalTokens))

	logger.Info("generateSummaries: generating session summary")
	session, err := s.chat.Query(ctx, fmt.Sprintf(sessionSummary, summary), Chat.MaxTokens)
	if err != nil {
		return bson.D{}, fmt.Errorf("error generating session summary: %w", err)
	}
	logger.Info("generateSummaries: session summary complete", zap.Int("TokensUsed", session.Usage.TotalTokens))

	tokenUsage[reports.CondensedSummary] = condensed.Usage.TotalTokens
	tokenUsage[reports.SessionSummary] = session.Usage.TotalTokens

	contentChan <- ContentChanPayload{Key: reports.CondensedSummary, Value: condensed.Content}
	contentChan <- ContentChanPayload{Key: reports.SessionSummary, Value: session.Content}

	return bson.D{
		{Key: reports.CondensedSummary, Value: condensed.Content},
		{Key: reports.SessionSummary, Value: session.Content},
	}, nil
}


// styleFromContentSection returns the style for the given content section.
func (r *ReportRequest) styleFromContentSection(contentSection string) (string, error) {
	switch contentSection {
	case reports.Subjective:
		return r.SubjectiveStyle, nil
	case reports.Objective:
		return r.ObjectiveStyle, nil
	case reports.AssessmentAndPlan:
		return r.AssessmentAndPlanStyle, nil
	case reports.Summary:
		return r.SummaryStyle, nil
	case reports.PatientInstructions:
		return r.PatientInstructionsStyle, nil
	default:
		return "", fmt.Errorf("error extracting style from content section: invalid content section: %s", contentSection)
	}
}

// styleFieldFromContentSection returns the style field for the given content section.
func styleFieldFromContentSection(contentSection string) (string, error) {
	switch contentSection {
	case reports.Subjective:
		return user.SubjectiveStyleField, nil
	case reports.Objective:
		return user.ObjectiveStyleField, nil
	case reports.AssessmentAndPlan:
		return user.AssessmentAndPlanStyleField, nil
	case reports.Summary:
		return user.SummaryStyleField, nil
	case reports.PatientInstructions:
		return user.PatientInstructionsStyleField, nil
	default:
		return "", fmt.Errorf("invalid content section: %s", contentSection)
	}
}

// contentFromContentSection returns the content for the given content section.
func (r *ReportRequest) contentFromContentSection(contentSection string) (string, error) {
	switch contentSection {
	case reports.Subjective:
		return r.SubjectiveContent, nil
	case reports.Objective:
		return r.ObjectiveContent, nil
	case reports.AssessmentAndPlan:
		return r.AssessmentAndPlanContent, nil
	case reports.Summary:
		return r.SummaryContent, nil
	case reports.PatientInstructions:
		return r.PatientInstructionsStyle, nil
	default:
		return "", fmt.Errorf("invalid content section: %s", contentSection)
	}
}
