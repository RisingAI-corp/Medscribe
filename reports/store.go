package reports

import (
	transcriber "Medscribe/transcription"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Constants for gender pronouns, report types, and SOAP sections.

HE and SHE are used to represent male and female gender pronouns.
Patient and Client are used to represent the types of individuals for whom a report is created.
Subjective, Objective, Assessment, and Summary represent the sections of a SOAP (Subjective, Objective, Assessment, and Plan) report.
*/
const (
	HE   = "HE"
	SHE  = "SHE"
	THEY = "They"
)

// Report Fields
const (
	Pronouns = "pronouns"

	Patient = "Patient"
	Client  = "Client"

	Subjective          = "subjective"
	Objective           = "objective"
	AssessmentAndPlan   = "assessmentAndPlan"
	Summary             = "summary"
	PatientInstructions = "patientInstructions"


	Loading = "loading"

	Content = "content"

	ContentData = "data"

	VisitType       = "visitType"
	PatientOrClient = "patientOrClient"

	IsFollowUp = "isFollowUp"

	Name = "name"

	ProviderID = "providerid"

	ID = "_id"

	TimeStamp = "timestamp"

	CondensedSummary = "condensedSummary"
	SessionSummary   = "sessionSummary"

	Transcript         = "transcript"
	DiarizedTranscript = "diarizedtranscript"

	readStatus = "readstatus"

	LastVisitID = "lastvisitid"

	Status = "status"

	tokens = "tokens"

	UsedDiarization = "useddiarizedtranscript"

)

type ReportContent struct {
	Data    string `json:"data"`
	Loading bool   `json:"loading"`
}

// RetrievedReportTranscripts
// includes multiple transcripts both unformatted and formatted(diarized)
type RetrievedReportTranscripts struct {
	ProviderID         string `json:"providerID"`
	Transcript string `json:"transcript"`
	DiarizedTranscript []transcriber.TranscriptTurn `json:"diarizedTranscript"`
	UsedDiarization bool `json:"usedDiarization"`
}

type Report struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProviderID          string             `json:"providerID"`
	Name                string             `json:"name"`
	TimeStamp           primitive.DateTime `json:"timestamp"`
	Duration            float64            `json:"duration"`
	Pronouns            string             `json:"pronouns"`
	IsFollowUp          bool               `json:"isFollowUp"`
	PatientOrClient     string             `json:"patientOrClient"`
	Subjective          ReportContent      `json:"subjective"`
	Objective           ReportContent      `json:"objective"`
	AssessmentAndPlan   ReportContent      `json:"assessmentAndPlan"`
	Summary             ReportContent      `json:"summary"`
	PatientInstructions ReportContent      `json:"patientInstructions"`
	CondensedSummary    string             `json:"condensedSummary"`
	SessionSummary      string             `json:"sessionSummary"`
	Transcript          string             `json:"transcript"`
	ReadStatus          bool               `json:"readStatus"`
	LastVisitID         string             `json:"lastVisitID"`
	Status              string             `json:"status"`
	UsedDiarizedTranscript bool `json:"usedDiarizedTranscript"`
}

type Reports interface {
	Put(ctx context.Context, name, providerID string, timestamp time.Time, duration float64, isFollowUp bool, pronouns string, lastVisitID string, usedDiarization bool) (string, error)
	Get(ctx context.Context, reportId string) (Report, error)
	GetAll(ctx context.Context, userId string) ([]Report, error)
	UpdateReport(ctx context.Context, reportId string, batchedUpdates bson.D) error
	Validate(report *Report) error
	Delete(ctx context.Context, reportId string) error
	GetTranscription(ctx context.Context, reportId string) (RetrievedReportTranscripts, error)
	MarkRead(ctx context.Context, reportId string) error
	MarkUnread(ctx context.Context, reportId string) error
	UpdateStatus(ctx context.Context, reportId string, status string) error
}

type reportsStore struct {
	client *mongo.Collection
}

func NewReportsStore(collection *mongo.Collection) Reports {
	return &reportsStore{client: collection}
}

