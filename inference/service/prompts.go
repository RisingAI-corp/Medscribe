package inferenceService

import (
	transcriber "Medscribe/Transcription"
	"Medscribe/reports"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)
type generatePromptConfig struct {
	transcript         string
	diarizedTranscript transcriber.TranscriptTurn
	useDiarization     bool
	targetSection      string
	context            string
	style              string
	providerName       string
	patientName        string
}

type regeneratePromptConfig struct {
	transcript         string
	diarizedTranscript transcriber.TranscriptTurn
	useDiarization     bool
	targetSection      string
	targetContent      string
	priorVisitContext  string
	providerName       string
	patientName        string
	reportUpdates      bson.D
}

// formatUpdateDetails formats the metadata updates for inclusion in the prompt.
func formatUpdateDetails(updates bson.D) string {
	var details strings.Builder
	details.WriteString("Required Metadata Updates:\n")
	for _, update := range updates {
		if value, ok := update.Value.(string); ok && value != "" {
			details.WriteString(fmt.Sprintf("- %s: %v\n", update.Key, value))
		}
	}
	if details.String() != "Required Metadata Updates:\n" {
		return "\n\nPlease also incorporate the following updates:\n" + details.String()
	}
	return ""
}

// GenerateReportContentPrompt creates a prompt for generating a NEW section of a report.
func GenerateReportContentPrompt(cfg generatePromptConfig) string {
    var taskDescription string

    switch cfg.targetSection {
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
You are an AI medical assistant acting as the provider. Your paramount responsibility is to generate each section of the clinical visit report with the utmost correctness, adhering meticulously to the Task Instructions provided below. This report may be presented as evidence in a legal setting.

--- BACKGROUND CONTEXT (FOR YOUR INFORMATION ONLY - DO NOT INCLUDE IN OUTPUT) ---
Patient Name: %s
Provider Name: %s
SOAP Section to Generate: %s
--- END BACKGROUND CONTEXT ---

--- TASK INSTRUCTIONS (Follow these instructions precisely and prioritize correctness) ---
%s

**CRITICAL OPERATING PRINCIPLES FOR REPORT GENERATION (LEGAL AND PROFESSIONAL STANDARD):**

* **Prioritize Correctness Above All:** Your primary directive is to ensure the absolute correctness of the generated content for each section. This requires a thorough understanding and strict adherence to the Task Instructions.
* **Diligent Review and Self-Correction:** Before generating the content for any section, carefully review the Task Instructions and mentally simulate the extraction and formatting process. Critically evaluate potential ambiguities or challenging parts of the transcript against the guidelines. If any part of your planned output deviates from the instructions, self-correct before generating the final text.
* **Meticulous Adherence to Guidelines:** Every aspect of the generated report, from the content to the formatting (headings, bullet points, etc.), MUST strictly comply with the explicit guidelines provided in the Task Instructions.
* **Legal Standard Reinforcement:** Remember that this report may be used in legal proceedings. Therefore, accuracy, clarity, and direct support from the transcript are non-negotiable. Errors or deviations from the Task Instructions could have serious consequences.
* **Contextual Accuracy:** Ensure that your interpretation of the transcript aligns with the overall context of a telehealth visit and standard medical practices, while always grounding your output in the specific details of the provided conversation.

**IMPORTANT ADDITIONAL INSTRUCTIONS FOR TRANSCRIPT ANALYSIS (incorporating legal and professional standards):**

* **Accent and Pronunciation Variations:** Use contextual understanding to interpret intended meaning, documenting standard medical terms when the intent is clear.
* **Anomaly Correction (Justified):** Correct clear transcription errors (e.g., "minute" to "mg") ONLY when the context overwhelmingly supports the correction. If ambiguity exists, document the transcription and note the uncertainty with your reasoning.
* **Contextual Interpretation (Evidence-Based):** Base interpretations firmly on the surrounding conversation. Avoid assumptions not directly supported by the transcript.
* **Completeness and Relevance (Legal Scrutiny):** Include all relevant information as per the Task Instructions.

--- END TASK INSTRUCTIONS ---

--- TRANSCRIPT (Analyze this transcript to perform the task with a focus on correctness) ---
%s
--- END TRANSCRIPT ---

GENERATE ONLY THE REQUIRED CLINICAL NOTE SECTION CONTENT for '%s' BASED ON THE TASK INSTRUCTIONS ABOVE, ensuring the highest degree of correctness through careful review and adherence to all guidelines.
***IMPORTANT: Your response MUST start directly with the narrative content for the requested section (%s). Do NOT include any section title or heading (like '%s:', 'Subjective:', 'Objective:', etc.) in your output. Your response should contain ONLY the narrative text.***

`, cfg.patientName, cfg.providerName, cfg.targetSection, taskDescription, cfg.transcript, cfg.targetSection, cfg.targetSection, cfg.targetSection)

    // --- Optional Style Integration ---
    if cfg.style != "" {
        prompt += "Integrate the style with the task description:\n" + cfg.style + "\n\n"
    }

    // --- Default Formatting and Warnings ---
    prompt += defaultReturnFormat + "\n\n" + defaultWarnings
    return prompt
}

// RegenerateReportContentPrompt creates a prompt for REWRITING an existing section based on metadata updates.
// MODIFIED to explicitly forbid section titles in the output.
func RegenerateReportContentPrompt(cfg regeneratePromptConfig) string {
	var taskDescription string

	switch cfg.targetSection {
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
	if strings.TrimSpace(cfg.priorVisitContext) != "" && cfg.priorVisitContext != "N/A" && !strings.Contains(strings.ToLower(cfg.priorVisitContext), "additional context") {
		prompt += fmt.Sprintf("**IMPORTANT CONTEXT FROM PREVIOUS VISIT:** %s\n\n", cfg.priorVisitContext)
		prompt += fmt.Sprintf("Use this as a reference when updating the %s section. If the content is vague or non-clinical (e.g., 'N/A', 'additional context needed'), ignore it completely and work only from the actual content and metadata updates provided.\n\n", cfg.targetSection)
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
	prompt += "Current SOAP Section to Rewrite: " + cfg.targetSection + "\n" +
		"Task Description Hint: " + taskDescription + "\n\n" +
		"Previous Content:\n" + cfg.targetContent + "\n\n" +
		"Required Metadata Updates:\n" + formatUpdateDetails(cfg.reportUpdates) + "\n\n"


	// --- OUTPUT FORMATTING RULE ---
	prompt += fmt.Sprintf("***IMPORTANT: Your final rewritten output MUST start directly with the narrative content for the %s section. Do NOT include any section title or heading (like '%s:', 'Subjective:', 'Objective:', etc.) in your output. Return ONLY the rewritten narrative text based on the updates and previous content.***\n\n", cfg.targetSection, cfg.targetSection) // Added emphasis and explicit instruction against titles

	// --- Default Formatting and Warnings ---
	prompt += defaultReturnFormat + "\n\n" + defaultWarnings
	return prompt
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
func GenerateLearnStylePrompt(targetSection, previous, current string) string {
	return fmt.Sprintf(LearnStylePromptTemplate, targetSection, previous, current)
}

