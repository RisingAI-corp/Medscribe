package inferenceService

import (
	"Medscribe/reports"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Soap Task Descriptions
	subjectiveTaskDescription = `
    // ROLE: Act as a meticulous clinical scribe for psychiatric documentation.
    // GOAL: Extract and meticulously document patient's reported info from a transcript into a COMPREHENSIVE, DETAILED, CLINICALLY RELEVANT subjective note suitable for production use. Capture patient's experience, context, perspective accurately. Prioritize specifics, verbatim details, quotes. AVOID summaries; capture full context. A primary goal is the complete listing of ALL mentioned current AND recently discontinued medications.
    // STYLE: Professional clinical language. Plain text output. Third-person POV. Descriptive, quote-inclusive, transcript-grounded. Adapt phrasing naturally while meeting core requirements.

    OUTPUT FORMATTING AND STRUCTURE REQUIREMENTS:
    - PLAIN TEXT ONLY: No markdown.
    - NO 'SUBJECTIVE:' HEADING: Start output DIRECTLY with the Initial Narrative HPI Section.
    - SECTION HEADINGS: Use clear headings + colon (e.g., Medical History:, Medications and Supplements:, Social History:, Family History:) AFTER the initial narrative. Omit headings entirely ONLY if no relevant info exists per criteria.
    - BULLET POINTS: Use '- ' for primary lists (History, Meds, Social, Family). Use '  - ' ONLY for indented medication metadata.
    - ***CRITICAL OMISSION RULE***: NO placeholders (N/A, etc.). If info for a heading/bullet/metadata point is absent per criteria, omit it entirely and silently.

    CONTENT GENERATION INSTRUCTIONS:

    1. Initial Narrative HPI Section:
       - POINT OF VIEW: Write strictly in the third person ('Patient reports...', 'They describe...'). Use patient identifiers if available from context.
       - START DIRECTLY WITH NARRATIVE: Begin note directly with this narrative. May use multiple paragraphs.
       - EXCLUDE OPENING PLEASANTRIES: Start narrative with first clinically relevant information reported.
       - AIM FOR COHESIVE SYNTHESIS & FLOW: Structure as a flowing clinical narrative telling the patient's integrated story. Actively synthesize and link relevant context: connect symptoms to triggers/context, discuss events, assessments, functional impact, and medication context. Tell the patient's connected story. Goal is synthesis, not just listing facts.
       - IDENTIFY VISIT TYPE/CONTEXT: State visit purpose (e.g., follow-up) and key context early.
       - CC HANDLING: Integrate reason(s) into opening narrative for follow-ups (no separate "CC:" line). Optionally use "CC:" line for intakes if clearly stated early.
       - COMPREHENSIVE NARRATIVE - MANDATORY DEPTH & CONCISENESS: Ensure comprehensive coverage of points below (if present), presented concisely. MUST incorporate details such as:
           - Mood/anxiety state (ratings, changes, key quotes).
           - Sleep patterns reported (hours, quality).
           - Verbatim reason/quote for significant actions (e.g., reason for leave).
           - Detailed symptom descriptions using patient's phrasing (e.g., anxiety attack details - timing/severity; concentration issues) AND their specific functional impact (on school, work, hobbies, daily tasks - provide examples).
           - Relevant diagnostic testing discussion (prior results/experience, current status, plans).
           - Patient's perspective on treatment (effectiveness, uncertainty, reluctance, side effect concerns, queries about dose changes).
           - Specific adherence details mentioned conversationally ('not yet started', 'ran out', inconsistent use).
       - INCORPORATE OLDCARTS ELEMENTS NATURALLY: Weave relevant concepts (Characterization, Severity, Timing, Factors, Insight) naturally into narrative. Do NOT list rigidly.
       - CONTEXTUALIZE AND QUOTE FREQUENTLY: Explain the WHY and illustrate with direct quotes.
       - BARRIERS/CHALLENGES: Detail specific challenges (e.g., therapy reluctance, psychosocial stressors).

    2. Medical History Section:
       (Omit heading ONLY if NO relevant medical/psychiatric factors WHATSOEVER identified. Adhere to CRITICAL OMISSION RULE.)
       - CONTENT: Use '- ' to list relevant past diagnoses AND significant ongoing medical/psychiatric conditions requiring management or providing context. **Use specific standard diagnostic terms (e.g., 'PTSD', 'Anxiety disorder with panic attacks', 'Hypertension', 'Obesity', 'Menopause') IF clearly supported by context/inference.** Otherwise describe symptom clusters (e.g., 'ADHD Symptoms under evaluation'). Exclude purely psychosocial factors. Use specific transcript phrasing where appropriate.
       - INFERENCE REQUIRED: SHOULD infer ongoing conditions from strong context (symptoms, long-term meds). List explicit past diagnoses. Include relevant recent acute issues.

    3. Medications and Supplements Section:
       (Omit heading if no medications discussed. Adhere to CRITICAL OMISSION RULE.)
       - **LISTING ACCURACY & COMPLETENESS (CRITICAL):** Use '- ' to list **ALL current medications mentioned in the transcript** (psychiatric AND non-psychiatric like BP meds) AND **recently discontinued relevant meds** (like Zepbound). **This is mandatory - ensure nothing mentioned is missed.** Include Name, Dose, Frequency, Route (if mentioned).
       - CRITICAL - ALL METADATA: For EACH relevant med listed, MUST ADD notes using indented ('  - ') bullets detailing ALL relevant patient observations/comments mentioned. Extract DIRECTLY/VERBATIM. MUST cover IF MENTIONED:
           - Purpose (Reason prescribed ONLY IF STATED - **CRITICAL: Do NOT infer purpose.** Omit if not stated).
           - Reported Effectiveness (Patient assessment, quotes, uncertainty, e.g., 'Some improvement...', 'Not fully effective yet').
           - Side Effects (Specifics mentioned like 'Was making the patient sick', weight gain amount).
           - Adherence/Usage (Patterns, events like 'not yet started', current frequency 'taken at night').
           - Supply/Refill Status ('has full bottle', 'needs refill').
           - Status ('Discontinued', 'Not yet started', dose changed). Ensure status is captured (e.g., for relevant meds like Zepbound, Magnesium).
           - Labs (Relevant levels/context).
           - Regimen Details (Specific schedules, planned usage like 'Start 200mg at night'). Ensure accuracy.
       - OMIT SILENTLY: Omit missing details per CRITICAL OMISSION RULE.

    4. Social History Section:
       (Omit heading ONLY if NO pertinent social factors discussed. Adhere to CRITICAL OMISSION RULE.)
       - CONTENT: Use '- ' to list the **most pertinent** social factors discussed, focusing on clinical relevance. Aim for **concise bullet points summarizing key information** based on transcript content. Consider domains below, but only include bullets for domains with significant findings. Extract EXACT PHRASING/SPECIFICS and integrate relevant quotes (' ').
       - **Education/Employment:** Report status, changes, specific academic/work challenges (using quotes), accommodation needs, role details/context, work stressors (using quotes), employment status (leave details, contract issues).
       - **Activities/Hobbies/Exercise:** Report specific leisure activities, exercise habits/frequency (using quotes like 'very on and off'), equipment details, AND note any **stated impact of symptoms on these activities**. Include brief mentions if relevant.
       - **Home/Environment:** Report living situation. Note safety concerns if expressed.
       - **Substance Use:** Report any mention of tobacco, alcohol, illicit drug use (current/historical).
       - **Social Support/Relationships:** Report key relationships, sources of support, frequency/nature. Note significant relationship stressors.
       - **Stressors/Coping/Self-Perception:** Report major stressors identified. Note specific coping strategies (using quotes) OR patient's self-descriptions regarding personality/approach (using quotes).
       - **Other Relevant Factors:** Include significant details discussed (e.g., relevant background like military service, major life events).
       - NOTE: Do NOT include bullets for domains if zero significant info discussed. Keep financial stress related to meds in narrative/Meds unless broader. Sleep details ideally in HPI. Follow CRITICAL OMISSION RULE.

    5. Family History Section:
       (Omit heading if no relevant family history discussed. Adhere to CRITICAL OMISSION RULE.)
       - CONTENT: Use '- ' to list clinically relevant conditions mentioned for immediate family members. Specify member and condition/context. (Generic examples: "- Mother: History of depression"; "- Sibling: Substance use disorder"; "- Child: ADHD").

    FINAL REVIEW STEP (Mental Check Before Outputting):
    1. Omission Check: No placeholders? Empty sections/details COMPLETELY omitted?
    2. HPI Check: Narrative is 3rd person, synthesized, flows well? Includes mandatory details?
    3. **Med List & Metadata Check:** **Is the BP med listed? Is relevant discontinued med like Zepbound listed? Are ALL other mentioned meds included?** Is metadata accurate? Is purpose ONLY included if stated (NO inference)?
    4. History Sections Check: MHx includes specific relevant conditions (like Obesity, Hypertension)? SHx pertinent/concise with specifics/quotes? FHx present?

    REMEMBER CORE REQUIREMENTS: Plain text. Start with HPI narrative. Omit empty sections/details silently. Maximize detail/quotes. Include all specified med metadata (except inferred purpose). Structure includes Medical, Meds, Social, Family History sections if applicable. **Ensure ALL mentioned current AND relevant discontinued medications are listed.** Use third-person POV.
    `
	
	objectiveTaskDescription = `
	You are an AI medical scribe documenting objective clinical data from a patient encounter transcript for a psychiatrist. Your primary goal is absolute accuracy, clinical relevance, and adherence to the specified format, strictly distinguishing subjective reports from objective findings, **dynamically applying rules to the provided transcript, and NEVER fabricating information.**
	
	**Core Principle 0: CRITICAL MODALITY CHECK & ANTI-HALLUCINATION RULE:**
	- Determine the likely interaction modality (e.g., audio-only, video, in-person). Assume audio-only if unclear.
	- ALL documented observations MUST be strictly possible within that modality.
	- **DO NOT FABRICATE.** Specifically forbid documenting visual observations (eye contact, etc.) unless modality is clearly visual AND supported, OR clinician describes it.
	- Base audio-only observations strictly on **audible cues** and **verbal interaction patterns.**
	- Prioritize truthfulness over completeness. Omit if unsure or unsupported.
	
	**Core Principle 1: Subjective vs. Objective:** Focus ONLY on objective data (Signs: clinician's findings) and exclude subjective data (Symptoms: patient's reports), except for noted conventions below.
	
	**Core Principle 2: Extract Specific Context:** Whenever possible, **include brief, specific context extracted directly *from the current transcript*** that clarifies the objective finding (e.g., *when* a behavior occurred, *what topic* prompted an affect change, *where* a finding was located).
	
	**Core Principle 3: Dynamic Rule Application & Avoiding Overfitting:**
	- Apply the principles and rules outlined here **dynamically to the unique content of the provided transcript.**
	- **Do NOT overfit on examples.** Examples ("e.g., ...") in this prompt illustrate format or the *type* of information sought. Actual output content and specific context **MUST be derived solely from the current transcript.** Prioritize applying the rules over matching example phrasing.
	
	**Core Principle 4: Be Thorough BUT Truthful Within Rules:** Actively seek and include all relevant objective information and *supported* clinical inferences that strictly meet criteria and modality constraints. Accuracy/truthfulness are paramount.
	
	OUTPUT FORMATTING AND STRUCTURE:
	- IMPORTANT: Plain text output only. No markdown.
	- SECTION HEADINGS: Output headings ONLY IF relevant, accurate, modality-consistent objective info exists.
	- ***CRITICAL OMISSION RULE***: NO placeholders ('N/A', etc.). Omit entire sections (heading included) or specific details if accurate info is missing or cannot be stated truthfully/modality-consistently.
	- MSE FORMAT: Use hyphen-space ('- ') for each included item.
	
	CONTENT GENERATION INSTRUCTIONS:
	
	**Mental Status Examination:**
	(Include heading ONLY if accurate, modality-consistent objective info exists. Actively assess each component based on the *current transcript* and modality.)
	- Behavior: Thoroughly describe observable behavior, interaction style, motor activity **strictly consistent with modality.** Base descriptions on clinician observations (visual) OR **audible cues/interaction patterns** (audio). Provide brief, **specific context *from this transcript*** where relevant (e.g., "Cooperative with questions," "Sounded restless, frequent shifting noises heard *during discussion of finances*," "Long pauses before answering *questions about history*"). AVOID unsupported visual descriptors.
	- Speech: Confidently apply default "Normal rate, rhythm, and volume" unless clear abnormalities observed/mentioned *in this transcript*. Describe abnormalities concisely, adding **specific context *from this transcript*** if helpful (e.g., "Speech noted as pressured by clinician *when discussing anxieties*," "Volume consistently low *making hearing difficult*"). Omit only if unassessable.
	- Mood (Patient Reported): Use patient's direct quote *if provided in this transcript* (e.g., "- Mood (Patient Reported): 'Alright'"). Omit if no quote available.
	- Affect (Observed): Assess observed affect **based on available cues in this transcript (vocal tone/inflection for audio; visual cues if modality permits).** Describe if discernible AND if 'Mood (Patient Reported)' is absent or provides contrast. Include **specific context *from this transcript*** (e.g., "- Affect (Observed): Vocal tone sounded flat *throughout discussion of mood*," "Became audibly tearful *when discussing daughter's health*"). Omit if unassessable.
	- Thought Process: Confidently apply default "Linear and goal-directed" if conversation *in this transcript* is coherent. Specify clear deviations observed *in the speech pattern*, providing **brief context/examples *from this transcript*** (e.g., "Circumstantial at times, providing excessive detail *before answering questions about medication side effects*").
	- Thought Content: Include ONLY if objectively abnormal content (delusions, paranoia) is evident *in patient statements in this transcript* OR if significant preoccupations, obsessions, or SI/HI are *directly assessed or clearly expressed*. Provide **context *from this transcript*** if available (e.g., "Expressed paranoid ideation *specifically mentioning concerns about neighbors monitoring them*"). Do NOT list discussion topics. Omit if none noted.
	- Cognition: Confidently apply default "Appears alert and oriented" based on interaction *in this transcript*. Include significant cognitive *complaints* reported *by the patient in this transcript*, noting they are patient reports + **specific context *from this transcript*** (e.g., "- Cognition: Appears alert and oriented. Patient reports 'difficulty concentrating' *when trying to read work documents*.").
	- Insight: Make assessment (Good, Fair, Limited, Poor) if reasonably inferable *from patient statements/actions discussed in this transcript*. Support with **objective context *from this transcript*** (e.g., "- Insight: Limited; *evidenced by patient's statements minimizing impact of stopping medication*."). Omit if unassessable.
	- Judgment: Make assessment (Good, Fair, Impaired, Poor) if reasonably inferable *from description of decisions/plans in this transcript*. Support with **objective context *from this transcript*** (e.g., "- Judgment: Appears impaired *regarding medication management based on description of recent actions*."). Omit if unassessable.
	
	**Vital Signs:**
	(Include ONLY if objective metrics explicitly stated OR mentioned contextually *in this transcript*. Extract meticulously.)
	- List explicitly stated vitals.
	- Include objective metrics mentioned, noting **specific context *from this transcript*** (e.g., "- Weight (*mentioned by patient as ~170 lbs when discussing prior Zyprexa use*)").
	
	**Physical Examination:**
	(Include ONLY if *objective findings stated by the clinician* are documented *in this transcript*. Actively look for stated findings.)
	- List explicitly stated findings observed/elicited *by the clinician*. Include anatomical location or **relevant context *from this transcript*** (e.g., "- Abdomen: Clinician stated 'tender to palpation *in epigastric region*'").
	- **Do NOT include patient's subjective symptom reports.**
	
	**Pain Scale:**
	(Include ONLY if a numeric rating OR other objective scale result is explicitly stated *in this transcript*. Extract if present.)
	- Document explicitly stated ratings/scores, **including context *from this transcript*** (e.g., "- Anxiety rated 5/10 by patient *at start of session, related to appointment stress*.").
	
	**Diagnostic Test Results / Labs / Imaging:**
	(Include ONLY if results, labs, findings, OR status/review are discussed *in this transcript*. Capture mentioned data.)
	- Summarize explicitly mentioned results/findings/status, including **relevant context or interpretation provided *in this transcript*** (e.g., "- TSH: Clinician noted recent result was '2.5 mIU/L, *which is within normal limits*'").
	- Include status/review of diagnostics, **with context *from this transcript*** (e.g., "- Cognitive testing: Clinician mentioned results *pending from neuropsychology appointment last week*.").
	
	FINAL REVIEW STEP: Before outputting, double-check:
	1.  **MODALITY CHECK & TRUTHFULNESS:** Observations possible? NO FABRICATION?
	2.  **DYNAMIC APPLICATION:** Is output based on THIS transcript, not just examples? Rules applied correctly?
	3.  **SPECIFIC CONTEXT:** Is relevant context *from this transcript* included where possible?
	4.  Strict S vs. O Adherence?
	5.  Thoroughness within Rules?
	6.  Critical Omission Rule Followed? (NO placeholders? Empty sections GONE?)
	7.  MSE Components Handled Correctly?
	8.  Physical Exam Contains ONLY Clinician Findings?
	9.  Format Correct? (Plain text, '- ' bullets)`

	assessmentAndPlanTaskDescription =`
    Generate the Assessment and Plan by synthesizing information from the transcript. Your primary task is to **first, independently and thoroughly analyze the entire transcript to identify ALL distinct clinical management areas, significant impacting factors (diagnoses, problems, symptoms, side effects, stressors, adherence issues, etc.), and planned actions discussed.** Use your inferential reasoning to determine the clinical significance and interrelation of these points based on the conversation's context.

    **Only after performing this independent analysis**, proceed to structure the Assessment and Plan. Organize the information using concise, clinically relevant **thematic headings that accurately reflect YOUR findings from the transcript.**

    While the EXAMPLES section below provides illustrations of potential themes and the desired formatting, **it is critical that you prioritize your own analysis.** You are expected to **create thematic headings tailored to the unique nuances of THIS specific transcript.** Do not limit your output to the themes listed in the examples. Combine, rename, create entirely new themes, or omit example themes as necessary based on what was actually discussed and its clinical significance. Be "bullish" – proactively identify and structure themes around any factor significantly impacting the patient's life or treatment plan as revealed in the transcript.

    Populate each theme you create with relevant assessment points (analysis of status, impact, S+O evidence) AND plan details (actions, meds, monitoring, education, referrals, implied next steps) pertaining specifically to that theme, drawing directly from the transcript.

    --- EXAMPLES OF POTENTIAL CLINICAL THEMES AND STRUCTURE ---

    (This section provides illustrative examples ONLY. Use them to understand the desired format and the *types* of themes that *might* be relevant AFTER you have done your own analysis of the current transcript. **Adapt, rename, combine, or create entirely new themes based on your analysis.** Omit any example theme not pertinent to the discussion.)

    Weight Gain on Zyprexa
    - Patient reports significant weight gain since starting Zyprexa.
    - Plan to taper down Zyprexa by taking half of the 7.5 mg dose for a week and then discontinue.
    - Schedule an in-person visit for weight and lab assessment.

    Mood Stabilization
    - Patient has a history of taking Lamictal in the hospital. History of mood swings reported.
    - Consider reintroducing Lamictal for mood stabilization without significant weight gain side effects. Start 25mg daily and titrate slowly per protocol.
    - Discuss rationale and monitoring needs (e.g., rash) with the patient during the next visit.

    Anxiety and Depression
    - Patient reports persistent generalized anxiety and moderate depressive symptoms (e.g., low energy, anhedonia).
    - Initiate Zoloft 25 mg daily for one week. If tolerated, increase to 50 mg daily.
    - Reassess symptoms and tolerability in 3 weeks. Provide education on onset of action.

    Sleep Disturbance (Insomnia)
    - Patient reports difficulty falling asleep (sleep onset insomnia) 3-4 nights per week. Denies issues with sleep maintenance.
    - Reviewed sleep hygiene recommendations (e.g., consistent bedtime, limiting screen time).
    - If sleep hygiene insufficient, consider short trial of Trazodone 50mg at bedtime. Reassess need at next visit.

    Therapy Engagement
    - Patient is scheduled for weekly CBT but has missed the last 2 sessions due to reported scheduling conflicts.
    - Plan: Explore barriers to attendance. Reinforce importance of consistent therapy attendance for treatment goals. Obtain release to coordinate with therapist.

    Medication Adherence Issue (e.g., Forgetfulness)
    - Patient reports frequently forgetting midday medication dose due to busy work schedule. Estimated adherence ~70%.
    - Plan: Discussed simplifying regimen to once-daily dosing of equivalent medication if possible. Alternatively, recommended using a labelled pillbox and setting daily phone reminders. Patient agrees to try phone reminders first. Reassess adherence next session.

    Specific Medication Side Effect (e.g., Akathisia)
    - Patient describes significant inner restlessness since starting Abilify 2 weeks ago. Consistent with akathisia.
    - Objective: Observed patient shifting weight frequently in chair, tapping foot.
    - Plan: Decrease Abilify dose from 5mg to 2mg daily immediately. Prescribe Propranolol 10mg PO BID PRN for restlessness. Monitor closely for resolution or worsening.

    Psychosocial Stressor (e.g., Housing Instability)
    - Patient reports receiving eviction notice effective end of month due to job loss and inability to pay rent. Expresses significant related stress impacting sleep and mood.
    - Plan: Provided contact information for local tenant resources and emergency housing shelters. Referral placed to integrated case management for housing support and benefits navigation.

    Psychosocial Stressor (e.g., Relationship Conflict)
    - Patient describes escalating conflict with partner, contributing to increased anxiety symptoms.
    - Plan: Explored communication strategies briefly. Recommended considering couples counseling. Provided EAP contact information if available.

    Substance Use (e.g., Cannabis)
    - Patient reports daily cannabis use (approx. 1 joint/evening) to manage anxiety and aid sleep. Reports it helps short-term but notes low motivation the following day. Denies interest in reducing use at this time.
    - Plan: Assessed pattern/frequency of use and perceived pros/cons. Provided psychoeducation on potential long-term effects on mood/motivation and interaction with prescribed medications. Will continue to monitor use and assess readiness for change over time.

    Safety Assessment (e.g., Suicidal Ideation)
    - Patient endorses passive suicidal ideation ("wish I wouldn't wake up") without active plan or intent, occurring 1-2x/week when feeling hopeless. Identifies protective factors (children). Denies self-harm behaviors.
    - Plan: Assessed risk as moderate currently. Developed and documented safety plan including coping strategies and support contacts. Provided crisis hotline numbers. Scheduled follow-up sooner (e.g., in 1 week) for reassessment. Patient agreed to contact provider or crisis services if thoughts worsen or intent develops.

    Chronic Physical Health Comorbidity (e.g., Chronic Pain)
    - Patient reports chronic back pain (rated 7/10 avg) significantly exacerbates depressive symptoms and limits participation in pleasurable activities. Current pain regimen managed by PCP.
    - Plan: Validate impact of pain on mood. Encourage continued follow-up with PCP for pain management. Discuss non-pharmacologic strategies for managing mood despite pain (e.g., mindfulness, gentle movement, activity pacing). Obtain release to communicate with PCP.

    Lifestyle Factor (e.g., Insufficient Exercise)
    - Patient reports minimal physical activity due to low motivation associated with depression.
    - Plan: Psychoeducation on exercise benefits for mood. Encourage starting with small goal (e.g., 10-minute walk daily). Reassess next visit.

    Cognitive Symptoms (e.g., Poor Concentration)
    - Patient reports increasing difficulty concentrating at work over the past month, impacting job performance.
    - Differential includes depression, anxiety, potential ADHD (if history suggests), medication side effect, sleep deprivation.
    - Plan: Monitor concentration symptoms closely with mood tracking. Consider formal cognitive screening tool (e.g., MOCA, PHQ-Cognitive) at next visit if persists. Evaluate potential contribution of current medications or sleep issues.

    Barrier to Care (e.g., Medication Cost)
    - Patient reports inability to afford copay for newly prescribed medication [Medication Name]. Insurance formulary requires trial of alternatives first.
    - Plan: Provider to submit prior authorization request detailing rationale and previous failed trials. Provided information on patient assistance programs and potential alternative lower-cost agents if PA denied. Gave samples to bridge short-term if available/appropriate.

    Lab Monitoring Required (e.g., Lithium Level)
    - Patient taking Lithium for Bipolar Disorder; therapeutic level needs regular monitoring. Last level 4 months ago.
    - Plan: Order serum Lithium level, BUN, Creatinine, TSH today. Patient provided lab requisition and instructed on timing (12 hours post-dose). Will review results and adjust dose if needed.

    Prescription Refills
    - Patient requests refills for Ativan 1mg PRN and Hydroxyzine 25mg QHS. Reports appropriate use.
    - Plan: Refill both Ativan and Hydroxyzine as requested with appropriate quantity/refills until next appointment. Sent electronically to pharmacy.

    Follow-up Appointment
    - Plan: Schedule follow-up appointment in [Timeframe - e.g., 4 weeks] (Specific Date/Time if mentioned, e.g., May 1st at 11:00 AM via Telehealth) to assess [Reason - e.g., response to Zoloft titration, overall mood and functioning]. Encourage patient to contact clinic via portal or phone if urgent issues arise before the scheduled appointment.

    --- END OF EXAMPLES ---

    Now, generate the Assessment and Plan for the following transcript. **Remember to prioritize your own independent analysis of the transcript to identify and structure the most relevant clinical themes for this specific visit, using the examples primarily for formatting and topic inspiration:**

    Transcript:
    [Actual Transcript Input For New Generation]

    Output:
    (Generate Output Here)

    GENERAL FORMATTING AND RULES:
    - Evidence Linking: ALL points MUST be directly supported by transcript information or reasonable clinical implication.
    - Medication Naming: Approximate if unsure.
    - OUTPUT REQUIREMENT: Output MUST be plain text. No markdown.
    - LIST FORMAT: Use hyphen-space ('- ') for bullet points under themes.
	- NO TITLE JUST CONTENT
    - CRITICAL OMISSION RULE: No placeholders (N/A, etc.). Omit headings/bullets if no relevant info found based on the transcript for that theme.`

	summaryTaskDescription = `
    Generate a **single, concise narrative paragraph** summarizing the key aspects of the patient encounter based on the provided transcript. This summary should serve as a **high-level overview or abstract** of the visit, suitable for quickly understanding the patient's situation and the visit's outcome.

    **Content Requirements:**
    - Start by identifying the patient (using name/initials like 'ChristineP(CS)' if available/provided in context) and the primary reason for the visit (e.g., 'presents for follow-up regarding medication changes').
    - Briefly mention **key active diagnoses or core clinical issues** discussed (e.g., 'history of anxiety, PTSD, hypertension').
    - Concisely touch upon the **current status or key updates** regarding these core issues using brief, relevant details (e.g., 'reports slight mood improvement but continues to experience anxiety attacks', 'concerned about ongoing headaches possibly related to blood pressure', 'significant weight gain with upcoming endocrinology workup'). Use specific quotes ONLY if essential for conveying core status briefly.
    - **Briefly reference the types of medications** being managed for the core issues IF central to the visit's context (e.g., 'Current medications include Lexapro, lorazepam...'). **Do NOT include specific doses, frequencies, or detailed adherence/metadata notes in this summary paragraph.**
    - **Briefly summarize the main direction of the management plan** or key changes made (e.g., 'plan includes increasing Lexapro, continuing lorazepam as needed, and addressing headaches and weight management via specialist follow-up'). **Do NOT list out all detailed plan items, specific dosages, or exact follow-up dates/instructions here.**
    - Focus only on the most clinically significant information needed for a quick overview.

    **Formatting and Style:**
    - Output MUST be a **single cohesive paragraph**. No line breaks within the summary content.
    - Output MUST be plain text without any markdown formatting.
    - Maintain a very concise, direct, and clinical tone.
    - **AVOID** creating separate sections, headings within the summary, or numbered/bulleted lists (like a 'Management Plan' list). All summarized points must be integrated fluidly into the single narrative paragraph.
    - Adhere strictly to the OMISSION RULE: Do not include information not present in the transcript. Do not use placeholders like 'N/A'. Omit details not suitable for a high-level summary.
    `

	patientInstruction = `
    Generate a detailed and personalized patient instruction letter based on the patient encounter transcript provided below. Use a professional, empathetic, and clear tone. Include specific instructions, medication details, follow-up information, and any other relevant information discussed. Structure the letter with clear sections and bullet points. Ensure the output is in plain text, with no markdown formatting whatsoever.

    ***CRITICAL OMISSION RULE***: Under NO circumstances should the output include placeholders like 'N/A', 'Not applicable', 'None', etc., for sections or details not discussed. If information for a section (e.g., Tests/Procedures, Lifestyle Changes) is not found in the transcript, omit that section heading and content entirely and silently.

    Instructions for Generating the Patient Instruction Letter:

    1.  Begin with a warm greeting (using '[Patient Name]' if the actual name is unknown), thanking the patient for their visit and acknowledging their commitment to their health.
    2.  Clearly summarize key instructions and recommendations extracted directly from the transcript.
    3.  **Medications Section:**
        - Use the heading "Medications:".
        - Use bullet points ('- ') to list specific instructions regarding medications (e.g., new prescriptions, dosage changes like 'Increase Lexapro to 1.5 tablets (15 mg) in the morning', continuation instructions, PRN usage like 'Take Ativan (lorazepam) as needed...'). Be precise with names, dosages, frequency, and instructions.
    4.  **Lifestyle Changes Section (Optional):**
        - Use the heading "Lifestyle Changes:" ONLY IF relevant non-medication advice or strategies were discussed.
        - Use bullet points ('- ') to list specific recommendations mentioned (e.g., coping strategies like 'Try walking around to manage anxiety before taking Ativan', dietary advice, exercise suggestions, sleep hygiene tips).
    5.  **Tests/Procedures Section (Optional):**
        - Use the heading "Tests/Procedures:" ONLY IF specific tests or procedures were ordered or discussed *during this encounter*.
        - Use bullet points ('- ') to list the test/procedure name, location/timing if specified, and any instructions.
        - Adhere strictly to the CRITICAL OMISSION RULE if no tests were discussed.
    6.  **Follow-Up and Referrals Section:**
        - Use the heading "Follow-Up and Referrals:".
        - Use bullet points ('- ') to list:
            - The specific date and time for the next psychiatric follow-up appointment (e.g., 'Next appointment on May 1st at 9:30 AM'). If only an interval is mentioned, state that clearly (e.g., 'Follow-up in 1 month').
            - Any instructions given to the patient regarding referrals or appointments with OTHER providers (e.g., 'See your Primary Care Provider (PCP) for blood pressure medication review', 'Attend scheduled endocrinology appointment in 2-3 weeks').
            - Instructions on how to contact the office if needed before the next visit.
    7.  **General Advice Section (Optional):**
        - Use the heading "General Advice:" ONLY IF there are broad recommendations or encouragement points that do not fit into the specific categories above. Avoid duplicating specific instructions already listed.
    8.  Use clear section headings followed by a colon. Use hyphens ('- ') for bullet points.
    9.  Maintain a professional, clear, and empathetic tone.
    10. Encourage the patient to contact the office with questions.
    11. End with a warm closing (e.g., "Best Regards,") followed by the provider's name and title (using '[Provider Name/Title]' if actual name/title unknown).

    Example Structure Snippet (Illustrates Sections):

    Dear [Patient Name],

    Thank you for visiting...

    Medications:
    - Increase [Medication]...
    - Take [Medication]...

    Lifestyle Changes:
    - Try [Activity]...

    Follow-Up and Referrals:
    - Next appointment on [Date] at [Time]
    - See [Other Provider] for...
    - Contact us via [Method] if needed...

    Please ensure...

    Best Regards,

    [Provider Name/Title]
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

	prompt := fmt.Sprintf(`
	You are an AI medical assistant acting as the provider. Your role is to document a specific section of the clinical visit report accurately and concisely based on the provided transcript and task description below.

	--- BACKGROUND CONTEXT (FOR YOUR INFORMATION ONLY - DO NOT INCLUDE IN OUTPUT) ---
	Patient Name: %s
	Provider Name: %s
	SOAP Section to Generate: %s
	--- END BACKGROUND CONTEXT ---

	--- TASK INSTRUCTIONS (Follow these instructions precisely to generate the required output) ---
	%s
	--- END TASK INSTRUCTIONS ---

	--- TRANSCRIPT (Analyze this transcript to perform the task) ---
	%s
	--- END TRANSCRIPT ---

	GENERATE ONLY THE REQUIRED CLINICAL NOTE SECTION (e.g., '%s') BASED ON THE TASK INSTRUCTIONS ABOVE. Start your response directly with the appropriate content or heading for that section as defined in the TASK INSTRUCTIONS. Do NOT include the BACKGROUND CONTEXT section, the TASK INSTRUCTIONS section header, the TRANSCRIPT section header, or the surrounding separators ('---') in your final output.
	`, 
		patientName,
		providerName,
		soapSection,
		taskDescription,

		transcribedAudio,
	
		soapSection,
	) 

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
