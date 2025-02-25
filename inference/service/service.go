package inferenceService

import (
	Chat "Medscribe/inference/store"
	"Medscribe/reports"
	Transcription "Medscribe/transcription"
	"Medscribe/user"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

var contentSections = []string{reports.Subjective, reports.Objective, reports.Assessment, reports.Planning, reports.Summary}

// InferenceService defines the methods for interacting with the inference service
type InferenceService interface {
	GenerateReportPipeline(ctx context.Context, report *ReportRequest, contentChan chan ContentChanPayload) error
	RegenerateReport(ctx context.Context, contentChan chan ContentChanPayload, report *ReportRequest) error
	LearnStyle(ctx context.Context, reportID, contentSection, content string) error
}

type inferenceService struct {
	reportsStore         reports.Reports
	transcriptionService Transcription.Transcription
	chat                 Chat.InferenceStore
	userStore            user.UserStore
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
func NewInferenceService(reportsStore reports.Reports, transcriptionService Transcription.Transcription, chat Chat.InferenceStore, userStore user.UserStore) InferenceService {
	return &inferenceService{
		userStore:            userStore,
		reportsStore:         reportsStore,
		transcriptionService: transcriptionService,
		chat:                 chat,
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
	ID                string
	PatientName       string
	AudioBytes        []byte
	TranscribedAudio  string
	ReportContents    []ReportContentSection
	ProviderID        string
	Timestamp         time.Time
	Duration          float64
	Updates           bson.D
	SubjectiveContent string
	ObjectiveContent  string
	AssessmentContent string
	PlanningContent   string
	SummaryContent    string
	user.Styles
}

func validateReportContents(reportContents *[]ReportContentSection) error {
	for _, report := range *reportContents {
		switch report.ContentType {
		case reports.Subjective, reports.Objective, reports.Assessment, reports.Planning, reports.Summary:
		default:
			return fmt.Errorf("%s is not a valid contentType", report.ContentType)
		}
	}
	return nil
}

//GenerateReportPipeline is a pipeline that generates a report based on the given audio bytes and report configuration.
//parameters:
// - ctx: The context of the request.
// - report: The report configuration.
// - contentChan: A channel to send content to the frontend.
//returns:
// - An error if the pipeline fails.
func (s *inferenceService) GenerateReportPipeline(ctx context.Context, report *ReportRequest, contentChan chan ContentChanPayload) error {
	defer close(contentChan)
	if err := validateReportContents(&report.ReportContents); err != nil {
		return fmt.Errorf("error validating report config: %w", err)
	}

	// transcribe audio
	transcribedAudio, err := s.transcriptionService.Transcribe(ctx, report.AudioBytes)

	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error transcribing audio: %w", err)
	}
	report.TranscribedAudio = transcribedAudio

	// create pre-configured report
	reportID, err := s.reportsStore.Put(ctx, report.PatientName, report.ProviderID, report.Timestamp, report.Duration, false, reports.THEY)
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error storing report: %w", err)
	}

	// send id to the client
	contentChan <- ContentChanPayload{Key: "_id", Value: reportID}

	//generate the report sections (subjective, objective, assessment, planning, summary)
	combinedUpdates, err := s.generateReportSections(ctx, report, contentChan, bson.D{})
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error generating report sections: %w", err)
	}

	// indicate to client that report finished generating
	contentChan <- ContentChanPayload{Key: reports.FinishedGenerating, Value: true}

	// batch update the report with the generated content
	combinedUpdates = append(combinedUpdates, bson.D{{Key: reports.FinishedGenerating, Value: true}}...)
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
	defer close(contentChan)
	if err := validateReportContents(&report.ReportContents); err != nil {
		return fmt.Errorf("error validating report config: %w", err)
	}

	if report.Updates == nil {
		return fmt.Errorf("RegenerateReport: no updates provided")
	}

	allowedKeys := map[string]bool{
		reports.Pronouns:        true,
		reports.VisitType:       true,
		reports.PatientOrClient: true,
		reports.IsFollowUp:      true,
	}

	for _, update := range report.Updates {
		if !allowedKeys[update.Key] {
			return fmt.Errorf("invalid update key: %s", update.Key)
		}
	}

	preUpdates := append(report.Updates, bson.D{{Key: reports.FinishedGenerating, Value: false}}...)
	if err := s.reportsStore.UpdateReport(ctx, report.ID, preUpdates); err != nil {
		return fmt.Errorf("RegenerateReport: error updating loading status before report regeneration: %w", err)
	}

	combinedUpdates, err := s.generateReportSections(ctx, report, contentChan, report.Updates)
	if err != nil {
		return fmt.Errorf("RegenerateReport: error generating report sections: %w", err)
	}

	contentChan <- ContentChanPayload{Key: reports.FinishedGenerating, Value: true}

	combinedUpdates = append(combinedUpdates, bson.D{{Key: reports.FinishedGenerating, Value: true}}...)
	if err := s.reportsStore.UpdateReport(ctx, report.ID, combinedUpdates); err != nil {
		return fmt.Errorf("RegenerateReport: error updating report after regeneration: %w", err)
	}

	return nil
}

