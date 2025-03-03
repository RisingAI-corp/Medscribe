package inferenceService

import (
	"Medscribe/reports"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const initialSystemPrompt = "You are an AI medical assistant that will be generating aspects of a report representing SOAP architecture: " +
	"Subjective, Objective, Assessment, and Planning. I will provide a transcript, and your task is to generate clear, " +
	"concise, and reliable outputs for the following fields: Subjective, Objective, Assessment, Planning, " +
	"IsFollowUp, IsReturning, Pronouns, and IsPatientOrClient. Accuracy and brevity are of supreme importance " +
	"in generating these components from the transcript."

const RegenerationInitialPrompt = "You are an AI medical assistant that will be generating aspects of a report representing SOAP architecture: " +
	"Subjective, Objective, Assessment, and Planning. I will provide a transcript, and your task is to generate clear, " +
	"concise, and reliable outputs for the following fields: Subjective, Objective, Assessment, Planning, " +
	"IsFollowUp, IsReturning, Pronouns, and IsPatientOrClient. Accuracy and brevity are of supreme importance " +
	"in generating these components from the transcript."

// GenerateReportContentPrompt generates a prompt for generating or regenerating a report section based on the SOAP architecture.
// It takes the transcribed audio, the specific SOAP section (Subjective, Objective, Assessment, Planning), the desired style,
// any updates to be incorporated, and the previously generated content (if any).
//
// Parameters:
// - transcribedAudio: The transcribed audio text to be used for generating the report content.
// - soapSection: The specific section of the SOAP report to be generated (e.g., Subjective, Objective, Assessment, Planning).
// - style: The desired writing style for the generated content.
// - updates: BSON document containing updates to be incorporated into the regenerated content.
// - content: The previously generated content for the specified SOAP section.
//
// Returns:
// - A string containing the generated prompt for the AI medical assistant to generate or regenerate the report content.
func GenerateReportContentPrompt(transcribedAudio, soapSection, style string, updates bson.D, content string) string {
	var taskDescription string

	switch soapSection {
	case reports.Subjective:
		taskDescription = "Extract patient-reported symptoms, concerns, and history relevant to the encounter."
	case reports.Objective:
		taskDescription = "Identify clinician-observed findings, including vitals, examination details, and test results."
	case reports.Assessment:
		taskDescription = "Summarize diagnosis or evaluation of the patient’s condition based on subjective and objective data."
	case reports.Planning:
		taskDescription = "Outline recommended treatment, follow-up instructions, and next steps for care."
	default:
		taskDescription = "Invalid SOAP section."
	}

	// Build update details if provided
	updateDetails := ""
	for _, update := range updates {
		if value, ok := update.Value.(string); ok && value != "" {
			updateDetails += "- " + update.Key + " updated to '" + value + "'\n"
		}
	}
	if updateDetails != "" {
		updateDetails = "\n\nPlease also incorporate the following updates:\n" + updateDetails
	}

	var prompt string

	// If both updates and previous content exist, generate a regeneration prompt
	if updateDetails != "" && content != "" {
		prompt = RegenerationInitialPrompt + "\n\n"
		if transcribedAudio != "" {
			prompt += "Here is the transcribed audio: \n" + transcribedAudio + "\n\n"
		} else {
			prompt += "No transcribed audio provided.\n\n"
		}
		prompt += "The following is the previously generated content for the " + soapSection + " section:\n" +
			content + "\n\n" +
			"Based on the updates provided below, please regenerate the content to reflect these changes:" +
			updateDetails + "\n\n" +
			"Task Details: " + taskDescription + "\n\n" +
			"Ensure that the regenerated content adheres to the SOAP framework and follows the given style: " + style + "\n\n" +
			"Ensure that the report is accurate, concise, and formatted in a structured way that is useful for medical documentation."
	} else {
		// Original prompt generation
		prompt = initialSystemPrompt + "\n\n"
		if transcribedAudio != "" {
			prompt += "Here is the transcribed audio: \n" + transcribedAudio + "\n\n"
		} else {
			prompt += "No transcribed audio provided.\n\n"
		}
		prompt += "Please generate the " + soapSection + " section of the report.\n" +
			"Task Details: " + taskDescription + "\n\n" +
			"Ensure that the generated content follows the SOAP framework and adheres to the given style: " + style + "\n\n" +
			"Ensure that the report is accurate, concise, and formatted in a structured way that is useful for medical documentation." +
			updateDetails
	}

	return prompt
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
