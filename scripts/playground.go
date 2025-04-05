package main

import (
	inferencestore "Medscribe/inference/store"
	dbhelper "Medscribe/scripts/dbHelper"
	promptDistiller "Medscribe/scripts/distelledPrompts"
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

const (
	apiKey      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ewogICAgInJvbGUiOiAiYW5vbiIsCiAgICAiaXNzIjogInN1cGFiYXNlIiwKICAgICJpYXQiOiAxNjc4MjYyNDAwLAogICAgImV4cCI6IDE4MzYxMTUyMDAwCn0.Dqu6ADslVdPp9DdlxMobJ6rAH_Uess-yDS4OXlHIlPk"
	bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYWwiOiJhYWwxIiwiYW1yIjpbeyJtZXRob2QiOiJvYXV0aCIsInRpbWVzdGFtcCI6MTc0MTk2MTIzMH1dLCJhcHBfbWV0YWRhdGEiOnsicHJvdmlkZXIiOiJnb29nbGUiLCJwcm92aWRlcnMiOlsiZ29vZ2xlIl19LCJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZW1haWwiOiJyaXNpbmdoaWxkYUBnbWFpbC5jb20iLCJleHAiOjE3NDI1NjYwMzAsImlhdCI6MTc0MTk2MTIzMCwicGhvbmUiOiIiLCJyb2xlIjoiYXV0aGVudGljYXRlZCIsInNlc3Npb25faWQiOiIxYTE2ODIwMC01MWJhLTQ5ZDQtYWJmMy0zYmI5Nzk5NDFhNzIiLCJzdWIiOiJkNWJjMzdjMC0yMzQ5LTRiOGItOGI0NS0yNjg5OThmMTBjOWYiLCJ1c2VyX21ldGFkYXRhIjp7ImF2YXRhcl91cmwiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NMc3FJV2RERkFfczZqR0RwR1p5WVg4aU03TDlaYUl4c3pydjFFMWJDdGJqejd4YXc9czk2LWMiLCJlbWFpbCI6InJpc2luZ2hpbGRhQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmdWxsX25hbWUiOiJSaXNpbmcgSGlsZGEiLCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJuYW1lIjoiUmlzaW5nIEhpbGRhIiwicGhvbmVfdmVyaWZpZWQiOmZhbHNlLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jTHNxSVdkREZBX3M2akdEcEdaeVlYOGlNN0w5WmFJeHN6cnYxRTFiQ3Riano3eGF3PXM5Ni1jIiwicHJvdmlkZXJfaWQiOiIxMDM5NTIyNTk4MTU2NDQ3MjIyNjEiLCJzcGVjaWFsdHkiOiJQc3ljaGlhdHJ5Iiwic3ViIjoiMTAzOTUyMjU5ODE1NjQ0NzIyMjYxIiwidXNlckNyZWF0ZWRBdCI6IjIwMjQtMDItMDhUMTU6MTM6MjYuMDE3MjAzIn19.Pmt8eh4CikI86iJ5wmXg4_WfolndXWHvid-GEOi8xEo"
)

var config *dbhelper.Config
var dbHelper *dbhelper.DatabaseHelper
var distiller promptDistiller.PromptDistiller

func init() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize configuration
	var err error
	config, err = dbhelper.ResourcesConfig()
	if err != nil {
		log.Fatalf("Error initializing resources: %v", err)
	}

	// Initialize MongoDB collections
	dbHelper, err = dbhelper.InitializeCollections(context.Background(), config)
	if err != nil {
		log.Fatalf("Error initializing collections: %v", err)
	}

	query := inferencestore.NewInferenceStore(config.OpenAIChatURL, config.OpenAIAPIKey)
	distiller = promptDistiller.PromptDistiller{
		FreedCollection:    dbHelper.GetFreedVisitsCollection(),
		AnalysisCollection: dbHelper.GetAnalysisCollection(),
		Query:              query,
	}
}

func processAndStoreDistilledReports(limit int) error {
	analyeses, err := distiller.DistillMultipleReports(context.Background(), limit)
	if err != nil {
		return fmt.Errorf("error distilling multiple reports: %w", err)
	}

	err = distiller.PutDistilledAnalysis(context.Background(), analyeses...)
	if err != nil {
		return fmt.Errorf("error putting analysis into collection: %w", err)
	}

	return nil
}