/* Put partially filled record into reports collection */
func (r *reportsStore) Put(ctx context.Context, name, providerID string, timestamp time.Time, duration float64, isFollowUp bool, pronouns string, lastVisitID string, usedDiarization bool) (string, error) {
	if name == "" {
		return "", errors.New("name cannot be an empty string")
	}

	if providerID == "" {
		return "", errors.New("providerId cannot be an empty string")
	}

	if duration <= 0 {
		return "", errors.New("duration must be greater than 0")
	}

	if pronouns != HE && pronouns != SHE && pronouns != THEY {
		return "", fmt.Errorf("pronouns must be either '%s', '%s', or '%s'", HE, SHE, THEY)
	}

	if lastVisitID != "" {
		_, err := primitive.ObjectIDFromHex(lastVisitID)
		if err != nil {
			return "", fmt.Errorf("invalid last visitID format: %v", err)
		}
	}

	// Initialize the Report struct
	report := Report{
		Name:                name,
		TimeStamp:           primitive.NewDateTimeFromTime(timestamp),
		Duration:            duration,
		ProviderID:          providerID,
		IsFollowUp:          isFollowUp,
		Pronouns:            THEY,
		Subjective:          ReportContent{Loading: true},
		Objective:           ReportContent{Loading: true},
		AssessmentAndPlan:   ReportContent{Loading: true},
		PatientInstructions: ReportContent{Loading: true},
		Summary:             ReportContent{Loading: true},
		LastVisitID:         lastVisitID,
		Status:              "pending",
		UsedDiarizedTranscript: usedDiarization,
	}

	insertResp, err := r.client.InsertOne(ctx, report)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v", err)
	}

	insertID, ok := insertResp.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("unexpected type for InsertedID: %T", insertID)
	}

	return insertID.Hex(), nil
}

/* Get retrieves a report by its unique identifier */
func (r *reportsStore) Get(ctx context.Context, reportId string) (Report, error) {
	objectID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		return Report{}, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{ID: objectID}
	projection := bson.M{Transcript: 0, DiarizedTranscript: 0} // Exclude transcript field
	opts := options.FindOne().SetProjection(projection)

	var retrievedReport Report
	err = r.client.FindOne(ctx, filter, opts).Decode(&retrievedReport)
	if err != nil {
		return Report{}, fmt.Errorf("failed to retrieve report: %v", err)
	}

	return retrievedReport, nil
}

func (r *reportsStore) UpdateStatus(ctx context.Context, reportId string, status string) error {
	if status != "pending" && status != "completed" && status != "failed" {
		return fmt.Errorf("status must be either 'pending', 'completed', or 'failed'")
	}
	objectID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}
	filter := bson.M{ID: objectID}
	update := bson.M{"$set": bson.M{status: status}}
	_, err = r.client.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update report status: %v", err)
	}

	return nil
}

/* GetTranscript retrieves the transcript for a report by its unique identifier */
func (r *reportsStore) GetTranscription(ctx context.Context, reportId string) (RetrievedReportTranscripts, error) {
	objectID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		return RetrievedReportTranscripts{}, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{ID: objectID}
	projection := bson.M{Transcript: 1, ProviderID: 1, UsedDiarization: 1,ID: 0} // Include only transcript and providerID fields
	opts := options.FindOne().SetProjection(projection)

	var partialReport struct {
		ProviderID string 
		Transcript string
		UsedDiarizedTranscript bool
	}

	err = r.client.FindOne(ctx, filter, opts).Decode(&partialReport)
	if err != nil {
		return RetrievedReportTranscripts{}, fmt.Errorf("failed to retrieve transcript: %v", err)
	}

	var transcriptTurns []transcriber.TranscriptTurn
	if partialReport.UsedDiarizedTranscript{
		err = json.Unmarshal([]byte(partialReport.Transcript), &transcriptTurns)
		if err != nil {
			return RetrievedReportTranscripts{}, fmt.Errorf("failed to unmarshal transcript: %v", err)
		}
	}
	fmt.Println("this right here is the partial report", partialReport)

	retrievedTranscript := RetrievedReportTranscripts{
		Transcript:         partialReport.Transcript,
		DiarizedTranscript: transcriptTurns,
		ProviderID: partialReport.ProviderID,
		UsedDiarization: partialReport.UsedDiarizedTranscript,
	}
	return retrievedTranscript, nil
}

