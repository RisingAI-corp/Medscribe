package inferenceService

import (
	Chat "Medscribe/inference/store"
	contextLogger "Medscribe/logger"
	"Medscribe/reports"
	reportsTokenUsage "Medscribe/reportsTokenUsageStore"
	transcriber "Medscribe/transcription"
	"Medscribe/user"
	"Medscribe/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// InferenceService defines the methods for interacting with the inference service
type InferenceService interface {
	GenerateReportPipeline(ctx context.Context, report *ReportRequest, w *utils.SafeResponseWriter) error
	RegenerateReport(ctx context.Context, report *ReportRequest,w *utils.SafeResponseWriter) error
	LearnStyle(ctx context.Context, providerID, contentSection, previous, content string) error
}

type inferenceService struct {
	reportsStore          reports.Reports
	transcriptionService  transcriber.Transcription
	chat                  Chat.InferenceStore
	userStore             user.UserStore
	reportTokenUsageStore reportsTokenUsage.TokenUsageStore
	diarization bool
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
func NewInferenceService(reportsStore reports.Reports, transcriptionService transcriber.Transcription, chat Chat.InferenceStore, userStore user.UserStore, reportTokenUsageStore reportsTokenUsage.TokenUsageStore, diarization bool) InferenceService {
	return &inferenceService{
		userStore:             userStore,
		reportsStore:          reportsStore,
		transcriptionService:  transcriptionService,
		chat:                  chat,
		reportTokenUsageStore: reportTokenUsageStore,
		diarization:           diarization,
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
	LastVisitID               string
	VisitContext              string
}

// CreateInitialReportEntry creates the initial report entry in the store.
func (s *inferenceService) createInitialReportEntry(ctx context.Context, report *ReportRequest) (string, error) {
	reportID, err := s.reportsStore.Put(ctx, report.PatientName, report.ProviderID, report.Timestamp, report.Duration, false, reports.THEY, report.LastVisitID, s.diarization)
	if err != nil {
		return "", fmt.Errorf("CreateInitialReportEntry: error storing report: %w", err)
	}
	return reportID, nil
}

func (s *inferenceService) processTranscript(ctx context.Context, reportRequest *ReportRequest) (string, error) {
	if s.diarization {
		return s.processWithDiarization(ctx, reportRequest)
	}
	return s.processWithoutDiarization(ctx, reportRequest)
}

func (s *inferenceService) processWithDiarization(ctx context.Context, reportRequest *ReportRequest) (string, error) {
	diarizedTranscript, err := s.transcriptionService.TranscribeWithDiarization(ctx, reportRequest.AudioBytes)
	if err != nil {
		return "", fmt.Errorf("error creating diarized transcript: %w", err)
	}

	logger := contextLogger.FromCtx(ctx)
	logger.Info("Generated Diarized transcript", zap.Any("diarizedTranscript", diarizedTranscript))

	diarizedTranscriptString,err := transcriber.DiarizedTranscriptToString(diarizedTranscript)
	if err != nil {
		return "", fmt.Errorf("error creating diarized transcript: %w", err)
	}
	diarizedToString, err := transcriber.CompressDiarizedText(diarizedTranscriptString)
	if err != nil {
		return "", fmt.Errorf("error compressing diarized Transcript: %w", err)
	}
	reportRequest.TranscribedAudio = diarizedToString
	return diarizedTranscriptString, nil
}

func (s *inferenceService) processWithoutDiarization(ctx context.Context, reportRequest *ReportRequest) (string, error) {
	transcript, err := s.transcriptionService.Transcribe(ctx, reportRequest.AudioBytes)
	if err != nil {
		return "", fmt.Errorf("error creating transcript: %w", err)
	}
	reportRequest.TranscribedAudio = transcript
	return transcript, nil
}

// RecordTokenUsage records the token usage for the report.
func (s *inferenceService) recordTokenUsage(ctx context.Context, reportID string, providerID string, tokenUsage *utils.SafeMap[int]) error {
	logger := contextLogger.FromCtx(ctx)
	reportIDtoPrimitive, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return fmt.Errorf("RecordTokenUsage: error converting reportID to primitive.ObjectID %w", err)
	}
	tokenEntry := reportsTokenUsage.TokenUsageEntry{
		ReportID:   reportIDtoPrimitive,
		ProviderID: providerID,
		Timestamp:  primitive.NewDateTimeFromTime(time.Now()),
		TokenUsage: tokenUsage.GetMap(),
	}
	if err := s.reportTokenUsageStore.Insert(ctx, tokenEntry); err != nil {
		return fmt.Errorf("RecordTokenUsage: error inserting token usage entry for report %s into store: %w", reportID, err)
	}
	logger.Info("Token usage recorded")
	return nil
}

// UpdateFinalReport updates the report in the store with the generated content and final status.
func (s *inferenceService) updateFinalReport(ctx context.Context, reportID string, combinedUpdates bson.D) error {
	updates := append(combinedUpdates,
		bson.E{Key: reports.Status, Value: "success"},
	)
	if err := s.reportsStore.UpdateReport(ctx, reportID, updates); err != nil {
		return fmt.Errorf("UpdateFinalReport: error updating report: %w", err)
	}
	return nil
}

// sendContentToFrontend writes the payload to the SafeResponseWriter and flushes it after encoding.
func sendContentToFrontend(w *utils.SafeResponseWriter, payload ContentChanPayload) {
    logger := contextLogger.FromCtx(context.Background())
    logger.Info("sendContentToFrontend: Encoding and sending payload to frontend", zap.String("key", payload.Key))

    w.Header().Set("Content-Type", "application/json")

    encoder := json.NewEncoder(w) // changed to use the safe writer.
    if err := encoder.Encode(payload); err != nil {
        logger.Error("sendContentToFrontend: Error encoding payload", zap.String("key", payload.Key), zap.Error(err))
        http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
        return
    }
	w.Flush()
}

// transcriptTurnsToBSONArray takes an array (slice) of transcriber.TranscriptTurn
func (s *inferenceService) GenerateReportPipeline(ctx context.Context, reportRequest *ReportRequest, w *utils.SafeResponseWriter) error {
	logger := contextLogger.FromCtx(ctx)

	var skipDefer bool
	// this will only run on failures as a less redundant way to mark a report as failed
	defer func() {
		if !skipDefer {
			logger.Error("GenerateReportPipeline: Updating Report Status to Failed", zap.String("report_id", reportRequest.ID))
			s.reportsStore.UpdateStatus(ctx, reportRequest.ID, "failed")
		}
	}()

	// Stage 1: Create pre-configured report
	logger.Info("Starting stage 1: creating pre-configured report")
	reportID, err := s.createInitialReportEntry(ctx, reportRequest)
	if err != nil {
		return err
	}
	sendContentToFrontend(w, ContentChanPayload{"_id", reportID})

	// Stage 2: Transcribe audio
	logger.Info("Starting stage 2: transcribing audio")
	rawTranscript, err := s.processTranscript(ctx, reportRequest)
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error creating transcript: %w", err)
	}

	var (
		transcript         string
		diarizedTurns      []transcriber.TranscriptTurn
	)

	if s.diarization {
		diarizedTurns, err = transcriber.StringToDiarizedTranscript(rawTranscript)
		logger.Info("Diarized Turns Generated Transcript", zap.Any("diarizedTurns", diarizedTurns))
		if err != nil {
			return fmt.Errorf("GenerateReportPipeline: error unmarshaling diarized transcript: %w", err)
		}
	} else {
		transcript = rawTranscript
	}

	// Send content to frontend
	sendContentToFrontend(w, ContentChanPayload{
		Key: reports.Transcript,
		Value: reports.RetrievedReportTranscripts{
			Transcript:       transcript,
			DiarizedTranscript: diarizedTurns,
			ProviderID:       reportRequest.ProviderID,
			UsedDiarization:  s.diarization,
		},
	})

	combinedUpdates := bson.D{
		{Key: reports.Transcript, Value: rawTranscript},
		{Key:reports.UsedDiarizationUpdateKey, Value: s.diarization},
	}

	// Stage 3: Generate report sections (SOAP + summary + patient Instructions)
	logger.Info("Starting stage 3: generating report sections")
	tokenUsage := utils.NewSafeMap[int]()
	contentUpdates, err := s.generateSoapSections(ctx, reportRequest, w,tokenUsage)
	if err != nil {
		return fmt.Errorf("GenerateReportPipeline: error generating report sections: %w", err)
	}

	combinedUpdates = append(combinedUpdates, contentUpdates...)

	// Stage 4: Record token usage
	logger.Info("Starting stage 4: recording token usage")
	if err := s.recordTokenUsage(ctx, reportID, reportRequest.ProviderID, tokenUsage); err != nil {
		logger.Error("GenerateReportPipeline: error recording token usage", zap.Error(err))
		// Decide if this error should fail the entire pipeline
		// For now, logging and continuing
	}

	sendContentToFrontend(w, ContentChanPayload{Key: reports.Status, Value: "success"})
	// Stage 5: Update the report with generated content
	logger.Info("Starting stage 5: updating report with generated content")
	if err := s.updateFinalReport(ctx, reportID, combinedUpdates); err != nil {
		return err
	}

	skipDefer = true
	return nil
}

// RegenerateReport regenerates the SOAP content based on key-value updates.
// probably will not make reportContents a pointer. it doesn't seem like it will have a high access pattern
func (s *inferenceService) RegenerateReport(
	ctx context.Context,
	reportRequest *ReportRequest,
	w *utils.SafeResponseWriter,
) error {
	logger := contextLogger.FromCtx(ctx)

	if reportRequest.Updates == nil {
		logger.Info("RegenerateReport: Regeneration aborted: no updates provided")
		return fmt.Errorf("RegenerateReport: no updates provided")
	}

	// Stage 1: Validate update keys
	logger.Info("Regenerating report: Validating update keys")
	allowedKeys := map[string]bool{
		reports.Pronouns:        true,
		reports.VisitType:       true,
		reports.PatientOrClient: true,
		reports.IsFollowUp:      true,
		reports.LastVisitID:     true,
	}
	for _, update := range reportRequest.Updates {
		if !allowedKeys[update.Key] {
			logger.Info("Regeneration aborted: Invalid update key encountered", zap.String("Key", update.Key))
			return fmt.Errorf("invalid update key: %s", update.Key)
		}
	}

	logger.Info("Regenerating report: Updating report with pre-generation state")
	preUpdates := append(reportRequest.Updates, bson.D{{Key: reports.Status, Value: "success"}}...)

	if err := s.reportsStore.UpdateReport(ctx, reportRequest.ID, preUpdates); err != nil {
		return fmt.Errorf("RegenerateReport: error updating loading status before report regeneration: %w", err)
	}

	// Stage 3: Regenerate SOAP sections
	logger.Info("Regenerating report: Generating report sections")
	tokenUsage := utils.NewSafeMap[int]()
	combinedUpdates, err := s.generateSoapSections(ctx, reportRequest, w,tokenUsage)
	if err != nil {
		return fmt.Errorf("RegenerateReport: error generating report sections while regenerating report: %w", err)
	}

	// Stage 4: Notify client and finalize
	sendContentToFrontend(w, ContentChanPayload{Key: reports.Status, Value: "success"})

	combinedUpdates = append(combinedUpdates, bson.D{{Key: reports.Status, Value: "success"}}...)
	logger.Info("Updating report with regenerated content")
	if err := s.reportsStore.UpdateReport(ctx, reportRequest.ID, combinedUpdates); err != nil {
		return fmt.Errorf("RegenerateReport: error updating report after regeneration: %w", err)
	}
	return nil
}

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
	response, err := s.chat.Query(ctx, "", learnStylePrompt, 100)
	if err != nil {
		return fmt.Errorf("LearnStyle: error querying for style: %w", err)
	}

	logger.Info("LearnStyle: updating style in user store")
	if err = s.userStore.UpdateStyle(ctx, providerID, styleField, response.Content); err != nil {
		return fmt.Errorf("LearnStyle: error updating style: %w", err)
	}
	return nil
}

