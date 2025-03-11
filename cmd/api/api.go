package main

import (
	"Medscribe/api/handlers/reportsHandler"
	userhandler "Medscribe/api/handlers/userHandler"
	"Medscribe/api/middleware"
	"Medscribe/api/routes"
	inferenceService "Medscribe/inference/service"
	inferencestore "Medscribe/inference/store"
	"Medscribe/reports"
	"Medscribe/user"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mockTranscriber struct{

}
func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return "I have a itchy throat and i feel really sick", nil
}

func main() {
	logger, err := zap.NewDevelopment() // Or zap.NewDevelopment() for development
	if err != nil {
		fmt.Println("Failed to initialize Zap logger:", err)
		return
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Error syncing logger: %v\n", err)
		}
	}()
	if err := godotenv.Load(".env"); err != nil {
		logger.Error("Error loading .env file", zap.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		logger.Error("MONGODB_URI environment variable is not set")
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Error("Failed to connect to MongoDB", zap.Error(err))
		return
	}

	if err = client.Ping(ctx, nil); err != nil {
		logger.Error("Failed to ping MongoDB", zap.Error(err))

		return
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(fmt.Sprintf("Critical error disconnecting client: %v", err))
		}
	}()

	db := client.Database(os.Getenv("MONGODB_DB"))
	userColl := db.Collection(os.Getenv("MONGODB_USER_COLLECTION_TEST"))
	reportsColl := db.Collection(os.Getenv("MONGODB_REPORT_COLLECTION_TEST"))

	userStore := user.NewUserStore(userColl)
	reportsStore := reports.NewReportsStore(reportsColl)

	// transcriber := azure.NewAzureTranscriber(
	// 	os.Getenv("OPENAI_API_SPEECH_URL"),
	// 	os.Getenv("OPENAI_API_KEY"),
	// )

	inferenceStore := inferencestore.NewInferenceStore(
		os.Getenv("OPENAI_API_CHAT_URL"),
		os.Getenv("OPENAI_API_KEY"),
	)

	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		&mockTranscriber{},
		inferenceStore,
		userStore,
	)

	userHandler := userhandler.NewUserHandler(userStore, reportsStore, logger)
	reportsHandler := reportsHandler.NewReportsHandler(reportsStore, inferenceService, userStore, logger)

	router := routes.EntryRoutes(routes.APIConfig{
		UserHandler:      userHandler,
		ReportsHandler:   reportsHandler,
		AuthMiddleware:   middleware.Middleware,
		LoggerMiddleware: middleware.LoggingMiddleware,
		Logger:           logger,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("PORT environment variable not set. Using default port", zap.String("port", port))
	}
	port = ":" + port

	logger.Info("Server listening on port", zap.String("port", port))
	err = http.ListenAndServe(port, router)
	if err != nil {
		logger.Error("Error starting server", zap.Error(err))
		return
	}
}