/* GetAll retrieves all the reports linked to a userId unique identifier */
func (r *reportsStore) GetAll(ctx context.Context, providerId string) ([]Report, error) {
    if providerId == "" {
        return []Report{}, errors.New("missing provider ID")
    }

    filter := bson.M{ProviderID: providerId}

    //  Set projection to exclude the "transcript" field.
    projection := bson.D{
        {Key: "transcript", Value: 0},
    }
    options := options.Find().SetSort(bson.M{TimeStamp: -1}).SetProjection(projection)

    retrievedReports := []Report{}

    cursor, err := r.client.Find(ctx, filter, options)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve reports: %v", err)
    }

    if err := cursor.All(ctx, &retrievedReports); err != nil {
        return nil, fmt.Errorf("failed to decode reports: %v", err)
    }

    return retrievedReports, nil
}


/* Delete removes a report by its unique identifier */
func (r *reportsStore) Delete(ctx context.Context, reportId string) error {
	objectID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{ID: objectID}
	result, err := r.client.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete report: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("report not found")
	}

	return nil
}
func (r *reportsStore) MarkRead(ctx context.Context, reportId string) error {
	objectID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{ID: objectID}
	update := bson.M{"$set": bson.M{readStatus: true}}

	result, err := r.client.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark report as read: %v", err)
	}
	fmt.Println(result.ModifiedCount, "modified count")
	return nil
}

func (r *reportsStore) MarkUnread(ctx context.Context, reportId string) error {
	objectID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{ID: objectID}
	update := bson.M{"$set": bson.M{readStatus: false}}

	result, err := r.client.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark report as unread: %v", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("report not found")
	}

	return nil
}

func validateUpdateDFS(updateMap map[string]interface{}, reportMap map[string]interface{}) error {
	// Loop through the keys in the updateMap and compare with reportMap
	for key, updateVal := range updateMap {
		reportVal, exists := reportMap[key]
		if !exists {
			return fmt.Errorf("key '%s' not found in reportMap", key)
		}

		// If both values are maps (indicating nested structures), recurse
		if isMap(updateVal) && isMap(reportVal) {
			// Recursively validate the nested maps
			if err := validateUpdateDFS(updateVal.(map[string]interface{}), reportVal.(map[string]interface{})); err != nil {
				return fmt.Errorf("error validating nested field '%s': %v", key, err)
			}
		} else {
			if !typesMatch(updateVal, reportVal) {
				return fmt.Errorf("type mismatch for key '%s': expected '%T', got '%T'", key, reportVal, updateVal)
			}
		}
	}
	return nil
}

// Helper function to check if the value is a map
func isMap(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}

func typesMatch(val1, val2 interface{}) bool {
	switch val1.(type) {
	case int:
		if _, ok := val2.(float64); ok {
			return true
		}
	case float64:
		if _, ok := val2.(int); ok {
			return true
		}
	}
	return reflect.TypeOf(val1) == reflect.TypeOf(val2)
}

// Function to convert the batched updates from bson.D to map[string]interface{}'
func bsonDToStringMap(bsonD bson.D) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, elem := range bsonD {
		// Reject dot notation in keys
		if strings.Contains(elem.Key, ".") {
			return nil, fmt.Errorf("dot notation is not allowed for key '%s'", elem.Key)
		}

		// Handle different BSON value types
		switch v := elem.Value.(type) {
		case string, int, int64, float64, bool:
			result[elem.Key] = v
		case bson.D:
			// Recursively handle nested bson.D
			nestedMap, err := bsonDToStringMap(v)
			if err != nil {
				return nil, err
			}
			result[elem.Key] = nestedMap
		default:
			return nil, fmt.Errorf("unsupported type for key '%s': %v type: %v", elem.Key,v,reflect.TypeOf(v))
		}
	}
	return result, nil
}

