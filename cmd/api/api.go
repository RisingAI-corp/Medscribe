package main

import (
	"Medscribe/api/handlers/reportsHandler"
	userhandler "Medscribe/api/handlers/userHandler"
	"Medscribe/api/middleware"
	"Medscribe/api/routes"
	"Medscribe/config"
	inferenceService "Medscribe/inference/service"
	inferencestorre "Medscribe/inference/store"
	contextLogger "Medscribe/logger"
	"Medscribe/reports"
	reportsTokenUsage "Medscribe/reportsTokenUsageStore"
	"Medscribe/transcription/azure"
	"Medscribe/user"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/genai"
	// For google.Credentials
	// For WithCredentialsJSON
)

type mockTranscriber struct{}

const sample1 = `Hello. Hi. Hi, how are you today? Good morning. Good morning. I'm doing well, thank you. And how have things going since we last spoke? Hello. Okay, so how things going since the last time we spoke? Okay. Any medical changes, emergency room, urgent care or doctor's visits I should be aware of? Okay. And how is your sleep, your mood and your appetite? Okay. How many hours of sleep would you say you get in a night? Okay. And if you have to rate your mood, 0 to 10, 10 being the best, how would you rate it? Fine. Okay. And last time we spoke about getting, calling the office to reschedule for therapy. Did that happen? To call the office to see if you can, right now you are not meeting with anybody, right? Okay, okay, good. How often do you do meet with her? And that is going okay? Okay. And what about the drug test? Did you get it done? Because the last time we did was, like I said, was high on alcohol. They said they don't have the paperwork, the lab. When we put it, we usually do it for six months, so they should have it. Okay, let me check. Okay, I just put it in again, okay? So when you have time, if you can get it done before the next appointment, okay? Okay, so add a 20 minute. twice a day it's gonna go to the same pharmacy okay all right I will send it over and we'll call the office and schedule again a month from today okay all right stay well bye bye`

func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return sample1, nil
}

func (m *mockTranscriber) TranscribeWithDiarization(ctx context.Context, audio []byte) (string, error) {
	return sample1, nil
}

type mockInferStore struct{}

func (m *mockInferStore) Query(ctx context.Context, request string, tokens int) (inferencestorre.InferenceResponse, error) {
	return inferencestorre.InferenceResponse{Content: "mock response"}, nil
}

func main() {
	log.Println("🚀 Starting up application")
	log.Println("⚡ ENV BEFORE .env load: PORT =", os.Getenv("PORT"))
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("❌ Critical error loading config: %v", err)
	}

	logger := contextLogger.Get(cfg.Env)
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("⚠️ Error syncing logger: %v\n", err)
		}
	}()
	logger.Info("✅ Configuration loaded", zap.String("env", cfg.Env))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("🌐 Connecting to MongoDB", zap.String("uri", cfg.MongoURI))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Fatal("❌ Failed to connect to MongoDB", zap.Error(err))
	}
	if err = client.Ping(ctx, nil); err != nil {
		logger.Fatal("❌ Failed to ping MongoDB", zap.Error(err))
	}
	logger.Info("✅ Connected to MongoDB")

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logger.Error("⚠️ Error disconnecting MongoDB client", zap.Error(err))
		}
	}()

	db := client.Database(cfg.MongoDBName)
	userColl := db.Collection(cfg.MongoUserCollection)
	reportsColl := db.Collection(cfg.MongoReportCollection) //TODO: Change back to MongoReportCollection

	// creating stores
	userStore := user.NewUserStore(userColl)
	reportsStore := reports.NewReportsStore(reportsColl)

	if cfg.Env == "production" {
		// in production we will use the metadata server to to leverage the cloud run service account to auth with vertex ai
		err := os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		if err != nil {
			logger.Warn("Error unsetting GOOGLE_APPLICATION_CREDENTIALS", zap.Error(err))
		} else {
			logger.Info("GOOGLE_APPLICATION_CREDENTIALS unset for production")
		}
	}

	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  cfg.ProjectID,
		Location: cfg.VertexLocation,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		logger.Fatal("❌ Failed to create gemini client", zap.Error(err))
	}

	GeminiInferenceStore, err := inferencestorre.NewGeminiInferenceStore(geminiClient)
	if err != nil {
		logger.Fatal("❌ Failed to create inference store", zap.Error(err))
	}

	reportsTokenUsage := reportsTokenUsage.NewTokenUsageStore(db.Collection(cfg.MongoReportTokenUsageCollection))

	//creating services
	azureTranscriber := azure.NewAzureTranscriber(cfg.OpenAISpeechURL, cfg.OpenAIDiarizationSpeechURL, cfg.OpenAIAPIKey)
	// geminiTranscriber := geminiTranscriber.NewGeminiTranscriberStore(geminiClient)
	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		// &mockTranscriber{},
		// geminiTranscriber,
		azureTranscriber,
		GeminiInferenceStore,
		// &mockInferStore{},
		userStore,
		reportsTokenUsage,
		false,
	)

	// instantiating api
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, logger, cfg.Env)
	userHandler := userhandler.NewUserHandler(userStore, reportsStore, *authMiddleware)
	reportsHandler := reportsHandler.NewReportsHandler(reportsStore, inferenceService, userStore, logger)

	router := routes.EntryRoutes(routes.APIConfig{
		UserHandler:        userHandler,
		ReportsHandler:     reportsHandler,
		AuthMiddleware:     authMiddleware.Middleware,
		MetadataMiddleware: middleware.MetadataMiddleware,
	})


	logger.Info("✅ starting up HTTP server", zap.String("port", cfg.Port))

	fullAddr := ":" + cfg.Port
	logger.Info("🌍 Binding to", zap.String("addr", fullAddr))
	err = http.ListenAndServe(fullAddr, router)
	if err != nil {
		logger.Fatal("❌ Error starting HTTP server", zap.Error(err))
	}
}
