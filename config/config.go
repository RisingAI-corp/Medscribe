package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                   string
	MongoURI              string
	MongoDBName           string
	MongoUserCollection   string
	MongoReportCollection string
	MongoFreedVisits      string
	MongoDistillAnalysis  string
	OpenAIChatURL         string
	OpenAISpeechURL       string
	OpenAIAPIKey          string
	DeepgramAPIKey        string
	DeepgramAPIURL        string
	JWTSecret             string
	FreedAuthToken        string
	Port                  string
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
		Env:                   env,
		MongoURI:              mongoURI,
		MongoDBName:           mongoDBName,
		MongoUserCollection:   mongoUserColl,
		MongoReportCollection: mongoReportColl,
		MongoFreedVisits:      mongoFreedVisits,
		MongoDistillAnalysis:  mongoDistillAnalysis,
		OpenAIChatURL:         openAIChatURL,
		OpenAISpeechURL:       openAISpeechURL,
		OpenAIAPIKey:          openAIKey,
		DeepgramAPIKey:        deepgramKey,
		DeepgramAPIURL:        deepgramURL,
		JWTSecret:             jwtSecret,
		FreedAuthToken:        freedToken,
		Port:                  port,
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