func main() {
	fmt.Printf(fmt.Sprintf("Emenike ", "emenike"))
	// processAndStoreDistilledReports(5)

	// analysis, err := distiller.GetTopKDistilledAnalysis(context.Background(),5)
	// if err != nil{
	// 	fmt.Println("error getting distilled report: %w",err)
	// }
	// fmt.Println(analysis)
	// reports, err := distiller.GetReports(context.Background(),10)
	// if err != nil{
	// 	fmt.Println("error getting report: %w",err)
	// }

	// ap := []string{}
	// for _, value := range reports{
	// 	ap  = append(ap , value.Summary)
	// 	fmt.Println(value.Summary)
	// 	fmt.Println("")
	// }

	// objective := []promptDistiller.DistilledAnalysis{}
	// for _, value := range analysis{
	// 	objective = append(objective, value.Objective)
	// }

	// fmt.Println(objective)

	// fmt.Println(analyses)

	// subjetives := []string{}
	// for _, value := range analyses{
	// 	subjetives = append(subjetives, value.Subjective.TeacherData)
	// }

	// fmt.Println(subjetives)
	// fmt.Println(analysis[0].Subjective.DistilledPrompt)

	// fmt.Println(analysis[0].Subjective.TeacherData)
	// fmt.Println(analysis[0].Subjective.DistilledOutput)

	// superPrompt, err := distiller.GenerateSuperPrompt(context.Background(),"subjective",subjectivePrompt)
	// if err != nil{
	// 	fmt.Println("error Generating super prompt: %w",err)
	// 	return
	// }
	// fmt.Println(superPrompt)

	// prompts, err := distiller.GenerateSuperSoapPrompts(context.Background(),analyses)
	// if err != nil{
	// 	fmt.Println(err)
	// }
	// 	var prompt = `
	// You are a medical provider documenting objective clinical data from a patient encounter transcript. Write the report clearly, concisely, and in plain text without markdown, incorporating specific details or examples directly from the patient's transcript to accurately reflect unique characteristics of the encounter. Completely omit any category that is neither explicitly mentioned nor confidently inferable.

	// Mental Status Examination:
	// - Behavior: Briefly describe patient's observable behavior, interaction style, and demeanor, integrating distinct and specific details or examples from the conversation if available (e.g., "Calm and cooperative, openly shared concerns about medication side effects," or "Withdrawn, minimal responses, appeared distracted during questioning"). Avoid overly generic phrases unless no distinctive details are provided.
	// - Speech: If explicitly mentioned or confidently inferred, briefly describe speech characteristics including pacing, tone, or clarity (e.g., "Speech rapid when discussing work stress," or "Quiet and hesitant when mentioning family conflicts"). Omit entirely if not clearly observed or documented.
	// - Mood: Use patient's direct quotes if provided, or succinctly describe mood based on patient's expressed emotions or tone (e.g., "Expressed feeling overwhelmed by family obligations," "Mood described as stable with improvements noted since last visit"). Avoid overly generic descriptions.
	// - Thought Process: Default to "Linear and goal-directed" if coherent; if not linear, briefly specify and describe clearly observed deviations using specific conversation examples (e.g., "Occasionally tangential, patient frequently shifted topics when discussing future plans," "Patient’s responses were focused but included excessive irrelevant details").
	// - Cognition: Succinctly state "Alert and oriented to conversation" unless explicit cognitive concerns (e.g., confusion, memory issues) are observed; provide a concise description and relevant examples if any cognitive issues are noted (e.g., "Mild confusion, difficulty recalling recent medication adjustments").
	// - Insight: Include if patient demonstrates clear awareness or understanding of their condition or treatment; support with specific examples or statements from the patient if possible (e.g., "Good insight demonstrated by proactive discussion of medication management," "Limited insight into the severity of reported anxiety symptoms").
	// - Judgment: Include if patient shows decision-making ability or planning that can be explicitly or implicitly inferred from the conversation; use examples from patient interactions if available (e.g., "Judgment appears fair; patient actively schedules follow-ups and adheres to medication despite reported side effects").

	// Vital Signs:
	// - Include explicitly stated vital signs (e.g., blood pressure, heart rate). Omit entirely if not explicitly stated.

	// Physical Examination:
	// - Include explicitly stated physical findings, briefly summarized by system with relevant details (e.g., "Tenderness noted in left ankle," "Clear lungs on examination"). Omit entirely if no physical exam mentioned.

	// Pain Scale:
	// - Clearly document numeric pain rating explicitly reported by the patient, with date and reference to scale (e.g., "Pain rated 8/10, described as severe and constant"). Omit entirely if not explicitly provided.

	// Diagnostic Test Results:
	// - Concisely summarize explicitly mentioned diagnostic tests or results (e.g., blood tests, imaging findings). Omit entirely if no test results mentioned.

	// Additional Relevant Information:
	// - Briefly document explicitly stated patient details directly relevant to clinical care decisions, medication management, or follow-up arrangements that have not been mentioned elsewhere. Avoid redundancy.

	// Provide all information concisely and specifically in plain text format without markdown, emphasizing distinct, transcript-specific details to avoid repetitive or overly generic documentation.
	// `

	// 	prompts := promptDistiller.Prompts{
	// 		Objective: prompt,
	// 	}
	// 	_, queries, err := distiller.Benchmark(context.Background(), prompts, 1)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println("----------")
	// 	fmt.Println("prompt: ", queries[0].Prompt)
	// 	fmt.Println("----------")
	// 	fmt.Println("transcript: ", queries[0].Transcript)
	// 	fmt.Println("----------")
	// 	fmt.Println("teacher data: ", queries[0].TeacherData)
	// 	fmt.Println("----------")
	// 	fmt.Println("distileld output from prompt: ", queries[0].DistilledOutputFromPrompt)
	// 	fmt.Println("----------")
	// 	fmt.Println("scores: ", queries[0].Scores)

}

