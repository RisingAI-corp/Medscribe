package main

import (
	"Medscribe/api/handlers/reportsHandler"
	userhandler "Medscribe/api/handlers/userHandler"
	"Medscribe/api/middleware"
	"Medscribe/api/routes"
	"Medscribe/config"
	inferenceService "Medscribe/inference/service"
	inferencestore "Medscribe/inference/store"
	"Medscribe/reports"
	"Medscribe/user"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mockTranscriber struct{}

func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return `How are things going with you since we last spoke?\nPretty good.\nAny changes in medical?\nNo.\n...`, nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Critical error loading config: %v", err))
	}
	var logger *zap.Logger
	if cfg.Env == "development" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(fmt.Errorf("failed to initialize Zap logger: %w", err))
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(fmt.Errorf("failed to initialize Zap logger: %w", err))			
		}
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Error syncing logger: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
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

	db := client.Database(cfg.MongoDBName)
	userColl := db.Collection(cfg.MongoUserCollection)
	reportsColl := db.Collection(cfg.MongoReportCollection)

	userStore := user.NewUserStore(userColl)
	reportsStore := reports.NewReportsStore(reportsColl)

	inferenceStore := inferencestore.NewInferenceStore(
		cfg.OpenAIChatURL,
		cfg.OpenAIAPIKey,
	)

	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		&mockTranscriber{},
		inferenceStore,
		userStore,
	)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, logger, cfg.Env)

	userHandler := userhandler.NewUserHandler(userStore, reportsStore, logger, *authMiddleware)
	reportsHandler := reportsHandler.NewReportsHandler(reportsStore, inferenceService, userStore, logger)



	router := routes.EntryRoutes(routes.APIConfig{
		UserHandler:      userHandler,
		ReportsHandler:   reportsHandler,
		AuthMiddleware:   authMiddleware.Middleware,
		LoggerMiddleware: middleware.LoggingMiddleware,
		Logger:           logger,
	})

	port := cfg.Port
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
