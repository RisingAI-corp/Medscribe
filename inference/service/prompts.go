package inferenceService

import (
	"Medscribe/reports"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// TODO: we don't need the section style
// GenerateReportContentPrompt creates a prompt for generating a NEW section of a report.
func GenerateReportContentPrompt(transcribedAudio, soapSection, style, providerName, patientName string, context string) string {
	var taskDescription string

	switch soapSection {
	case reports.Subjective:
		taskDescription = subjectiveTaskDescription
	case reports.Objective:
		taskDescription = objectiveTaskDescription
	case reports.AssessmentAndPlan:
		taskDescription = assessmentAndPlanTaskDescription
	case reports.PatientInstructions:
		taskDescription = patientInstruction
	case reports.Summary:
		taskDescription = summaryTaskDescription
	default:
		taskDescription = "Invalid SOAP section."
	}

	// --- Prompt Construction ---
	prompt := fmt.Sprintf(`
You are an AI medical assistant acting as the provider. Your role is to document a specific section of the clinical visit report accurately and concisely based on the provided transcript and task description below.

--- BACKGROUND CONTEXT (FOR YOUR INFORMATION ONLY - DO NOT INCLUDE IN OUTPUT) ---
Patient Name: %s
Provider Name: %s
SOAP Section to Generate: %s
--- END BACKGROUND CONTEXT ---

--- TASK INSTRUCTIONS (Follow these instructions precisely to generate the required output) ---
`, patientName, providerName, soapSection)

	// --- Optional Context from Previous Visit ---
	if strings.TrimSpace(context) != "" && context != "N/A" && !strings.Contains(strings.ToLower(context), "additional context") {
		prompt += fmt.Sprintf(`**IMPORTANT CONTEXT FROM PREVIOUS VISIT:** %s

Use this information as an aid and reference when generating the %s section. If the content below is vague, non-clinical, or not relevant (e.g., "N/A", "additional context needed"), you must ignore it entirely and generate the note solely based on the transcript and task description below.

`, context, soapSection)
	}

	// --- Core Task Description ---
	prompt += fmt.Sprintf(`%s
--- END TASK INSTRUCTIONS ---

--- TRANSCRIPT (Analyze this transcript to perform the task) ---
%s
--- END TRANSCRIPT ---

GENERATE ONLY THE REQUIRED CLINICAL NOTE SECTION CONTENT for '%s' BASED ON THE TASK INSTRUCTIONS ABOVE.
***IMPORTANT: Your response MUST start directly with the narrative content for the requested section (%s). Do NOT include any section title or heading (like '%s:', 'Subjective:', 'Objective:', etc.) in your output. Your response should contain ONLY the narrative text.***

`, taskDescription, transcribedAudio, soapSection, soapSection, soapSection) // Added emphasis and explicit instruction against titles

	// --- Optional Style Integration ---
	if style != "" {
		prompt += "Integrate the style with the task description:\n" + style + "\n\n"
	}

	// --- Default Formatting and Warnings ---
	prompt += defaultReturnFormat + "\n\n" + defaultWarnings
	return prompt
}

// RegenerateReportContentPrompt creates a prompt for REWRITING an existing section based on metadata updates.
// MODIFIED to explicitly forbid section titles in the output.
func RegenerateReportContentPrompt(previousContent string, soapSection, exampleStyle string, updates bson.D, context string) string {
	var taskDescription string

	switch soapSection {
	case reports.Subjective:
		taskDescription = subjectiveTaskDescription
	case reports.Objective:
		taskDescription = objectiveTaskDescription
	case reports.AssessmentAndPlan:
		taskDescription = assessmentAndPlanTaskDescription
	case reports.Summary:
		taskDescription = summaryTaskDescription
	case reports.PatientInstructions:
		taskDescription = patientInstruction
	default:
		taskDescription = "Invalid SOAP section."
	}

	// --- Prompt Construction ---
	prompt := "You are an AI medical assistant acting as the provider, responsible for **fully rewriting** a clinical SOAP note section (Subjective, Objective, Assessment, Planning, Summary, or Patient Instructions) " +
		"to ensure consistency between the provided metadata updates and the existing content. Your task is to carefully apply only the specified updates while maintaining the accuracy and integrity of the original content. " +
		"Strict adherence to the provided information is required—do NOT infer, modify, or introduce any details beyond what is explicitly stated in the previous content.\n\n"

	// --- Optional Context from Previous Visit ---
	if strings.TrimSpace(context) != "" && context != "N/A" && !strings.Contains(strings.ToLower(context), "additional context") {
		prompt += fmt.Sprintf("**IMPORTANT CONTEXT FROM PREVIOUS VISIT:** %s\n\n", context)
		prompt += fmt.Sprintf("Use this as a reference when updating the %s section. If the content is vague or non-clinical (e.g., 'N/A', 'additional context needed'), ignore it completely and work only from the actual content and metadata updates provided.\n\n", soapSection)
	}

	// --- Handling Insufficient Content ---
	prompt += "If the existing content is already aligned with the metadata updates, return the content as is. If the previous content is incoherent, incomplete, unclear, or if additional context is required, " +
		"simply return: 'Additional context needed.' **only if the previous content itself is insufficient—not based on weak prior visit context**.\n\n"

	// --- Metadata Update Instructions ---
	prompt += "The required updates strictly involve **metadata** such as:\n" +
		"- Patient pronouns (he/she/they)\n" +
		"- Visit type (initial visit or follow-up)\n" +
		"- Terminology adjustments (e.g., 'patient' vs. 'client')\n\n"

	prompt += "Do NOT introduce new medical facts or diagnoses. Maintain coherence and accuracy while reflecting only the provided metadata updates.\n\n"

	// --- Core Task and Formatting Instructions ---
	prompt += "Current SOAP Section to Rewrite: " + soapSection + "\n" +
		"Task Description Hint: " + taskDescription + "\n\n" +
		"Previous Content:\n" + previousContent + "\n\n" +
		"Required Metadata Updates:\n" + formatUpdateDetails(updates) + "\n\n"

	// --- Optional Style Matching ---
	if exampleStyle != "" {
		prompt += "Ensure the regenerated content closely matches this style:\n" + exampleStyle + "\n\n"
	}

	// --- OUTPUT FORMATTING RULE ---
	prompt += fmt.Sprintf("***IMPORTANT: Your final rewritten output MUST start directly with the narrative content for the %s section. Do NOT include any section title or heading (like '%s:', 'Subjective:', 'Objective:', etc.) in your output. Return ONLY the rewritten narrative text based on the updates and previous content.***\n\n", soapSection, soapSection) // Added emphasis and explicit instruction against titles

	// --- Default Formatting and Warnings ---
	prompt += defaultReturnFormat + "\n\n" + defaultWarnings
	return prompt
}


func formatUpdateDetails(updates bson.D) string {
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

const LearnStylePromptTemplate = `You are an AI medical assistant tasked with analyzing and refining the writing style of a clinical report.

Content Section: %s

Previous version:
%s

Current version:
%s

Analyze the key stylistic differences between the previous and current versions. Identify specific changes in:
- Tone (e.g., clinical vs. narrative, patient quotes, emotional context)
- Structure (e.g., bullet points vs. full sentences, paragraph flow)
- Vocabulary (e.g., descriptive phrasing, use of time/location details)
- Formatting (e.g., labeled sections, sentence complexity)
- Level of detail (e.g., expanded context, inclusion of patient history)

Extract a set of **style recommendations** that will guide future content generation to match the **current** style as closely as possible. These recommendations should be clear, structured, and actionable so that other sections can be written consistently in the same manner.`

// GenerateLearnStylePrompt constructs a prompt for the LearnStyle function.
func GenerateLearnStylePrompt(contentSection string, previousContent string, currentContent string) string {
	return fmt.Sprintf(LearnStylePromptTemplate, contentSection, previousContent, currentContent)
}
