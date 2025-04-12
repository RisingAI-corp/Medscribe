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
const sample1 = `I had the honor of getting the Saturdays.
I was like, are you sure?
Okay.
But yeah, everything on this end seems to be going.
Seems to be going.
I told you about disability, correct?
Yeah.
How's that going?
So I had had that.
So just this last Tuesday I had, let me see if I have my volume on so I can hear you.
Are you there?
Yeah, I'm here.
All right, I hear you better now.
So I had my call with Disability on Tuesday, and just so happened I happened to be talking to somebody from Lowell, where I'm at, like I'm in Billerica, but still, I got to tell them kinda something.
So it wasn't like me explaining it to someone from Texas, do you get what I'm saying?
They actually, he could understand me, so when I got on the file on my online account,
Out of 5 steps, I was already step 3 because of my pre-existing, even though they denied my appeal, I'm reapplying, I'm already through to 5 steps, so I would think they'd be reaching out for you soon.
Okay.
That's good.
But I'm very, very happy.
I'm very, that takes a huge weight off my shoulders, but I have a huge other thing that I've been dealing with.
Okay, what is that?
I found that my testosterone
It's only been, so a normal male, say, my age, would be, like, somewhere between 600 and 1100 on the higher end, maybe.
I think it's, like, it might be 350, but it's 1100 on the high end.
That means, like, these cops you see walking around huge.
Or doctors, or doctors, doesn't mean that they're naturally on TRT, they could just be normally a high level 1100 looking swole because they work out.
Now me, I should be somewhere in 350 or 600 to 1100, I know 1100 is the high end, but mine have only been 14.
Wow.
Now, Ms.
Hilda, last week I didn't even talk to you.
On Saturday, I got in a little bit of a fight with my girl because we went out and I was looking.
Excuse my French, because I speak freely, and even though I'm a white person, I don't believe in racism.
But I dress like a, you know what I mean?
That's how I am.
I dress like a young kid.
I knew it going on with my girl, but I didn't really realize when I'm going into the store how much I look different than her.
I usually try to collab kinda on my on how much I do it with her.
You gotta get what I'm saying?
Yeah.
So this goes back to the story.
Don't worry.
So coming home arguing about something.
Well, what we were arguing about in the toe to toe with like the power struggle in the house, not listening to me.
At that point, my testosterone must have dropped so bad, because I didn't get these levels until the other day, that I was in bed from Saturday — I slept from one to the next morning, got up at 7 like I usually do, was in bed by 3 Sunday, slept all the way to Monday, and didn't get up until halfway through Monday when my probation officer popped up on me.
So what I'm saying is that my testosterone, so that would give the reason that I'm thinking, I beat around the bush a little, sorry, my testosterone.
That's why I think that I have these significant downs.
Okay.
Well, are they treating it now?
Are they giving you anything for it?
Yeah.
They've given me a shit ton of blood work so far, excuse my French, but they don't know what to do with me because
In 2013, when I went to a different doctor, they seen something in the CAT scan or an MRI on my pituitary.
But didn't you know how you are so so follow me if you can you know how a CAT scan cuts like a cut a deck a deck of cards?
Am I talking correctly?
Yes.
Okay, so it's cutting the deck, let's just say your brain, you're looking at the top, the left side of the deck, the right side of the deck, when they were doing it, no, the right side of the deck, and then chopping that side in half the opposite direction, when they saw it one way, the abnormally or abnormalities, they didn't see it in the other direction, do you copy what I'm saying?
Yes.
Okay, so now they just put me on testosterone back in 2014 for about 3 years.
I'm getting shot, shot, shot.
No blood work is shown.
They never tested my blood or nothing.
Maybe once or twice.
I think I've had 3 tests by them, but I was on test with them for 2 years.
Why?
I don't know.
I had to call my ex-girlfriend because she was in the office when I heard I had something wrong with my pituitary.
Now this is way back when I was using.
Yeah.
So I reached out and she said something, so the doctors don't know what to do with me right now and on Wednesday they're going to have a huge meeting about me.
Well, they have to figure it out because somehow it has to be managed, it has to be taken care of.
Can I ask you something just as a person that I've talked to for a while and I don't have people?
Now if the endocrinologist, back in 14, we have like a 10 to 11 depending on the state for BAC.
I'm not looking for that, but I'm saying if the doctor in 2013 put me on testosterone for something that they saw in a CAT scan, they should have followed it up with another CAT scan, correct?
Yeah, they should have maybe reached out to you.
So my doctor didn't see the 2nd one, she only saw one when I was in her office the other day for one.
But for two, at such a young age, even if it was dip, because they were blaming all my problems on my substance abuse when I was really asshole, but my dopediction was maybe a tenth of everything a day.
You know what I mean?
It wasn't big.
It was the crimes I was committing.
So now I'm asking, I'm getting off track, but
Could they shut my casties off back then, you think?
I've looked everywhere online, I've tried to reach Miss Hilda, I'm trying on my own and I have no one to ask.
I've read so much stuff I don't even know, my girlfriend wants to say I don't know what to do, I don't know, I don't want to ask the wrong person.
So your question is what, should they have continued to give you the shots or what?
They were giving me the shots, they haven't, not yet.
They haven't made a move on me yet, back then.
I was getting the shots every other week from them.
Okay, and why did you stop the shot?
I, well, if I put you, if you'd had the minute, do you have a minute?
I'll just pull up the graph and tell you the date.
I have it all on my phone.
I forget, it was like July of 13, 14.
Okay, so 2014.
I remember, I remember.
2014 was the last time you got the shot?
2013 or 14 was the last time?
No, it was, like, 16, the last time I got my shot.
It was, like, 9 years ago.
Okay, why weren't you getting it?
Why weren't you getting it since 9 years ago?
I don't know.
They just dropped me, and I don't know why, but also at the same time...
Miss Hilda is, I've had the same doctor, no matter where I was, and now I've been home since 2009 from the States, so I've been in Massachusetts, I haven't been nowhere else.
So when I came home, Doctor Lee was my 1st doctor.
I've went so many different places depending on the doctor and if I vibe with him, like I vibe with you.
And I've always had him listed as my primary, to the point where just last year, I was talking to him about my
Excuse my French my limp dick issue.
He's okay with just throwing me on sedelafil just last year with an open script as of right now.
I can't even eat those because they don't work.
I'm not face to face so I can't read your face.
You understand what I'm saying?
I understand what you're saying.
I think the most important thing right now is
Yes, I think that is, don't even divide your mind or your mind should be focusing on this or that.
I think the most important thing is, is focusing on them, making sure that they fix what is going on right now.
Then other things can follow up later.
You know, I think that is the most, keep to your appointments.
Yes.
Yeah, make sure you keep to your appointment if they asked you to go to see a specialist, to go to your lab, whatever, get it done.
Can I say something?
Like, from the bottom of my heart, I know you have a lot of kids, we've been talking for a little while now, like, it really makes me feel...
at ease that I get to hear from you because nobody, I don't talk to my mom and Silda, I don't talk to nobody, the only one that talks is the girls at the bar at Elm Tree in the dark, and you, and you know what I mean?
So it means a lot to me, like if I could talk to you right now, because I'm thinking the same way, because I'm like, look, I need it taken care of, I can't get out of bed, like,
I can't get out of bed.
And so she's like, well, we don't know what to do with you.
Keep taking your omeprazole for your swallowing problem.
That could all be led to, she's like, and they're trying to say that my swallowing problem is healthy because I just had a barium swallow.
Yeah.
What was the result of that?
They haven't even told me the results from that.
The marmoterygram I had.
Okay.
Where they went up my nose and down my throat with a tumor that cancels the barium swallow.
But she's not even, they haven't even touched on that, so I don't know.
Okay, well, you're doing what you're supposed to do, don't you?
This is a new doctor, like, this is a new doctor, like, you're my new doctor.
This doctor came in the same time I switched everything at the same time, because of where I was going.
They always looked at me like an addict, like I couldn't get out of that shell, you understand what I mean?
Well, I'm glad that- I hate the boy, I hate- No, no, no, it's not.
Boy, you're sad.
I know you're banging people out on Saturday.
I hate the boy you're sad as is.
No, it's okay.
I'm glad that at least that is being taken care of, that is being discussed.
I just- So that- And like, in all in all, my wife's telling- My girl's telling me, like, I haven't been able to hold nothing down for a job.
Like, could it be real ditched to where, like- Cause, listen, a big thing I read, and I don't know if it's true, but if you don't have your ADD in check and your test goes real low, you could be suicidal.
All these years, I've never had my ADD in check until just recently.
It's just, like, it's mind-boggling to me.
It's having me read more.
I don't know if this was something natural that was supposed to be brought on to me.
Well, one thing about the Google thing is that you shouldn't be on it a lot because everybody's experience is different, okay?
People have different body types.
Our body is not the same.
The way you react to some medications is not the same way somebody has reacted to it, and what you might be reading is somebody that had a negative side effect or a negative experience, they put it on Google.
So be very careful because if you're thinking about it like that, you know, it might worry you a lot.
I prefer going to my doctor, discussing with my doctor.
and getting the feedback from them than getting it off the Internet, because everyone is in there, everyone is talking about their negative or their positive experience, most, 80 percent of the time, the only thing you're going to hear is all the bad things about it.
Because people that had a good experience, they don't come to the Internet to talk about it most of the time, but what you see all the time,
is people that have a negative experience, they will come to the Internet to talk more about it.
And most of the time, people will just take that and it becomes their problem.
So just be careful the way you do that.
I can totally agree with what you're saying on other forums when I'm reading other interviews, but I get what you're saying.
I, you know, I don't, like I said, I was telling my girl's son, he's 13, he's here today, and I'm like, dude, I'm not hopping on you, you're like my little dude, you're like my homeboy, I don't even look at you like my girl, you know what I mean?
Like, I don't look at it, because you always think I'm hopping on him, but I'm like, dude, I'm not hopping on you, people hop on me.
I was like, that's all, I just have no one to talk to, so it's kind of odd, I just try to read as much as I can during the days.
Yeah, all right.
Yeah, so that's about it, Miss.
Oh, thank you for that.
I appreciate the talk.
I guess I look at you like my family, in a way.
That is fine.
So, other than that, everything going well, right?
Everything is going superbly well.
Okay, very good.
And what do you need for medication today?
I need my Adderall, and let me see right now, because I have them right here.
I don't know if I'll make it 2 weeks on my Lexapro.
I'll make it 2 weeks on Lexapro.
I think I have the TENS, if anything.
Okay.
So I'm all right there.
All right.
So the only thing we're doing today is the Adderall 15 mg 3 times a day, right?
Yes, ma'am.
And you want to go to where?
Where do you want me to send it to?
North Belrica CVS, please.
Okay.
Alright, I will send it there.
Yeah.
Okay, what did you want me to schedule you?
Let me know what I have.
What did you want?
Okay, let me see.
Depends on what I have.
Okay, so, um...
No, I don't have anything on the 25th, because if it's Friday, it's supposed to be on the 25th.
I don't have anything on the 25th, but I don't mind doing on the 23rd, which is Wednesday.
I'm available on Wednesday.
That's fine, if you like to, yeah, that's fine.
No, no, no, sorry, sorry.
I also know my Saturdays.
Saturdays is just, it's just, it's different for me to hear from a doctor, and I get talking like this, you know what I mean?
Yeah, so do you want me to do Saturday or do you want me to do Wednesday?
What do you want?
You can do Wednesday.
Okay, let's see.
All right, I'll do Wednesday, okay?
2 weeks on Wednesday, thank you.
Yeah, all right, stay well.
Take care.
All right, bye-bye.`


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