func contentPromptFunc(
	transcript string,
	targetSection string,
	context string,
	style string,
	providerName string,
	patientName string,
	content string,
	updates bson.D,
) string {
	if content == "" {
		cfg := generatePromptConfig{
			transcript:         transcript,
			targetSection:      targetSection,
			context:            context,
			style:              style,
			providerName:       providerName,
			patientName:        patientName,
		}
		return GenerateReportContentPrompt(cfg)
	}

	cfg := regeneratePromptConfig{
		transcript:         transcript,
		targetSection:      targetSection,
		targetContent:      content,
		priorVisitContext:  context,
		providerName:       providerName,
		patientName:        patientName,
		reportUpdates:      updates,
	}
	return RegenerateReportContentPrompt(cfg)
}

func (s *inferenceService) generateSectionPipeline(
	ctx context.Context,
	systemPrompt,
	queryMessage, 
	field string,
	tokenUsage *utils.SafeMap[int],
	aggregator func(...bson.E),
	writer *utils.SafeResponseWriter,
) error {
	logger := contextLogger.FromCtx(ctx)

	// Stage 1: Query chat model
	logger.Info("generateReportSection: querying chat model", zap.String("Section", field))
	response, err := s.chat.Query(ctx, systemPrompt, queryMessage, Chat.MaxTokens)
	if err != nil {
		return fmt.Errorf("error generating report section: %w", err)
	}

	// Stage 2: Record token usage
	tokenUsage.Set(field+"Tokens", response.Usage.TotalTokens)

	// Stage 3: Send content to frontend
	sendContentToFrontend(writer, ContentChanPayload{Key: field, Value: response.Content})

	// Stage 4: Aggregate updates
	aggregator(bson.E{
		Key: field,
		Value: bson.D{
			{Key: reports.ContentData, Value: response.Content},
			{Key: reports.Loading, Value: false},
		},
	})
	return nil
}

