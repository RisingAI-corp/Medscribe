package integrationtests

import (
	"Medscribe/api/handlers/reportsHandler"
	userhandler "Medscribe/api/handlers/userHandler"
	"Medscribe/api/middleware"
	"Medscribe/api/routes"
	inferenceService "Medscribe/inference/service"
	inferenceStore "Medscribe/inference/store"
	"Medscribe/reports"
	transcriber "Medscribe/transcription"
	"Medscribe/transcription/azure"
	"Medscribe/user"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// constants for cookie names
const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

type TestEnv struct {
	Router           *chi.Mux
	UserStore        user.UserStore
	ReportsStore     reports.Reports
	MongoClient      *mongo.Client
	UserColl         *mongo.Collection
	ReportsColl      *mongo.Collection
	Transcriber      transcriber.Transcription
	InferenceStore   inferenceStore.InferenceStore
	InferenceService inferenceService.InferenceService
}

func SetupTestEnv() (*TestEnv, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI_DEV")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(os.Getenv("MONGODB_DB"))
	userColl := db.Collection(os.Getenv("MONGODB_USER_COLLECTION_DEV"))
	reportsColl := db.Collection(os.Getenv("MONGODB_REPORT_COLLECTION_DEV"))

	userStore := user.NewUserStore(userColl)
	reportsStore := reports.NewReportsStore(reportsColl)

	transcriber := azure.NewAzureTranscriber(
		os.Getenv("OPENAI_API_SPEECH_URL"),
		os.Getenv("OPENAI_API_KEY"),
	)

	inferenceStore := inferenceStore.NewInferenceStore(
		os.Getenv("OPENAI_API_CHAT_URL"),
		os.Getenv("OPENAI_API_KEY"),
	)

	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		transcriber,
		inferenceStore,
		userStore,
	)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set in environment")
	}

	authMiddleWare := middleware.NewAuthMiddleware(jwtSecret, logger, "dev")
	userHandler := userhandler.NewUserHandler(userStore, reportsStore, logger, *authMiddleWare)
	reportsHandler := reportsHandler.NewReportsHandler(reportsStore, inferenceService, userStore, logger)

	router := routes.EntryRoutes(routes.APIConfig{
		UserHandler:      userHandler,
		ReportsHandler:   reportsHandler,
		AuthMiddleware:   authMiddleWare.Middleware,
		LoggerMiddleware: middleware.LoggingMiddleware,
		Logger:           logger,
	})

	return &TestEnv{
		Router:           router,
		UserStore:        userStore,
		ReportsStore:     reportsStore,
		MongoClient:      client,
		UserColl:         userColl,
		ReportsColl:      reportsColl,
		Transcriber:      transcriber,
		InferenceStore:   inferenceStore,
		InferenceService: inferenceService,
	}, nil
}

// CleanupTestData removes all test data from the database
func (env *TestEnv) CleanupTestData() error {
	ctx := context.Background()

	if _, err := env.UserColl.DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("failed to clean user collection: %v", err)
	}

	if _, err := env.ReportsColl.DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("failed to clean reports collection: %v", err)
	}

	return nil
}

// Disconnect closes the MongoDB connection
func (env *TestEnv) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := env.MongoClient.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}

	return nil
}

// Helper function to create a test user
func (env *TestEnv) CreateTestUser(name, email, password string) (string, error) {
	ctx := context.Background()
	return env.UserStore.Put(ctx, name, email, password)
}

// Helper function to create a test report
func (env *TestEnv) CreateTestReport(providerID string) (string, error) {
	ctx := context.Background()
	return env.ReportsStore.Put(ctx, "Test Report", providerID, time.Now(), 30, false, "HE","")
}

// GetTestUser fetches a user by ID
func (env *TestEnv) GetTestUser(userID string) (user.User, error) {
	ctx := context.Background()
	return env.UserStore.Get(ctx, userID)
}

// GetTestReport fetches a report by ID
func (env *TestEnv) GetTestReport(reportID string) (reports.Report, error) {
	ctx := context.Background()
	return env.ReportsStore.Get(ctx, reportID)
}

// GetUserByAuth fetches a user by email and password
func (env *TestEnv) GetUserByAuth(email, password string) (user.User, error) {
	ctx := context.Background()
	return env.UserStore.GetByAuth(ctx, email, password)
}

// GetAllReports fetches all reports for a user
func (env *TestEnv) GetAllReports(userID string) ([]reports.Report, error) {
	ctx := context.Background()
	return env.ReportsStore.GetAll(ctx, userID)
}

// GenerateJWT generates a JWT token for the given userID and adds it to the request's cookies.
func (env *TestEnv) GenerateJWT(req *http.Request, userID string) (*http.Request, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set in environment")
	}

	// Create the token
	claims := &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("error signing token: %v", err)
	}

	refreshTokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("error signing token: %v", err)
	}

	// Add the token to the request's cookies
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: tokenString,
		Path:  "/",
	})
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshTokenString,
		Path:  "/",
	})

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx), nil
}
