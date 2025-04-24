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
I just needed some things to talk about that's on my chest.
Okay, yeah.
Tell me what's on your chest.
Couple things.
Every time I try to do a goal, set off for a goal.
For a goal?
A goal.
Okay.
I don't complete it.
The babies are 1st.
That's the truth.
Okay.
All right.
Let me get a little history from you.
Have you been diagnosed with any condition from any providers?
Yes.
I got pain through my whole body.
I got nerve damage.
So medical nerve damage?
Yes, I got medical nerve damage to my whole body.
I have spasms.
Sometimes I crawl.
You crawl?
Yeah, I crawl sometimes.
Did you get a... was it... were you involved in an accident?
I see that you have a cane.
What happened?
I threw up.
Okay.
Okay.
When was this, when did you have 2 strokes?
About 5 years ago.
Alright, have you been hospitalized for any psychiatric condition?
Yeah.
Where were you hospitalized?
I've been to the one in Burkman.
Yes, I was hospitalized.
Oh, the one in Brockton or Brookline?
Brookline.
And why were you hospitalized at the Brookline?
Do you know why you ended up there?
How long ago was this?
3 years ago?
Okay.
For medications, I see that the doctor that you see here is giving you Remeron to take for sleep.
30 milligram.
We have nightmares.
What is the nightmares?
Did you have any trauma?
Can you tell me a little bit about, no, you don't have to go into detail.
Tell me what you want me to know.
I don't know how to say it.
High-strung.
Nervous because my nerves start shaking all the time.
I have high anxiety.
When I go to sleep and wake up, I wake up with high anxiety.
I know this one.
Okay, let me see what else.
You look very nice today.
This is your 1st time seeing me.
I just wanted to say you look very nice.
Thank you.
Thank you.
Any substance use history?
Do you go by Daniel or Lawrence?
Lawrence.
Any substance use history including alcohol and marijuana?
I smoked marijuana for my birthday.
For your birthday?
It helps?
Yeah, it does.
How much are you smoking?
Because I wasn't eating.
I lost over, at one time I lost over 40-something pounds.
I was very sickly.
I've been getting a little better since I moved into this apartment.
There's no apartment now.
Okay, so apart from marijuana, do you use any other drug?
No, nothing.
Okay.
Any surgical history, like surgeries?
Yeah.
Why did you have the surgery?
I had spinal neck surgery.
I went under what I called an R3.
All right, so you're the doctor you've seen, Doctor Aziz, Fareez Aziz?
My doctor is Olivia.
Olivia?
The doctor who's giving you the Remeron?
Oh yeah.
Why is this doctor here?
The DMC.
DMC, okay.
Where do you live?
Randolph?
You're coming all the way from Lynn?
Why?
There's no... So how did you get here?
You drive?
Do you live alone?
Are you married?
Any children?
And you have support system like family that you can talk to, except your wife.
Do you do any therapy?
I used to.
Not anymore?
Why did you stop?
I stopped because every time I go with 2 or 3 times, they take me off the schedule for 6 months.
It don't make sense.
Did you find it helpful?
Was it helping you in any way?
Do you want to restart?
Are you sure?
So you have your PCP here, right?
You're not working, I assume.
I know that, but sometimes you just have to ask some people.
I gotta make some money.
I got a lot to take in.
Okay, so just confirming the medications, Lawrence, the Remeron you take at bedtime, do you still take the trazodone also?
And sertraline, you take 50 for your anxiety.
Does it help?
Nope.
It makes me have higher anxiety.
Sertraline?
And there's a lot of anxiety going on here.
I've been really bad lately.
When I started that medicine, that's when I started waking up in the morning like I had a stroke.
So are you sleeping well with the Remeron and the trazodone?
I still sleep about 2 days out.
Because of the stroke, it limits us to the medications that we can give you.
Some medications will raise your blood pressure, but it's not going to be healthy for you.
So it's... Because nothing helps, nothing works.
And the doctor told me I did really good for someone who had 2 strokes.
Because sometimes your arm gets like this.
Like my friend, he had 2 strokes.
Me, I had one stroke.
And his arm is like this.
Yeah.
Yeah.
Try the number 50.
Okay.
One to 200 milligrams.
Okay.
Try the number.
That was over like 5 years ago.
What?
This drum.
Oh, yeah.
You know what, it's a good name.
I'm in a lot better shape now.
Yeah, it did not affect your speech.
I can hear you clearly.
I'm in a lot better shape than what I was.
Are you, do you experience any depression?
I'm always depressed.
I stop looking at depression.
I guess it's because I want more.
Well, it must be depression because my family used to be depressed.
The medications you're on, I don't like it because they're going to make you fall.
You're not supposed to be combining the Remeron and the trazodone together.
It will lead to a lot of falls, and I don't want that.
What time do you take the Remeron?
You take the 30 milligrams at night?
The young actor. But I need something to help me sleep.
I can't sleep.
So when I, when I, when I engage with marijuana, it helps, it helps me eat and sleep.
And it also brings my anxiety down.
Because I have high anxiety.
Where do you have the pain?
You said you have chronic pain.
I never had back pain before, but now I do.
The pain is like this, you know, it goes from my buttocks all the way up to my neck, and then from my neck all the way down to my buttocks.
It's like awful.
Okay, so the Remeron, right, it's a medication that if you are on the lower dose, it helps you sleep.
If you go up higher, you can't sleep.
No, so the dose that you're on is high, 30.
So what I'm gonna do today is reduce it to 15.
Yep, the lower the dose,
The sleepier you get with that medication.
So the higher dose is preventing you from sleeping.
If you're on a higher dose, it activates.
So you are awake, you are alert, you get sleep.
So that's not, you're gonna lower that day and the settling.
And not being able to sleep, it feels like my whole body's heavy.
Like sometimes my body's heavy because I can't sleep.
I asked my therapist to leave a message for my doctor, and she did.
My doctor wrote me a letter from my landlord, and they told me to come pick you up, but they didn't tell me where to pick you up from.
So, you know, I don't even know who to ask.
Do I go to Family Medicine, too, and ask for the lemon?
So, which doctor wants this?
Is this your primary care?
Primary care, no.
Yes, you have to go there and ask them for it.
Which pharmacy do you get your medications?
Yerke.
Yerke?
Why didn't you get your remedy here?
Prefer Yerke?
Yerke, the lemon isn't anything else, though.
Plus I live in Lynn.
Okay, so it's closer there to coming down here, yeah.
Plus about an hour drive.
Okay.
And I think my primary care doctor wanted me to draw blood.
Hmm, okay.
So I might have to go draw blood.
I think I have come an appointment with her.
I'm not sure on the date.
Let me see.
Can I check your blood pressure real quick?
I need just one.
Just one.
Just one.
Same here.
No, you don't have to get up.
I just need your sleep, that's all.
Oh, okay.
It's easier for you to get it up.
That's one of 3 things we just worked on.
I bundled up.
I didn't know how.
You know what?
I'm gonna be safe.
Safe, yeah.
Sorry.
Because I can't- If somebody walks by and they see- Even if they sneeze, I catch it.
You catch it, okay.
So, I'm just trying to protect me and my body.
Yeah, that's what you should do.
Oh boy, I'm putting you through hell, huh?
Are these too long?
Okay.
Okay.
You're scaring me.
It's okay, it's okay, I'm kidding.
I'm playing with you.
Try not to talk.
Can you move here, please?
I was said, no, you have to focus.
You told me to read this, but that was fine.
I didn't need it.
I'll say that the person is the president.
Oh, really?
Okay.
I didn't know that.
Thank you so much.
Do you take blood pressure medications?
You should drink water.
I've been drinking water all morning.
You have to drink your blood pressure.
If your blood pressure is low, you're going to feel dizzy and you're going to fall.
You have to drink.
Okay, so we'll do that.
So what we're going to do today, Lawrence, is decrease the Remeron and the sertraline.
You said it doesn't help with your anxiety anyways.
So what we're going to do is taper down, you know how we taper down medication?
So right now you're taking one, is it one tablet you're taking?
Take half, half for one week.
And then after one week you can stop it.
Yes.
So take half for one week.
Starting, no, what time do you take it?
Essentially.
You're supposed to take it in the morning.
If you take it at nighttime, you can't sleep.
The Zoloft, the medication for anxiety, no, that's why you're not sleeping.
Take it in the morning.
So take half in the morning for one week and then after one week you can stop it.
Starting from tomorrow.
Starting from tomorrow, yes.
So you're gonna take the half for 7 days.
After 7 days, one week, you stop it.
And then when you stop it, today I'm gonna start you on some medication called Cymbalta.
The reason why I think Cymbalta is good for you is it's gonna help you with anxiety, it's gonna help you with depression, and is safe for somebody who has heart stroke, and also it would help you a little bit for your pain.
But I'm not treating your pain, but it helps.
That's Gabapentin in my mouth.
Gabapentin in your mouth?
Really?
I take 2 in the morning.
Oh, I didn't even notice that you take Gabapentin.
Okay, let me see.
I take 2 in the morning.
What is the dose?
400.
I think in Greece it's an egg, but I'm not sure.
But it was at 400. It doesn't help to take 2 in the morning, 2 in the afternoon, and 2 at night.
When I ran the bottle, it said, Take lightly.
Now, he's all off to take it in the morning, the settling.
Now, let me see what he says on the thing, because you shouldn't take it.
Oh, the gabapentin, you take 400 milligrams 3 times a day.
You take it 3 times?
The 400 milligrams?
So 400 milligrams, you take 2 in the morning.
You're taking more than you should be taking.
No, if you do that and you run out, you won't get any because your insurance won't pay for it.
It says here, 400 milligrams, take 3 times a day.
So you should take one in the morning, one in the afternoon, one at night.
That's all you should take.
So if you're taking 2, if it's not helping, you have to talk to your doctor.
What you're doing is you're going to run out of the medication and it's a controlled substance.
When you do for refill, they won't give it to you.
Don't do that.
Just do take one in the morning.
Yeah.
3 times a day.
That's what it says.
Yeah.
So that's why you have to have the discussion with your doctor.
400 milligram.
If I see you next time and you're still taking the gabapentin 2 in the morning, 2 in the afternoon, 2 at night, I won't continue you on the Cymbalta because that would be too much.
Remember, you've had stroke.
So too many medications that have interaction in your brain would lead you to have stroke again.
So always stick with whatever they're telling you.
If you want to make any adjustment or changes, check with your doctor to make sure it's safe because some of the medications, when they cross,
They interact and they make you sick.
So always check with your doctor, check with me if you wanna decrease a medication, if you wanna increase it, if you wanna make any changes, check with your doctor.
The gabapentin starting from tonight, take one in the morning, one in the afternoon, one at night.
Don't go too, that's too much.
So I'll give you the symbol, I'll start you on 20 milligram.
20 milligram, take once a day in the morning.
And the other medication, the sertraline.
Take half for one week and then you stop.
Can you tell me exactly what I said?
I want to make sure you understand.
You told me to take one every day.
Now I said it's true.
Take one in the morning, take one in the afternoon, take one at night.
Then you said the insomnia pill.
The what?
For insomnia.
Insomnia, yeah.
The insomnia one I am going to give you.
Are you going to go to the pharmacy today?
No, if you send it to Jackie, they send it to me.
Oh, okay.
So do you have enough at home, the insomnia one, the 30 mg?
I have two.
You have two much?
Yeah.
They gave you a lot?
Yeah, they gave me a lot.
I had like three bottles of them.
I don't need that much.
Okay, so the Remeron.
Remeron one is the one that starts with M, mirtazapine.
The metazepines take half.
I don't know how they look.
Maybe, I don't know.
I haven't seen it, I didn't take it, so I don't know.
So the one that you take at nighttime, take instead of the full one, take half of it.
Do you have a pill cutter?
Okay, so I'll send you the 15 milligrams then, so you don't have to cut it.
Okay, and the anxiety one, what did I tell you about it?
Take half for 7 days and you stop.
And then when you stop, I'm giving you.
You want me to write it down for you?
No, you write it so you understand whatever you write.
You sure?
Yeah.
So one, mirtazapine.
Take 15 mg at bedtime.
Is that more than a calorie?
Yeah, take that one, the morning one.
Cut it in half.
Take half for 7 days and then you stop.
Okay.
Take citrulline, 25 mg every morning for 7 days.
And then stop.
And start the new medication, duloxetine.
Where does it take that?
You can start taking it.
Anytime you get it, you can start taking it.
What's that gonna help?
Anxiety, depression, and pain.
The pain is not a lot, but it helps a little bit.
Start duloxetine, 20 milligrams.
Dissolve the sertraline after 7 days.
Just stop it.
I don't want you to take that past 7 days in addition to the new medication.
Okay, so after 7 days of going down to 25, stop it and continue.
So let me make a follow-up.
I'll check in with you to make sure that you did everything I told you.
Do you want to come here or you want to telephone?
Yeah, we can be doing telephone any coming, telephone any coming.
I will check in with you on Tuesday, May 13th.
When do you want me to call you?
Afternoon or morning?
Are you a late sleeper?
I can call you at 9.30.
You know what can help me with high blood pressure?
That's the new medication I'm giving you.
Yeah.
My cat can help me too.
Your what?
Oh, you have a cat?
Oh, nice.
That's nice.
I heard they can be very feisty.
Yeah.
They can be fun too.
They're fun.
Don't they hate?
I've been seeing videos of them hating.
She does it to me all the time, like, I don't know what that means, but I guess it means they're hungry.
That's the story.
Oh.
Okay, so I'm sending the new medication.
I want to make sure you don't take sertraline past 7 days in addition to the new medication.
Okay.
If your wife is not sure, have her call me.
Okay, have her call and I will explain to her, okay?
So I sent the new medication, yeah, also let me write something, the trazodone, I want you to stop the trazodone.
The trazodone and the Remeron, there's no need for that.
I want to make sure you want a very small medication. Because of interaction, stop the trazodone.
Do you take the trazodone every night?
5, 3, 4.
Okay, so if you take it every night, what I need you to do is, how many do you take every night?
Just one, okay.
The Remeron, the lower dose of the temazepam, will help you.
So you were all set.
I made that appointment for you in 3 weeks.
I will call just to make sure you did everything.
And at that point, if you stopped settling, I will increase the new medication for you.
But I started you very low.
And when I see you, I'll go up on it for you.
Okay.
`
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
	// azureTranscriber := azure.NewAzureTranscriber(cfg.OpenAISpeechURL, cfg.OpenAIDiarizationSpeechURL, cfg.OpenAIAPIKey)
	// geminiTranscriber := geminiTranscriber.NewGeminiTranscriberStore(geminiClient)
	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		&mockTranscriber{},
		// geminiTranscriber,
		// azureTranscriber,
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