// getUpdateValue returns the content of the given field in the combined updates or an empty string if the field is not found.
func getSectionValue(combinedUpdates bson.D, field string) string {
	for _, update := range combinedUpdates {
		if update.Key == field {
			for _, v := range update.Value.(bson.D) {
				if v.Key == reports.ContentData {
					return v.Value.(string)
				}
			}
		}
	}
	return ""
}

// generateReportSections generates all sections of the report concurrently.
// It serves as a helper function for both generateReportPipeline and regenerateReport.
func (s *inferenceService) generateSoapSections(
	ctx context.Context,
	reportRequest *ReportRequest,
	w *utils.SafeResponseWriter,
	tokenUsage *utils.SafeMap[int],
) (bson.D, error) {
	logger := contextLogger.FromCtx(ctx)
	logger.Info("SOAP: starting concurrent section generation")

	g, ctx := errgroup.WithContext(ctx)
	var m sync.Mutex

	combinedUpdates := bson.D{}
	aggregateUpdates := func(update ...bson.E) {
		m.Lock()
		combinedUpdates = append(combinedUpdates, update...)
		m.Unlock()
	}

	stitchSystemPrompt := func (subSystemPrompt string) string{
		return fmt.Sprintf("%s\n%s\n%s\n%s", baseSystemPrompt, subSystemPrompt,defaultReturnFormatSystemPrompt, defaultWarningsSystemPrompt)
	}

	// Generate Subjective Section
	g.Go(func() error {
		contentPrompt := contentPromptFunc(
			reportRequest.TranscribedAudio,
			reports.Subjective,
			reportRequest.VisitContext,
			reportRequest.SubjectiveStyle,
			reportRequest.ProviderName,
			reportRequest.PatientName,
			reportRequest.SubjectiveContent,
			reportRequest.Updates,
		)
		err := s.generateSectionPipeline(ctx, stitchSystemPrompt(subjectiveTaskDescription), contentPrompt, reports.Subjective, tokenUsage, aggregateUpdates, w)
		if err != nil {
			return fmt.Errorf("error generating report section: %w", err)
		}
		return nil
	})

	// Generate Objective Section
	g.Go(func() error {
		contentPrompt := contentPromptFunc(
			reportRequest.TranscribedAudio,
			reports.Objective,
			reportRequest.VisitContext,
			reportRequest.ObjectiveStyle,
			reportRequest.ProviderName,
			reportRequest.PatientName,
			reportRequest.ObjectiveContent,
			reportRequest.Updates,
		)
		err := s.generateSectionPipeline(ctx, stitchSystemPrompt(objectiveTaskDescription), contentPrompt, reports.Objective, tokenUsage, aggregateUpdates, w)
		if err != nil {
			return fmt.Errorf("error generating report section: %w", err)
		}
		return nil
	})

	// Generate Assessment and Plan Section
	g.Go(func() error {
		contentPrompt := contentPromptFunc(
			reportRequest.TranscribedAudio,
			reports.AssessmentAndPlan,
			reportRequest.VisitContext,
			reportRequest.AssessmentAndPlanStyle,
			reportRequest.ProviderName,
			reportRequest.PatientName,
			reportRequest.AssessmentAndPlanContent,
			reportRequest.Updates,
		)
		err := s.generateSectionPipeline(ctx, stitchSystemPrompt(assessmentAndPlanTaskDescription),contentPrompt, reports.AssessmentAndPlan, tokenUsage, aggregateUpdates, w)
		if err != nil {
			return fmt.Errorf("error generating report section: %w", err)
		}
		return nil
	})
	
	// Generate Patient Instructions Section
	g.Go(func() error {
		contentPrompt := contentPromptFunc(
			reportRequest.TranscribedAudio,
			reports.PatientInstructions,
			reportRequest.VisitContext,
			reportRequest.PatientInstructionsStyle,
			reportRequest.ProviderName,
			reportRequest.PatientName,
			reportRequest.PatientInstructionContent,
			reportRequest.Updates,
		)
		err := s.generateSectionPipeline(ctx, stitchSystemPrompt(patientInstruction),contentPrompt, reports.PatientInstructions, tokenUsage, aggregateUpdates, w)
		if err != nil {
			return fmt.Errorf("error generating report section: %w", err)
		}
		return nil
	})


	// Generate Summary and Sub-Summaries
	g.Go(func() error {
		contentPrompt := contentPromptFunc(
			reportRequest.TranscribedAudio,
			reports.Summary,
			reportRequest.VisitContext,
			reportRequest.SummaryStyle,
			reportRequest.ProviderName,
			reportRequest.PatientName,
			reportRequest.SummaryContent,
			reportRequest.Updates,
		)
		err := s.generateSectionPipeline(ctx, stitchSystemPrompt(summaryTaskDescription),contentPrompt, reports.Summary, tokenUsage, aggregateUpdates, w)
		if err != nil {
			return fmt.Errorf("error generating report section: %w", err)
		}

		summary := getSectionValue(combinedUpdates, reports.Summary)
		if summary == "" {
			return fmt.Errorf("error generating report sub-summaries due to no summary generated: %w", err)
		}

		// Generate condensed and session summaries
		err = s.generateSummaries(ctx, summary, tokenUsage, aggregateUpdates, w)
		if err != nil {
			return fmt.Errorf("error generating report section: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	logger.Info("SOAP: all sections generated successfully")
	return combinedUpdates, nil
}

// generateSummaries generates the condensed and session summaries based on the given summary.
// The token usage of the summaries is recorded in the tokenUsage map.
// The aggregator function is called with the generated summaries to update the report.
func (s *inferenceService) generateSummaries(
	ctx context.Context,
	summary string,
	tokenUsage *utils.SafeMap[int],
	aggregator func(...bson.E),
	writer *utils.SafeResponseWriter,
) error {
	logger := contextLogger.FromCtx(ctx)

	// Generate condensed summary
	logger.Info("generateSummaries: generating condensed summary")
	condensed, err := s.chat.Query(ctx, condensedSummary, summary, Chat.MaxTokens)
	if err != nil {
		return fmt.Errorf("error generating condensed summary: %w", err)
	}

	// Generate session summary
	logger.Info("generateSummaries: generating session summary")
	session, err := s.chat.Query(ctx, sessionSummary, summary, Chat.MaxTokens)
	if err != nil {
		return fmt.Errorf("error generating session summary: %w", err)
	}

	// Record usage and send to client
	tokenUsage.Set(reports.CondensedSummary, condensed.Usage.TotalTokens)
	tokenUsage.Set(reports.SessionSummary, session.Usage.TotalTokens)

	sendContentToFrontend(writer, ContentChanPayload{Key: reports.CondensedSummary, Value: condensed.Content})
	sendContentToFrontend(writer, ContentChanPayload{Key: reports.SessionSummary, Value: session.Content})


	// Update report with summaries
	aggregator(
		bson.E{
			Key: reports.CondensedSummary,
			Value: condensed.Content,
		},
		bson.E{
			Key: reports.SessionSummary,
			Value:reports.ContentData,
		},
	)

	return nil
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