// func main() {

// 	scraper := promptDistiller.FreedScraper{
// 		API_KEY:      apiKey,
// 		ACCESS_TOKEN: bearerToken,
// 		BASE_URL: "https://secure.getfreed.ai/",
// 	}

// 	visits, err := scraper.GetRecentVisits(30)
// 	if err != nil {
// 		fmt.Println("Error in GetRecentVisits:", err)
// 	}
// 	fmt.Println(visits)

// md, err := scraper.GetVisitMetadataIDs([]string{visits[0]})
// if err != nil {
// 	fmt.Println("Error in GetRecentVisits:", err)
// } else {
// 	fmt.Println("Recent Visits:", md)
// }

// payload := promptDistiller.IDs{
// 	IDs: []string{"f26d93a0-c151-4970-89f5-17808b1be7c3"},
// }
// txt, err := scraper.GetSectionText(payload)
// if err != nil {
// 	fmt.Println("Error in GetRecentVisits:", err)
// } else {
// 	fmt.Println("Recent Visits:", txt)
// }
// }

// // func main(){

// // 	if err := godotenv.Load(".env"); err != nil {
// // 		fmt.Errorf("error loading .env file: %v", err)
// // 		return
// // 	}

// // 	freedCollection := os.Getenv("MONGODB_FREED_VISITS")

// // 	config, err := dbhelper.ResourcesConfig()
// // 	if err != nil{
// // 		fmt.Println("error initializing reasources: %w", err)
// // 		return
// // 	}

// 	// Connect to MongoDB
// 	// clientOptions := options.Client().ApplyURI(config.MongoURI)
// 	// client, err := mongo.Connect(context.Background(), clientOptions)
// 	// if err != nil {
// 	// 		fmt.Println("error connecting to MongoDB:", err)
// 	// 		return
// 	// }
// 	// defer client.Disconnect(context.Background())

// 	// collection := client.Database(config.MongoDBName).Collection(freedCollection)

// 	// query := inferencestore.NewInferenceStore(config.OpenAIChatURL,config.OpenAIAPIKey)

// 	// distiller := promptDistiller.PromptDistiller{Collection: collection,Query: query}
// 	// reports, err := distiller.GetReports(context.Background(), 2)
// 	// if err != nil{
// 	// 	fmt.Println("error getting reports: ", err)
// 	// }

// 	// analysis, err := distiller.DistillReportPrompts(context.Background(),reports[0])
// 	// if err != nil{
// 	// 	fmt.Println("error occured distilling report prompts ", err)
// 	// 	return
// 	// }
// 	// fmt.Println(analysis.Subjective.DistilledPrompt)

