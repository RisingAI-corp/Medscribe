package inferenceService

import (
	"Medscribe/reports"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Soap Task Descriptions
	subjectiveTaskDescription = `
	Extract and summarize the patient's reported information from the following transcript, ensuring a structured, comprehensive, and clinically relevant format. This should be formatted into clearly labeled sections with full details and no markdown formatting in the output.
	
	General Guidelines:
	1. Begin with a free-flowing narrative summarizing the patient’s concerns, symptoms, and the visit purpose (HPI):
	   - Identify the type of visit: intake, follow-up, or transition of care.
	   - Provide context for medication use: If the reason for treatment is mentioned, explain why the patient is using certain medications based on their medical history or stated symptoms.
	   - Preserve direct quotes when the patient expresses significant emotions, distress, or relief (e.g., medication access, cravings, or symptom relief).
	   - If the patient discusses barriers to obtaining medications, document the exact reason and prior provider interactions. Do NOT paraphrase excessively—preserve details about insurance, debts, and provider changes.
	   - Use past tense for events that have already occurred and present tense for current symptoms or ongoing issues.
	
	2. Structured Format:
	   - Do NOT title the first paragraph—allow it to flow naturally.
	   - Use clearly labeled sections for medical history, medications, and social history.
	   - Do NOT summarize financial stress under Social History unless it is a broad, ongoing issue. If it is strictly related to medication access, keep it within the HPI or medication section.
	
	3. Medication Documentation:
	   - Always list current and past medications with adherence details.
	   - If the patient obtained medication from non-traditional sources (e.g., the street), document why and how it affected their health or stress levels.
	   - Ensure the patient’s reasoning for medication use is explicitly stated.
	
	4. Social History:
	   - Include only explicitly mentioned social factors. Do NOT move situational financial struggles into Social History unless they are a recurring or broad issue.
	
	Missing Information Handling:
	- If specific details such as the patient's name, gender, adherence details, or medical history are not explicitly mentioned in the transcript, **do not request additional context. Simply omit those details and proceed with the available information.**
	`

	objectiveTaskDescription = `
You are a medical provider documenting objective clinical data from a patient encounter transcript. Write the report clearly, concisely, and in plain text without markdown, incorporating specific details or examples directly from the patient's transcript to accurately reflect unique characteristics of the encounter. Completely omit any category that is neither explicitly mentioned nor confidently inferable.

(REQUIREMENTS: its critical here that if you can not clearly infer aspects of the mental status examination then don't include it. It doesn't serve a purpose to explicitly state Not mentioned. If not mentioned -> infer. if can not be inferred then omit the filed  )

Mental Status Examination:
- Behavior: Briefly describe patient's observable behavior, interaction style, and demeanor, integrating distinct and specific details or examples from the conversation if available (e.g., "Calm and cooperative, openly shared concerns about medication side effects," or "Withdrawn, minimal responses, appeared distracted during questioning"). Avoid overly generic phrases unless no distinctive details are provided.
- Speech: If explicitly mentioned or confidently inferred, briefly describe speech characteristics including pacing, tone, or clarity (e.g., "Speech rapid when discussing work stress," or "Quiet and hesitant when mentioning family conflicts"). Omit entirely if not clearly observed or documented.
- Mood: Use patient's direct quotes if provided, or succinctly describe mood based on patient's expressed emotions or tone (e.g., "Expressed feeling overwhelmed by family obligations," "Mood described as stable with improvements noted since last visit"). Avoid overly generic descriptions.
- Thought Process: Default to "Linear and goal-directed" if coherent; if not linear, briefly specify and describe clearly observed deviations using specific conversation examples (e.g., "Occasionally tangential, patient frequently shifted topics when discussing future plans," "Patient’s responses were focused but included excessive irrelevant details").
- Cognition: Succinctly state "Alert and oriented to conversation" unless explicit cognitive concerns (e.g., confusion, memory issues) are observed; provide a concise description and relevant examples if any cognitive issues are noted (e.g., "Mild confusion, difficulty recalling recent medication adjustments").
- Insight: Include if patient demonstrates clear awareness or understanding of their condition or treatment; support with specific examples or statements from the patient if possible (e.g., "Good insight demonstrated by proactive discussion of medication management," "Limited insight into the severity of reported anxiety symptoms").
- Judgment: Include if patient shows decision-making ability or planning that can be explicitly or implicitly inferred from the conversation; use examples from patient interactions if available (e.g., "Judgment appears fair; patient actively schedules follow-ups and adheres to medication despite reported side effects").

(some requirements for below. If the patient is not explicitly mentioned, omit entirely, but if mentioned and is still inclusive still specify it)

Vital Signs:
- Include explicitly stated vital signs (e.g., blood pressure, heart rate). Omit entirely if not explicitly stated.

Physical Examination:
- Include explicitly stated physical findings, briefly summarized by system with relevant details (e.g., "Tenderness noted in left ankle," "Clear lungs on examination"). Omit entirely if no physical exam mentioned.

Pain Scale:
- Clearly document numeric pain rating explicitly reported by the patient, with date and reference to scale (e.g., "Pain rated 8/10, described as severe and constant"). Omit entirely if not explicitly provided.

Diagnostic Test Results:
- Concisely summarize explicitly mentioned diagnostic tests or results (e.g., blood tests, imaging findings). Omit entirely if no test results mentioned.

Additional Relevant Information:
- Briefly document explicitly stated patient details directly relevant to clinical care decisions, medication management, or follow-up arrangements that have not been mentioned elsewhere. Avoid redundancy.

Provide all information concisely and specifically in plain text format without markdown, emphasizing distinct, transcript-specific details to avoid repetitive or overly generic documentation.
`

	assessmentAndPlanTaskDescription = `
	This section synthesizes subjective and objective evidence from the transcript to document concisely the patient's clinical issues in order of importance. Clearly differentiate patient-reported symptoms (subjective) from clinical observations or objective findings. Title each issue simply by its clinical name, even if inferred from medication or clearly described symptoms without explicitly stating it's inferred.

	Important guidelines:
	- If a diagnosis or condition is explicitly discussed, document it clearly. If medications, symptoms, or treatment strongly imply a diagnosis (e.g., Adderall implying ADHD, lorazepam implying anxiety), document the diagnosis clearly without stating it's inferred.
	- Each condition must directly link to medications or treatment plans explicitly mentioned in the transcript.
	- If a condition or its management is not explicitly mentioned or clearly implied by medications or symptoms, omit it entirely.
	- make sure to actually think before putting a medication down because some providers may have different accents so approximate which medications are most likely the provider and patient are referring to. Make sure to take fully into account the the context of the conversation in which it is used in.
	
	Structure:
	[Condition Name]
	- Briefly summarize current patient status, including symptoms, medication efficacy, adherence, side effects, or explicitly stated patient concerns. Use specific examples or context from the transcript when relevant. Clearly distinguish subjective patient reports from objective clinical observations.
	- Plan: Clearly document specific next steps as provided, including medication details (exact dosages, frequency), referrals, recommended tests, monitoring, lifestyle advice, or patient instructions. Include follow-up details explicitly stated or confidently infer a standard clinical interval if not specified (e.g., "Follow-up in 1 month").
	
	Follow-Up Appointments:
	- Clearly indicate scheduled follow-ups as explicitly stated. If the follow-up timing is implied without an explicit date, state a standard clinical interval clearly (e.g., "Follow-up in 1 month (standard interval)").
	
	Additional Information:
	- Briefly include explicitly stated details not captured elsewhere that directly affect patient care (pending diagnostic tests, significant life events, upcoming travel, or environmental factors impacting health or treatment).
	
	Summary:
	- Briefly reiterate the key clinical priorities and immediate actionable steps explicitly discussed to ensure clarity and patient understanding.
	
	Provide the response in plain text without markdown formatting.
	`

	summaryTaskDescription = `

	Write a concise, narrative-style visit summary based on the provided transcript. Clearly identify each clinical issue discussed, briefly summarizing the patient's reported symptoms, clinical findings, diagnostic outcomes, and relevant personal or psychosocial context. Include specific management plans such as medication adjustments (with dose and frequency if mentioned), follow-up appointments (with specific timing if given), referrals, or lifestyle recommendations. Maintain brevity and clarity throughout, avoiding redundancy. Omit any conditions or plans not explicitly discussed or clearly inferable from the transcript. Ensure the summary is direct, reads naturally, and includes clinically relevant details without unnecessary verbosity. `

	patientInstruction = `
	Generate a detailed and personalized patient instruction letter based on the patient
	encounter transcript provided below. Use a professional, empathetic, and clear tone. Include
	specific instructions, medication details, follow-up information, and any other relevant
	information discussed during the encounter. Structure the letter with clear sections and
	bullet points, similar to the provided example. Ensure the output is in plain text, with no
	markdown formatting whatsoever. If the transcript mentions important categories not covered by
	Medications, Tests/Procedures, Follow-Up, or General Advice, create a new section for those
	categories. If a category is not explicitly mentioned in the transcript, omit it entirely. It doesn't look aesthetic if you include a section followed by lack of information
	
	Instructions for Generating the Patient Instruction Letter:
	
	1. Begin with a warm greeting, thanking the patient for their visit and acknowledging their
	   commitment to their health.
	2. Clearly summarize key instructions and recommendations discussed during the encounter,
	   including specific details.
	3. Include specific details about medications, such as names, dosages, frequency, and
	   instructions for use, in a dedicated Medications section.
	4. Provide information about any tests or procedures ordered, including where and when they
	   should be performed, in a dedicated Tests/Procedures section.
	5. Specify any follow-up appointments explicitly mentioned in the transcript, including exact
	   dates, times, and locations.
	   - If the provider implies or mentions a follow-up without giving a specific timeline,
		 confidently infer a standard follow-up interval based on typical clinical practice or
		 past appointments (e.g., "Follow-up in one month (standard interval)"). Clearly state
		 that this timeline is inferred in the Follow-Up section.
	6. Offer general advice and recommendations for managing symptoms or improving health, in a
	   dedicated General Advice section.
	7. Use clear section headings and bullet points (using hyphens "-") to organize information for
	   clarity, similar to the provided example letter.
	8. Maintain a professional and empathetic tone throughout the letter.
	9. Encourage the patient to contact the office with any questions or concerns.
	10. End with a warm closing, such as "Best regards" or "Warm regards," followed by your
		professional title.
	
	Example Structure to Follow:
	
	Dear %s,
	
	Thank you for visiting us today. We appreciate your commitment to maintaining your health and
	addressing your concerns promptly. Here is a summary of the key instructions from our
	consultation:
	
	Medications:
	- Medication 1: Dosage, Frequency, Instructions
	- Medication 2: Dosage, Frequency, Instructions
	
	Tests/Procedures:
	- Test/Procedure Name: Location, Date/Time, Instructions
	
	Follow-Up:
	- Date, Time, Location, Purpose (or clearly inferred timeline if not explicitly stated)
	
	General Advice:
	- Advice 1
	- Advice 2
	
	[New Section Name (if applicable)]:
	- Information 1
	- Information 2
	
	Please ensure to follow these instructions carefully. We are here to support you in your health
	journey. Do not hesitate to reach out if you have any questions or concerns.
	
	Best regards,
	
	%s
	  `

	//Default Return format
	defaultReturnFormat = "Return Format (Always Adhere to These Rules):" +
		" Responses must be strictly plain text, suitable for direct display in a textbox." +
		" Never use markdown formatting (no \"*\", \"-\", \"#\", or \"---\")." +
		" Do not return a response as if you were responding based off of a question." +
		" All your queries will be aggregated and decorated from the client. Do not give any indication you were prompted iteratively." +
		" If this task description cannot be answered or there is a gross amount of information missing to answer the task description, just remove it."
	//Default Warnings
	defaultWarnings = "Warnings:" +
		" output needs to be plain text absolutely no markdown" +
		" Always rely exclusively on the provided transcript without assumptions or inference beyond clearly available context." +
		"if the transcript is inconclusive meaning no relevant dat is enough to even generate the report then just reply back with N/A"
)