// Function to apply the updates to the Report struct
func applyUpdatesToReport(updateMap map[string]interface{}, report *Report) error {
	marshalledData, err := json.Marshal(updateMap)
	if err != nil {
		return fmt.Errorf("error marshaling update data: %v", err)
	}

	err = json.Unmarshal(marshalledData, report)
	if err != nil {
		return fmt.Errorf("error unmarshaling update data into report: %v", err)
	}
	return nil
}

// UpdateReport handles updates for any field in the report after validation
func (r *reportsStore) UpdateReport(ctx context.Context, reportId string, updates bson.D) error {
	// Validate that the reportId is not empty and there are updates
	if reportId == "" {
		return errors.New("reportId cannot be empty")
	}

	if len(updates) == 0 {
		return errors.New("cannot apply zero updates to report")
	}

	// Convert BSON to map
	updateMap, err := bsonDToStringMap(updates)
	if err != nil {
		return fmt.Errorf("error converting BSON to map: %v", err)
	}

	// Initialize a report with default values for validation
	report := Report{
		ID:              primitive.NewObjectID(),
		ProviderID:      "for validation purposes",
		Name:            "for validation purposes",
		TimeStamp:       primitive.NewDateTimeFromTime(time.Now()),
		Duration:        60,
		Pronouns:        HE,
		IsFollowUp:      false,
		PatientOrClient: Patient,
		Subjective: ReportContent{
			Data:    "for validation purposes",
			Loading: false,
		},
		Objective: ReportContent{
			Data:    "for validation purposes",
			Loading: false,
		},
		AssessmentAndPlan: ReportContent{
			Data:    "for validation purposes",
			Loading: true,
		},
		Summary: ReportContent{
			Data:    "for validation purposes",
			Loading: false,
		},
		PatientInstructions: ReportContent{
			Data:    "for validation purposes",
			Loading: false,
		},
		SessionSummary:     "for validation purposes",
		CondensedSummary:   "for validation purposes",
	}

	reportJSON, err := json.Marshal(&report)
	if err != nil {
		return fmt.Errorf("error marshalling report %v", err)
	}

	var reportMap map[string]interface{}
	if err := json.Unmarshal(reportJSON, &reportMap); err != nil {
		return fmt.Errorf("error marshalling report %v", err)
	}

	// Validate the fields of the updateMap
	if err := validateUpdateDFS(updateMap, reportMap); err != nil {
		return err
	}

	// Apply updates to the report
	if err := applyUpdatesToReport(updateMap, &report); err != nil {
		return err
	}

	// Validate the updated report
	err = r.Validate(&report)
	if err != nil {
		return fmt.Errorf("error validating report: %v", err)
	}

	// Convert reportId to ObjectId
	objectId, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	// Perform the update in MongoDB
	result, err := r.client.UpdateOne(ctx, bson.D{{Key: ID, Value: objectId}}, bson.D{{Key: "$set", Value: updates}})
	if err != nil {
		return fmt.Errorf("error updating the report field in MongoDB: %v", err)
	}

	// If no document was matched, return an error
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %s", reportId)
	}

	return nil
}

func (r *reportsStore) Validate(report *Report) error {
	if report.ProviderID == "" {
		return errors.New("ProviderID cannot be empty")
	}

	if report.Name == "" {
		return errors.New("name cannot be empty")
	}

	if report.Duration < 1 {
		return fmt.Errorf("duration must be greater than 0, got: %f", report.Duration)
	}

	if report.Pronouns != HE && report.Pronouns != SHE && report.Pronouns != THEY {
		return fmt.Errorf("pronouns must be either '%s', '%s', or '%s'", HE, SHE, THEY)
	}

	if report.PatientOrClient != Patient && report.PatientOrClient != Client {
		return fmt.Errorf("PatientOrClient must be either '%s' or '%s'", Patient, Client)
	}

	if report.SessionSummary == "" {
		return errors.New("sessionSummary cannot be empty")
	}

	if report.CondensedSummary == "" {
		return errors.New("condensedSummary cannot be empty")
	}

	if report.LastVisitID != "" {
		_, err := primitive.ObjectIDFromHex(report.LastVisitID)
		if err != nil {
			return fmt.Errorf("invalid last visitID format: %v", err)
		}
	}

	return nil
}
