package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                       string
	MongoURI                  string
	MongoDBName               string
	MongoUserCollection       string
	MongoReportCollection     string
	MongoReportTestCollection string
	MongoReportTokenUsageCollection     string
	MongoFreedVisits          string
	MongoDistillAnalysis      string
	OpenAIChatURL             string
	OpenAISpeechURL           string
	OpenAIAPIKey              string
	GeminiAPIKey            string
	VertexLocation string
	ProjectID string
	GoogleApplicationCredentialsFileContent string
	DeepgramAPIKey            string
	DeepgramAPIURL            string
	JWTSecret                 string
	FreedAuthToken            string
	Port                      string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	env, err := getEnvStrict("ENVIRONMENT", "development")
	if err != nil {
		return nil, err
	}
	isProd := strings.ToLower(env) == "production"

	mongoURI, err := getEnvStrictConditional("MONGODB_URI", "MONGODB_URI_DEV", isProd)
	if err != nil {
		return nil, err
	}
	
	if isProd{
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS","/var/run/secrets/google_application_credentials")
	}

	mongoDBName, err := getEnvStrict("MONGODB_DB", "")
	if err != nil {
		return nil, err
	}
	mongoUserColl, err := getEnvStrict("MONGODB_USER_COLLECTION", "")
	if err != nil {
		return nil, err
	}
	mongoReportColl, err := getEnvStrict("MONGODB_REPORT_COLLECTION", "")
	if err != nil {
		return nil, err
	}
	mongoReportTestColl, err := getEnvStrict("MONGODB_REPORT_TEST_COLLECTION", "")
	if err != nil {
		return nil, err
	}

	mongoReportTokenUsageColl, err := getEnvStrict("MONGODB_REPORTS_TOKEN_USAGE_STORE", "")
	if err != nil {
		return nil, err
	}

	mongoFreedVisits, err := getEnvStrict("MONGODB_FREED_VISITS", "")
	if err != nil {
		return nil, err
	}
	mongoDistillAnalysis, err := getEnvStrict("MONGODB_DISTILL_ANALYSIS", "")
	if err != nil {
		return nil, err
	}
	openAIChatURL, err := getEnvStrict("OPENAI_API_CHAT_URL", "")
	if err != nil {
		return nil, err
	}
	openAISpeechURL, err := getEnvStrict("OPENAI_API_SPEECH_URL", "")
	if err != nil {
		return nil, err
	}
	openAIKey, err := getEnvStrict("OPENAI_API_KEY", "")
	if err != nil {
		return nil, err
	}

	geminiApiKey, err := getEnvStrict("GEMINI_API_KEY", "")
	if err != nil {
		return nil, err
	}
	projectID, err := getEnvStrict("GCP_PROJECT_ID", "")
	if err != nil {
		return nil, err
	}

	vertexLocation, err := getEnvStrict("VERTEX_LOCATION", "")
	if err != nil {
		return nil, err
	}
	deepgramKey, err := getEnvStrict("DEEPGRAM_API_KEY", "")
	if err != nil {
		return nil, err
	}
	deepgramURL, err := getEnvStrict("DEEPGRAM_API_URL", "")
	if err != nil {
		return nil, err
	}
	jwtSecret, err := getEnvStrict("JWT_SECRET", "")
	if err != nil {
		return nil, err
	}
	freedToken, err := getEnvStrict("FREED_AUTH_TOKEN", "")
	if err != nil {
		return nil, err
	}
	port, err := getEnvStrict("PORT", "8080")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Env:                       env,
		MongoURI:                  mongoURI,
		MongoDBName:               mongoDBName,
		MongoUserCollection:       mongoUserColl,
		MongoReportCollection:     mongoReportColl,
		MongoReportTokenUsageCollection: mongoReportTokenUsageColl,
		MongoReportTestCollection: mongoReportTestColl,
		MongoFreedVisits:          mongoFreedVisits,
		MongoDistillAnalysis:      mongoDistillAnalysis,
		OpenAIChatURL:             openAIChatURL,
		OpenAISpeechURL:           openAISpeechURL,
		OpenAIAPIKey:              openAIKey,
		GeminiAPIKey:            geminiApiKey,
		VertexLocation:            vertexLocation,
		ProjectID:                 projectID,
		DeepgramAPIKey:            deepgramKey,
		DeepgramAPIURL:            deepgramURL,
		JWTSecret:                 jwtSecret,
		FreedAuthToken:            freedToken,
		Port:                      port,
	}

	return cfg, nil
}

func getEnvStrict(key, fallback string) (string, error) {
	val := os.Getenv(key)
	if val == "" && fallback == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	} else if val == "" {
		return fallback, nil
	}
	return val, nil
}

func getEnvStrictConditional(prodKey, devKey string, isProd bool) (string, error) {
	if isProd {
		return getEnvStrict(prodKey, "")
	}
	return getEnvStrict(devKey, "")
}
