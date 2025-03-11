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
	planningTaskDescription = "Outline recommended treatment, follow-up instructions, and next steps for care."
	summaryTaskDescription = "Summarize the key points of the patient's history, symptoms, and examination findings."

	//Default Return format
	defaultReturnFormat = "Return Format (Always Adhere to These Rules):" + 
	" Responses must be strictly plain text, suitable for direct display in a textbox." + 
	" Never use markdown formatting (no \"*\", \"-\", \"#\", or \"---\")." + 
	" Do not include headers, special characters, or extraneous punctuation." + 
	" Ensure all responses are brief, precise, and directly related to the requested component."

	//Default Warnings
	defaultWarnings = "Warnings:" + 
           " Never provide information or details outside of what's explicitly requested." +
           " Always rely exclusively on the provided transcript without assumptions or inference beyond clearly available context." +
		   " if more context is needed do not hallucinate and just specify that more context is needed"
)


func GenerateReportContentPrompt(transcribedAudio, soapSection, style string, updates bson.D, content string) string {
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

	prompt := "You are an AI medical assistant tasked with generating clinical notes strictly based on provided transcripts. " +
		"Accuracy and brevity are paramount. Never include information not explicitly found in the provided transcript. " +
		"If you lack sufficient information, clearly indicate that more context is required.\n\n" +

		"Your current task (" + soapSection + "): " + taskDescription + "\n\n" +

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
	
	prompt := "You are an AI medical assistant tasked with regenerating a clinical SOAP note section (Subjective, Objective, Assessment, Planning). " +
		"Your task is to clearly, accurately, and concisely regenerate the provided content, strictly incorporating the provided updates. " +
		"Do NOT infer or add information not explicitly provided. If details are unclear or missing, state explicitly that more clarification or context is needed.\n\n" +

		"Current SOAP Section: " + soapSection + "\n" +
		"Task Description: " + taskDescription + "\n\n" +

		"Previous Content:\n" + previousContent + "\n\n" +

		"Updates to Incorporate:\n" + formatUpdateDetails(updates) + "\n\n"

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