// LearnStyle learns the style from the given report and content section.
func (s *inferenceService) LearnStyle(ctx context.Context, provideID, contentSection, content string) error {
	if content == "" {
		return errors.New("cannot learn from empty content")
	}

	styleField, err := styleFieldFromContentSection(contentSection)
	if err != nil {
		return fmt.Errorf("LearnStyle: invalid content section%w", err)
	}

	learnStylePrompt := GenerateLearnStylePrompt(contentSection, content)
	response, err := s.chat.Query(ctx, learnStylePrompt, 100)
	if err != nil {
		return fmt.Errorf("LearnStyle: error querying for style: %w", err)
	}

	if err = s.userStore.UpdateStyle(ctx, provideID, styleField, response); err != nil {
		return fmt.Errorf("LearnStyle: error updating style: %w", err)
	}
	return nil
}

// generateReportSections generates all sections of the report concurrently.
// It serves as a helper function for both generateReportPipeline and regenerateReport.
func (s *inferenceService) generateReportSections(ctx context.Context, report *ReportRequest, contentChan chan ContentChanPayload, updates bson.D) (bson.D, error) {
	g, ctx := errgroup.WithContext(ctx)
	updatesChan := make(chan bson.E, len(report.ReportContents))

	combinedUpdates := bson.D{}
	doneAggregationChan := make(chan struct{})

	go func() {
		for {
			update, ok := <-updatesChan
			if !ok {
				doneAggregationChan <- struct{}{} // kill signal
				return
			}
			combinedUpdates = append(combinedUpdates, update)
		}
	}()

	for _, section := range contentSections {
		g.Go(func() error {
			style, err := report.styleFromContentSection(section)
			if err != nil {
				return fmt.Errorf("invalid content Section: %w", err)
			}
			content, err := report.contentFromContentSection(section)
			if err != nil {
				return fmt.Errorf("invalid content Section: %w", err)
			}
			contentPrompt := GenerateReportContentPrompt(report.TranscribedAudio, section, style, updates, content)
			return s.generateReportSection(ctx, contentPrompt, section, contentChan, updatesChan)
		})
	}

	if err := g.Wait(); err != nil {
		return bson.D{}, fmt.Errorf("failed to generate report sections: %w", err)
	}

	close(updatesChan)

	// Wait for the aggregator goroutine to finish processing updates.
	<-doneAggregationChan

	return combinedUpdates, nil
}

// generateReportSection generates a single section of the report.
func (s *inferenceService) generateReportSection(ctx context.Context, queryMessage string, field string, contentChan chan ContentChanPayload, aggregateUpdatesChan chan<- bson.E) error {
	response, err := s.chat.Query(ctx, queryMessage, Chat.MaxTokens)
	if err != nil {
		return fmt.Errorf("error generating report section: %w", err)
	}
	contentChan <- ContentChanPayload{Key: field, Value: response}

	// Send the update to the updates channel.
	aggregateUpdatesChan <- bson.E{
		Key: field,
		Value: bson.D{
			{Key: reports.ContentData, Value: response},
			{Key: reports.Loading, Value: false},
		},
	}

	return nil
}



// styleFromContentSection returns the style for the given content section.
func (r *ReportRequest) styleFromContentSection(contentSection string) (string, error) {
	switch contentSection {
	case reports.Subjective:
		return r.SubjectiveStyle, nil
	case reports.Objective:
		return r.ObjectiveStyle, nil
	case reports.Assessment:
		return r.AssessmentStyle, nil
	case reports.Planning:
		return r.PlanningStyle, nil
	case reports.Summary:
		return r.SummaryStyle, nil
	default:
		return "", fmt.Errorf("invalid content section: %s", contentSection)
	}
}

// styleFieldFromContentSection returns the style field for the given content section.
func styleFieldFromContentSection(contentSection string) (string, error) {
	switch contentSection {
	case reports.Subjective:
		return user.SubjectiveStyleField, nil
	case reports.Objective:
		return user.ObjectiveStyleField, nil
	case reports.Assessment:
		return user.AssessmentStyleField, nil
	case reports.Planning:
		return user.PlanningStyleField, nil
	case reports.Summary:
		return user.SummaryStyleField, nil
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
	case reports.Assessment:
		return r.AssessmentContent, nil
	case reports.Planning:
		return r.PlanningContent, nil
	case reports.Summary:
		return r.SummaryContent, nil
	default:
		return "", fmt.Errorf("invalid content section: %s", contentSection)
	}
}
