package inferenceService

import (
	"Medscribe/reports"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const(
	// Soap Task Descriptions
	subjectiveTaskDescription = "Extract patient-reported symptoms, concerns, and history relevant to the encounter."
	objectiveTaskDescription = "Identify clinician-observed findings, including vitals, examination details, and test results."
	assessmentTaskDescription = "Summarize diagnosis or evaluation of the patient’s condition based on subjective and objective data."
	planningTaskDescription = "Outline recommended treatment from the transcript (if provided), follow-up instructions, and next steps for care."
	summaryTaskDescription = "Summarize the key points of the patient's history, symptoms, and examination findings."

	//Default Return format
	defaultReturnFormat = "Return Format (Always Adhere to These Rules):" + 
	" Responses must be strictly plain text, suitable for direct display in a textbox." + 
	" Never use markdown formatting (no \"*\", \"-\", \"#\", or \"---\")." + 
	" Do not include headers, special characters, or extraneous punctuation." + 
	" Ensure all responses are brief, precise, and directly related to the requested component." +
	" Do not return a response as if you were responding based off of a question." + 
	" All your queries will be aggregated adn decorated from the client. Do not give any indication you were prompted iteratively"

	//Default Warnings
	defaultWarnings = "Warnings:" + 
           " Never provide information or details outside of what's explicitly requested." +
           " Always rely exclusively on the provided transcript without assumptions or inference beyond clearly available context." +
		   " if more context is needed do not hallucinate and just specify that more context is needed"
)


func GenerateReportContentPrompt(transcribedAudio, soapSection, style string) string {
	var taskDescription string

	switch soapSection {
	case reports.Subjective:
		taskDescription = subjectiveTaskDescription
	case reports.Objective:
		taskDescription = objectiveTaskDescription
	case reports.Assessment:
		taskDescription = assessmentTaskDescription
	case reports.Planning:
		taskDescription = planningTaskDescription
	case reports.Summary:
		taskDescription = summaryTaskDescription
	default:
		taskDescription = "Invalid SOAP section."
	}

	prompt := "You are an AI medical assistant responsible for generating precise and concise clinical notes based solely on the provided transcript. " +
    "Strict adherence to the transcript is required—do not include information that is not explicitly stated. " +
    "If critical details are missing, clearly indicate that additional context is necessary and why it is needed.\n\n" +
    
    "Do not infer treatment plans or SOAP elements unless they are explicitly present in the audio transcript.\n\n" +

    "Current Task (" + soapSection + "): " + taskDescription + "\n\n" +

    "Transcript:\n" + transcribedAudio + "\n\n"

	if style != "" {
		prompt += "Follow this style strictly:\n" + style + "\n\n"
	}

	prompt += defaultReturnFormat + "\n\n" + defaultWarnings

	return prompt
}


func RegenerateReportContentPrompt(previousContent string, soapSection, exampleStyle string, updates bson.D) string {
	var taskDescription string

	switch soapSection {
	case reports.Subjective:
		taskDescription = subjectiveTaskDescription
	case reports.Objective:
		taskDescription = objectiveTaskDescription
	case reports.Assessment:
		taskDescription = assessmentTaskDescription
	case reports.Planning:
		taskDescription = planningTaskDescription
	case reports.Summary:
		taskDescription = summaryTaskDescription
	default:
		taskDescription = "Invalid SOAP section."
	}

	prompt := "You are an AI medical assistant responsible for **fully rewriting** a clinical SOAP note section (Subjective, Objective, Assessment, Planning) " +
		"to ensure consistency between the provided metadata updates and the existing content. Strict adherence to the provided information is required—" +
		"do NOT infer, modify, or introduce any details that are not explicitly stated in the previous content. Your task is to apply only the given updates and restructure the content " +
		"if the metadata (e.g., pronouns, visit type, patient/client designation) conflicts with the original wording.\n\n" +

		"If the existing content is already aligned with the metadata updates, return the content as is. If the previous content is incoherent, incomplete, unclear, or if it requests additional context, " +
		"simply return: 'Additional context needed.' and specify why it is needed.\n\n" +

		"The required updates consist strictly of **metadata** such as:\n" +
		"- Patient pronouns (he/she/they)\n" +
		"- Visit type (initial visit or follow-up)\n" +
		"- Terminology adjustments (e.g., 'patient' vs. 'client')\n\n" +

		"These updates do NOT introduce new medical details, symptoms, diagnoses, or treatment plans. Your only modifications should ensure the content remains **coherent and aligned with the updated metadata** while preserving its original meaning.\n\n" +

		"Current SOAP Section: " + soapSection + "\n" +
		"Task Description: " + taskDescription + "\n\n" +

		"Previous Content:\n" + previousContent + "\n\n" +

		"Required Metadata Updates:\n" + formatUpdateDetails(updates) + "\n\n"

	if exampleStyle != "" {
		prompt += "Ensure the regenerated content closely matches this style (if no style provided, disregard this):\n" + exampleStyle + "\n\n"
	}

	prompt += defaultReturnFormat + "\n\n" + defaultWarnings

	return prompt
}

func formatUpdateDetails(updates bson.D) string{
	updateDetails := ""
	for _, update := range updates {
		if value, ok := update.Value.(string); ok && value != "" {
			updateDetails += "- " + update.Key + " updated to '" + value + "'\n"
		}
	}
	if updateDetails != "" {
		updateDetails = "\n\nPlease also incorporate the following updates:\n" + updateDetails
	}
	return updateDetails
}


const LearnStylePromptTemplate = `You are an AI medical assistant tasked with learning the writing style of an existing report.
Content Section: %s

The content of this section is as follows:
%s

Analyze the text above and extract the key stylistic elements, including tone, vocabulary, sentence structure, and formatting conventions.
Provide a summary of these stylistic elements that can be used to guide future content generation so that it matches the original style as closely as possible.`

// GenerateLearnStylePrompt constructs a prompt for the LearnStyle function.
func GenerateLearnStylePrompt(contentSection string, content string) string {
	return fmt.Sprintf(LearnStylePromptTemplate, contentSection, content)
}
