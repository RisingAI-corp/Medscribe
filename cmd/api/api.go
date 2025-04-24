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
	transcriber "Medscribe/transcription"
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

const sample1 = `
I'm doing okay, doing okay.
Just checking in to see how things are going.
Has there been, do you have any concerns with your mood, sleep, appetite, anything at all, stresses?
No, but however, I haven't taken my medication.
Which one?
Which medication?
The injection, there's only one medication I'm on.
Okay, let me get into your meds, because I'm not in your meds that much.
Oh, yes, just the injection.
I thought you were due for it sometime this week, no?
Yes, I was due for it since last week.
I had an appointment to take it.
When I went to get the medication, they said I didn't have any refill left.
You could have, when it happens like that next time, you know, you do fill it at a Cordman Pharmacy, right?
Yeah, they could have they did put in a refill request for but I saw it 2 days ago 2 days ago, but then when I looked at it, I figured that you were coming.
I was seeing you today.
So I just didn't put in the refill.
I just wanted to make sure the dose is fine everything and put it through but it didn't come through to me last week.
Yeah.
Yeah.
Listen, I left a message for you.
I don't know if you received it.
This is going to be a problem.
Did you do it on my chart?
No, I told them to leave a message for me.
They said they would tell you and then you would call me.
They never told me anything.
Because some days I work from home, like today, yesterday, Mondays and Tuesdays, I work from home.
I'm only in clinic on Wednesdays and Thursdays.
And because I'm new there, having been there, I think I started back in early January, I'm still not settled.
There's a lot of things that is still not put in place well.
So I apologize for that.
I'll put it in right away for you.
Are you feeling any different because you have missed a week?
I don't anticipate that, but just making sure that you're doing okay.
I've just been feeling more tired than usual.
Okay.
I'm just going to do this real quick before you go to sleep.
So I put in a refill right away.
I know you said you're feeling tired, but besides that, is there any mood fluctuations, any concern, any hallucinations, anything like that?
No.
Okay.
All right.
How about your sleep and appetite?
I've been sleeping a lot.
My appetite is good.
Okay.
Is that your baby's crying?
Yeah.
Okay.
All right, so the refill, I'll put it in now and I'll make sure I put a lot of refills so this doesn't happen again.
But in the near future, if it happens again, if you go to the pharmacy, sometimes the pharmacy will not send me that message.
You can come upstairs and, you know, usually where you check in, you can tell them to let me know that you went there to fill your meds and it's not there or, you know, whatever is going on, just leave them a message and I would get that right away.
Okay.
How about therapy?
How is therapy going?
Um, actually, it's been going good.
I have an appointment, I think, next week.
Next week?
Okay.
So, I will put in the refill for you.
When can you get it?
When can you, um, can you schedule to go get it anytime this week?
Yes.
All right, so it will be there.
I would suggest that you call them right away just so you can get an early appointment.
I don't want you to get off of it for too long before restarting it.
The program will be good.
Okay.
So yeah, just continue therapy.
If you notice any difference in your medication or if the tiredness continues after you get back on the meds, that's something to look out for and let me know.
Maybe we need to make some revision.
With your meds, but other than that I hope you just get back on it continue therapy and I think because of this med issue I would want to check in with you.
Let me actually make the next appointment.
Do you want to come in in person next time?
So I can see you in the 1st week of June, which is gonna be like 6 weeks from now.
So maybe a week when you get your... Were you saying something?
No, it's my daughter talking.
All right, so let's do June 4th.
Let me see, June 4th, is that not a holiday?
Yeah, this is July.
Okay, so June 4th, I do have,
930 open, I have 1130 open, I have noon open, I have 1230 open, actually I have, what time?
930, okay.
All right, so that is all set.
If you happen to have any issue with the meds, I'm sending it right away as soon as I hang up with you.
Just go upstairs to the 3rd floor, let them know, and tell them to message me and I'll get in touch with the pharmacist to figure out what's going on.
Okay.
All right.
Okay.
So you take care and I will see you in what, around 5, 6 weeks?
Yes.
I'll see you.
Between the time period, if anything changes, don't wait until that time period.
You can always switch that appointment to a closer date if you're having any issues.
Okay.
No problem.
All right.
You take care then.
Bye-bye.
Thank you.
It's not going to work.
So you're taking the Abilify injection every 28 days.
How is that helping?
That helps a lot.
I feel pretty good.
It's just that I was having a hard day for some days.
That's why I feel kind of tired.
But I plan to schedule and get it sometime this week.
I'm sleeping fine.
I'm eating okay.
And things are pretty good.`
func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return sample1, nil
}

func (m *mockTranscriber) TranscribeWithDiarization(ctx context.Context, audio []byte) ([]transcriber.TranscriptTurn, error) {
	return []transcriber.TranscriptTurn{}, nil
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
		//&mockTranscriber{},
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
