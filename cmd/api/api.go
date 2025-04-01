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
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mockTranscriber struct{}

func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return `All right. How are you? That's good. Okay. Okay. Okay. All right. How many sections do you have? How many sections do you have? Is that just a one or multiple sections? One more section, okay. How did the first one go? Yeah. Uh-huh. Yeah? There is, yeah, mm-hmm. Yeah, yeah, it makes it worse. Yeah, mm-hmm. Okay. Yeah. Oh. Okay. So we're talking about another month or plus. to be complete. Yeah, yeah. So if you're still going to have a section with her, review, talk to your family, review and write up the report. So it's going to take a while. And when that is done and you get, I think they will maybe after the report, they're going to print it out and give it to you. So if you could send it to the office, drop it off or fax it in, whichever way you want us to get it, that would be great. So I will look into it and see, okay, this is what is the, you know, the diagnosis and you hold on to it and, you know, we go from there, okay? All right. Okay. No, no, we, I want you to keep it because Yeah, we're going to make a copy. You keep the original one. Anytime you're sharing it, giving it to maybe your primary care doctor, to whoever, you want to keep a copy of it. It's very important that you have a copy. Yep. Now tell them to give you the hard copy. Because the thing with having the hard copy is that if there's no fire or water damage and you know where you kept it, you can always go and get it. But one thing with all these electrical things, you know, I give you a password or whatever. If you lose it, that's it. Or something wiped it out and that's it. It's completely gone. Yeah. Oh, forget it. And sometimes they, where, you know, the clinic, they close. There's no way you can get in contact with anybody. They're no more in business. That's it. You know. So I think that is the most thing because some of the patients that I'm getting that are new patients, they're coming to me because where they're going, they close. And we cannot be able to get that information from them because nobody's there anymore. And that's it. Yeah, exactly. No, no, no, they have a lot to worry about. That is the least of their worries. So no, they're not going to do that. Okay. All right. So you mentioned that you are on three weeks leave from your job and that is due to the anxiety you said? Mm-hmm. Mm-hmm. Okay. Okay. Mm-hmm. Okay. Yeah. Yeah. And also to reduce the anxiety because it's very stressful to, you know, engage with people like, you know, you are a customer service provider. Are you dealing with a lot of, some of them angry, some of them frustrated, some of them, you know, very nasty people, unfortunately. And when you are doing something where they feel like, oh my God, he doesn't or she know what he's doing, get me somebody or something like that. So if your job will accommodate you and not maybe in a different department where you don't have to deal with customers, you know, maybe review what other people are doing, you know, the different department there. Is that, is that like an option for you there? Yeah. Yeah. Mm-hmm. Yeah. Yeah, job description. Okay. You don't get paid. Yeah, yeah. Okay. Mm-hmm. No, yeah. I want to say because of your contractor, then your job responsibility is going to stick to it. Okay. Okay. Okay. Yeah, yeah, yeah. Yeah. It is. It is. It is. But well, just do what you can, but your well-being is also very important. If, you know, the customer service is also a trigger that increases the anxiety, then it's something that you should be looking towards into channeling your job description maybe to a different, choosing a different job. You know, because yeah, we want to make the money, but at the end of the day, we are so miserable. Then somehow the money becomes useless to us, you know? Mm-hmm. I know. Mm-hmm. Why exactly? Yeah. No, no, no, no. You know, yeah, as a kid, your brain does not proceed to think about other things. Your brain only wants what I want right away, right? So the immediate notes, that's what your brain, you know, able to compromise at a time or function with. But now as an adult, okay, I have this, I have that. If I want to buy this, I have money to do that, you know, but now where is the happiness? Why do I still feel like I'm even better off without it, you know? So Yeah. So even, it's not, it's not, it's not worth it. It's not worth it. Because at the end of the day, you want to not just feel happy that people are looking at your face, you're smiling, but feel happy deep inside that you are happy. You know, you're waking up, you're looking forward to the day. I'm happy to be here. I'm happy to engage in whatever that I'm doing. Right. So even if you're taking like a maybe 20% cut from what you're making right now and you have the peace of mind, you have the joy that you wake up in the morning, you're looking forward to your day, I think it's worth it. No. Yeah. No. Mhm. Yeah. Yeah. Wow. Wow. You see this, this Army thing, the whole thing is so complicated. The way they deal with I don't even want to, I don't even want to think about it because you can tell me about it from now to next till the heaven comes down. I'm not even going to be able to understand it. You know, the rules and regulations there, you have, I mean, it's a lot. It's a lot to understand. And sometimes you even look into it, it's like, does this even make sense that this rule have to be here? This has to be there, you know? Yeah. Yeah, exactly. They never take it away. They just leave it there. And they think, yeah, yeah, yeah. And it may be in a whole five years, 10 years, right? Then something will happen. They will go dust it up and say, well, this is the rule, you know, but it's not needed because of out of 100, the rules is used 1% of the time does not mean that it should be there. Wow. Wow. Really? Isn't that? Somebody like me, they will fire me a long time ago because I'm not even going to remember to put it on or to take it out. Yeah, I thought that what they would have done is make it like, if you want it, we make it available to you, right? If you want to wear it when you are out, you wear it. If you don't feel like wearing it, you don't have to, but here it is that we have it for you if you need it. It's by choice, optional. Wow. Uh-huh. Uh-huh. Oh yeah. Once you become a soldier, you will never live that the character, the behavior, the personality stays with you forever. You know? Oh yeah. You do it without even knowing that you're doing it. Oh. Yeah. Wow. Yeah, yeah. Okay. All right, so the medications you are doing well on it. You said that your symptoms are better than the last visit. Mm-hmm. Okay. Mm-hmm. Hmm. Okay. Okay. You still go to the gym? On and off. Okay. Yeah. Okay. Okay. At home. Okay. That's good. All right. So I see that you are out of most of the medications. I see you forgot to schedule. Okay. Yeah. Oh. Okay. No problem. And Klonopin, you're still taking it, right? Okay. Okay. Okay. Well, you have it twice a day if it's needed, right? So if you... Okay. Let's see. Okay, you can hear me? Okay. All right. So you're taking it once a day, but you know, We keep it as twice a day if you need it if you don't need it that is fine But I want you to have it in case if you need it so you don't have to Struggle with that. Okay. Do you need the lamotigine today? Okay, so you need a refill you need a refill on it Okay So the clonopin, you have enough that will last you till in the next 30 days. Okay. Okay. All right, so what I'm going to be sending today, you say you don't need to clone a pin? Okay. Okay. All right, so I will send the Vivian. We're doing the Vivian. We're doing the Vivian 15 milligram once a day. They are Bellified, 5 milligrams once a day. The Zoloft, we are doing a 50 milligram. You take two once a day. The Lamotigine, 200 milligram, you take two tablets, making it 400 a day, right? All right. And everything is still going to go to the same pharmacy? All right, then. What? Oh, you didn't take it today? Oh, okay. All right, I will send all. And now you're going to call them. Don't forget, okay? So call them, schedule them again a month from today. If anything changes, call us and let us know, okay? All right, stay well. All right, bye-bye.`, nil
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
	reportsColl := db.Collection(cfg.MongoReportCollection)

	userStore := user.NewUserStore(userColl)
	reportsStore := reports.NewReportsStore(reportsColl)

	inferenceStore := inferencestore.NewInferenceStore(
		cfg.OpenAIChatURL,
		cfg.OpenAIAPIKey,
	)

	// transcriber := azure.NewAzureTranscriber(
	// 	cfg.OpenAISpeechURL,
	// 	cfg.OpenAIAPIKey,
	// )
	// fmt.Println(cfg.OpenAIChatURL,"checking")


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
