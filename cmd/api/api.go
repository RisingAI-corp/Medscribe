package main

import (
	"Medscribe/api/handlers/reportsHandler"
	userhandler "Medscribe/api/handlers/userHandler"
	"Medscribe/api/middleware"
	"Medscribe/api/routes"
	"Medscribe/config"
	inferenceService "Medscribe/inference/service"
	inferencestorre "Medscribe/inference/store"
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
)

type mockTranscriber struct{}

const sample1 = `Good.
Hi, Miss Joanne.
How are you?
I'm okay.
Thank you.
That's good.
Happy new month.
Happy new month to you, too.
Okay.
No problem.
I've been taking it easy on myself.
That's good.
Awesome.
Yeah.
Okay.
Okay.
I've been doing walking, making walk, you know, that too.
I've been doing that around the building, you know, so my fears could go away.
I've been doing that.
Okay.
That's good.
I love that.
We got little Chloe here with me.
Oh, really?
Yeah, little dog.
She's so cute.
How old is she again?
She's gonna be two.
Wow.
Okay, that's good.
That's good, that's awesome.
Yeah, she don't shed, she don't, she's good, she's good.
Okay, all right, that's awesome.
So no panic attacks, no depressive episodes?
Not this month, no.
Okay, all right, awesome.
How about your mood, your mood okay?
You know, my moods always go up and down.
That's something normal, I think.
Us women, we go through that, you know?
Of course, that's we women, always going through that.
Yeah, like we want to get the satisfaction that at one point we were getting and we're not getting it today.
We gotta give it to ourselves.
Yes.
Yes.
Okay.
All right.
So how about your medication?
Everything?
Everything's working.
Everything is working.
Okay.
All right.
I will send you all of your medication to your pharmacy, the gabapentin, lorazepam, you know, hydralazine, mirtazapine and duloxetine.
Okay.
And your Suboxone.
Okay.
Thank you very much.
All right.
Yes, I will see you next month on, um, let me tell you, um, that would be one, 2, 3 July, May 3rd at 9 a. m.
Yeah, we're good.
Alright, see you, okay?
Nice talking with you, honey.
Alright.
Have a nice weekend.
And you too, bye.
Okay.`

func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return sample1, nil
}

type mockInferStore struct{}

func (m *mockInferStore) Query(ctx context.Context, request string, tokens int) (inferencestorre.InferenceResponse, error) {
	return inferencestorre.InferenceResponse{Content: "mock response"}, nil
}

func main() {
	log.Println("🚀 App is starting...")
	log.Println("⚡ ENV BEFORE .env load: PORT =", os.Getenv("PORT"))
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Critical error loading config: %v", err)
	}

	var logger *zap.Logger
	if cfg.Env == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("❌ Failed to initialize zap logger: %v", err)
	}
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

	userStore := user.NewUserStore(userColl)
	reportsStore := reports.NewReportsStore(reportsColl)

	inferenceStore := inferencestorre.NewInferenceStore(
		cfg.OpenAIChatURL,
		cfg.OpenAIAPIKey,
	)

	reportsTokenUsage := reportsTokenUsage.NewTokenUsageStore(db.Collection(cfg.MongoReportTokenUsageCollection))

	transcriber := azure.NewAzureTranscriber(
		cfg.OpenAISpeechURL,
		cfg.OpenAIAPIKey,
	)

	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		// &mockTranscriber{},
		transcriber,
		inferenceStore,
		userStore,
		reportsTokenUsage,
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

	// Ensure we're reading the PORT env var properly
	port := cfg.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
		logger.Warn("PORT not set; defaulting to 8080")
	}
	logger.Info("✅ Ready to start HTTP server", zap.String("port", port))

	fullAddr := ":" + port
	log.Printf("🌍 Binding to %s", fullAddr)
	err = http.ListenAndServe(fullAddr, router)
	if err != nil {
		logger.Fatal("❌ Error starting HTTP server", zap.Error(err))
	}
}
