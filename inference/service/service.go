package inferenceService

import (
	Chat "Medscribe/inference/store"
	"Medscribe/reports"
	Transcription "Medscribe/transcription"
	"Medscribe/user"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	fmt.Println("this is trasncription ", s.transcriptionService)
	defer close(contentChan)

	// create pre-configured report
	reportID, err := s.reportsStore.Put(ctx, report.PatientName, report.ProviderID, report.Timestamp, report.Duration, false, reports.THEY, report.LastVisitID)
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error storing report: %w", err)
	}
	contentChan <- ContentChanPayload{Key: "_id", Value: reportID}

	// transcribe audio
	transcribedAudio, err := s.transcriptionService.Transcribe(ctx, report.AudioBytes)
	fmt.Println(transcribedAudio, "check check check")

	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error transcribing audio: %w", err)
	}
	report.TranscribedAudio = transcribedAudio

	//generate the report sections (subjective, objective, assessment, planning, summary)
	combinedUpdates, err := s.generateSoapSections(ctx, report, contentChan, bson.D{})
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error generating report sections: %w", err)
	}
	// indicate to client that report finished generating
	contentChan <- ContentChanPayload{Key: reports.FinishedGenerating, Value: true}

	// batch update the report with the generated content
	combinedUpdates = append(combinedUpdates, bson.E{Key: reports.FinishedGenerating, Value: true}, bson.E{Key: reports.Transcript, Value: transcribedAudio})
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

	combinedUpdates, err := s.generateSoapSections(ctx, report, contentChan, report.Updates)
	if err != nil {
		return fmt.Errorf("RegenerateReport: error generating report sections while regenerating report: %w", err)
	}

	contentChan <- ContentChanPayload{Key: reports.FinishedGenerating, Value: true}

	combinedUpdates = append(combinedUpdates, bson.D{{Key: reports.FinishedGenerating, Value: true}}...)
	if err := s.reportsStore.UpdateReport(ctx, report.ID, combinedUpdates); err != nil {
		return fmt.Errorf("RegenerateReport: error updating report after regeneration: %w", err)
	}

	return nil
}

// LearnStyle learns the style from the given report and content section.
func (s *inferenceService) LearnStyle(ctx context.Context, providerID, contentSection, previous, current string) error {
	if current == "" {
		return errors.New("cannot learn from empty content")
	}

	styleField, err := styleFieldFromContentSection(contentSection)
	if err != nil {
		return fmt.Errorf("LearnStyle: invalid content section%w", err)
	}

	learnStylePrompt := GenerateLearnStylePrompt(contentSection, previous, current)
	response, err := s.chat.Query(ctx, learnStylePrompt, 100)
	if err != nil {
		return fmt.Errorf("LearnStyle: error querying for style: %w", err)
	}

	if err = s.userStore.UpdateStyle(ctx, providerID, styleField, response); err != nil {
		return fmt.Errorf("LearnStyle: error updating style: %w", err)
	}
	return nil
}

// generateReportSections generates all sections of the report concurrently.
// It serves as a helper function for both generateReportPipeline and regenerateReport.
func (s *inferenceService) generateSoapSections(ctx context.Context, reportRequest *ReportRequest, contentChan chan ContentChanPayload, updates bson.D) (bson.D, error) {
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
		g.Go(func() error {
			style, err := reportRequest.styleFromContentSection(section)
			if err != nil {
				return fmt.Errorf("invalid content Section: %w", err)
			}
			content, err := reportRequest.contentFromContentSection(section)
			if err != nil {
				return fmt.Errorf("invalid content Section: %w", err)
			}
			contentPrompt := ""
			if reportRequest.TranscribedAudio != "" { // if there is no transcript then we are regenerating report
				contentPrompt = GenerateReportContentPrompt(reportRequest.TranscribedAudio, section, style, reportRequest.ProviderName, reportRequest.PatientName, reportRequest.VisitContext)
			} else {
				contentPrompt = RegenerateReportContentPrompt(content, section, style, reportRequest.Updates,reportRequest.VisitContext)
			}

			text, err := s.generateReportSection(ctx, contentPrompt, section, contentChan)
			if err != nil {
				return fmt.Errorf("error generating report section: %w", err)
			}
			if section == reports.Summary {
				summaries, err := s.generateSummaries(text, contentChan)
				if err != nil {
					return fmt.Errorf("GenerateReport: error generating report sections while regenerating report: %w", err)
				}
				aggregateUpdates(summaries...)
			}
			aggregateUpdates(bson.E{Key: section, Value: bson.D{{Key: reports.ContentData, Value: text}, {Key: reports.Loading, Value: false}}})
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return bson.D{}, fmt.Errorf("failed to generate report sections: %w", err)
	}
	close(updatesChan)

	return combinedUpdates, nil
}

// generateReportSection generates a single section of the report.
func (s *inferenceService) generateReportSection(ctx context.Context, queryMessage string, field string, contentChan chan ContentChanPayload) (string, error) {
	response, err := s.chat.Query(ctx, queryMessage, Chat.MaxTokens)
	if err != nil {
		return "", fmt.Errorf("error generating report section: %w", err)
	}

	contentChan <- ContentChanPayload{Key: field, Value: response}

	// Send the update to the updates channel.
	return response, nil
}

func (s *inferenceService) generateSummaries(summary string, contentChan chan ContentChanPayload) (bson.D, error) {
	condensed, err := s.chat.Query(context.Background(), fmt.Sprintf(condensedSummary, summary), Chat.MaxTokens)
	if err != nil {
		return bson.D{}, fmt.Errorf("error generating condensed summary: %w", err)
	}

	session, err := s.chat.Query(context.Background(), fmt.Sprintf(sessionSummary, summary), Chat.MaxTokens)
	if err != nil {
		return bson.D{}, fmt.Errorf("error generating session summary: %w", err)
	}

	contentChan <- ContentChanPayload{Key: reports.CondensedSummary, Value: condensed}
	contentChan <- ContentChanPayload{Key: reports.SessionSummary, Value: session}
	return bson.D{{Key: reports.CondensedSummary, Value: condensed}, {Key: reports.SessionSummary, Value: session}}, nil
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