const condensedSummary = `
You are an AI medical assistant tasked with generating a concise medical history based on the patient's report. 
The patient's detailed summary is as follows:

%s

If the summary is "N/A" or contains no meaningful medical information, respond only with: N/A

Otherwise, summarize the patient's past medical history, including major conditions, and list relevant medications with dosages. 
Keep the summary brief and to the point, no more than 15 words.

Example format: "Past medical history of [condition], currently managed with [medication1], [medication2], [medication3]."
`


const sessionSummary = `
You are an AI medical assistant tasked with generating a short description of the main topic or focus of a clinical session. 
The description should be concise and capture the key focus of the session in **a few words** (e.g., "Anxiety and medication discussion", "Follow-up on treatment plan").

If the input is "N/A" or provides no meaningful information, respond only with: N/A

otherwise proceed with the following:
The summary should:
- Be no more than **a few words**.
- Focus on the **primary subject of the session** (e.g., a condition, treatment, or discussion).
- Avoid unnecessary details or lengthy explanations.

Example input:  
The patient discussed her ongoing issues with anxiety, the impact of recent medication changes, and her challenges with managing side effects.  

Example output:  
Anxiety and medication discussion

%s
`


// TODO: we don't need the section style
func GenerateReportContentPrompt(transcribedAudio, soapSection, style, providerName, patientName string) string {
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

	prompt := "You are an AI medical assistant acting as the provider, documenting the visit report after reviewing the transcript of the encounter. " +
		"Your task is to generate precise and accurate clinical notes that reflect exactly what was stated in the transcript, as if you were the provider writing them up. " +
		"Strict adherence to the transcript is required—do not add, infer, or assume any information beyond what is explicitly stated. " +
		"If critical details are missing, clearly indicate that additional context is necessary and explain why it is needed.\n\n" +

		"Current Task (" + soapSection + "): " + taskDescription + "\n\n" +

		"Patient Name: " + patientName + "\n" +
		"Provider Name: " + providerName + "\n\n" +

		"Transcript:\n" + transcribedAudio + "\n\n"

	if style != "" {
		prompt += "Integrate the style with the task description:\n" + style + "\n\n"
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
	case reports.AssessmentAndPlan:
		taskDescription = assessmentAndPlanTaskDescription
	case reports.Summary:
		taskDescription = summaryTaskDescription
	case reports.PatientInstructions:
		taskDescription = patientInstruction
	default:
		taskDescription = "Invalid SOAP section."
	}

	prompt := "You are an AI medical assistant acting as the provider, responsible for **fully rewriting** a clinical SOAP note section (Subjective, Objective, Assessment, Planning) " +
		"to ensure consistency between the provided metadata updates and the existing content. Your task is to carefully apply only the specified updates while maintaining the accuracy and integrity of the original content. " +
		"Strict adherence to the provided information is required—do NOT infer, modify, or introduce any details beyond what is explicitly stated in the previous content. " +
		"Your role is to ensure that the documentation remains precise, professional, and aligned with the updated metadata.\n\n" +

		"If the existing content is already aligned with the metadata updates, return the content as is. If the previous content is incoherent, incomplete, unclear, or if additional context is required, " +
		"simply return: 'Additional context needed.' and specify why it is needed.\n\n" +

		"The required updates strictly involve **metadata** such as:\n" +
		"- Patient pronouns (he/she/they)\n" +
		"- Visit type (initial visit or follow-up)\n" +
		"- Terminology adjustments (e.g., 'patient' vs. 'client')\n\n" +

		"These updates do NOT introduce new medical details, symptoms, diagnoses, or treatment plans. Your modifications should solely ensure the content remains **coherent, accurate, and aligned with the updated metadata** while preserving its original medical meaning.\n\n" +

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
