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

	// For google.Credentials
	// For WithCredentialsJSON
	"google.golang.org/genai"
)

type mockTranscriber struct{}
const sample1 = `Hello.
How you doing?
Okay.
Excuse me.
Now I can hear you have to be logged in and out to be able to get to be able to get this.
How are you?
Uh, not great.
Okay.
What's going on?
Um, my dad passed away.
Oh, oh my word.
So sorry to hear that.
When did this happen?
Thursday.
In the hospital?
Yeah.
So since the last time we spoke he has been in the hospital or he came home he went back?
He came home January 2nd and then he went back.
Came on January 2nd?
He came back January 2nd and then he went back in the hospital March 7th.
Okay.
So since March 7th he's been in the hospital?
Yeah.
So what changed?
It was just, uh, he was just progressing and getting worse.
Oh, wow.
And then, um, uh,
The Friday previous, I guess my mom and dad decided to change to hospice care.
So he died at home?
No.
The VA hospital decided it was best to not even move him to the Brockton VA.
They were going to do the hospice care at West Roxbury.
They thought it would be too much to even move him.
So they did the comfort care at the West Roxbury on the med floor at West Roxbury.
Wow, wow, wow.
I am so sorry.
Oh God.
So sorry, so sorry.
I'm so sorry.
Even though, you know, you know that it's not like getting, um,
Better, because once they start talking hospice, it means that, you know, it's the end of life.
Once that word hospice is mentioned, it's like, hmm, but somehow you still have that, even though that little tiny bit of hope that maybe somehow, you know, closely there will be some different outcome, even though with hospice, because, you know,
I've worked in a hospice, though I didn't last long there because I couldn't deal with it.
I've worked like in a nursing home, visiting nurses that they give people, like, 6 months to go, and come one year, they're still here.
And then they take them off the hospice, and sometimes they go for another 2 years, even, before 3 years before they go.
Even with the hospice, sometimes it's not like that, you know, so because of that, you still have that little hope that maybe there will be some, a little bit better outcome that will give that person a little longer to stay.
Then sometimes the disease just rampages the body that
You rather even have them go and have peace and some rest instead of looking at them and there's nothing you can do, you're just kind of keeping them there because you don't want them to go, but at the same time they are going through hell, you know?
So it's like, which one is the lesser evil here to deal with, you know?
So it's a very hard thing to let go, extremely hard.
But just know that it's something you have no control over.
Absolutely you cannot fix it, even though, no matter how much you try, that is something that once it happens, there's no going back.
There's no going back.
And I wish this
All the scientists, all the discovering that we're doing, I wish there's a way that somehow you can stop it.
But there's absolutely no way to do that.
Absolutely no way.
There was a friend of mine that lost her daughter.
That is when I really knew that, well, I know it, but it really kind of hammered it home very well to me that when somebody close that eyes forever, there's absolutely nothing that we can do to bring that person back.
Because this, my friend was like, if she could go to the end of the world to bring this child back, she would have, if T. S. could wake up a dead person.
She doesn't need 2 tiers to do, 2 people's tiers to do, that she already has is enough, but regardless of everything, there's no going back.
So you look at it in a way that there's nothing you could do.
So remember the good memories.
That is what you're left with, the good memories.
Try your best to bring everybody together because
When something like this happens in a family, sometimes it can bring families together, sometimes it can tear families apart.
So, that is one thing that you guys should be very, very careful of, because different people will want to do things differently.
You know, oh, that wouldn't want that, yes, that want this, what are you talking about?
Then you are sad already, you are depressed already, then a little thing,
can just mess things up, you know?
So just be very, very mindful of that.
I've seen where it happens a lot.
It happens more when you have only one parent left and that parent left, that parent dies, then there's no person to imitate or say, okay, this is how I want to do it, right?
So it happens more there where it's just left for the children to decide what to do with the remaining parents.
So it becomes very difficult.
But the good thing is that your mother is here, so most of the thing is going to be almost based on her decision, right?
And whatever that is going to happen, she has to okay it.
But at the same time, you guys should be very careful.
Well, thankfully, he decided everything as far as...
His whole services and everything, songs, the psalms, the hymns, everything.
He already paid for everything already, so everything's completely decided.
Like he knew of the debt and he's dying and he prepared everything how he wanted to go.
He did this a couple years ago.
Oh, wow.
Yeah, I think that is good.
I have a friend that did that.
She bought her symmetry like 10, 15 years ago.
I'm like, at the time she told me that I was like, are you telling me something that I don't know?
Are you, did the doctors give you bad news that you're going to die maybe in the next six months?
She was like, Oh no, no, I've seen where ETS family apart.
I don't want that.
And she said, because I have a very large family.
I don't want everybody to start putting in the opinions and then everybody's going to, it's going to destroy the family.
So she 15 years ago, and she's still here going strong.
And I'm like, you have yours already bought the plot everything.
Yeah.
Yeah.
And she planned everything.
What she does is like every year, she kind of, maybe she has a different idea.
She kind of updates this update that yeah.
Yeah.
My dad, uh, well, my dad was in a family plot since,
long before I even met my mother.
Wow!
Yeah, since 1963, they've had a family plot.
Family plot, okay.
Okay.
That's good.
That's good.
Yeah, so... So he already knew where he was going.
Hmm.
Hmm, okay.
So the only thing left is for you guys to carry it out.
Yeah, and my mother asked my husband and I to be pallbearers.
And I think my other brothers as well.
I don't know if my sister is gonna be or not, I'm not sure.
So when is the funeral?
The week's Friday, the funeral's on Saturday.
And Saturday was supposed to be his 93rd birthday.
90 what? 93.
And he's gonna get buried on his birthday?
Wow.
Wow.
Hmm.
I think when if, you know, look back, he would say, 93 years old is pretty a good year.
You know, somehow I pray to make it to 90.
You know.
It's something that, like I said, you hold onto the memories.
I think that is what is gonna get you through it.
Hold onto the memories and let the memories bring so much joy and comfort to you that when you remember things with him, the way he laughed, the way he cracks up jokes, the way, you know, he would say something and when you ask him again, he might be like, oh, no, I didn't say that.
When did we have that discomfort?
So let those happy times be the things that carry you through.
Okay.
That's what I think was like, so difficult about that, that last week in hospice, because he had a sense of humor the entire time, like the things that would come out of his mouth.
Oh my God.
And then that last, that last week, like, like he couldn't get anything out of his mouth, you know, like it was just like facial expressions and it was just so painful.
Yeah. To not, like, be able to communicate with them while you're, like, here is sense of humor, because that was, like, the one thing that was comforting that whole time, was instilling the sense of humor.
That is the most difficult part of it, you know?
You're looking at this person that has all these qualities, and they're just laying there, helpless, and somehow I feel like the things they want to say, they still, it's, like, trapped inside.
And it cannot come up, and that's why sometimes they use the facial expression, because I know what I'm saying, but I can't voice it out for you to hear it.
So now they kind of... Like he was trying to like mouth things, and it's like, I could tell he was like, like saying bye to me, like, and trying to mouth it.
But nothing was coming out and I just wanted to hear it like one more time.
Yeah, yeah.
And you always want to hear it.
Even if he said that the last time, you still want to hear it.
There's never the last time.
With loved ones, there's never the last time.
You know, there's always reason you want to hear it again the last time.
But there will never be the last time.
That's the thing.
See, so...
Hold on to the good memories.
Hold on, hold on to it.
It's gonna be hard.
And the healing hasn't started yet.
So until the funeral is done, you know, because right now he's not buried, everything is still going very raw.
So after the funeral and as these weeks goes by, it will start kind of,
kind of registering somewhere in your brain that, Oh my God, he's not here because sometimes you want to, when you call, you want him to pick up.
You want him to, you know, speak in his usual way, the way he speaks to you when you call him.
But it is, um, it is, it is how life is.
So I'm so sorry that you're going through that.
There's some things within that will remain with you forever because I remember my, um,
My dad, even up to today, I still remember how my father calls me.
I mean, there's some certain things that, it's been like 16 years now because my daughter, he died in February, my daughter was born in June.
There's a way he calls my name that sometimes I hear, he's so clear, it's like his voice.
I mean, it's like no one else calls me like that but him.
I don't know if it's,
People call me the same name, right, but there's a way he calls it that is different when other people calls me that name.
So sometimes when I hear, I can't, I don't know if I'm imagining it or I'm hearing it, but it sounds so clear to me sometimes.
I'm like, okay.
You know?
So you just hold onto the good memory.
There's so many unique things about him that will bring you comfort, sense of joy, when you actually take a moment and think about it.
Okay.
I'm sorry to hear that.
All right.
How was your trip to Georgia?
It was, it was good.
Like I, like I said to you, like I almost didn't want to go because of like everything with him, but it was, it was good to like, be able to like see friends and laugh, you know, even though I had this like anxiety in the back of my head, you know, with everything going on with him.
Yeah.
Okay.
The weather was nice except for one day it rained, but, but other than that, it was, it was a good time.
And then, you know, I come back and, and it was like, I, you know, it was like, it had been a year since I'd been to Georgia cause things were just so, you know, chaotic, like being up at the hospital, like and everything, you know?
Okay, good.
And, um, how is your sleep?
Oh God.
Awful.
I just can't breathe, you know?
I know you have issues with sleep and something that will make it worse as I would say has happened.
So it makes it more, um, more difficult.
Yeah.
And I've been just having these like crazy, crazy dreams, like,
That just, like, wake me up and then I can't fall- then I can't fall back to sleep and then I do fall back to sleep but it's only, like, an hour at a time and... Yeah, the last- the last couple weeks have been really difficult.
Difficult, okay.
And, uh, how is your appetite?
Eating okay?
I don't know, maybe once a day?
Gotta be very careful, though.
I know.
Gotta be careful you don't, um... The last thing you need now is to get sick.
You don't want that at all.
I know.
Like, I'm hungry right now, which is a good thing, so... After we finish this, I'll go eat something.
Okay.
Alright, good.
And, um...
What about any medical doctors visits, any study on any new medications?
No, nothing new right now.
All I did was my mammogram.
I don't even remember if I looked to see what the results were.
I think that, you know, I don't even, like, so many things are just like an absolute blur to me.
in the last, like, couple of weeks?
Well, for the mammogram, I really believe that if there is some concern or some problems, they would give you a call immediately.
Yeah.
So, sometimes when I don't hear from them, I don't usually sweat it.
I don't usually feel... I feel... Sometimes I feel okay.
And if you have my chat, you can log in there and...
See what is going on, you know, I do that.
Yeah, they send the results and yeah, right.
Like no news is good news.
Yes.
So I think that is the best thing.
If you feel like, okay, I've not heard from them.
I don't know what is going on.
You can just log in there and see what the result of whatever the test they've scheduled for you or you've done and go from there.
But I also will last thing on my mind and I just was like, you know, I did have like a colonoscopy scheduled, but I rescheduled it because I didn't know what was going on.
And I had this, the, the growth that came back after the basal cell, um, I had, they were going to do another biopsy next week.
And I was like, I was like, so I rescheduled that this is before he had passed away.
You know what I mean?
So I was like,
You know, I don't know what's going on.
So I rescheduled that.
And then, you know, I just, I, like, I just started like rescheduling a bunch of things, everything except for that herpian acupuncture, I rescheduled.
So I was like, those are the only 2 things I can handle right now, you know?
Yeah.
Sometimes when you have, um, you have a lot going on.
The only the best thing to do is to pick out what is the most important thing that I'm going to attend to.
So you don't get so overwhelmed, you know, and you go with that, whatever you feel like, this is something that cannot wait.
This is the immediate thing that I need to handle.
Then you pick that up with and go, I think that is the best way to do it.
Because if you try to.
Combine everything, do everything at the same time, it's gonna be too overwhelming for you, it's gonna be way too much.
Yeah, and I've been extremely overwhelmed, I've had multiple daily anxiety attacks, couple panic attacks, obviously crying every day.
When I did get the call when it happened,
I probably screamed like a 2-year-old, like, from my gut, like, on the kitchen floor, like, probably, probably alerted every dog in the neighborhood that came out, like, it was just so bad.
Yeah, it's not, it's never a good news to, to get to, it's never a good news, you know, there's never anything good
about somebody getting the news that their loved one passed away.
I'm sorry to hear that.
It's, I don't know, it's like I said to my husband, it's like, I'm almost like tired of people like calling and asking if I'm okay.
And it's like, what kind of question is that?
Of course I'm not okay.
Like, you know, like, and I just want to like hide and isolate and just be left alone.
Yeah.
Well, sometimes people calling is for them also to have the closure, to say, okay, so this I'm hearing is real.
So I'm confirming from the daughter, from the son, from the loved one that this is happening, you know, this happened.
So sometimes it's a way of them to, to put like, this is it for them, for them to also know that this, and I'm sure if they didn't call that, if they didn't call,
They will not feel good about themselves about it.
And sometimes- Yeah, and I know that it's out of the kindness of their heart or anything like that, but my heart is just shattered right now and I can't handle much.
I know, yeah.
My siblings all wanna go over to my parents' house today and I don't even wanna get out of my pajamas right now.
I just don't, I just wanna...
Kind of go through some boxes and find some old photos and just stay home.
Like, I just, I don't know, I just can't do it today. Well, do it when you feel comfortable.
It might be too much for you to handle.
You don't need to be there when everybody's there.
And what you need to do is just go at your own pace.
And if you have, you know, the time, you call your mother.
I know she's too, she's hurting also.
It's her husband, somebody that she's home.
I'm sure most of you guys have moved out and it's just the 2 of them.
So it's going to hit her so much.
You know, at the end of the day, everybody has, you have your husband, everybody has their ones, their wife, their children that will comfort them.
And the only person she's been with is your, your father.
So I think, um, you know, keep close touch, you know, call her mom, how are you doing this day?
What do you need?
Um, I can come and spend time and something like that.
Not just today, not just tomorrow, it's going to be something that is going to be ongoing.
Because loneliness can really, really make her decline.
Yeah, I know.
Okay, so that is another thing that you guys have to be extremely mindful of, how you guys take care of her in terms of not to let her be on her own for so, like, weeks and all that.
The broken heart syndrome?
Oh yeah, oh yeah.
You know, sometimes when you hear it, it's like, what does that mean?
But it's a lot, I've seen it a lot.
Yeah, I mean, she's been with him 3 quarters of her life, you know, so.
That's the only thing she knows.
And the things, it's like when I was in school, I've been in school for so long, and when I finish school, it's like I don't know what to do with myself.
I feel like something is missing.
I feel like I should be doing something.
I feel like, Oh my God, why, why am I not doing something I supposed to do?
So now you have, she has to start learning to live without him.
There are some certain things that is like a routine for her.
She wakes up in the morning, she knows how to take care of him.
She knows how to do this, to do that.
And now that thing that defined her for so long is no longer there.
So there's this vacuum there.
That, of course, she has to find something to fill it.
If she's not filling it or doing something that will take her attention from that, she's going to be very, very depressed.
And she has PTSD as well.
Okay, so that would make it even worse for her.
Yeah.
Okay, so this is where you guys come in and help her with that, okay?
All right.
Thankfully, she has 5 kids, so I think today I'll sit at home with my PTSD.
I don't think it would be helpful for her.
Yeah.
Okay.
Good.
And let's see, so how other things, apart from going to have the biopsy again and manage with sleep, but at the same time, you just have to also take care of yourself.
You need to be
Um, strong enough health, healthy, and take care of yourself to be able to help your mother in this difficult time.
The last thing she needed is for you to, to break down, you know?
So remember that your wellbeing is extremely important and you've done so, you know, you've made some progress.
So, um, I wouldn't want this.
the death of your father to drag you back to square one.
Definitely he wouldn't want that for you, okay?
I know.
Okay, so that's why I just want you to be mindful of that and also please continue meeting with your therapist.
It's extremely important that you continue to meet with her, okay?
I know.
Every
Every Monday I meet with her.
Okay.
All right.
And, um, so you needed to medication today?
Yeah.
Okay.
Did you get, did you get my text messages?
Oh yes.
I forgot about that.
I got the text message about, um, sending the information to them.
You know that every time the back and forth, back and forth with these people.
Um, to send this or to send that, I know that, uh, when you forwarded me those information, I called my husband to follow up.
And, um, when I got your message again, he said, um, he did.
So, um, they said, what are they saying?
I don't know.
I haven't checked back in with them since, uh, all this stuff has, um, gone on, but they said they just needed the, the information.
From you from, uh, basically all because the last thing they got was, um, December 2024.
So basically they needed everything from this year, I guess, all the office notes from this year.
And then, um, they also needed, uh, the, uh, information from my, uh, from my therapist that I see weekly.
So those, those are the 2 things, um, that they needed.
They did, they did the same game all over again.
They cut me off without notifying me.
So when, uh, the.
looking for this to end because I wouldn't, uh, this thing is becoming so cumbersome and so annoying.
And that's what I said to him, I said, if I can work on my mental health, you know, 50% of my mental health is dealing with these people.
You know what I mean?
You know, there's always, um, at my, I was just talking to my husband when I got your message, I'm like, are you sending the payment?
Because
Every time they're requesting for that, each of the pages, there's how much we're supposed to include for them to pay, because that is how it is.
Even I remember when the 1st time we started it, they actually sent us like a page or a form to fill out and see how much we usually do for transfer of record.
And my husband was like, he has never even filled any of the form and done any of those things.
So it's like, I don't know what they want more from me.
I don't know, they just want every.
Every month I do a visit, I should send them the notes, so it's becoming... Well, that's what I'm saying to them, like, you should have told me, like, you know, I'm like, I just had this fixed in January, and then you cut it off in February, I was like, why didn't you tell me, like, that you wanted this information, like, or you want it every month or every 3 months, you just go and cut it off and you don't even say anything, like, when I just went through this in January, like, what's wrong with you people?
I'm like, I'm tied to these games, like, I'm trying to work on my mental health and you're making it worse.
Like, enough.
It's awful.
You know, and then to deal with that on top of, like, right before I went to Georgia and then everything with my dad, it was just like, it just, it's been absolutely brutal, this last, like,
Month and a half has been absolutely brutal on me.
The combination of everything.
I don't know what else, you know, it's like.
I don't know if they're thinking that maybe they're thinking that, you know, the situation is going to change and you get bigger and get something on the note that will say, OK, now patient is doing this patient.
Uh, no longer need medication, patient, I don't know, maybe that's what they are looking for, but uh, unfortunately, whatever you reported to me is whatever is gonna be on the note there, so.
And all they do is just keep making me go backwards with the games.
Well, what they want is, um, I think, uh... They wanna break, they wanna break me down.
Uh, well, it's, it's...
It's not even breaking you down.
They want to frustrate and frustrate to the extent that whatever that is needed, you're not going to, you know, you're giving up that.
No, I'm done with this.
You know, this is going to, was it 23 years now that since actually since I started with you, it's been this way.
So that is that one of the last while I don't want to deal with, I have some patients that they will come here and they will tell me, Oh,
I'm meeting with you because I need the letter.
I'm like, no, I don't want to do that.
I don't do it.
I've turned so many patients away because it's something that is never ended, you know?
It's something that is never ended.
They're requested for peer review that you're not even getting paid for.
They're requesting for your time to print out all this record.
It's a lot.
It's a lot.
I know.
I don't want to deal with this.
I just want to stop feeling
Like so overwhelmed and having like all these anxiety attacks, like over small things, big things, you know, obviously recently, but like even small things or want not isolate.
And I just, well, just, just, you know, that, um, it's not going to keep going this way.
One day this thing is going to come to an end and.
It's usually not going to come to an end that you're going to be happy about.
So it's something that you just need to prepare yourself.
One of my, a very good friend of mine with Boston Medical, even though right there he got, she got hurt at the job and
Still, finally, after 4 years or so, they closed the case, they went to court, she had a good lawyer, they went to court and everything, and nothing was given, nothing, because the reason why they would have the best lawyer, paid the best lawyer, what they would give to you, they would pay to a lawyer for you not to win the case.
So this is not something they are going to keep doing.
It's something that one day is going to come to an end.
So the goal here is for you to be prepared that they're not going to keep it going like this for a long time.
I just want to not feel like my mental health-wise, though.
My goal was to get my mental health better and check, to be able to go back to work and not feel so overwhelmed.
Dealing with live electricity and having panic and anxiety attacks is not a good recipe. That is the thing, the thing is that they will tell you that you have what we call a chronic depression or a chronic anxiety, right?
And this is what it is.
So meaning that you're not going to get better is something that you're going to manage for the rest of your life.
Just like somebody will have diabetes or somebody will have high blood pressure, it doesn't go away.
Once you have high blood pressure, it's going to be there for the rest of your life.
The only thing is that you manage it with medication and live your life, right?
So they're going to tell you that, well, you're not going to get better.
You've been on this for the past 2 years, for the past 3 years, for the past 4 years.
So you're not going to get better.
You have what is called chronic anxiety or chronic depression.
That is something that is going to be with you for the rest of your life.
So what is it going to be?
You cannot go back to work because of this anxiety that will never be cured, but you have to manage it with
medication, you keep going to therapy.
So are you going to come back to work with this anxiety that you are managing with medication and depression?
Or are you going to stay away because it's not going to get better.
So as time goes on, based on what I have dealt with with insurance and all that is going to come to that because they are not going to keep you for the next 5 years, 6 years, 10 years for you to get, they're not going to do that.
So,
one day that is what is going to happen because I've seen it a lot.
Okay.
So it's just for you to, I don't want to live like that either.
You know what I mean?
Like, you know, I want to eventually get some sort of like better grasp as to where I am right now.
Well, that, that is the goal.
The goal is to accept you and,
Be able to accept what it is and make peace with what it is and be able to move on and live your life as normal as possible, right?
So the thing with them is that, like, you know, the information they are sending to me, they want me to put a date on it, which is that is the way they usually do it.
All the ones I've done, put a date on it.
Are you saying?
in the next 2 months, in the next 6 months, in the next one year.
When do you- Like a broken bone, like when they've healed?
Yeah, so I'm like, no, it's not something that you can say, okay, in 6 months, the patient is going to be well enough to do this or well enough to do that.
And again, like I said, it's a chronic, it's something that will always be there, you manage it with medication.
That is, sometimes you, it's like when you have a history of cancer, right, for example, you have a history of cancer,
It will never go back like, I never had cancer, right?
It's always part of your medical history when somebody says to you, what kind of medical problem you have?
Oh, I had cancer, right?
It's part of your medical record.
So just like anxiety, depression is going to be part of you for the rest of your life.
So that is where it comes in that, okay, this is a chronic thing that is going to be with you for the rest of your life.
And they will be like, well, we're not going to keep waiting for you for it to go away before you return back to work, because it's not something that will ever go away.
So, so you know that the insurance, they will never want to keep paying for not getting that service they're paying for.
And they will tell you, well, it is what it is, and they didn't cause it, they didn't do whatever.
Just like when people are involved in a car accident and the car insurance settled them and they go their way, whatever is the outcome of that accident is no longer the problem.
They've given you what they need to give you and that is the case closed.
That is how insurance work.
Their job is to save as much money for their client as possible.
Okay.
They didn't have all these lawyers, they didn't hire them to just come there.
and get paid and they complain to lose money.
No, their job is to save money as much as possible.
So I just want you to be aware of that, that one day they're going to tell you, no, they will cut off the check.
They're not going to do it anymore.
And then they will say, if you want to go to court, wherever you want to go, they are ready to meet you anyway you want to call them.
That is how they operate.
Okay.
All right.
So for medications today everything is going to go to the same pharmacy, right?
Yes.
Okay.
When, when, when did, I don't even remember when I texted you that.
Maybe I can see.
When did you send the, oh, that was back at the.
Well, I'm going to talk to my husband.
I don't know why.
Oh, it was right before I went to Georgia.
Yeah.
Okay.
So I'm going to talk to him.
Almost a month ago, because they said they have to have everything by, um, I think the end of my calendar by the end of, uh, uh, or the middle of this week.
Hmm.
Okay.
I will tell him, um, I don't know to check if he had a confirmation of that.
Because I know that every time we are sending something, it's always, they didn't get it, they didn't get it, we have to resend again and again and again and sometimes talk to them multiple times and get different people before they will finally say they got it.
Since we started this with them, I've never sent it once, it went through, there's always a lot of back and forth, back and forth for it.
Yeah.
So I will ask him again, if
confirmation.
If not, I would say send it again because it's it's always that is how it is.
It's they never go through one time.
It's always again and again and then they confirm that they got it and I remember that last time the last one we sent, there was a lot of phone calls that went between my husband and because I would if he was for me doing it, I would not be able to do it because I
I have a lot going on, I will not be calling and calling and leaving messages and calling, I couldn't do it.
Yeah, no, I get super frustrated when it comes to stuff like that.
Okay, so I will follow up with him and see if, if not, if they say they didn't receive it, I would say send it again, send it again and maybe Monday follow up with them and he will also follow up to see if they got it, okay?
Okay.
So I send the Zoloft 100 mg once a day and the Ativan 0.5 once a day.
I send all that to CVS on Washington Street.
And I will schedule you again in a month.
So the next one is going to be, let's see, I'll put it like for 12, 1230 on the 3rd of May.
1230?
Mm-hmm.
May 3rd.
Yeah.
Oh, is that Kentucky Derby Day?
I forget.
I have no idea.
All right.
OK.
All right.
Stay well, OK?
I'm trying.
Try your best.
I think that's the best you can do.
That's what you can do about what is going on right now.
OK.
All right, all right, stay well, okay?
All right, thank you.
All right, bye-bye.
Bye-bye.`


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
	// gptInferenceStore := inferencestorre.NewGPTInferenceStore(
	// 	cfg.OpenAIChatURL,
	// 	cfg.OpenAIAPIKey,
	// )

	var geminiClient *genai.Client
	var clientConfig *genai.ClientConfig
	

	clientConfig = &genai.ClientConfig{
		Project:  cfg.ProjectID,
		Location: cfg.VertexLocation,
		Backend:  genai.BackendVertexAI,
	}
	
	geminiClient, err = genai.NewClient(ctx,clientConfig)
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
		// gptInferenceStore,
		GeminiInferenceStore,
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
