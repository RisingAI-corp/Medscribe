package inferenceService

import (
	"Medscribe/reports"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Soap Task Descriptions
subjectiveTaskDescription = `
    Extract and summarize the patient's reported information from the following transcript, ensuring a structured, comprehensive, and clinically relevant format. Act as an efficient AI medical scribe processing the transcript for a psychiatrist. The output should be plain text with clearly labeled sections and no markdown formatting.

    **Core Principle for Missing Information:**
    - **IMPORTANT:** If, after analyzing the transcript, no relevant information (either explicitly stated or reasonably inferred where specifically allowed, like in Social History) can be found for an entire section (e.g., Medications, History) OR a specific subsection heading within History (e.g., Surgical History, Family History), **COMPLETELY OMIT** that specific heading or subheading from the final output. Do not include headers for sections or subsections devoid of content based on the transcript.

    **Structured Sections:**

    **1. Chief Complaint (CC) and History of Present Illness (HPI):**

        **Chief Complaint (CC):**
        - Identify and state the primary reason(s) for the visit as reported by the patient.
        - Present this as a concise statement(s) or phrase(s), like a title for the visit (e.g., "Follow-up for depression," "Worsening anxiety," "Medication management"). Use the patient's own words if they provide a clear, brief reason.
        - If multiple complaints are offered, list the main ones briefly.
        (If no CC can be discerned, omit the "Chief Complaint (CC):" heading).

        **History of Present Illness (HPI):**
        - (If no HPI details beyond CC are available, this section might be omitted or very brief).
        - **Opening Statement:** Attempt to begin with a one-line summary including patient's age, sex/gender identity (if mentioned), and the chief complaint/reason for visit. (e.g., "Patient is a 35-year-old individual presenting for follow-up of anxiety."). If age/sex are not mentioned, start with the presentation reason and context.
        - **Elaborate on CC using OLDCARTS (Where Applicable):** Structure the narrative elaboration of the CC by addressing relevant components (Onset, Location, Duration, Characterization, Alleviating/Aggravating factors, Radiation, Temporal factors, Severity).
        - **Flexibility in Applying OLDCARTS:** Adapt elements based on the CC type (e.g., focus on Characterization, Severity, Duration, Aggravating/Alleviating for mood disorders; omit Location/Radiation).
        - **Contextual PMH:** Briefly weave in relevant Past Medical History only if it directly informs the HPI narrative and is mentioned in the transcript.
        - **Quote Usage for Clarity and Impact:** Use direct patient quotes strategically to illustrate key aspects (Characterization, Severity, patient experience). Prioritize clarity and impact over excessive quotation.
        - **Narrative Flow:** Maintain a logical narrative flow, integrating OLDCARTS elements into readable sentences.
        - **Focus on Clarity:** Ensure documentation is clear, concise, and captures essential details.
        - **Separate Interval Updates:** Keep HPI focused on the patient's story. Do not include routine provider checklist answers here.

    **2. Medications:**
        (Omit this entire section including the heading if no medications are discussed).
        - List all current medications mentioned, including name, dosage, and frequency if specified.
        - Document reported adherence details for each medication (e.g., "taking as prescribed," "reports missing 2 doses last week," "stopped taking [med name]").
        - Include any relevant past medications mentioned.
        - If the patient discusses obtaining medications from non-traditional sources (e.g., "street," friend), document this, including any stated reasons or impact.
        - Note the patient's stated reason for taking specific medications if mentioned (e.g., "takes [med name] for anxiety").

    **3. History:**
        (Omit this entire section including the heading if NO history details - Medical, Surgical, Family, or Social - are discussed).

        **Medical History (PMH):**
        (Omit this subheading if no relevant PMH is discussed).
        - Identify and list established, chronic, or significant past medical conditions mentioned.
        - Focus on diagnosed conditions stated or referenced.
        - Prioritize listing actual diagnoses over simple negative interval updates.
        - List conditions clearly and concisely.

        **Surgical History:**
        (Omit this subheading if no surgical history is discussed).
        - Identify and list any past surgical procedures mentioned.
        - Include type of surgery, year/date (if specified), and surgeon (if specified).
        - Present clearly (e.g., "Appendectomy (approx. 1998)").

        **Family History:**
        (Omit this subheading if no *pertinent* family history is discussed).
        - Identify mentions of pertinent medical conditions in family members.
        - Focus on hereditary conditions or those explicitly discussed as relevant to the patient. Avoid exhaustive lists.
        - Specify family member, condition, and relevant age details (if specified and relevant, e.g., "Father - MI at age 45").
        - Present clearly.

        **Social History:**
        (Omit this subheading if NO information across ANY HEEADSSS categories is found or inferred).
        - Analyze for information pertinent to the patient's social context using the HEEADSSS framework (Home, Education/Employment, Eating, Activities, Drugs/Substance Use, Sexuality, Suicide/Mood, Safety).
        - Include explicit statements and reasonably inferred details based on strong contextual clues (as allowed per original prompt).
        - **Structure Output by Category:** For EACH HEEADSSS category where relevant information *IS* found, include that specific category heading (e.g., "Home:", "Drugs/Substance Use:") followed by the concise summary for that category.
        - **Omit Empty HEEADSSS Categories:** If no information (explicit or inferred) is found for a *specific* HEEADSSS category (e.g., no mention related to 'Eating'), completely omit that specific category heading (do not write "Eating: Not discussed").
        - Include the HEEADSSS guidance details below for reference during processing:
            * Home: Living situation, relationships, safety, weapons.
            * Education/Employment: School/job status, performance, safety, plans, IEPs.
            * Eating: Habits, body image, weight goals/history, dieting behaviors, disordered eating signs.
            * Activities: Hobbies, fun, peers, clubs/teams, exercise, driving.
            * Drugs/Substance Use: Personal/peer use (tobacco, alcohol, other), frequency, context, consequences (CRAFFT), impaired driving exposure.
            * Sexuality: Dating, activity, orientation, partners, safety, coercion, STI/pregnancy history.
            * Suicide/Mood: Sadness, hopelessness, self-harm (thoughts/actions), running away (focus on baseline/historical context).
            * Safety: Integrate safety across categories; note explicit risks/violence.
`
	objectiveTaskDescription = `
	You are an AI medical scribe documenting objective clinical data from a patient encounter transcript for a psychiatrist. Write the report clearly, concisely, and in plain text without markdown, using the specific formatting guidelines below.

	**Core Objective Focus:**
	- Focus ONLY on objective, observable data (signs) in this section. Distinguish these from the patient's subjective reports (symptoms). (e.g., Subjective: "stomach pain"; Objective: "abdominal tenderness").
	- Base all observations and inferences DIRECTLY on the transcript content (what was said, how it was said, provider observations explicitly stated, or physical findings explicitly mentioned). Inference should be conservative and directly tied to clear transcript evidence.

	**Handling of Missing Information:**
	- For the main sections **Vital Signs, Physical Examination, Pain Scale, Diagnostic Test Results, Relevant Medication Details:** If no information is explicitly mentioned in the transcript for that category, explicitly state "[Category Name]: Not documented in the transcript."
	- For the **Mental Status Examination (MSE):**
		- If NO relevant observations can be made or confidently inferred across ALL MSE subcategories, omit the entire "Mental Status Examination (MSE):" heading.
		- For individual MSE subcategories (Behavior, Speech, Mood/Affect, etc.): If no relevant observation or confident inference can be made for a specific subcategory, **OMIT that specific subheading** (e.g., do not include "Insight:" if insight cannot be assessed).

	**Formatting:**
	- Use colons after main section headings (e.g., "Vital Signs:").
	- Within the Mental Status Examination section, use bullet points ("-" ) before each assessed category (e.g., "- Behavior: ...").

	**Structured Sections:**

	**(Omit heading entirely if no MSE observations possible)**
	**Mental Status Examination (MSE):**

		**(Omit bullet point/subheading if no relevant observation/inference)**
		- **Behavior:** Describe observable behavior, interaction style (e.g., cooperative, guarded), psychomotor activity (e.g., restless, calm), and demeanor based *directly* on transcript evidence. **Support with specific examples whenever possible** (e.g., "Cooperative with evaluation process," "Remained seated calmly throughout interview").

		**(Omit bullet point/subheading if speech cannot be assessed)**
		- **Speech:** Describe discernible characteristics (rate, volume, rhythm, tone, clarity) based on *audible qualities*. If speech appears within normal limits and no abnormalities are noted, state "Normal rate and fluency" or similar. **Support descriptions of abnormalities with specific examples or context** (e.g., "Speech occasionally slowed, seemed to search for words").

		**(Omit bullet point/subheading if affect/mood report cannot be discerned/inferred)**
		- **Mood/Affect:** Describe observed **Affect** (e.g., "Affect bright and reactive," "Affect constricted"). Include patient's self-reported **Mood** *using a direct quote* if provided clearly during the interaction (e.g., Patient stated mood was "'Not as bad as last time'"). **Support affect description with specific observations** (e.g., "Affect congruent with stated mood").

		**(Omit bullet point/subheading if thought process cannot be assessed)**
		- **Thought Process:** Assess organization/flow. Default to "Linear and goal-directed" if conversation is coherent. If deviations noted (tangential, circumstantial, etc.), describe clearly. **Support deviations with specific examples from the conversation** (e.g., "Thought process linear and goal-directed").

		**(Omit bullet point/subheading if cognition cannot be assessed)**
		- **Cognition:** Assess based *only* on interaction. Default to "Appears alert and oriented" if appropriate. Include *patient's reports* about their own cognition mentioned during the encounter (e.g., memory, attention) and any *observed* deficits. **Support with specific examples** (e.g., "Patient reported 'poor short-term memory' during recent testing, especially when anxious," "Appeared oriented to person, place, and situation").

		**(Omit bullet point/subheading if insight cannot be assessed)**
		- **Insight:** Assess patient's understanding of their situation/treatment based on statements. State level (e.g., Good, Fair, Limited). **Support assessment with specific examples, reasoning, or direct quotes demonstrating this understanding** (e.g., "Insight appears fair; patient recognizes need for testing to 'proceed forward properly'").

		**(Omit bullet point/subheading if judgment cannot be assessed)**
		- **Judgment:** Assess decision-making/planning based *only* on statements/reported actions in transcript. State level (e.g., Fair, Impaired). **Support assessment with specific examples of plans or decisions discussed** (e.g., "Judgment appears fair; patient discussed considering job change or further education as future options").


	**Vital Signs:** [If not mentioned, state "Not documented in the transcript."]
	- [List explicitly stated vitals here, e.g., BP 120/80, HR 70]

	**Physical Examination:** [If not mentioned, state "Not documented in the transcript."]
	- [List explicitly stated findings here, e.g., "Neuro exam: CN II-XII intact"]

	**Pain Scale:** [If not mentioned, state "Not documented in the transcript."]
	- [List explicitly stated pain rating here, e.g., "Patient reports anxiety 5/10"]

	**Diagnostic Test Results:** [If not mentioned, state "Not documented in the transcript."]
	- Include mention of tests ordered, ongoing, recently completed, OR specific results discussed. (e.g., "Cognitive testing ongoing, first session completed (approx. 3 hours duration)," "Recent TSH reported as 2.1").

	**Relevant Medication Details:** [If not mentioned, state "Not documented in the transcript."]
	- [Include any objective details mentioned related to medications, e.g., "Observed patient taking morning dose," "Injection site without redness." Note: Patient self-report belongs in Subjective].
	`


assessmentAndPlanTaskDescription = `
    You are an AI medical scribe synthesizing subjective and objective information from a patient encounter transcript to create a clinically reasoned Assessment and Plan (A/P) for a psychiatrist. The output must be plain text, well-organized, concise yet comprehensive, and use a problem-based structure. Avoid markdown formatting.

    **Structure:**
    - Organize the A/P by clinical problem (diagnosis, significant side effect, major psychosocial issue discussed). List problems addressed during the visit.
    - For each problem identified:
        1. State the **Problem** clearly as a heading.
        2. Underneath, provide a brief **Assessment** summarizing the relevant subjective and objective findings and patient status pertaining *only* to that problem. Use concise bullet points. This should synthesize information, not just repeat S/O data. Include status (e.g., improved, stable, worsening), contributing factors, or key complexities.
        3. Following the Assessment, detail the **Plan** specifically for that problem. Use concise bullet points outlining actions.

    **Content Guidance:**

    **Problem Identification:**
    - Identify all distinct clinical problems, diagnoses, or significant issues addressed during the encounter (e.g., Major Depressive Disorder, Anxiety Disorder, Medication Side Effect [specify, e.g., Weight Gain], Insomnia, Occupational Stress, Substance Use). List them, potentially in order of importance if discernible.

    **Assessment (per Problem):**
    - Synthesize key subjective reports (patient statements, symptom reports) and objective findings (MSE observations, test mentions, vital signs if relevant) that relate *directly* to this specific problem.
    - Briefly summarize the current status of the problem (e.g., "Symptoms improved since last visit," "Stable on current regimen," "Experiencing increased stress related to work").
    - Mention relevant contributing factors discussed (e.g., "Anxiety exacerbated by customer interactions," "Weight gain potentially related to [Medication Name]").

    **Plan (per Problem):**
    - Detail specific actions, monitoring, and follow-up related to this problem. Include:
        * **Medications:** Specify actions - continue, initiate, adjust dose, discontinue. List the *current medication regimen* concisely under the primary problem it treats (e.g., "Continue Zoloft 150 mg daily (2 x 75 mg tablets), Abilify 5 mg daily..."). Mention refills if addressed.
        * **Testing/Monitoring:** Include plans for diagnostic tests (e.g., "Await cognitive testing results," "Order TSH"), or monitoring parameters (e.g., "Monitor weight at next visit," "Continue monitoring mood symptoms").
        * **Referrals/Consultations:** List any referrals made or planned (e.g., "Referral to therapy," "Consult cardiology").
        * **Patient Education/Counseling:** Mention key education points or counseling provided (e.g., "Discussed sleep hygiene," "Counseled on potential side effects of [Medication Name]," "Discussed stress management techniques").
        * **Other Interventions/Actions:** Include any other steps (e.g., "Continue current exercise regimen," "Reassess work situation upon return from leave").

    **Overall Follow-up:**
    - Conclude the entire A/P section with a general statement about the planned follow-up interval (e.g., "Follow up in 4 weeks for medication management.").

    **Emphasis:**
    - Focus on clinical relevance. Identify and assess the problems that were actually addressed or impact the patient's current psychiatric status.
    - Ensure the Plan logically follows from the Assessment for each problem.
    - Be concise but include necessary specifics (like medication doses, specific tests, key patient concerns synthesized in assessment).
    - Use bullet points for assessment and plan details under each problem heading for clarity.
`

	summaryTaskDescription = `

	Write a concise, narrative-style visit summary based on the provided transcript. Clearly identify each clinical issue discussed, briefly summarizing the patient's reported symptoms, clinical findings, diagnostic outcomes, and relevant personal or psychosocial context. Include specific management plans such as medication adjustments (with dose and frequency if mentioned), follow-up appointments (with specific timing if given), referrals, or lifestyle recommendations. Maintain brevity and clarity throughout, avoiding redundancy. Omit any conditions or plans not explicitly discussed or clearly inferable from the transcript. Ensure the summary is direct, reads naturally, and includes clinically relevant details without unnecessary verbosity. `

	patientInstruction = `
    Generate a detailed and personalized patient instruction letter based on the patient
    encounter transcript provided below. Use a professional, empathetic, and clear tone. Include
    specific instructions, medication details, follow-up information, and any other relevant
    information discussed during the encounter. **Scan the entire transcript carefully to extract relevant details, even if mentioned conversationally.** Structure the letter with clear sections and
    bullet points ("-"), similar to the provided example. Ensure the output is in plain text, with no
    markdown formatting whatsoever. **The final output string MUST use %s placeholders for the patient name (in the greeting) and provider name/title (in the closing), which will be substituted later programmatically.**

    **Handling Missing Information & Sections:**
    - Include headings for **Medications, Tests/Procedures, Follow-Up, and General Advice** *only if relevant information pertaining to that category was clearly discussed* in the transcript. Summarize the information found under the appropriate heading.
    - Omit a heading completely if the topic was not addressed at all in the transcript.
    - If the transcript mentions important instructions/categories not covered by the standard sections (e.g., specific reporting requirements, non-medical advice strongly emphasized), create a relevant new section heading for those items.
    - **Placeholders:** Ensure the output uses **%s** for the patient's name in the greeting and for the provider's name/title in the closing. Do not attempt to fill these in.

    Instructions for Generating the Patient Instruction Letter:

    1. Begin with the exact greeting line: **"Dear %s,"**.
    2. Clearly summarize key instructions and recommendations discussed during the encounter, extracting specific details even from conversational parts.
    3. **Medications Section:** If medications were discussed, include this section heading ("Medications:"). Detail medication names, dosages **(as stated accurately in the transcript)**, frequency, specific instructions regarding refills (which meds need refill, which don't), and any relevant notes on usage (like Klonopin PRN), using bullet points.
    4. **Tests/Procedures Section:** If tests were discussed, include this section heading ("Tests/Procedures:"). Detail the status (e.g., ongoing, ordered, completed), any instructions for the patient (e.g., scheduling, sending reports), and the type of test if mentioned, using bullet points.
    5. **Follow-Up Section:** If follow-up was discussed, include this section heading ("Follow-Up:"). Specify the timeframe mentioned (e.g., "one month from today"). Include the method for scheduling if stated (e.g., "Please call the office to schedule"). Use bullet points if listing multiple follow-up items.
    6. **General Advice Section:** If general advice was given (e.g., related to lifestyle, stress management, coping, exercise), include this section heading ("General Advice:"). Summarize the key advice points discussed, using bullet points.
    7. Use clear section headings followed by a colon (e.g., "Medications:", "Follow-Up:") and bullet points (using hyphens "-") for lists within sections.
    8. Maintain a professional and empathetic tone throughout the content generated between the placeholders.
    9. Include a concluding sentence encouraging the patient to contact the office with any questions or concerns before the closing.
    10. End with the exact closing lines:
        **Best regards,**

        **%s**

    Example Structure (for AI reference, ensure final output matches this structure with %s):

    Dear %s,

    Thank you for visiting us today. We appreciate your commitment to maintaining your health and
    addressing your concerns promptly. Here is a summary of the key instructions from our
    consultation:

    Medications:
    - [Details extracted from transcript]
    - [Details extracted from transcript]

    Tests/Procedures:
    - [Details extracted from transcript]

    Follow-Up:
    - [Details extracted from transcript]

    General Advice:
    - [Details extracted from transcript]

    [New Section Name (if applicable)]:
    - [Details extracted from transcript]

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

If the input is "N/A" or provides no meaningful information or more information is required, respond only with: N/A. We don't want to give any indication that we need more information just put N/A

Otherwise, summarize the patient's past medical history, including major conditions, and list relevant medications with dosages. 
Keep the summary brief and to the point, no more than 15 words.

Example format: "Past medical history of [condition], currently managed with [medication1], [medication2], [medication3]."
`


const sessionSummary = `
You are an AI medical assistant tasked with generating a short description of the main topic or focus of a clinical session. 
The description should be concise and capture the key focus of the session in **a few words** (e.g., "Anxiety and medication discussion", "Follow-up on treatment plan").

If the input is "N/A" or provides no meaningful information or more information is required, respond only with: N/A. We don't want to give any indication that we need more information just put N/A

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
	"Patient Name: " + patientName + "\n" +
	"Provider Name: " + providerName + "\n\n" +
	"The following transcript is from a clinical session between the provider and the patient:\n\n" +
	"Current Task (" + soapSection + "): " + taskDescription + "\n\n" +
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