// 	// analyses, err := distiller.DistillMultipleReports(context.Background(), 2)
// 	// if err != nil{
// 	// 	fmt.Println("error occured distilling report prompts ", err)
// 	// 	return
// 	// }
// 	// fmt.Println(analyses)

// }

// func main(){

// 	if err := godotenv.Load(".env"); err != nil {
// 		fmt.Errorf("error loading .env file: %v", err)
// 		return
// 	}

// 	freedCollection := os.Getenv("MONGODB_FREED_VISITS")

// 	config, err := dbhelper.ResourcesConfig()
// 	if err != nil{
// 		fmt.Println("error initializing reasources: %w", err)
// 		return
// 	}

// 	// Connect to MongoDB
// 	clientOptions := options.Client().ApplyURI(config.MongoURI)
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 			fmt.Println("error connecting to MongoDB:", err)
// 			return
// 	}
// 	defer client.Disconnect(context.Background())

// 	// Get the collection
// 	collection := client.Database(config.MongoDBName).Collection(freedCollection)

// 	distiller := promptDistiller.PromptDistiller{Collection: collection}
// 	reports, err := distiller.GetReports(10)
// 	if err != nil{
// 		fmt.Println("error getting reports: ", err)
// 	}

// 	fmt.Println(reports)

// }

// func main(){

// 	if err := godotenv.Load(".env"); err != nil {
// 		fmt.Errorf("error loading .env file: %v", err)
// 		return
// 	}

// 	scraper := promptDistiller.FreedScraper{
// 		API_KEY:      apiKey,
// 		ACCESS_TOKEN: bearerToken,
// 		BASE_URL: "https://secure.getfreed.ai/",
// 	}

// 	teacherData,err := scraper.ScrapeReports(2)
// 	if err != nil{
// 		fmt.Println("error sraping: %w", err)
// 		return
// 	}

// 	fmt.Println(teacherData)

// freedCollection := os.Getenv("MONGODB_FREED_VISITS")

// config, err := dbhelper.ResourcesConfig()
// if err != nil{
// 	fmt.Println("error initializing reasources: %w", err)
// 	return
// }

// // Connect to MongoDB
// clientOptions := options.Client().ApplyURI(config.MongoURI)
// client, err := mongo.Connect(context.Background(), clientOptions)
// if err != nil {
// 		fmt.Println("error connecting to MongoDB:", err)
// 		return
// }
// defer client.Disconnect(context.Background())

// // Get the collection
// collection := client.Database(config.MongoDBName).Collection(freedCollection)

// // Insert the teacherData
// for _, data := range teacherData {
// 	_, err := collection.InsertOne(context.Background(), data)
// 	if err != nil {
// 			fmt.Println("error inserting data:", err)
// 			return // or continue, depending on your error handling preference
// 	}
// }

// fmt.Println("Successfully inserted teacherData into MongoDB.")

// }

// func main(){
// 	config, err := dbhelper.ResourcesConfig()
// 	if err != nil{
// 		fmt.Println("error initializing reasources: %w", err)
// 		return
// 	}

// 	dbHelper, err := dbhelper.InitializeCollections(context.Background(),config)
// 	if err != nil{
// 		fmt.Println("error intializating collections: %w", err)
// 		return
// 	}

// 	query := inferencestore.NewInferenceStore(config.OpenAIChatURL,config.OpenAIAPIKey)
// 	if err != nil{
// 		fmt.Println("failed to setup inference store: ", err)
// 		return
// 	}

// 	scraper := promptDistiller.FreedScraper{
// 		API_KEY:      apiKey,
// 		ACCESS_TOKEN: bearerToken,
// 	}

// 	teacherData,err := scraper.ScrapeReports(10)
// 	if err != nil{
// 		fmt.Println("error sraping: %w", err)
// 		return
// 	}
// 	distiller := promptDistiller.PromptDistiller{Query:query, Collection: dbHelper.GetAnalysisCollection()}
// 	for key,value := range teacherData{
// 		analysis, err := distiller.Distill(key, "subjective",value.Subjective,value.Transcript)
// 		if err != nil{
// 			fmt.Println("there was an error during distillation: ", err)
// 			return
// 		}

// 		// fmt.Println(analysis.DistilledPrompt)
// 		if err := distiller.PutDistilledAnalysis(context.Background(),analysis); err != nil{
// 			fmt.Println("there was an error when putting the analysis in collection: ", err)
// 			return
// 		}

// 		return
// 	}
// }
