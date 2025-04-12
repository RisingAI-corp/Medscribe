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
const sample1 = `Hi, Jen.
Hi.
Oh, hold on one second.
Mm-hmm.
Here she is.
OK.
Good.
How are you?
I can't hear you.
What can I tell you?
How you doing today?
I can't hear you.
You can't?
Oh, now I can, now I can.
Oh, okay.
So I do meditation with a mask at night and it has like the meditation music in it, so it was connecting to my, I'm sorry about that.
It's okay.
It's okay.
How are you today?
Good.
Good.
All right.
How is your medication?
How is the Vyvanse, the 20 mg?
Um, it's okay.
I'm still finding, like, halfway through the day I'm, like, meh, but, um, but it's okay.
I'm trying to take it as close to, like, noon as I can.
Okay.
All right.
So, so with that being said, I know that I increased it last week for that 20, so you want us to sit on that 20 for, like, a month to see how that goes or what?
What's your plan?
Um, if we could do, like,
If we can do it in one, um, I just, because I'm feeling it right now, like the feeling of not being, um, working, like towards the end of the day, like, I'm a little nervous of that.
Okay.
My daughter's getting married and stuff.
I know.
So if you're kind of busy, like in the afternoons and you know, so I would prefer like.
2 weeks, like maybe 2 weeks.
I'm going to try 2 weeks.
Yeah.
Okay.
All right, so, because I still, you know, the reason why I'm asking you that I didn't want to force something to you, I just want you to make the decision yourself, you know?
Yeah, I was going to say a week, or if you're going to do a month, I was going to say a week, even like an hour, like I used to do, like in the afternoon.
No, I'm not going to do a month.
I throw that out to you just to let you know that you can still go to a month.
So if you want us to continue the weekly, it's fine.
If you want us to do 2 weeks, you know, to try how you do 2 weeks, it's also fine.
So it's kind of, you have that, you know, window.
Can we try one week just so I can make sure that I can, because if it's not working, I feel like to go up a little bit more.
the afternoon at all.
Either one is fine with me, because I know both did work at one point.
It works for me.
I just want you to find that.
Yeah, this is it for me.
So I don't want to push it and we go over the board.
All right.
So that's awesome.
So I will still continue the one week.
At least we do next week and see how that goes, how your medication goes.
If you still don't find it helpful, we can increase it.
You know, because it's a capsule, so we can't do that half, so it has to be, I think, 20, then the next one would be, I think, 30 or 40, one of those, so there is no, like, 25 or something like that, okay?
And I know, because, you know, you have this wedding, big wedding planning, could be another thing that making not to fill the effect of the doses that you have, because, you know, come on, we are busy, we are everywhere.
The concentration is just not there, like I went over, I was fine for like a good portion of the day when I went over last weekend to help her and then like I started to like kind of just go, I'm getting tired, I'm getting tired, like just not.
So at what time do you take the Vyvanse?
I'm trying to take it closer to noon now because I'm noticing it's wearing off.
In the afternoon, I'm trying to get up by 10 to noon.
Okay.
All right.
So, even though I might not be like having the concentration, I might have like the little bit of like the caffeine feeling type of thing that you get from me.
Okay.
All right.
That's awesome.
Okay.
So, um, I will tell you, you know, I know you don't like too much medication.
Me too.
I don't like too much medication.
You know, Adderall is more stronger, lasts longer than the Vyvanse, but you have been on Vyvanse for years, so I don't want you to jump from one medication to another, but I just want you to know that you have a room just in case, you know, if the Vyvanse is not working, we can throw in like a 5 mg of Adderall or 10, you know, to compensate it if at the long run the Vyvanse is not working, so we can do that.
But for now, we stick with the Vyvanse, you know, if it doesn't work, then we can
Do the other route too, okay?
All right.
So for your son, I spoke, I forgot to call you during the week.
I spoke, yeah, let me give you the doctor's number.
So her name is Doctor Jenny.
Jenny, her phone number is 860-834-7537.
Just tell her that Doctor Chichi gave you her information.
She was asking me about your insurance, your son's insurance and all that.
I told her that I'd take your insurance, but I wasn't on my computer, so I couldn't tell her the name.
But you call her, then you 2 talk, and I think she'll have to run your insurance to make sure that she accepts it, then she'll take him, okay?
All right, so I will reschedule you, I'm so sorry, do you have any other questions or concerns?
No, not right now, I don't.
Okay, so I will reschedule you the 19th, April 19th at 10?
Yeah, 10, perfect.
Okay, all right, I will send you your prescription right now to your pharmacy, okay?
All right.
Thank you very much.
All right, thank you, bye, bye Jen.
I was just reading it off the top of my head.`


func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return sample1, nil
}

type mockInferStore struct{}

func (m *mockInferStore) Query(ctx context.Context, request string, tokens int) (inferencestorre.InferenceResponse, error) {
	return inferencestorre.InferenceResponse{Content: "mock response"}, nil
}

func main() {
	log.Println("🚀 Starting up application")
	log.Println("⚡ ENV BEFORE .env load: PORT =", os.Getenv("PORT"))
	cfg, err := config.LoadConfig()
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


	
	if(cfg.Env == "production") {
		// in production we will use the metadata server to to leverage the cloud run service account to auth with vertex ai
		err := os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		if err != nil {
			logger.Warn("Error unsetting GOOGLE_APPLICATION_CREDENTIALS", zap.Error(err))
		} else {
			logger.Info("GOOGLE_APPLICATION_CREDENTIALS unset for production")
		}
	}

	geminiClient, err := genai.NewClient(ctx,&genai.ClientConfig{
		Project:  cfg.ProjectID,
		Location: cfg.VertexLocation,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		logger.Fatal("❌ Failed to create gemini client", zap.Error(err))
	}
	
	GeminiInferenceStore,err := inferencestorre.NewGeminiInferenceStore(
		geminiClient,
	)
	if err != nil {
		logger.Fatal("❌ Failed to create inference store", zap.Error(err))
	}

	reportsTokenUsage := reportsTokenUsage.NewTokenUsageStore(db.Collection(cfg.MongoReportTokenUsageCollection))

	//creating services
	transcriber := azure.NewAzureTranscriber(
		cfg.OpenAISpeechURL,
		cfg.OpenAIAPIKey,
	)
	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		// &mockTranscriber{},
		transcriber,
		GeminiInferenceStore,
		// &mockInferStore{},
		userStore,
		reportsTokenUsage,
	)

	// instantiating api
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, logger, cfg.Env)
	userHandler := userhandler.NewUserHandler(userStore, reportsStore, *authMiddleware)
	reportsHandler := reportsHandler.NewReportsHandler(reportsStore, inferenceService, userStore, logger)

	router := routes.EntryRoutes(routes.APIConfig{
		UserHandler:      userHandler,
		ReportsHandler:   reportsHandler,
		AuthMiddleware:   authMiddleware.Middleware,
		MetadataMiddleware: middleware.MetadataMiddleware,
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
