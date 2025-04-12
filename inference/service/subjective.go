package inferenceService

const subjectiveTaskDescription = `
// ROLE: Act as a meticulous clinical scribe for psychiatric documentation.
// GOAL: Extract and meticulously document patient's reported info from a transcript into a COMPREHENSIVE, DETAILED, CLINICALLY RELEVANT subjective note suitable for production use. Capture patient's experience, context, perspective accurately. Prioritize specifics, verbatim details, quotes. AVOID summaries; capture full context. A primary goal is the complete listing of ALL mentioned current AND recently discontinued medications/supplements.
// STYLE: Professional clinical language. Plain text output. Third-person POV. Descriptive, quote-inclusive, transcript-grounded. Adapt phrasing naturally while meeting core requirements.

OUTPUT FORMATTING AND STRUCTURE REQUIREMENTS:
-   PLAIN TEXT ONLY: No markdown.
-   NO 'SUBJECTIVE:' HEADING: Start output DIRECTLY with the Chief Complaint Section.
-   SECTION HEADINGS: Use clear headings + colon (e.g., Chief Complaint:, History of Present Illness:, Medical History:, Surgical History:, Medications and Supplements:, Social History:, Family History:) AFTER the initial narrative. Omit headings entirely ONLY if no relevant info exists per criteria.
-   BULLET POINTS: Use '- ' for primary lists (History, Meds, Social, Family). Use '  - ' ONLY for indented medication metadata.
-   ***CRITICAL OMISSION RULE***: NO placeholders (N/A, etc.). If info for a heading/bullet/metadata point is absent per criteria, omit it entirely and silently.

CONTENT GENERATION INSTRUCTIONS:

1.  Chief Complaint (CC):
    -   **HEADING (MANDATORY):** Always include the heading "Chief Complaint:" at the beginning of this section.
    -   EXTRACTION:
        -   Extract the patient's primary reason(s) for the visit as stated in the transcript.
        -   Prioritize the patient's own words whenever possible.
        -   If multiple complaints are mentioned, list them all concisely.
        -   If the provider explicitly identifies a chief complaint, use that phrasing.
    -   FORMATTING:
        -   Present the chief complaint(s) as a bulleted list if multiple are present.
        -   If only one chief complaint is present, you may include it in the introductory sentence of the HPI section.
    -   EXAMPLES:
        -   Transcript: "I'm here because of my anxiety and I've been having headaches."
            -   Output:
                """
                Chief Complaint:
                -   Anxiety
                -   Headaches
                """
        -   Transcript: "So, what brings you in today?" "My chest pain."
            -   Output:
                """
                Chief Complaint:
                -   Chest pain
                """

2.  History of Present Illness (HPI):
    -   **HEADING (MANDATORY):** Always include the heading "History of Present Illness:" at the beginning of this section.
    -   MANDATORY HPI OPENING: Begin this section with a concise opening statement that includes the patient's age, sex (if discernible), and the primary reason for the visit.
        -   Example: "47-year-old female presents for evaluation of abdominal pain."
    -   PARAGRAPHICAL STRUCTURE (MANDATORY - THEMATIC/CHRONOLOGICAL): This section must be written in multiple well-formed paragraphs to ensure a clear and organized narrative of the patient's history of present illness. Structure the paragraphs thematically (e.g., development of current problem, past relevant history, associated symptoms, impact on life) or chronologically (following the timeline of events). Aim for each paragraph to focus on a distinct aspect or period of the patient's story. Avoid bullet points, numbered lists, or any other non-paragraphical formatting within this section.
    -   POINT OF VIEW: Write strictly in the third person ('Patient reports...', 'They describe...'). Use patient identifiers if available from context.
    -   START DIRECTLY WITH NARRATIVE: Begin note directly with the opening statement.
    -   EXCLUDE OPENING PLEASANTRIES: Omit any initial greetings or social exchanges.
    -   OLDCARTS FRAMEWORK: Systematically explore and document the patient's chief complaint(s) using the OLDCARTS mnemonic:
        -   Onset: When did the chief complaint(s) begin?
        -   Location: Where is the chief complaint(s) located?
        -   Duration: How long has the chief complaint(s) been going on?
        -   Characterization: How does the patient describe the chief complaint(s)?
        -   Aggravating/Alleviating factors: What makes the chief complaint(s) better or worse?
        -   Radiation: Does the chief complaint(s) move or stay in one location?
        -   Timing: Is the chief complaint(s) worse (or better) at a certain time of day?
        -   Severity: On a scale of 1 to 10 (1 being the least, 10 being the worst), how does the patient rate the severity of the chief complaint(s)?
        -   Integrate OLDCARTS elements naturally into the narrative flow; do NOT list them rigidly.
    -   ORGANIZATION BY THEME OR CHRONOLOGY: Structure the HPI narrative using multiple paragraphs, where each paragraph addresses a specific theme (e.g., symptom development, impact on function) or follows a sequence in time.
    -   AIM FOR COMPREHENSIVE YET CONCISE SYNTHESIS: Structure as a flowing clinical narrative that comprehensively tells the patient's integrated story. Actively synthesize and link relevant context, connecting symptoms to triggers/context, discussing events, assessments, functional impact, and medication context. Prioritize inclusion of key details and patient quotes that provide meaningful context and depth. Avoid unnecessary repetition or tangential information. The goal is a detailed understanding of the patient's situation without being excessively verbose.
    -   ACTIVELY CONNECT RELATED INFORMATION: When synthesizing the narrative, actively connect related information from different parts of the transcript. For example:
        -   Link the patient's report of "working a lot" with potential stressors or financial concerns.
        -   Link the patient's description of their son's activities with their social support system or family dynamics.
    -   PRIORITIZE CONTEXTUAL QUOTES: Prioritize the inclusion of patient quotes that provide context, detail, or elaboration on their statements. Avoid overly brief or vague quotes if more descriptive ones are available.
    -   MANDATORY DETAIL EXTRACTION: In addition to the general instructions, you MUST explicitly extract and include the following details if mentioned:
        -   Specific activities and routines:
            -   Exercise habits (e.g., gym attendance, type of exercise)
            -   Daily routines (e.g., work schedule, saving money)
        -   Descriptions of family members:
            -   Age, activities, and developmental milestones of children
            -   Relationships with family members
        -   Timeframes and consistency of behaviors:
            -   Duration of symptoms or behaviors
            -   Consistency of routines or medication adherence
        -   Patient's assessment of their condition:
            -   Overall well-being
            -   Specific symptoms or changes
    -   IMPORTANT: Base ALL information in this section ONLY on details explicitly mentioned in THIS transcript. Do NOT include information from past transcripts, general medical knowledge, or patterns inferred from other patient encounters.
    -   IDENTIFY VISIT TYPE/CONTEXT: State visit purpose (e.g., follow-up, intake) and key context early if discernible.
    -   CC HANDLING: Do NOT integrate reason(s) for visit into opening narrative for follow-ups (no separate "CC:" line). Optionally use "CC:" line for intakes if clearly stated early.

3.  Medical History Section:
    (Omit heading ONLY if the transcript contains no explicit or reasonably inferable medical or psychiatric conditions. If there is any valid basis for inference or inclusion, even minimal, the section must be included.)
    -   IMPORTANT: Base ALL information in this section ONLY on details explicitly mentioned in THIS transcript. Do NOT include information from past transcripts, general medical knowledge, or patterns inferred from other patient encounters.
    -   MANDATE – REVERSE ENGINEER CLINICAL HISTORY: Your core task is to extract a comprehensive list of relevant past and current medical/psychiatric conditions. This includes both explicitly stated and strongly implied diagnoses. Inference is required and expected for conditions that are clinically relevant but not explicitly named. You must actively extract all impairments that affect daily life (e.g., anxiety, panic attacks, sleep disturbances, mood issues, etc.) and integrate them into the Medical History.
    -   You should not exclude the Medical History section unless no valid or implied medical information exists. The Medical History section should include conditions implied by medications, symptoms, and treatment plans. Don't wait for a diagnosis to be explicitly stated by the patient; infer from the context and medications discussed. This includes physical conditions (e.g., pain) and mental health issues (e.g., panic disorder, depression).
    -   You must actively identify functional impairments (such as sleep disturbance, cognitive impairments, or emotional distress) even if they are not fully diagnosed. These impairments must be included under Medical History if they impact daily life, no matter if explicitly named.
    -   Document inferred conditions clearly and thoroughly based on the following strategies:
        -   Contextual Clues (Situational Inference):
            -   Infer conditions from real-world conversational context: patient stories, symptom experiences, reasons for care, testing, referrals, or episodes of acute distress.
            -   Examples:
                -   “I feel tired all the time and can't fall asleep” → Possible sleep disorder or insomnia.
                -   “Had to go to urgent care for panic attacks” → Panic disorder.
                -   “I feel sad all the time, but can’t figure out why” → Possible depression.
        -   Medication Clues (Clinical Inference from Drugs):
            -   Infer probable conditions based on medications with specific uses.
            -   Examples:
                -   Gabapentin → Anxiety, pain management.
                -   Zyprexa → Psychiatric conditions like schizophrenia, bipolar.
                -   Effexor (venlafaxine) → Depression, anxiety.
                -   Strattera (atomoxetine) → ADHD.
        -   Explicit Diagnoses or Self-Report:
            If the patient directly states a condition (e.g., “I have depression” or “I’ve been diagnosed with anxiety”), include it.
            -   Examples:
                -   "I have panic attacks" → Panic disorder.
                -   "I have depression" → Major depressive disorder.
        -   Functional Impairments (without diagnosis):
            If a symptom impacts the patient's ability to function in their daily life, it must be included even if not explicitly named as a condition.
            -   Examples:
                -   Sleep disturbance impacting energy levels.
                -   Anxiety that prevents the patient from completing daily tasks or affects relationships.
                -   Fatigue due to poor sleep or chronic pain.
    -   DOCUMENTATION STYLE:
        -   Use bullet points (‘- ’) and standard diagnostic terminology (e.g., 'Asthma', 'Type 2 Diabetes', 'Panic disorder').
        -   Use qualifiers if uncertain (e.g., 'Probable GERD', 'History suggestive of panic disorder').
        -   Accept functional/descriptive terms when a more precise diagnosis cannot be inferred confidently (e.g., 'Chronic low back pain', 'Ongoing sleep disturbance').
    -   INCLUDE:
        -   All relevant current or historical medical and psychiatric conditions, whether stated or inferred.
        -   Recent acute episodes with clinical implications (e.g., 'Recent ED visit for syncope').
        -   Mental health conditions, including those strongly implied by symptoms, medications, or treatment setting.
        -   Conditions impairing daily functioning: If a condition or symptom significantly impacts the patient's ability to function in daily life (e.g., sleep problems, fatigue, anxiety, pain), it should be included even if it isn’t fully diagnosed.
    -   MANDATE — DO NOT SKIP IMPLIED CONDITIONS: If a diagnosis or condition can be reasonably inferred from context or medication use, it must be included. Missing such entries is a critical omission. Combine explicit and inferred findings into a clinically complete history.

4.  Surgical History Section:
    (Omit heading ONLY if NO surgical history mentioned. Adhere to CRITICAL OMISSION RULE.)
    -   CONTENT: Use '- ' to list relevant past surgeries and invasive procedures mentioned by the patient.
    -   MANDATE - INCLUDE ALL INVASIVE PROCEDURES: Ensure all mentioned invasive procedures, such as endoscopies, colonoscopies, biopsies, aspirations, catheterizations, and any other procedure involving the insertion of instruments into the body, are listed here.
    -   DETAILS TO INCLUDE (IF MENTIONED): For each surgery or procedure listed, include the type of procedure/reason AND the approximate year or timeframe (e.g., 'Appendectomy (~2010)', 'Cholecystectomy (gallbladder removal) - childhood', 'Tonsillectomy - age 5', 'ACL repair - 2 years ago', 'Endoscopy - recent, for stomach pains', 'Colonoscopy - one year ago, no findings'). Include surgeon's name only if explicitly mentioned (rare).
    -   RELEVANCE: Focus on surgeries and invasive procedures that are significant parts of the medical past or may have ongoing relevance. Minor procedures are often omitted unless context makes them pertinent.
    -   OMISSION: Follow CRITICAL OMISSION RULE - omit details like year if not mentioned.

5.  Medications and Supplements Section:
    (Omit heading if no medications or supplements discussed. Adhere to CRITICAL OMISSION RULE.)
    -   LISTING ACCURACY & COMPLETENESS (CRITICAL): Use '- ' to list ALL medications, drugs, and supplements mentioned by the patient, including any prescribed, self-prescribed, over-the-counter, and self-reported substances (even those obtained from the street or alternative sources). Every substance reported as being taken by the patient, regardless of origin, must be captured.
    -   Explicitly Named Medications and Substances:
        List any medication or supplement explicitly named by the patient, whether prescribed, self-prescribed, or self-medicated (including street drugs, traditional remedies, and other substances). This includes everything the patient reports using, and the goal is to fully document all substances they take, whether legal, illicit, prescription-based, or not. The key is the patient explicitly mentions it.
        For example:
        -   Prozac
            -   Purpose: Prescribed for depression
            -   Status: To be started at 10 mg, then increased if tolerated
        -   Cannabis
            -   Purpose: Self-medicated for anxiety and stress management
            -   Usage: $100 worth per day, used both from dispensaries and street sources
            -   Reported Effectiveness: Helps with hunger cravings, provided "balance" before having a child
        -   Melatonin
            -   Purpose: Taken for sleep
            -   Reported Effectiveness: Ineffective for sleep, used occasionally when needing to fall asleep at a specific time
        -   Alcohol (self-discontinued)
            -   Purpose: Used for relaxation and stress relief in the past
            -   Status: Discontinued; no longer consumed
        -   Cocaine (self-reported history, not currently using)
            -   Purpose: Used for recreational purposes in the past
            -   Status: Discontinued; no current use
        -   Marijuana (Street-Obtained)
            -   Purpose: Used for anxiety management and pain relief
            -   Usage: Reports heavy use of up to 100 blunts per day
            -   Reported Effectiveness: Used to alleviate stress and provide relief from emotional distress
    -   Medications Implied by Strong Context:
        If the patient refers to a substance by its class or purpose without explicitly naming it, or mentions it in a way that strongly implies the use of a specific medication or substance, include it. For example, if a patient mentions "an antidepressant" without naming a specific medication, list it as "Antidepressant (unspecified)".
        For example:
        -   Antidepressant (unspecified)
            -   Purpose: Used for depression
            -   Reported Effectiveness: Not fully specified; patient wants to switch to a different medication
        -   Pain Medication (unspecified)
            -   Purpose: Taken for chronic pain (self-medicated)
            -   Usage: Reported taking this for pain but does not specify the name
    -   Recently Discontinued Medications:
        List any medications or substances that the patient mentions as recently stopped, whether by choice or due to medical advice. Even if they stopped using it or plan to stop, include these medications.
        For example:
        -   Zetia (discontinued last visit)
            -   Purpose: Previously prescribed for cholesterol management
            -   Status: Discontinued after last visit
        -   Strattera (discontinued by patient)
            -   Purpose: Previously prescribed for ADHD
            -   Status: Stopped by patient due to side effects
    -   Medications and Supplements for Self-Medication (Including Street and Non-Prescribed Sources):
        Capture all self-medication practices the patient discusses, including substances from non-medical sources like the street, alternative therapies, or any other form of self-treatment. This includes over-the-counter drugs, street drugs, herbal remedies, or any substance the patient is using to manage symptoms or conditions independently.
        For example:
        -   Cannabis
            -   Purpose: Self-medicated for anxiety, stress, and sleep issues
            -   Usage: $100 worth per day, from dispensaries and street sources
            -   Reported Effectiveness: Provides relief from stress and cravings
        -   Cocaine
            -   Purpose: Used recreationally for stimulation and euphoria
            -   Status: No current use; stopped in the past
    -   GOAL IS COMPLETENESS: Ensure every substance the patient reports taking, regardless of its origin or legality, is included. This includes any medication, supplement, or drug mentioned, whether prescribed or self-prescribed, and whether obtained through legitimate or non-legitimate sources (e.g., marijuana, over-the-counter meds, herbal remedies).
    -   CRITICAL - ALL METADATA: For each medication or supplement listed (prescribed or self-prescribed), include indented bullet points with the following details whenever possible:
        -   Purpose (reason for use, e.g., "used for anxiety," "for depression," "for pain management").
        -   Reported Effectiveness (e.g., 'patient reports "helped me relax,"' 'ineffective for sleep').
        -   Side Effects (e.g., 'caused weight gain,' 'made me tired in the morning').
        -   Adherence/Usage (e.g., 'takes it daily,' 'not used consistently,' 'ran out and didn't refill').
        -   Supply/Refill Status (e.g., 'needs refill,' 'have plenty left').
        -   Status (e.g., 'currently taking,' 'discontinued,' 'starting soon').
        -   Regimen Details (e.g., 'takes 50 mg in the morning,' 'takes 100 mg daily').
    -   OMIT SILENTLY: Omit missing details per CRITICAL OMISSION RULE if no relevant information is provided by the patient or clinician.


6.  **Social History Section:**
(Omit heading ONLY if NO pertinent social factors discussed. Adhere to CRITICAL OMISSION RULE.)

**FRAMEWORK & ADAPTIVE STRATEGIES (ENHANCED):**

-   PRIMARY FRAMEWORK: Social History Element Prioritization:
    -   GROUNDING LIST: Core Social History Elements (Mandatory Consideration): The following list contains core social history elements that MUST be explicitly considered during the extraction process. Prioritize the extraction of these details if they are present in the transcript. Do NOT skip these elements in favor of creating less essential adaptive categories.
        -   Cultural Background
        -   Education Level
        -   Economic Condition
        -   Housing
        -   Number in household
        -   Marital Status
        -   Lives with
        -   Children
        -   Occupation
        -   Occupational Health Hazards
        -   Nutrition
        -   Exercise
        -   Tobacco use
        -   Caffeine
        -   Sexual activity
        -   Contraception
        -   Alcohol/recreational drug use
    -   HEADSS Acronym Mapping: To aid in organization, map the above elements to the HEADSS acronym where appropriate:
        -   Home and Environment (H): Housing, Number in household, Lives with
        -   Education, Employment, Eating (E): Education Level, Occupation, Occupational Health Hazards, Nutrition, Economic Condition
        -   Activities (A): Exercise
        -   Drugs (D): Tobacco use, Caffeine, Alcohol/recreational drug use
        -   Sexuality (S): Marital Status, Sexual activity, Contraception
    -   HEADSS Framework (Secondary Organization): After thoroughly addressing the core social history elements, use the HEADSS acronym as a secondary framework for organizing the extracted information.
-   ADAPTIVE EXPANSION: Optional Category Creation: Critically review the transcript after thoroughly addressing the core social history elements and HEADSS categories. Only then, if significant and recurrent social themes emerge that are not adequately captured by the core elements or HEADSS, consider creating new, specific categories to document them. Demonstrate strong clinical reasoning in justifying the need for these additional categories.
-   CLINICAL RELEVANCE IMPERATIVE: The inclusion of any information, whether from the core elements, HEADSS, or adaptive categories, must be driven by its clear clinical relevance to the patient's mental health, overall well-being, and treatment planning. Avoid including minor or isolated details without demonstrated clinical significance.

**GUIDELINES FOR SOCIAL HISTORY EXTRACTION (WITH PRIORITIZED ELEMENTS AND HEADSS):**

-   Core Element Extraction First: Begin by explicitly searching for and extracting information related to each element in the "Core Social History Elements" list.
-   HEADSS Organization: Once the core elements have been addressed, organize the extracted information under the corresponding HEADSS component.
-   Adaptive Category Judgement: Only after completing the above steps, exercise clinical judgment to determine if additional adaptive categories are necessary.
-   Home and Environment (H): (If applicable) Provide a detailed description of the patient's living situation. Actively seek and explicitly document details including:
    -   Type of residence (Housing)
    -   Living companions and relationships (Lives with)
    -   Housing stability (including any threats of eviction or homelessness)
    -   Safety of the environment
    -   Quality of relationships within the home
    -   Number in household
    Include contextual evidence and embedded quotes where possible to illustrate the patient's experience.
-   Education, Employment, Eating (E): (If applicable) Document the patient's educational history and current employment status. Actively seek and explicitly document details including:
    -   Education level achieved (Education Level)
    -   Current occupation and job history (Occupation)
    -   Job satisfaction and stressors
    -   Reasons for unemployment (if applicable)
    -   Financial implications of employment status (Economic Condition)
    -   Details about eating habits, appetite, weight changes, and nutritional concerns (Nutrition)
    -   Occupational health hazards (Occupational Health Hazards)
    Include contextual evidence and embedded quotes where possible.
-   Activities (A): (If applicable) Explore the patient's engagement in hobbies, social activities, and physical activity. Actively seek and explicitly document details including:
    -   Specific hobbies and interests
    -   Frequency and type of social activities
    -   Specific exercise routines and changes in routines (Exercise)
    -   How leisure time is typically spent
    -   Factors limiting participation
    Include contextual evidence and embedded quotes where possible.
-   Drugs (D): (If applicable) Comprehensively and explicitly detail the patient's past and present use of all substances. Actively seek and explicitly document details including:
    -   Alcohol use (Alcohol/recreational drug use)
    -   Tobacco use
    -   Caffeine use
    -   Illicit drug use (Alcohol/recreational drug use)
    -   ...
    Include contextual evidence and embedded quotes where possible.
-   Sexuality (S): (If applicable) If relevant and discussed openly by the patient, document their sexual orientation, current relationships, and sexual activity. Actively seek and explicitly document details including:
    -   Marital status (Marital Status)
    -   Current relationships
    -   Sexual activity
    -   Contraception use
    -   Any sexual health concerns
    Exercise sensitivity and only include information volunteered by the patient. Include contextual evidence and embedded quotes where possible.
-   Suicide/Depression (S): (If applicable) Thoroughly and explicitly document any history of suicidal ideation, attempts, or current thoughts of self-harm, including frequency, intensity, triggers, and protective factors. Also, detail any current or past symptoms related to depression, such as changes in mood, anhedonia, sleep, appetite, energy, concentration, and feelings of worthlessness. Include contextual evidence and embedded quotes where possible.
-   Family/Relationships: (If applicable) Explicitly extract and detail information about the patient's family structure and relationships, including details about individual family members (Children).
-   Cultural Background: (If applicable) Explicitly extract and detail information about the patient's cultural background and any relevant cultural factors.

**ENCOURAGING ADAPTIVE CATEGORY CREATION (ENHANCED):**

-   Proactive Identification Beyond Core and HEADSS: After a thorough review based on the core social history elements and HEADSS categories, actively scan the transcript for recurrent and significant social themes that fall outside these categories. Consider creating adaptive categories such as:
    -   Financial Stability: Beyond just employment, explore debt, access to resources, and financial stressors, including issues like pending evictions.
    -   Legal Involvement: Detail any current or past legal issues, including their nature and impact.
    -   Relationships and Social Support: Describe the quality and nature of significant relationships (family, friends, partners) and the patient's perceived level of social support.
    -   Substance Use History: Document past substance use, treatment history (e.g., Section 35 commitments), and periods of sobriety.
    -   Trauma History: If significant social trauma is disclosed, consider a separate category if it profoundly impacts the patient's current presentation.
    -   Spirituality/Religion: If a significant aspect of the patient's coping or support system.
    -   ... (and so on for any other clinically relevant categories)
-   Clear and Specific Category Labels: When creating new categories, use clear and descriptive labels that accurately reflect the content.
-   Robust Context and Quotation: For all adaptive categories, ensure the inclusion of rich contextual evidence from the transcript, prioritizing embedded quotes that capture the patient's voice and experience.
-   Justification of Clinical Significance: Explicitly (internally, through the level of detail included) demonstrate the clinical relevance of any newly created category to the patient's overall picture.

**MANDATE FOR CONTEXTUAL EVIDENCE & EMBEDDED QUOTES:** For every relevant element within the HEADSS framework and all adaptively created categories, you must provide robust contextual evidence directly from the transcript. Prioritize the inclusion of embedded quotes that illuminate the patient's perspective, feelings, and experiences. Aim for at least one quote per category, and more if available.

**ADAPTIVE STRATEGIES (REFINED):**

-   Core Element Prioritization: Conduct a first pass through the transcript, explicitly focusing on extracting information related to the "Core Social History Elements" list.
-   Systematic HEADSS Review: After addressing the core elements, conduct a second pass, organizing the extracted information under the appropriate HEADSS component.
-   Targeted Secondary Scan: Only after completing the above steps, perform a focused third pass, specifically seeking out recurring and clinically significant social themes not covered by the core elements or HEADSS.
-   Quote Prioritization: Actively identify and extract direct patient quotes during all passes, noting their relevance to specific categories (core elements, HEADSS, or adaptive).
-   Synthesize and Organize: Group related information and quotes under the appropriate category.
-   Clinical Judgment in Categorization: Continuously apply clinical reasoning to determine the significance of social information and the appropriateness of creating new categories.

**Format for Social History Output (Organized by HEADSS and Adaptive Categories):**

-   Home and Environment: [Detailed description with contextual evidence and embedded quotes where possible]
-   Education, Employment, Eating: [Detailed description with contextual evidence and embedded quotes where possible]
-   Activities: [Detailed description with contextual evidence and embedded quotes where possible]
-   Drugs: [Detailed description with contextual evidence and embedded quotes where possible]
-   Sexuality: [Detailed description with contextual evidence and embedded quotes where possible (if relevant)]
-   Suicide/Depression: [Detailed description of relevant history and current status, with contextual evidence and embedded quotes where possible]
-   Family/Relationships: [Detailed description with contextual evidence and embedded quotes where possible]
-   Cultural Background: [Detailed description with contextual evidence and embedded quotes where possible]
-   Financial Stress: [Detailed description with contextual evidence and embedded quotes where possible] *(Example Adaptive Category)*
-   Legal Issues: [Detailed description with contextual evidence and embedded quotes where possible] *(Example Adaptive Category)*
-   Relationships and Social Support: [Detailed description with contextual evidence and embedded quotes where possible] *(Example Adaptive Category)*
-   Substance Use History: [Detailed description with contextual evidence and embedded quotes where possible] *(Example Adaptive Category)*
-   ... (and so on for any other clinically relevant categories)

---

7.  Family History Section:
(Omit heading if no relevant family history discussed. Adhere to CRITICAL OMISSION RULE.)
-   IMPORTANT: ONLY include family history details EXPLICITLY mentioned within THIS transcript. Do NOT include information from past transcripts, general medical knowledge, or patterns inferred from other patient encounters.
-   CONTENT: If, and ONLY if, a family member's condition is directly stated in this transcript, include it using the following format: '- Family Member: Condition (Quote, if available)'.
-   When possible, use direct quotes from the transcript to indicate the source of the family history information (e.g., '- Mother: 'She has high blood pressure').
-   Specify member and condition/context (e.g., "- Mother: History of depression"; "- Parents: History of drug use ('My mom and my daddy did do drugs')").

FINAL REVIEW STEP (Mental Check Before Outputting):
1.  Omission Check: No placeholders? Empty sections/details COMPLETELY omitted? It's OKAY for some sections to be brief or absent if the transcript lacks relevant information.
2.  HPI Check: Narrative is 3rd person, synthesized, flows well? Includes mandatory details? Quotes used effectively?
3.  Med List & Metadata Check: Are ALL mentioned/implied meds/supplements listed (psych and non-psych, current & recent D/C)? Is metadata accurate? Is purpose ONLY included if stated (NO inference)? Adherence/Status clear?
4.  History Sections Check: MHx includes inferred conditions based on meds/context? Is Surgical History included if mentioned? SHx pertinent/concise with specifics/quotes reflecting style guidance? FHx present if mentioned?
5.  POV/Format Check: Third-person? Plain text? Starts directly with HPI? Correct bulleting?

REMEMBER CORE REQUIREMENTS: Plain text. Omit empty sections/details silently. Maximize detail/quotes. Include all specified med metadata (except inferred purpose). Structure includes Medical, Surgical, Meds, Social, Family History sections if applicable. Ensure ALL mentioned current AND relevant discontinued/implied medications/supplements are listed. Use third-person POV. Apply inference rules for Medical History. Capture specific social details concisely, including the Social History section if any relevant details are present. Capture Surgical History if mentioned.
`;