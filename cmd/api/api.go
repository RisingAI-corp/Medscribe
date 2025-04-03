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

const sample1 = `I'm good.
All right.
And, uh, how things going since, um, the last time we spoke, since last month?
Um, pretty okay.
I did the blood work for the... Okay.
When did you get it done?
I think it was 2 weeks ago or less.
Okay.
Let's check that.
Any medical changes I should be aware of?
No.
Okay.
There's just a steroid in my ears that I'm supposed to take that my doctor said it could, like, make my mood a little bit weird.
Hmm.
What are you taking this steroid for?
It's for my ears.
Okay, like you have an ear infection, ear pain, what?
It was like clogged because I had the flu.
Okay.
I don't know if it pulls up in my chart.
How long are you going to be on it?
I think it's 6 days.
Okay.
Do you think that would affect me a lot?
No, but you need it though.
So you take it, that's why I asked you for how many days?
You take it, if it's going to affect you, it's not going to be that much, and it's going to, your body's going to get rid of it, and you're going to be back to your normal base, okay?
How many days have you, how many days has it been since you started taking it, or how many days is left for you to complete it?
I didn't start it yet.
Oh, you've not studied it.
Okay.
Yeah.
How many times a day?
Once a day?
I think it's like once a day, then by each day it goes up one time, like the 6th day I take 6 or something.
Yeah.
Something like that.
Okay.
Okay, and how is school going?
Yeah, so, um, I managed to do a few things for school, but I just have a really hard time concentrating and memorizing, not memorizing, but like remembering what I'm studying.
If I read something, I have to read it like thousands of times to actually get it.
Okay.
So I was looking for another ADHD assessment, but because the one I did, it was like one hour long and it was just one day, which I feel like isn't the usual.
So I was looking to get something more in-depth to see if I have ADHD and if I could benefit from medication?
Well, right now I am prescribing you, based on the symptoms that you are reporting, I'm prescribing you 60 mg of Strattera for the ADHD.
And, you know, I have quite a lot of patients on it, and they're doing absolutely well on it.
But the thing that I'm worried about is that, the thing I'm more worried about is that, is that something that is actually lack of motivation?
or something that you don't have interest to do, because sometimes I've had some patients that believe that the reason why they're not doing well in school is because they need a stimulant to help them.
And successfully, they got it from whoever prescribed it for them, and they're still having the same problem.
So sometimes it's not because
Medication is not going to read that book for you.
The medication is not going to do things.
The only thing that medication can do for you is that you have the willing and you have the motivation to do it.
And then you are doing it and you're reading it, but you are not sustaining those information or you are not sustaining your attention long enough.
That's where the medication can help.
Okay, so if you're looking into to get more in-depth, like you said, assessment for ADHD, you can.
If you look for whatever they can obtain it, you can.
But my biggest concern is that we're not just dealing with poor attention.
I think there's other things that is there.
Are you still seeing your therapist?
Yep.
The one in Brazil or the one here or both of them?
The one there yesterday.
Okay.
So this is what, you know, where both of you should actually get into some of the main reasons.
Okay.
If it's, um, if it's the lack of the motivation or lack of, you know, knowing
Lack of responsibility, because this has to do a lot with responsibilities also, and also what you enjoy doing.
Sometimes when you enjoy doing some certain things, you see yourself doing it, getting, you know, so interested in doing it, looking forward to doing it, you know?
So for example, I, um, see some kids that when it comes to games,
They don't lose focus when it comes to game, they can play game for hours.
And that when they're playing that game, they sustain their attention in that game.
They don't lose their attention playing game.
Why?
Because this is what they're enjoying to do.
They find so much joy in doing that.
So it's not like when they're playing that game, they're losing their attention.
No.
So, but when you move them from that game to go study, they cannot concentrate because it's not something that they're enjoying to do.
Or something they're not getting pressure from, they're not getting enjoyment from doing it.
So they're not focused on it.
Okay, so if you wanna watch a movie, some people, they see themselves, they can stay and watch a movie from the beginning to the end.
And they can watch like two or three movies in a day without losing their focus from it.
Okay, well, without changing 5 channels, they can't decide what they want.
They are doing this, they jump to this one.
But they can stay and watch a home movie and really enjoy it.
So now the question becomes, why are you not losing focus on watching a movie that you can stay for 2 hours movie without flipping the channel?
But you cannot spend 30 to an hour or 2 hours to read a book and to be able to take an exam.
So, uh huh.
Um, I'm not like, it's not just school that I'm having a hard time concentrating.
Even like reading for fun or watching TV, I just can't sit down and do it.
I just can't pay attention.
And with school, at least I feel a lot of difficulty.
Like I'm pushing myself lately.
When I was tired, I tried to do and I did everything, but it took me so long to do something so simple because I had to reread it over and over again.
Something that was really easy and I could have done easily, but I just couldn't do it.
Okay.
And the Strattera, you don't think that is doing much for you, right?
So the Strattera, you feel like it's not doing anything?
Hi.
Hi, okay.
Wait, Wi-Fi went down.
Yeah, so the Strattera, do you think that the Strattera is doing something or you feel like it's of no use to even continue to take it?
I feel like it's better than nothing.
Okay.
Is there like a higher dose or something or is that like the recommended dose?
Well, I think the highest dose for Septra is about 70 mg.
Oh.
Yeah.
I've been taking it, like, for a while, and I don't feel a lot of improvement.
Yeah.
Like I haven't been able to do a lot of stuff that I did for fun in a couple months because I just can't sit down and do it and concentrate on it.
It's not even sitting down and doing it.
It's not like finding motivation to do it because I've been doing like everything I have to do lately.
I've been really responsible with everything I have to do.
It's just, uh, I can't make my brain focus on something.
Okay. Okay, so, um, if you can, because the last, um, the last ADHD evaluation that you did, they did not, that report, they did not, uh, diagnose you with, um, they did not diagnose you with ADHD, so maybe... Yeah, he did, like, the testing he did during the report, it said that I get, like, 5 criteria out of six.
So I'm not sure why he didn't give me a diagnosis.
I know it's not only the criteria that meet the diagnosis, but like I scored like a one percentile in memory and I feel that like that's not very normal.
Okay, so let me, I will increase the Strattera from the 60 to 70, if I will say it here.
No, I think it's 80.
Let's see.
Okay, so I don't know where you're calling, you know, to do this.
You can also call the office and get some of the list of places that you can call and see if you have a place you've already or if your primary care doctor can refer you to a different place.
I know that most of my patients that had it done, they usually do it for
Some would do it for 4 hours, some would do it 4 hours in different days altogether, like 8 hours.
Yeah, mine was like an hour long, so that's why I don't feel like it was very thorough.
Okay.
All right.
Also, I was looking to get accommodations.
Okay.
Okay.
At my school now, because I changed schools.
You changed schools, okay.
There's a form that I think you might need to fill out.
Okay.
Where should I send that?
Yeah, I think that's a form that they sent, something like that.
So I will look at it and see if I have to fill it up, anything like that.
Okay.
All right, so for medication today, what do you need?
I think it's the vilazodone.
Let me check.
I think it's the vilazodone and the abilify.
Okay, so you don't need the Depakote?
No, I think I just picked it up.
Okay.
And the Strattera.
Okay, and you don't need the Lamotrigine?
No, I think I stopped that one.
Why did they have it here?
Okay.
That is the Lamotrigine.
Yeah, I think I stopped that one.
Okay.
I just received yes to put it back.
Okay.
Let me delete it from the system here.
Okay.
And, um, the buspirone, you don't need it?
No.
Okay.
The billifier, you need it 10 mg once a day, right?
Yep.
Uh, the Depakote, you don't need it today?
No.
You'll have enough for the next one month?
I think I have like a full bottle, so I think so.
Okay.
Let me just pull up your lab and see what is in it.
Okay.
So the Depakote is 22.one.
The normal level is 50 to 100.
So that is low.
They only gave me the Depakote level.
I want the whole bloodwork.
They only gave me.
The only check for Depakote level.
Oh my God.
Okay.
So currently the Depakote, you are taking 250.
That's why it's low.
You're taking a very low dose, but at the same time, we are now just going to increase it because it's low.
We are going to increase it based on what you're reporting because
I have some patients that are on the same low dose and they're doing very well and their mood is stable, so we leave it at that.
But what it tells us is that we have room to increase it.
That is not above the normal level, right?
So if you feel like the current dosage is managing your mood well, you're stable on it, we keep it at that.
If it's not, we have room to increase it.
Okay.
So how are you doing in general?
If you have to rate your mood 0 to 10, 10 being the best, how would you rate it?
Um, like a seven.
Like a seven?
Okay.
And how many hours of sleep would you say you get in a night?
Probably around eight to nine.
Okay.
Okay.
So we keep the medication as it is now and, um,
We keep reevaluating every month and see where we need to do anything with, okay?
And I will try to fill out the form by the end of this week.
You should be getting it.
And if you get the Neurosoc evaluation done and just inform them that they need to fax the results or they can give it to you.
You can drop it off at the office.
Some people will say, No, I don't want to fax it.
Or if you insist, they will charge you money to do that.
But you can just pick it up and drop it at the office for us, okay?
Okay.
All right, any question or any other concern you have?
Someone at the Disability Center of my school, they said there's a website that you can get a psych evaluation done.
I don't take that.
Yeah, I thought so.
I don't take that, I have some.
People will go and get it done and come, take the pill like maybe 180, 170, get it done and my hypon doesn't take that.
Okay.
They want more in depth where you can sit and do some, they will do like assessment, give you some things to do at the office.
That's how they can observe you and see if you have that or not.
Okay.
Okay.
Um, for the accommodationss, I was just going to ask for more time, like on assignments.
Okay.
Because it took longer to complete them.
Okay.
So you need more time assist you on the assignments and on the exams, something like that?
Yeah.
Okay.
Just the same thing as before.
Okay.
Have I done that for you before?
Yeah, for my other university.
Oh, okay.
I will see.
Okay, I will look if I can pull it up, okay?
So I can just model whatever is on it, okay?
Okay.
All right, stay well, okay?
No problem.
All right, bye-bye.`
const sample2 = `
All right, how are you?
Not as bad as last time we talked, so that's good.
That's good, okay.
I'm also, though, I'm on a 3-week unpaid medical leave right now to try and kind of reset from all the anxiety from work that was happening.
So I'm doing that, and I did my 1st of the diagnostic testing yesterday.
Okay.
with the South Coast Diagnostics.
So... How many sessions do you have?
Sorry?
How many sessions do you have?
Is that just one or multiple sessions?
So I have one more session.
One more session, okay.
Yeah, so the next, so this one was about 3 hours and the next one will be about 2 and a half hours.
Okay.
So... How did the first one go?
Lot of
Tests, a lot of questions.
The interesting part was the cognitive tests.
That's what I started, that I'll be finishing in the next one.
They make you feel stupid eventually.
I mean, I know they tell you at the start, like, you're not meant to get them right, but it makes you realize how
pour your short-term memory, especially since I was anxious.
Yeah.
Yeah.
So until I started really hearing stuff over and over, that's when it finally clicked.
And I think it was just, my anxiety was kind of fading and I was getting more into the motion of actually trying to memorize what we're saying or recall whatever.
Um, but I have no idea.
when I'll get any answers back about anything.
Um, so I don't know where it's going to go.
Um, and there was zero information provided back to me, but I know it was mainly the doctor assessing me verbally.
And then I think she's going to take part of the week to review my tests and then I'll have the follow on and then another week or review or two, I think she's going to call my family.
and talk to them about some of my childhood stuff.
Oh.
And then she'll write a whole report, I guess, eventually.
That'll come my way.
Okay, so we're talking about another month or plus.
Probably, yeah.
To be complete.
Yeah, if I had to guess.
Yeah, yeah.
So if you're still gonna have a session with Revue, talk to your family, Revue, and write up the report.
So it's gonna take a while.
And when that is done, and you get, I think they would maybe after the report, they're going to print it out and give it to you.
So if you could send it to the office, drop it off or fax it in, whichever way you want us to get it, that would be great.
So that we'll look into it and see, okay, this is what
is the diagnosis and you hold onto it and we go from there, okay?
Yeah, that works.
I think it's better in your hands than it is in mine.
No, I want you to keep it because we're going to make a copy.
You keep the original one.
Anytime you're sharing it, giving it to maybe your primary care doctor or to whoever, you want to keep a copy of it.
Okay, very important that you have a copy.
Yeah, the last copy medical records I got, they gave me the wrong password to the CD.
Yeah.
Um, and when I tried to get the password, they lost all the records.
No, tell them to give me the hard copies.
You know, because the thing with having the hard copy is that if there is no fire or water damage, and you know where you kept it, you can always go and get it.
But one thing with all these electrical things, you know, I give you a password or whatever, if you lose it, that's it, or something, wipe it out, and that's it, it's completely gone.
Yeah, oh, forget it.
And sometimes the, you know, the clinic, they close.
There's no way you can get in contact with anybody if they're no more in business.
That's it, you know.
So I think that is the most thing, because some of the patients that I'm getting that are new patients,
They're coming to me because we're going to close.
And we cannot be able to get that information from them because nobody's there anymore.
And that's it.
Yeah.
And a lot of people don't think to ask for hard copies before they leave.
I'm sure if the place is going out of business, they're not going to be sending out hard copies.
No, no.
They have other things to worry about.
They have a lot to worry about.
That is the least of their worries.
So no, they're not going to do that.
Okay.
All right, so you mentioned that you are on 3-weeks leave from your job?
Yes.
And that is due to the anxiety you said?
Yeah, the anxiety.
I guess just the consistent drain from the interactions I was having at work.
So it's my attempt with my company at a reset.
So hopefully when I go back, the customer relations and interactions don't push me over the edge as much as they have been.
If that doesn't work, then I'll probably have to start looking for a new job or think about going back to school.
But this is basically just my attempt to hope it's just bad burnout or something that's making it so much worse.
Okay.
So in the diagnostic testing is the kind of way to, I guess, guarantee that I proceed forward properly.
Yeah.
So I know I'm handling everything at least the way I should be on the job front and personally, but also on the medical side.
Yeah.
And, and also to, um, reduce the anxiety because it's very stressful to,
you know, engage with people like, you know, you are a customer service provider, are you dealing with a lot of some, some of them angry, some of them frustrated, some of them, you know, very nasty people, some unfortunately, and when you are doing something with the feel like, Oh my God, she, he doesn't know she knows what she's doing, get me somebody or something like that.
So, um, if your job could accommodate you and not
Maybe in a different department where you don't have to deal with customers, you know, maybe review what other people are doing, you know, in a different department there.
Is that, is that like an option for you there?
So my boss has tried to make my like projects more internal focused, so they have tried to help.
Um, it's just, I'm on the service desk, so the core function of my responsibility is to
work with people who have problems on their computers.
And since I'm a contractor for the government, I have to follow my job description rule.
If it's not fit, then I don't get paid for the contract.
It's semantics, personally, but they care a lot about it.
I understand.
Because you're a contractor, then your job responsibility, you're going to stick to it.
OK.
Yeah.
So the contract's up for renegotiation in July.
But right now, the government obviously is having a bunch of issues, so I don't even know if the contract will stay.
So I might be without a job in July either way, depending on how that all plays out.
But wouldn't be the worst thing, because I think I need to try to stay away from customer-facing jobs in general.
It's just tough to find those.
Sometimes it is, it is, it is, but well, just do what you can, but your wellbeing is also very important.
If, you know, the customer service is also a trigger that increases the anxiety, then, um, is something that you should be looking towards into channeling your job description, maybe to a different, um, choosing a different job.
You know, because yeah, we want to make the money, but at the end of the day, we are so miserable, then somehow the money becomes useless to us, you know?
Yeah.
I grew up in a very, my family was very poor when I grew up.
So like I always thought money was the solution to everything growing up because as a kid, you know, I just thought we didn't have money and it caused problems.
So eventually I made enough and I'm like, well, I'm making enough.
Why am I still not?
Why is nothing great?
That's when it finally clicked, like, oh, okay, now I get what Aldo's saying.
No, no, you know, as a kid, your brain does not proceed to think about other things.
Your brain only wants what I want right away, right?
Yeah.
So the immediate needs, that's what your brain, you know, able to compromise at a time or function with.
But now as an adult, okay, I have this, I have that.
If I want to buy this, I have money to do that, you know, but now where is the happiness?
Why do I still feel like I'm even better off without it?
You know?
So, um, yeah, so even it's not, it's not, it's not worth it.
It's not worth it because at the end of the day, you want to not just feel happy that people are looking at your face, you're smiling, but feel happy deep inside that you are happy.
You know, you you're waking up, you're looking forward to the day.
I'm happy to be here, I'm happy to engage in whatever that I'm doing, right?
So even if you're taking like maybe 20% cut from what you're making right now, and you have the peace of mind, you have the joy that you wake up in the morning, you're looking forward to your day, I think it's worth it.
Yeah, that would be nice.
Yeah.
I used to have a job like that, which is very difficult to get back into.
And it was my military job, so doing it as a civilian is a little bit more challenging to get in.
Just because, obviously, the military sends you right there.
But with civilian jobs, you still have to be hired by the Navy, and then the Navy has to decide, is that position going to go there? And then whether the the command usually has to pay for the civilian personnel from the Navy So they have to be willing to basically buy you Wow in a weird way See this this amethyst and the whole thing is so complicated the way they deal with I don't even wanna I don't even wanna think about it because you can tell me about it from now till next Till the heaven comes down.
I'm not even gonna be able to understand it
You know, the rules and regulations there, you have, I mean, it's a lot, it's a lot to understand.
There is a ton.
Too much sometimes.
And sometimes you even look into it, it's like, does this even make sense that this rule have to be here, this has to be there, you know?
Yeah, a lot of rules just never left.
They were applied when it was needed, but then when it wasn't needed anymore, they never took it away.
They never took it away, they just leave it there.
It became a tradition and not a rule.
Yeah, yeah, yeah.
And maybe in a whole 5 years, 10 years, right, then something will happen, they will go and dust it off and say, well, this is the rule, you know?
But it's not needed.
Because out of 100, the rules is used one percent of the time, does not mean that it should be there.
Exactly, yeah.
But yet they keep them.
It's like wearing your cover outside in the military.
If you go outside and you're in uniform, you have to wear some kind of hat, whether your command ball cap or your normal like naval cap or something.
And that all started off just as people used to complain that their heads were getting sunburned out in the sun while they were working in the military.
So they got them all.
something to wear on their head.
And they made it a regulation for a time being to make it so they had to wear it on their head when they're outside, so they didn't get sunburned.
Well, that rule eventually turned into a tradition, kind of, where now it doesn't matter what job you're in or what you do, you have to wear your ball cap outside if you're in your uniform.
Really?
Or you always have to take it off inside because it's technically disrespectful to wear them inside.
So, like, you have to always put them on and off every time you go inside or out.
Somebody like me, they will fire me a long time ago, because I'm not even gonna remember to put it on or to take it off!
Yeah, it's still weird, the little things that... Yeah, I thought that what they would have done is make it like, if you want it, we make it available to you, right?
If you wanna wear it when you're out, you wear it.
If you don't feel like wearing it, you don't have to.
But here it is, that we have it for you if you need it.
It's by choice, optional.
That would be nice.
Wow.
I have that option now.
It's kind of ingrained in me now, though, like I won't wear hats inside most of the time.
It's been so, it was drilled into me so much.
Once you become a soldier, you will never leave that character, the behavior, the personality stays with you forever.
You know, they'll always be that side.
Oh yeah.
You do it without even knowing that you're doing it.
Yup.
You know, a lot of the training I've, I've run into situations where I've always thought like when I was younger, like I stumbled upon someone that I thought was dead at one point.
Thankfully they were not.
Um, but in those situations you never know, like, am I going to freeze up?
Am I going to run away and do something?
And thankfully I just went on autopilot and followed the training they gave me and everything worked out.
It was interesting.
She was drunk and collapsed and passed out by herself on the side of the road.
Um, and we ended up having to follow her because she decided to try to run away from us instead of staying on the ground when she was hurt.
But yeah, stuff like that.
But I was like, okay, well, at least I know in the moment I handle the situation well, whether or not I can cope with it mentally after.
But in the moment, at least I go on autopilot and do what I need to do.
All right, so the medications you are doing well on it, you said that your symptoms are better than the last visit?
Yes, I would say they're better.
I don't know if it was going up in Zoloft that did it or if it was the Abilify.
Obviously, I don't want to stay on the Abilify forever.
I've been eating less and I've still been gaining a little bit of weight, but I'm not
I'm still, I believe, only 5 pounds above where I was when I started taking it, so it's still very reasonable.
I just know I should be losing weight, but with how much I'm eating, I shouldn't be gaining.
You still go to the gym?
It's still very on and off, but I still do go.
I actually have, I don't know if you can see it, but right down
There, yeah, that little piece is actually part of a home gym kit that has like up to 200 pounds of resistance I can use with it.
Um, so I'm trying to use that too.
So even if I don't feel socially free, I can still work out at home.
Okay.
That's good.
All right.
So I see that, uh, we are out of most, uh, you are out of most of the medications.
I see you forgot to schedule.
Yup.
Yeah, I was very caught up with the diagnostics one, and then got the emails for those, thinking it was an appointment for this.
And then yesterday, that's when I found out that I had that appointment for the thing, and I was like, okay, I gotta rush there.
And then by the time I got out of that, I completely forgot to even call the appointment in the 1st place.
And Klonopin, you're still taking it, right?
Yep, yeah, I still take.
It's usually only once a day now at this point.
Well, you have it twice a day if it's needed, right?
I'm sorry, I gotta change the battery in my headset.
It's about to die.
One second.
Okay, I can hear you again.
Okay, you can hear me?
Yeah, my headset gives me like a 10-second warning once I die.
Okay.
All right.
So, you're taking it once a day, but, you know, we keep it as twice a day if you need it.
If you don't need it, that is fine, but I want you to have it in case if you need it, so you don't have to struggle with that.
Okay, do you need the lamotrigine today?
On lamotrigine, I have one for tomorrow still.
In Zoloft, I have a little bit for tomorrow.
Okay, so you need a refill?
You need a refill on it?
Yeah, I need a refill on all the medications, I think, except the Klonopin.
Okay.
So the Klonopin, you have enough that will last you until in the next 30 days?
If I keep taking one, yes.
Okay.
Yeah, I should have more than enough of that.
Okay.
I had it right here with me, so I just wanted to double-check that.
All right, so what I'm going to be sending today, you say you don't need the Klonopin?
Nope.
Yeah, it's the lamotrigine.
the Zoloft and the Abilify.
Okay.
I'm going to either be out of or I am out of.
Okay.
All right.
So I will send the Vyvanse.
Uh, we're doing the Vyvanse.
Oh, and Vyvanse, yes.
Thank you.
We're doing the Vyvanse, 15 mg once a day, the Abilify, 5 mg once a day, the Zoloft, we're doing a 15 mg, you take 2 once a day.
Um,
The lamotrigine 200 mg, you take 2 tablets, making it 400 a day, right?
Yes.
All right.
And everything is still going to go to the same pharmacy?
Yep.
All right.
Well, today is the day without Vyvanse, because I completely forgot I didn't even get to take it.
Oh, you didn't take it today?
No, today is the 1st day I ran out of the Vyvanse.
Oh.
It's kind of funny that I forgot that I ran out of Vyvanse.
Okay.
All right, I will send all, and now you're going to call them, don't forget, okay?
So call them, schedule them again a month from today.
If anything changes, call us and let us know, okay?
Alright, stay well.
Alright, bye-bye.ss`

const sample3 = `
How are you?
I'm good, how are you?
I'm doing well, thank you.
So just a follow-up on the changes that we made the last visit.
So have you been, has there been any mood swings, mood changes, any fluctuations at all?
A little bit better the last couple of weeks, but not, you know, it's a little bit better.
A little bit, yeah.
Okay, that is expected, but remember I had told you it might take 46 weeks?
till it's kicking back.
The headaches, though, has that subsided a little?
They're coming, though, but I'm not sure, you know, I have blood pressure, too, so it might be because of the blood pressure.
Oh, that is, that is, so are you taking any blood pressure medication?
Yes, you are, I turned it on, yes, yes, yes.
But when was the last time you saw your PCP?
Maybe the dose needs to be adjusted, because I see that you're only on, what, 25 mg?
Yeah, so maybe I think you should if it's not going away after the switch, I would suggest that you make a follow up for them to kind of reevaluate.
Sometimes the medication needs to kind of needs to be increased, you know, needs to be adjusted, which you never know until maybe you get seen.
So I think I think you should do that.
Other than that, how is your mood?
How would you describe it?
Is it is it is it?
You know, something that is anything going on?
Is there any triggers, anxiety, any depressive symptoms?
Sure, yeah.
Hello?
Yes, I'm here.
Yeah, I don't know, I had a really bad anxiety attack.
I had a attack last night, but it wasn't as bad, but on Saturday night, maybe Saturday, I was sleeping and I jumped out of my sleep and I just couldn't calm it down for a while.
Is that new or it's happened before?
Hmm.
Okay.
Sometimes it happens when I'm sleeping.
Sometimes it happens when I'm, you know, I just can't figure out what's triggering it.
I just, and I think that's what irritates it either more sometimes, like when I'm up during the day and I had really bad anxiety or, um, you know, my PTSD kicks in, I'm raising three of my grandsons and one of them, I mean, three of my grandchildren and one of them, the nine year old, he has really bad ADHD.
So when he starts bouncing around and getting all over the place, I noticed that sets it off.
But that's sometimes, but sometimes I don't know what sets it off.
I don't know what triggers it or like, I even try to like sometimes go back on my steps and like, it's like, okay, what did I just do to, you know what I'm saying?
But- Yeah.
So what- Especially when I'm sleeping.
Usually when you're sleeping.
So when that happens, what do you do?
Do you take the lorazepam and when you, if yes, does it help when you do?
I don't take it right away.
I get up and I try to walk around and walk it off.
And then if I see it's still doing it, I'll take one point five.
And then like the other night, what I did was I did dishes and then I folded laundry.
Um, I opened the, you know, I sat on the front porch for a little bit.
I live in Randolph now.
Oh yeah.
Yeah.
I live in Randolph.
I can barely get out of my house.
I'm so afraid of that.
And on the ring, the neighbors, like some coyote, last week, a coyote grabbed somebody's dog and took it away.
Wow.
I'm like, oh my God.
That's crazy.
Crazy, crazy, crazy.
Yeah, so I was just, I would suggest that, you know, the lorazepam, just take it as you're doing.
Just don't take it if you don't need it.
Just so you, you know.
Yeah.
calm down before I jump, like, you know what I'm saying?
Even with any medication, usually anytime that's switched, I don't really, like, I'm trying my best to, like, I listen to a lot of gospel music too when I sleep.
So I try to, like, calm myself before, but sometimes it's just so overbearing.
Okay.
Yeah.
So I would suggest that, you know, when that happens, that's why you have the Ativan.
It happens, just use it.
It's calming, it's going to calm you down.
But Lexapro, I would suggest that we go up, we started at 10, maybe we should increase it just so you're kind of having that anxiety controlled throughout the day.
But I know you really don't like medication, so it's something that you...
Because I'm raising 3 kids, I don't want to look like I'm out of it.
I don't want to feel like I'm out of it.
That is totally understandable, especially if your grandson has an ADHD.
You have to be very careful and be alert all the time.
The older one, she's so teenage, she's sneaky.
I get so upset with it.
Their mom is my daughter and she's living her best life.
She has a home relationship.
She knows a lot about the kids, and I'm just like, I'm over here struggling.
Are you the primary caregiver of them?
I legally adopted them.
Oh, okay.
Wow.
When she had the baby, the 9-year-old, she had got postpartum depression and got really, she developed a really bad drinking problem, so she's a really severe alcoholic.
Okay, okay, so you're just helping out, making sure the kids are safe and good.
That's good.
Well, you see, I took them away, I took them, and I took them with the intention that my daughter
would get herself together, because I had debris, though.
My kids are all grown.
So I thought she would get herself together, and she never did.
She still drinks a lot.
She's hospitalized at least once a month from her drinking.
Oh, man.
Yeah, so that's a lot of burden on me.
That's a lot, yeah.
Another reason why my anxiety, my PTSD kicks in, because I don't know how to explain it.
I'm not angry, I'm just hurt that it was my daughter that hurt my grandkids.
And seeing them without their parents, because their father didn't step up either.
So, you know, it's just a little, it's hard, it's sad for their birthdays and all these roll around, and you know, their mother's always promising them stuff and then she never gets it because she'll go and buy stuff for herself or she'll go drink and then she'll disappear.
That is tough, and that caregiver strain on you can trigger your anxiety.
Now I understand where that panic is probably coming from.
Are you doing any form of therapy or are you interested?
Am I doing what?
Any therapy, like someone to talk to, like sometimes they do it virtual with you, you don't have to- No, it's a chance for me too, sometimes, because I don't really- I feel like when I keep talking about it, like- You're triggered more or you get- Yeah, yeah, exactly, it just brings up more.
More, okay.
Yeah.
Yeah, but that- Because I don't wanna- I don't wanna grow- I don't wanna, um, become angry, you know what I'm saying?
Mm-hmm, mm-hmm.
And it's just hard, you know, and my daughter's like,
You know, I, when I, when I tell her she can do it, you know, it's been 9 years.
You could have stopped drinking by now, but you chose not to.
You know, you're in a whole relationship.
How can you be in a relationship by not taking care of your kids?
Does she want care?
Like, you know, the substance use department in, um, uh, Codman.
She gets everything.
She, she tried, she, listen, she, she, um, is very good at maneuvering the system.
She likes to play the victim.
She tells her bad, took her kids away from her.
Listen.
I raised 6 kids, I do not want no more.
It's good that you're stepping in to find a good home for these, for your grandkids, but it's also a lot on you.
I understand where you're coming from when it comes to medication, because you have to be alert all the time.
You cannot be like, you're out of it, overly medicated.
But the Lexapro is not, it's not sedating.
It's more of, yeah, it has this calming effect.
It's just going to keep you stabilized.
You know, even if you take it at night, you're not going to sleep.
You're supposed to take it in the morning.
So, um, we can, we can increase it if you're open to it.
If you feel like you're okay, you're okay with what you're taking now to how are you sleeping well?
I thought we started magnesium.
Let me see.
Yeah, I didn't, I didn't take it.
I didn't do anything.
Did you even try?
Because taking that in addition to the gabapentin helps a lot.
All right, I'll try it.
Yeah, try it.
Try it.
All right.
I didn't even prescribe it.
I think they've given it to you a while, say back in July.
Yeah.
Yeah, you never took it?
So now I think they started, let me look at it, 400 mg, take 200 by month.
Okay, so start off with half of the 400.
Don't start with the 400 because you don't know how it's gonna affect you.
You don't wanna be out of it like you said.
So take half.
Take half and save one of the gabapentin for nighttime.
You take one at night anyways, right?
Okay.
The gabapentin, right?
Yeah.
Yeah, so take one of the gabapentin, the nighttime dose, in addition to half of the magnesium.
You should be able to sleep.
Okay.
Yeah, try that.
Let's see how that goes during the next follow-up.
Appetite, any issue with that?
I mean, I really don't eat.
I mean, I'm having issues with my weight.
I don't know if it's because of menopause or something's going on.
I have to go see endocrinology in 2 weeks, I think, or 3 weeks.
Are you losing weight?
I was trying.
I was on the shot, but the shot wasn't working.
It's making me sick.
Oh, yeah.
Is it what, Wegovy or the other one?
It was Zepbound.
Yeah, it's making a lot of people sick.
A lot of people that I've seen said, okay.
What?
How much do you weigh now?
Roughly average?
Over 300, like 320, 330.
Oh, okay.
So the shot would have been helpful, but it's not for everyone.
Yeah, it wasn't working.
And then she was supposed to switch me over to, um, we'll go be, but then she wouldn't approve it.
I don't know.
The doctor was being weird, so I was just like, you know what?
Why don't you go to the weight management clinic in BMC?
That's where I'm going to see the immunology specialist.
Yeah, go there.
They think it might be my thyroid.
I don't know.
If you have the hypothyroid, then that will make you gain weight a lot. Yeah, because I don't eat, like yesterday I had a yogurt and I had a chicken sandwich, like a homemade chicken sandwich because I'm out eating out.
Like I had made chicken breast the night before and I made a chicken sandwich with just the chicken and the bread, that was it.
I didn't have like no extra nothing on it.
So you shouldn't really be gaining any weight.
Huh?
I said, so you shouldn't be gaining any weight.
Yeah.
But I'm not losing weight.
Yeah, yeah.
You know, I watched those shows 600,000 more.
I see how much I'm like, I'm really be sick if I eat that much, but I don't eat like that.
I'm just like, why am I not losing women?
I go, you know, it could be your thyroid.
It could be right.
Um, but that's a good starting point.
Sometimes the PCPs are not comfortable with all these medications, but with the weight management, they specialize in it and people do have a lot of success.
So, um,
Yeah, I hope you're able to find something that's a good fit for you there.
Maybe sometimes just phentermine alone can do it.
Have you tried that?
Tried what?
Phentermine.
Vitamins?
No, phentermine, phentermine.
What is that?
It's a medication they usually give for weight.
It helps your weight too, but yeah.
You know, until they kind of run blood tests, sometimes it's not good for your liver and all that.
So when you go there, they'll do, they'll do all the necessary blood work and find something that would help you lose the weight.
Yeah.
Okay.
Well, I appreciate it.
All right.
So try for the sleep, for the sleep, try the magnesium and one, um, one, um, one tap one, let me see the gabapentin, see the nighttime dose, take it in addition to.
half of the magnesium, get some good sleep.
Sometimes getting a good sleep can also change a lot of things.
Like if you feel better, your anxiety is well controlled.
So let's start off from there.
Um, so for the Lexapro, let me see, do you want me to increase it for you?
Okay.
The only thing is, um, what I can do, do you like splitting pills or you want me to,
It's fine, I have a splitter, I have a splitter.
Okay, so if you have a splitter, then I would, let me go in Lexapro, where is it?
Okay, so I'll change it to one and a half tablet, which will be 15 milligram in the morning.
Okay.
Do you understand that?
Yes, ma'am.
Okay, so one-and-a-half in the morning, and hopefully by the time we get to that 6 weeks, you should get some relief from it.
And I will check in with you.
Let me make that follow-up appointment.
Do you want to come in?
It depends on what time the time of day is.
Okay, what time usually works for you?
Nine-thirty would be good.
I dropped my grandson, I leave him at 820, so from there to there, it's a little bit traffic, so I think I get there by 930.
But if you do me a favor from there at 930, can you call me just to check on me?
What did you say?
I said if I don't show up at 930, can you call me at 930?
Because I have really bad, like,
Um, I always like get up and I'm like, Oh, I'm going to do this today.
And then I'll send something just like, I think what they do is they call the night before the day before to remind you.
No, I know that, but I'm saying that.
So like, I'll say I get up and I'm sitting in my, my, the school and I'm waiting for him, waiting for him.
And then, I don't know, something like my anxiety will kick in and I'll just go home.
Like just to be in the house.
Like, you know what I'm saying?
If I don't come, if you
Okay.
All right.
But I'm gonna, I'm gonna do my best to film.
Okay.
No worries.
Okay.
Uh, follow up.
Okay.
So let's do a month.
Okay.
One.
So that one, 234.
Okay.
So how about May 1st?
Yes.
There's the May 1st.
Yes.
And I do have... Is Cinco de Mayo, are you gonna bring a drink?
I'll try my best.
Okay, so I have 930, I have 1130.
Which one?
We'll do 930.
Okay.
Because if I'm already out, I think I can do it.
Okay.
You see what I'm saying?
If I go out, if I go out, come home, I won't go back out.
Okay.
All right.
So that should be all set for May 1st at...
What is it?
May 1st at 930, yes.
Okay.
So I'll see if anything changes, you know, you know where to find me.
Send me a message on my chat call or call the main number and I would, uh, I would call you back.
Okay.
All right.
You're welcome.
You too.
Take care.
Yes.
Bye-bye.
`
const sample4 = `Good afternoon.
I'm Naomi Hilda, nice to meet you.
I know you were seeing Doctor Rison.
Have you been?
I'm alright.
Are you tearing?
I have dry eyes.
It's always like this?
Have you talked to your PCP?
It's allergies and I was told by the eye doctor that I should use the synthetic eye drops.
Change covers?
Oh, so maybe they won't cover.
It's expensive, I've used it before.
Anyways, how are you, and how have you been?
Alright.
Oh, how far, how far are you?
Well, I lived down the street, but I was filling out paperwork and making sure that my daughter was okay before I left the house.
And then, like, before I knew it, it went from 107 to 730, but... And so, as I'm on my way over here and I'm calling, and I'm like... You could have rescheduled, right?
Yeah, well, I've rescheduled all the time.
Yeah, listen, this is already a reschedule.
Oh, okay, okay.
All right.
Call and let them know that I'm on the way, I'm just gonna be a little bit late.
And the lady's on the phone talking about something.
You know you have the 15-minute grace period?
Well, I'm in my 12 minutes of the 15 minutes, I'm like 4 minutes away.
Yeah, but you know, you're not listening to me, you know what, let me just hang up.
Because by the time you get to understand what I'm saying, I'm gonna be walking through the building.
Luckily I didn't have anyone after you, so it worked out well.
Right.
Yeah.
So tell me, what's going on?
Anxiety?
Depression?
I see that you're prescribed Zoloft by Doctor Isom, which you take a 100 every day.
Does it help?
You take it every day?
Why?
Because I was taking it, and it wasn't doing nothing.
But...
You wouldn't know until you stop.
How long ago did you, did you stop?
What is a while ago?
You just stopped it like that.
Just stop taking it because that medication, you can't, you have to slowly, if you stop it, well.
What symptoms were you having that made, because it's prescribed for both anxiety and depression, was yours anxiety or more of depression?
Both.
Both?
How long did you take it for before you realized it wasn't doing anything?
Almost a month.
What's your name again, Christina?
Christina, it takes 4 to 6 weeks, sometimes 8 weeks for you to notice that it's working or not.
We're gonna have to restart it.
We're gonna have to restart it.
If you've stopped for a month, you're going to start off with 50 for a week and then continue to a hundred.
You're going to have to take that for 6 weeks.
If it doesn't help after 6 weeks, then there's stuff that I can add or increase it.
That medication, the highest you can go is 200.
You only have a hundred and you stop.
So if you had come in or you had reached out, she probably would have increased or added something more.
Tell me how you feel on a daily basis, just so I would know if you're sick.
Like, dry.
No energy?
What do you do during the day?
Are you working?
No?
So what do you do in the day?
When you wake up, go to school, what else do you do?
Your daughter doesn't go to school?
How old is she?
Oh, sure.
She goes, she goes to work and her job has provided her to sign up for the program to help her get her license.
So she does that on Tuesdays and Thursdays.
Okay.
When you say your daughter, I thought she's little.
17, she's a big woman.
That's your only child?
That's your only child?
Yeah.
Um, how do you sleep at night?
How many hours do you get?
Are you sleeping a lot or are you getting less sleep?
I'm like sleeping a lot, but it's like an interval.
So like if I go, if I happen to fall asleep at like nine and then wake up at like three, then I'm up to about like 10.
And usually what, what makes you wake up?
What makes you, is it noise or is just?
Is your mind going?
You can't sleep again?
You have ADHD?
I don't see it.
When did you get diagnosed?
So it's not in the end, like that's not a part of something I talk about.
I don't, I don't see that.
But seeing you just, um, the, like the man start working on a piece of scar because I do my braids myself, like.
All by yourself, even the back?
How do you do the back?
I always say people that are witches, they do that because there's only,
It's impossible, I can't even do it on your head, and I can't imagine you turning your head to... I don't know, I have a friend who does it too, I call her Rich.
I say, you must be rich.
I've been doing hair since like 12.
Wow, you're lucky, because it's damn expensive.
Sure is, like this bad boy would cost me a couple of pennies.
Yes, my girl.
Prices, just like four, five, six dollars a peg.
Wow, so you did all this by yourself?
How long does it take you?
A couple of days.
And just like, what is it, like, last week, like every year, almost every night of last week, I just do a couple of years, do a couple of days, because these have been in my head for a while.
For other people?
Yeah.
On your own head?
I mean, I do it to, like, my daughter.
I had a client and then I used to do my nieces hair.
But like last week I was doing it like often, like just refreshing up the braids that were in the middle of my hair.
They had been in here for like a month.
The new growth was out of braids here.
Hmm.
Okay.
So let's get back to this meds here.
The depression and the anxiety, which one is more persistent to you during the day, like on a daily basis?
Do you get depression?
How often do you get depression and anxiety?
Almost every day, but by night, I don't know, because I stay in the house to myself.
But then when I do get outside, you know, there's always
Are you in therapy?
Who do you talk to?
Let me see.
Oh, okay.
This person here, Tamara, okay.
How often do you see Tamara?
Every week?
Yeah, every other couple weeks ago.
Okay.
Um, so why are you not working?
You just don't want to or you can't?
Cause my daughter has type one diabetes.
Oh, type one.
Okay.
And it's like every time that I give a job, she forgets that she has diabetes.
Does she?
And it's like, I can't.
5 people, so she forgets that she has diabetes and stops taking her medicine and ends up in the hospital and the child fires me for being in the hospital with her.
So, why get a job if she says I'm insane, forget it, and then I get fired.
So how are you surviving now?
Any past psychiatric admissions?
How many?
One?
How did you- why did you get hospitalized?
Were you having any suicidal thoughts, ideation, thoughts of death? But there was just a whole lot of stress going on here.
And my assistant's trying to tell me fuck me is a more difficult word.
So before I say fuck me, while everybody else and everything around me is saying fuck me, let me go check in and see what's up with them.
You know?
How long did you stay there?
8 weeks.
Any substance use history, including alcohol, marijuana, drugs on the street?
What have you been using?
Marijuana and alcohol.
No, no heroin, cocaine, or nothing?
I can't do that.
My mom and my daddy did do drugs.
I'm not trying to be like, I guess you could call me an addict, whatever, because of how much I didn't drink alcohol.
But I'm trying to be an addict like that.
Do you have any medical history?
So this eye is always like this?
No, now I'm, like, high level.
Oh, okay.
I was gonna be like, this is, this is serious.
This is the worst I've seen.
You know, I see people with it, but yours is more... No, no, it's not that.
The 1st game we, it was the... The, okay.
Yeah, yeah, yeah.
When I talked about this, I was emotional.
Okay.
Do you have the energy to, to get up, wake up and stuff?
Like when you're, do you have the desire to get things done and all that?
Like I make stuff happen, but do I want to?
Do I have to?
Yes.
So.
Sure.
Sure.
There you go.
You can leave it there.
I'm sorry, I just, uh... How bad is the alcohol use?
Do you drink every day?
I can afford to, sir.
Okay.
It's not problematic, it's not?
Okay.
So for now you live with your daughter alone, right?
My daughter and my 3 dogs.
And 3 dogs.
Oh, you're both able to feed 3 dogs?
Huh, so you're not broke.
Dog food costs a lot, I heard.
3, not one!
You're going to my pantry, sir.
Oh, I see.
Okay, so that's not bad.
All right, um, I'm glad you're in therapy.
The sertraline, you took it for a month, you didn't give it time to work.
Some of the medications, I'll tell you, next time you're on medication, you want to stop, you want to get help to kind of slowly get off of it.
Don't just completely stop.
You were lucky that when you stopped it, you didn't feel like your brain was zapping.
Like sometimes you get that feeling with that medication.
If you stop, luckily your dose was not a lot.
It was just a hundred.
So maybe you got away with it.
But if you're on like 200 milligrams and you stop it right away, properly, you're not going to feel good.
You're not going to feel good.
Okay.
So we're going to try another SSRI.
We won't try the same thing if you didn't really feel anything, but although you didn't give it time to work,
It's the same class of medication.
It helps with anxiety, helps with depression, but it's more activating.
Activating means it will give you the energy to get up and kind of desire to work on and get things done.
So I think we should do Prozac.
I would start you off with that.
If you tolerate it, I will increase it.
If not, there are options that we can take a look at later.
Sleep.
Have you tried melatonin for sleep?
I mean, for now, how many hours beside getting up in between?
Do you feel well rested when you wake up?
Okay, so that, like, if you feel fine, then let's sleep, just sleep for now until you can sleep at all.
Then maybe we can start melatonin or something all natural.
So which pharmacy do you use?
Downstairs?
Okay, so I'll put in Prozac for you.
I'll start you with a very low dose to see if your body can handle it, you can tolerate it.
If you do, then I can increase it for you.
So let me see here.
So just like the other one, don't give up and say it's not working and stop.
It takes at least a good 6 weeks.
Sometimes I tell patients, give it 2 months before you start complaining or give up on it.
Take it consistently every morning for 6 to 8 weeks.
Then once we reach the max, you don't feel good, then we know it didn't work.
But one week, 2 weeks, even one month is not going to do it.
After 6 weeks, that's when you start to notice it.
some effects from it.
So we'll do a Prozac, give it 30.
Okay, so I sent that for you.
Let me make an appointment with you in 3 weeks to see that you're tolerating it.
There's no side effect, then I can increase.
The dose I gave you today is really not a big dose, a very small baby dose.
This is because I'm making sure that you don't have a reaction to it.
You will tolerate it, and then I'll see you in 3 weeks.
If you're fine, no side effect, then I'll increase it for you, and then we'll wait for the period, time period, okay?
Would you like to come in?
Yeah.
Don't do much.
Great.
That'll be 3 weeks from today.
I have 12, I have 1230, I have 1130 too.
What was for you?
All right.
Do you get a text reminder on your phone?
Okay, so 1230, March 5th, which is 3 weeks from today, I'll see you, but I'll tell you this.
If you take the medication and for some reason you get side effects like rash or anything that you feel like this is new, this is weird, this is completely new before I started the medication, don't continue to take it.
Stop it and call and let me know and I'll call you and guide you as to what to do from there, okay?
So you don't have to wait for the 3 weeks if you're not doing too good on it, if you're having side effects, but if you feel fine, then continue to take it.
Even if you don't notice a change in the way you feel, continue to take it.
It will take a good 6 weeks to a month for you to start noticing
a change in the way you feel.
Okay.
A card.
So, as you can see, they moved everything around because they changed the form.
So I came into there and I can find anything.
Oh, there's some here.
It's not my name on it, but I can remind you.
Do you want the appointment on it?
Or just my phone number?
Oh, no, I don't have that.
Yeah, there's a phone number here.
But I don't have a direct one, so you'll have to go through the main switch vault and do let me know if there's anything going on.
All right.
Nice meeting you, and I will check in with you in 3 weeks.
And we will come up with a plan to make you feel a lot better.
All right.
You take care.
No, don't worry.
Don't worry about it.
See you later.
You're welcome.
All the way.`
func (m *mockTranscriber) Transcribe(ctx context.Context, audio []byte) (string, error) {
	return sample4, nil
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

	transcriber := azure.NewAzureTranscriber(
		cfg.OpenAISpeechURL,
		cfg.OpenAIAPIKey,
	)
	// fmt.Println(cfg.OpenAIChatURL,"checking")


	inferenceService := inferenceService.NewInferenceService(
		reportsStore,
		transcriber,
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
