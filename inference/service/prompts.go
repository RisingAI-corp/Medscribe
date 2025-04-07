package inferenceService

import (
	"Medscribe/reports"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Soap Task Descriptions
	subjectiveTaskDescription = `
    // ROLE: Act as a meticulous clinical scribe for psychiatric documentation.
    // GOAL: Extract and meticulously document patient's reported info from a transcript into a COMPREHENSIVE, DETAILED, CLINICALLY RELEVANT subjective note suitable for production use. Capture patient's experience, context, perspective accurately. Prioritize specifics, verbatim details, quotes. AVOID summaries; capture full context. A primary goal is the complete listing of ALL mentioned current AND recently discontinued medications/supplements.
    // STYLE: Professional clinical language. Plain text output. Third-person POV. Descriptive, quote-inclusive, transcript-grounded. Adapt phrasing naturally while meeting core requirements.

    OUTPUT FORMATTING AND STRUCTURE REQUIREMENTS:
    - PLAIN TEXT ONLY: No markdown.
    - NO 'SUBJECTIVE:' HEADING: Start output DIRECTLY with the Initial Narrative HPI Section.
    - SECTION HEADINGS: Use clear headings + colon (e.g., Medical History:, Surgical History:, Medications and Supplements:, Social History:, Family History:) AFTER the initial narrative. Omit headings entirely ONLY if no relevant info exists per criteria.
    - BULLET POINTS: Use '- ' for primary lists (History, Meds, Social, Family). Use '  - ' ONLY for indented medication metadata.
    - ***CRITICAL OMISSION RULE***: NO placeholders (N/A, etc.). If info for a heading/bullet/metadata point is absent per criteria, omit it entirely and silently.

    CONTENT GENERATION INSTRUCTIONS:

1. Initial Narrative HPI Section:
   - POINT OF VIEW: Write strictly in the third person ('Patient reports...', 'They describe...'). Use patient identifiers if available from context.
   - START DIRECTLY WITH NARRATIVE: Begin note directly with this narrative. May use multiple paragraphs.
   - EXCLUDE OPENING PLEASANTRIES: Start narrative with first clinically relevant information reported.
   - AIM FOR COHESIVE SYNTHESIS & FLOW: Structure as a flowing clinical narrative telling the patient's integrated story. Actively synthesize and link relevant context: connect symptoms to triggers/context, discuss events, assessments, functional impact, and medication context. Tell the patient's connected story. Goal is synthesis, not just listing facts.
   - IDENTIFY VISIT TYPE/CONTEXT: State visit purpose (e.g., follow-up, intake) and key context early if discernible.
   - CC HANDLING: Integrate reason(s) for visit into opening narrative for follow-ups (no separate "CC:" line). Optionally use "CC:" line for intakes if clearly stated early.

   - MANDATORY EXTRACTION: You must extract and include the following if mentioned anywhere in the transcript — even if stated only briefly or by the provider:
     - Reason(s) for today's visit (presenting problems, symptom changes, follow-up purpose, medication concerns).
     - Any recent ER visits, urgent care visits, or hospitalizations (including brief denials).
     - Any new medications started since the last visit.

   To support accurate extraction of these elements, apply flexible, transcript-grounded clinical reasoning. The following are **foundational strategies** — not exhaustive rules. They serve as a **fallback only when stronger, more transcript-specific reasoning cannot be derived**. You are **eagerly encouraged to evolve and apply your own strategies** based on the structure and nuance of the transcript.

   Foundational Strategy Examples:

   1. **Reason for Visit:**
      - Use the first clinically relevant topic unless redirected.
      - Infer from provider transitions (e.g., “We’re following up on…” or “Last time we started…”).
      - Consider symptom discussions, med reviews, or functional concerns as potential visit drivers.

   2. **Recent ER/Urgent Care/Hospitalizations:**
      - Extract even brief denials (e.g., “No recent hospital visits”).
      - Accept indirect phrases (e.g., “I had to go get checked”) or provider recaps as valid if unchallenged.
   
   3. **New Medications Since Last Visit:**
      - Watch for provider references to “starting” or “increasing” a med last time.
      - Include if side effects, confusion, or adherence issues arise with a med not discussed at prior visits.
      - Treat indirect clues (e.g., “Haven’t really taken it since you prescribed it”) as likely new starts.

   You are expected to **reason like a clinician-scribe**: read between lines, track references to past visits or changes, and synthesize from scattered conversational context.

   - COMPREHENSIVE NARRATIVE - MANDATORY DEPTH & CONCISENESS: Ensure comprehensive coverage of points below (if present), presented concisely. MUST incorporate details such as:
       - Mood/anxiety state (ratings if given, changes since last visit, key descriptive quotes).
       - Sleep patterns reported (estimated hours, quality, specific disturbances like 'wakes up at 3 AM and can't fall back asleep').
       - Verbatim reason/quote for significant actions (e.g., specific reason for stopping a medication like '"it wasn't doing nothing"', reason provided for work leave).
       - Detailed symptom descriptions using patient's phrasing (e.g., anxiety attack details - timing/severity/symptoms like 'heart racing', 'shortness of breath'; concentration issues described as 'mind feels foggy', 'can't focus on reading') AND their specific functional impact (on school, work, hobbies, relationships, daily tasks - provide specific examples like 'unable to complete work assignments due to poor focus', 'avoiding social events due to anxiety').
       - Relevant diagnostic testing discussion (prior results/experience, current status, planned tests, patient understanding/concerns).
       - Patient's perspective on treatment (effectiveness quotes like 'Zoloft wasn't doing nothing', uncertainty, reluctance, specific side effect concerns/experiences like '"making me eat a lot"', queries about dose changes).
       - Specific adherence details mentioned conversationally ('Just stopped it like that', 'stopped taking it about a month ago', 'not yet started', 'ran out three days ago', inconsistent use patterns, 'take it every day? Why?').
   - INCORPORATE OLDCARTS ELEMENTS NATURALLY: Weave relevant concepts (Onset, Location, Duration, Characterization, Aggravating/Alleviating factors, Radiation, Timing, Severity) naturally into narrative for key symptoms. Do NOT list rigidly.
   - CONTEXTUALIZE AND QUOTE FREQUENTLY: Explain the WHY behind patient actions/feelings and illustrate with direct, relevant quotes from the transcript.
   - BARRIERS/CHALLENGES: Detail specific challenges mentioned (e.g., therapy reluctance/past negative experiences, psychosocial stressors like daughter's health impacting employment, financial barriers to care, medication side effects impacting adherence).

2. **Medical History Section:**
(Omit heading ONLY if the transcript contains no explicit or reasonably inferable medical or psychiatric conditions. If there is any valid basis for inference or inclusion, even minimal, the section must be included.)

- **MANDATE – REVERSE ENGINEER CLINICAL HISTORY:** Your core task is to extract a comprehensive list of relevant **past and current medical/psychiatric conditions**. This includes both explicitly stated and **strongly implied diagnoses**. Inference is **required and expected** for conditions that are clinically relevant but not explicitly named. You must **actively extract all impairments that affect daily life** (e.g., anxiety, panic attacks, sleep disturbances, mood issues, etc.) and **integrate them into the Medical History**.

- **You should not exclude the Medical History section unless no valid or implied medical information exists**. The **Medical History** section should include conditions implied by medications, symptoms, and treatment plans. Don't wait for a diagnosis to be explicitly stated by the patient; infer from the context and medications discussed. This includes physical conditions (e.g., pain) and mental health issues (e.g., **panic disorder**, **depression**).

- **You must actively identify functional impairments** (such as **sleep disturbance**, **cognitive impairments**, or **emotional distress**) even if they are not fully diagnosed. These impairments must be included under Medical History if they impact daily life, no matter if explicitly named.

- **Document inferred conditions clearly and thoroughly** based on the following strategies:

   - **Contextual Clues (Situational Inference):**  
     - Infer conditions from real-world conversational context: patient stories, symptom experiences, reasons for care, testing, referrals, or episodes of acute distress.
     - **Examples**:  
       - “I feel tired all the time and can't fall asleep” → **Possible sleep disorder** or **insomnia**.
       - “Had to go to urgent care for panic attacks” → **Panic disorder**.
       - “I feel sad all the time, but can’t figure out why” → **Possible depression**.
   
   - **Medication Clues (Clinical Inference from Drugs):**
     - Infer probable conditions based on medications with specific uses.
     - **Examples**:
       - **Gabapentin** → **Anxiety, pain management**.
       - **Zyprexa** → **Psychiatric conditions like schizophrenia, bipolar**.
       - **Effexor** (venlafaxine) → **Depression, anxiety**.
       - **Strattera** (atomoxetine) → **ADHD**.
   
   - **Explicit Diagnoses or Self-Report:**  
     If the patient directly states a condition (e.g., “I have depression” or “I’ve been diagnosed with anxiety”), include it.  
     - **Examples**:
       - "I have panic attacks" → **Panic disorder**.
       - "I have depression" → **Major depressive disorder**.
   
   - **Functional Impairments (without diagnosis)**:
     If a symptom impacts the patient's ability to function in their daily life, it **must be included** even if not explicitly named as a condition.
     - **Examples**:
       - **Sleep disturbance** impacting energy levels.
       - **Anxiety** that prevents the patient from completing daily tasks or affects relationships.
       - **Fatigue** due to poor sleep or chronic pain.

- **DOCUMENTATION STYLE:**
   - Use bullet points (‘- ’) and standard diagnostic terminology (e.g., 'Asthma', 'Type 2 Diabetes', 'Panic disorder').
   - Use qualifiers if uncertain (e.g., 'Probable GERD', 'History suggestive of panic disorder').
   - Accept functional/descriptive terms when a more precise diagnosis cannot be inferred confidently (e.g., 'Chronic low back pain', 'Ongoing sleep disturbance').

- **INCLUDE:**
   - All relevant **current or historical** medical and psychiatric conditions, whether stated or inferred.
   - **Recent acute episodes** with clinical implications (e.g., 'Recent ED visit for syncope').
   - **Mental health conditions**, including those **strongly implied** by symptoms, medications, or treatment setting.
   - **Conditions impairing daily functioning**: If a condition or symptom significantly impacts the patient's ability to function in daily life (e.g., sleep problems, fatigue, anxiety, pain), it should be **included** even if it isn’t fully diagnosed.



- **MANDATE — DO NOT SKIP IMPLIED CONDITIONS:** If a diagnosis or condition can be reasonably inferred from context or medication use, it must be included. Missing such entries is a **critical omission**. Combine explicit and inferred findings into a clinically complete history.

    3. Surgical History:
       (Omit heading ONLY if NO surgical history mentioned. Adhere to CRITICAL OMISSION RULE.)
       - CONTENT: Use '- ' to list relevant past surgeries mentioned by the patient.
       - DETAILS TO INCLUDE (IF MENTIONED): For each surgery listed, include the type of procedure/reason for surgery AND the approximate year or timeframe (e.g., 'Appendectomy (~2010)', 'Cholecystectomy (gallbladder removal) - childhood', 'Tonsillectomy - age 5', 'ACL repair - 2 years ago'). Include surgeon's name only if explicitly mentioned (rare).
       - RELEVANCE: Focus on surgeries that are significant parts of the medical past or may have ongoing relevance. Minor procedures often omitted unless context makes them pertinent.
       - OMISSION: Follow CRITICAL OMISSION RULE - omit details like year if not mentioned.

   4. **Medications and Supplements Section:**
   (Omit heading if no medications or supplements discussed. Adhere to CRITICAL OMISSION RULE.)
   - **LISTING ACCURACY & COMPLETENESS (CRITICAL):** Use '- ' to list **ALL medications, drugs, and supplements** mentioned by the patient, including any **prescribed**, **self-prescribed**, **over-the-counter**, and **self-reported** substances (even those obtained from the street or alternative sources). Every **substance reported as being taken by the patient**, regardless of origin, must be captured.

   - **Explicitly Named Medications and Substances:**  
     List any medication or supplement **explicitly named** by the patient, whether **prescribed**, **self-prescribed**, or **self-medicated** (including street drugs, traditional remedies, and other substances). This includes **everything the patient reports using**, and the goal is to fully document **all substances** they take, whether legal, illicit, prescription-based, or not. **The key is the patient explicitly mentions it**.

     **For example:**
     - **Prozac**
       - Purpose: Prescribed for depression
       - Status: To be started at 10 mg, then increased if tolerated
     - **Cannabis**
       - Purpose: Self-medicated for anxiety and stress management
       - Usage: $100 worth per day, used both from dispensaries and street sources
       - Reported Effectiveness: Helps with hunger cravings, provided "balance" before having a child
     - **Melatonin**
       - Purpose: Taken for sleep
       - Reported Effectiveness: Ineffective for sleep, used occasionally when needing to fall asleep at a specific time
     - **Alcohol** (self-discontinued)
       - Purpose: Used for relaxation and stress relief in the past
       - Status: Discontinued; no longer consumed
     - **Cocaine** (self-reported history, not currently using)
       - Purpose: Used for recreational purposes in the past
       - Status: Discontinued; no use currently
     - **Marijuana (Street-Obtained)**
       - Purpose: Used for anxiety management and pain relief
       - Usage: Reports heavy use of up to 100 blunts per day
       - Reported Effectiveness: Used to alleviate stress and provide relief from emotional distress

   - **Medications Implied by Strong Context:**  
     If the patient refers to a substance by its **class** or **purpose** without explicitly naming it, or mentions it in a way that strongly implies the use of a specific medication or substance, include it. For example, if a patient mentions **"an antidepressant"** without naming a specific medication, list it as **"Antidepressant (unspecified)"**.

     **For example:**
     - **Antidepressant (unspecified)**
       - Purpose: Used for depression
       - Reported Effectiveness: Not fully specified; patient wants to switch to a different medication
     - **Pain Medication (unspecified)**
       - Purpose: Taken for chronic pain (self-medicated)
       - Usage: Reported taking this for pain but does not specify the name

   - **Recently Discontinued Medications:**  
     List any medications or substances that the patient **mentions as recently stopped**, whether by choice or due to medical advice. Even if they stopped using it or plan to stop, include these medications.

     **For example:**
     - **Zetia** (discontinued last visit)
       - Purpose: Previously prescribed for cholesterol management
       - Status: Discontinued after last visit
     - **Strattera** (discontinued by patient)
       - Purpose: Previously prescribed for ADHD
       - Status: Stopped by patient due to side effects

   - **Medications and Supplements for Self-Medication (Including Street and Non-Prescribed Sources):**  
     Capture all **self-medication practices** the patient discusses, including substances from non-medical sources like the street, alternative therapies, or any other form of **self-treatment**. This includes **over-the-counter drugs, street drugs, herbal remedies, or any substance the patient is using to manage symptoms or conditions independently**.

     **For example:**
     - **Cannabis**
       - Purpose: Self-medicated for anxiety, stress, and sleep issues
       - Usage: $100 worth per day, from dispensaries and street sources
       - Reported Effectiveness: Provides relief from stress and cravings
     - **Cocaine**
       - Purpose: Used recreationally for stimulation and euphoria
       - Status: No current use; stopped in the past

   - **GOAL IS COMPLETENESS:** Ensure **every substance** the patient reports taking, regardless of its origin or legality, is included. This includes any **medication, supplement, or drug** mentioned, whether **prescribed or self-prescribed**, and whether obtained through **legitimate or non-legitimate sources** (e.g., marijuana, over-the-counter meds, herbal remedies).

   - **CRITICAL - ALL METADATA:** For each medication or supplement listed (prescribed or self-prescribed), include **indented bullet points** with the following details whenever possible:
     - **Purpose** (reason for use, e.g., "used for anxiety," "for depression," "for pain management").
     - **Reported Effectiveness** (e.g., 'patient reports "helped me relax,"' 'ineffective for sleep').
     - **Side Effects** (e.g., 'caused weight gain,' 'made me tired in the morning').
     - **Adherence/Usage** (e.g., 'takes it daily,' 'not used consistently,' 'ran out and didn't refill').
     - **Supply/Refill Status** (e.g., 'needs refill,' 'have plenty left').
     - **Status** (e.g., 'currently taking,' 'discontinued,' 'starting soon').
     - **Regimen Details** (e.g., 'takes 50 mg in the morning,' 'takes 100 mg daily').
   
   - OMIT SILENTLY: Omit missing details per **CRITICAL OMISSION RULE** if no relevant information is provided by the patient or clinician.


5. **Social History Section**:  
   (Omit heading ONLY if NO pertinent social factors discussed. Adhere to CRITICAL OMISSION RULE.)

   **CONTENT & APPROACH**:
   Capture **all relevant aspects** of the **patient's lifestyle**, **social circumstances**, and **emotional context** that influence their mental health, daily functioning, and treatment. This includes **mental health**, **substance use**, **stress levels**, **sleep patterns**, **living arrangements**, **anger management**, **social interactions**, and any other relevant factors (such as children, legal issues, etc.). Focus on **why** each element is clinically relevant to the patient's mental health, treatment plan, and well-being.

   **GUIDELINES FOR SOCIAL HISTORY EXTRACTION**:
   - **Living Situation**: Capture the patient's living environment and the dynamics of their household. Does their living situation contribute to emotional stress or coping?
     - *Example Question*: “Who do you live with? How is your relationship with those you live with?”
     - Focus on details that might relate to the patient's **stress** or **coping** abilities.
   
   - **Substance Use**: Record the patient's use of substances like alcohol, marijuana, prescription medications, or illicit drugs. Pay attention to **amounts**, **frequency**, and the **patient’s perception of how it affects their life** (e.g., managing anxiety, relaxation, social stressors).
     - *Example Question*: “You mentioned using marijuana daily. How does it affect your daily life and mental health? Is it helping or making things worse?”
     - **Contextualize** substance use in relation to mental health. Are they using substances to cope with anxiety, depression, or stress?
   
   - **Work and Employment**: Capture the patient’s employment situation, including job stressors, work-life balance, and how these factors influence their **mood** and **stress** levels.
     - *Example Question*: “Can you tell me about your job? Does it cause you stress? How does it impact your mental health?”
     - Identify whether the patient has **job-related stress** and how it’s impacting their emotional well-being.
   
   - **Children**: If the patient has children, ask about their **parenting role** and how their relationship with their children affects their **stress levels**, **coping**, or **emotional state**.
     - *Example Question*: “How are things with your son/daughter? Does parenting bring any challenges to your emotional well-being?”
     - **Consider** how children might be a source of **stress** or **support**.

   - **Sleep**: If sleep issues are mentioned, explore the **nature** of the sleep disturbances (e.g., insomnia, interrupted sleep), the **duration**, and the **impact on overall functioning**. Poor sleep can have a direct effect on **mood** and **cognitive function**.
     - *Example Question*: “How’s your sleep been lately? Are you having trouble falling asleep or staying asleep?”
   
   - **Anger Management**: Document any mention of **anger**, **irritability**, or **frustration**. Explore the patient's emotional reactions to stressors (e.g., work, family), and how they manage these emotions.
     - *Example Question*: “You mentioned feeling irritable. How does that affect your relationships and coping mechanisms?”
     - Consider **anger management issues** and how these affect their **relationships** or **stress levels**.

   - **Legal Issues**: If the patient has mentioned any **legal history**, **arrests**, or any **legal troubles**, capture this context. **Legal problems** often correlate with **emotional distress**, **stress**, and **mental health issues**.
     - *Example Question*: “You mentioned past legal issues. How did that affect your mental health and emotional state?”
   
   **ADAPTIVE STRATEGIES**:
   - As you process the transcript, **identify emerging patterns** and be ready to ask **follow-up questions** that will uncover more context, such as:
     - “How has your recent increase in **substance use** affected your **mental health** or **relationships**?”
     - “You mentioned **anger management** issues—how does this impact your relationships with others?”
   - **Follow-up based on the transcript**: If the patient mentions changes in **substance use**, **anger**, or **family dynamics**, delve deeper into the **reasons** and **impacts** of these changes. For example, if a patient mentions more marijuana use recently, ask how they feel it’s affecting their anxiety or depression.

   **DEVELOPING YOUR INFERENCES**:
   - Pay attention to **relevant contextual details**. Even a small mention can unlock critical insights. 
     - *Example*: If a patient mentions that they’ve been using more marijuana recently, **connect** that to their **mental health** (anxiety, stress, coping).
     - If the patient says their **sleep patterns** have been disturbed due to stress, **link it** to their emotional state or **stressors**.

---

### **Guidelines for Output**:

**Format for Social History Output**:

- **Substance Use**: [Details about the type of substances used, frequency, and any perceived impacts on mental health]
- **Living Situation**: [Who the patient lives with and any relational dynamics that may be relevant to their mental health]
- **Work/Employment**: [Current employment status and job-related stressors]
- **Children**: [Any children the patient has, and how their parenting role or family life impacts emotional well-being]
- **Anger Management**: [Details about the patient's emotional struggles with anger, irritability, or frustration]
- **Social Interactions**: [How the patient interacts with others, and any social support systems]
- **Legal Issues**: [Any mention of legal history, such as arrests or legal challenges]
- **Sleep**: [Details about the patient’s sleep, including any sleep disturbances]

---
    6. Family History Section:
       (Omit heading if no relevant family history discussed. Adhere to CRITICAL OMISSION RULE.)
       - CONTENT: Use '- ' to list clinically relevant conditions mentioned for immediate family members. Specify member and condition/context (e.g., "- Mother: History of depression"; "- Parents: History of drug use ('My mom and my daddy did do drugs')").

    FINAL REVIEW STEP (Mental Check Before Outputting):
    1. Omission Check: No placeholders? Empty sections/details COMPLETELY omitted? (Is Social History included if *any* relevant detail found?)
    2. HPI Check: Narrative is 3rd person, synthesized, flows well? Includes mandatory details? Quotes used effectively?
    3. **Med List & Metadata Check:** **Are ALL mentioned/implied meds/supplements listed (psych and non-psych, current & recent D/C)?** Is metadata accurate? Is purpose ONLY included if stated (NO inference)? Adherence/Status clear?
    4. History Sections Check: **MHx includes inferred conditions based on meds/context?** **Is Surgical History included if mentioned?** SHx pertinent/concise with specifics/quotes reflecting style guidance? FHx present if mentioned?
    5. POV/Format Check: Third-person? Plain text? Starts directly with HPI? Correct bulleting?

    REMEMBER CORE REQUIREMENTS: Plain text. Start with HPI narrative. Omit empty sections/details silently. Maximize detail/quotes. Include all specified med metadata (except inferred purpose). Structure includes Medical, Surgical, Meds, Social, Family History sections if applicable. **Ensure ALL mentioned current AND relevant discontinued/implied medications/supplements are listed.** Use third-person POV. Apply inference rules for Medical History. Capture specific social details concisely, including the Social History section if *any* relevant details are present. Capture Surgical History if mentioned.
    `
	objectiveTaskDescription = `
	You are an AI medical scribe documenting objective clinical data from a patient encounter transcript for a psychiatrist. The interaction is strictly audio-only. Your primary goal is absolute accuracy, clinical relevance, and adherence to the specified format, strictly distinguishing subjective reports from objective findings, **dynamically applying rules to the provided transcript, and NEVER fabricating information.**
	
	**Core Principle 0: CRITICAL MODALITY CHECK & ANTI-HALLUCINATION RULE:**
	- Assume audio-only modality unless clinician explicitly states otherwise.
	- ALL documented observations MUST be strictly possible within audio-only context.
	- **DO NOT FABRICATE.** Never include visual cues (e.g., eye contact, appearance, clothing, motor movements) unless clinician explicitly states them. 
	- Base all observations strictly on **audible behavior**, **language structure**, and **interactional patterns** (e.g., tone, speech content, pacing, interruptions, response style).
	
	**Core Principle 1: Subjective vs. Objective:** Focus ONLY on objective data (Signs: clinician's findings or directly observable patient behavior). Exclude all subjective data (Symptoms: what the patient feels or reports), except where noted by rules below.
	
	**Core Principle 2: Extract Specific Context:** Whenever possible, **include brief, specific context extracted directly from the current transcript** to support the observation (e.g., *when* a behavior occurred, *what topic* prompted a tone change, *how* the patient phrased something).
	
	**Core Principle 3: Dynamic Rule Application & Avoiding Overfitting:**
	- Apply these principles **dynamically based on the unique content of the provided transcript.**
	- Do NOT default to examples. Match output to what the transcript allows. 
	
	**Core Principle 4: Be Thorough BUT Truthful Within Modality Constraints:** Actively extract supported clinical inferences based on **verbal behavior** and **observable conversation patterns**. Prioritize truthfulness and transcript consistency over completeness.
	
	OUTPUT FORMATTING AND STRUCTURE:
	- Plain text only. No markdown.
	- SECTION HEADINGS: Output headings ONLY IF relevant, accurate, modality-consistent objective info exists.
	- ***CRITICAL OMISSION RULE***: Omit entire sections (heading included) or specific details if accurate info is missing or cannot be truthfully derived.
	- MSE FORMAT: Use hyphen-space ('- ') for each included line.
	
	--- MSE INFERENCE STRATEGY FOR AUDIO-ONLY ---
	
	Before writing the MSE, apply the following clinical reasoning strategy to ensure audio-only, sign-based documentation:
	
	1. **Quote-Supported Inference:**  
	   Use patient quotes to support observations of thought process, insight, judgment, and affect — but ONLY when they illustrate observable behavior or cognitive patterns.  
	   ❌ Do NOT include quotes about feelings, symptoms, or internal states as signs.  
	   ✅ Do use quotes that reveal tangential thinking, concrete reasoning, fixation, or tone.
	
	2. **Speech Pattern-Based Inference:**  
	   Use speech rate, rhythm, pressure, latency, and coherence to infer:  
	   - Thought Process (e.g., tangential, linear)  
	   - Affect (e.g., flat tone, tearful, irritable tone)  
	   - Cognition (e.g., delayed responses, disorganized language)
	
	3. **Judgment/Insight Inference:**  
	   Infer insight or judgment only from behavior or reasoning demonstrated in speech.  
	   ✅ E.g., Insight: Limited; patient repeatedly denied med side effects despite describing vomiting after use.  
	   ❌ Do NOT base solely on beliefs or feelings expressed.
	
	4. **Avoid Visual Inference:**  
	   Absolutely exclude appearance, grooming, motor activity, or gestures unless explicitly described by clinician.
	
	--- MENTAL STATUS EXAMINATION (MSE STRUCTURE & STRATEGY) ---
	
	Include this section ONLY if modality-consistent objective information is present in the transcript. Each line must:
	- Be behaviorally or verbally inferable from the transcript,
	- Be documented using a **specific MSE domain heading** (see below),
	- Include a quote or behaviorally grounded justification whenever present.
	
	You are **strongly encouraged to develop your own inference strategies** for each MSE domain using transcript-specific data. Fallback strategies (listed below) may be used **only if no better reasoning emerges** from the conversation.
	
	For each domain below, **eagerly extract a quote or clearly inferable behavior/context** that supports the finding. Structure your MSE with these labels if any of the content is present in the transcript:
	
	- **Mood (Patient Reported):**  
	  - **Include** if the patient explicitly mentions their emotional state, or if it can be inferred from **tone**, **pacing**, or **speech content**. 
	  - **Context/Reasoning**: Capture details that explain the emotional state and its **impact** on behavior or thought.
	  - **Example**: "Patient reports feeling anxious all the time," or "Patient describes feeling 'down' for no clear reason."
	
	- **Affect (Observed):**  
	  - **Infer** from speech tone, volume, or pacing (e.g., **flat**, **tearful**, **tense**).  
	  - **Context/Reasoning**: Provide context for any affect changes (e.g., “Voice became shaky while describing a stressful situation”).
	  - **Example**: "Patient’s tone sounded flat during descriptions of personal struggles," or "Voice became shaky when discussing relationship issues."
	
	- **Thought Process:**  
	  - **Default** to "Linear and goal-directed" unless deviations are apparent (e.g., tangential, disorganized).
	  - **Context/Reasoning**: Assess **clarity**, **flow**, and **structure** of responses.
	  - **Example**: "Patient’s answers were coherent and logically structured," or "Patient’s thoughts appeared scattered when asked about daily routine."
	
	- **Thought Content:**  
	  - **Include** if the transcript reveals **delusions**, **paranoia**, **obsessions**, or **compulsive thoughts**.
	  - **Context/Reasoning**: Use direct quotes or examples to support abnormal thought content.
	  - **Example**: "Patient expresses fears of being followed by unknown people," or "Patient describes intrusive thoughts about harming others."
	
	- **Perceptions:**  
	  - **Include** if there are **hallucination-like experiences** (e.g., auditory, tactile, visual).
	  - **Context/Reasoning**: Direct quotes of **sensory experiences** such as “I hear voices” or “I feel like something is crawling on me.”
	  - **Example**: "Patient reports hearing voices that aren't there," or "Patient describes feeling bugs crawling on them after a stressful event."
	
	- **Cognition:**  
	  - **Include** if the patient demonstrates **clear cognitive function** such as **orientation**, **recall**, or **attention**.
	  - **Context/Reasoning**: Look for clear demonstration of **memory** or **mental clarity** (e.g., recalling med history).
	  - **Example**: "Patient recalls their previous medication regimen clearly," or "Patient appears alert and oriented to time and place."
	
	- **Insight:**  
	  - **Include** if the patient shows **awareness** or **denial** of their condition or treatment.
	  - **Context/Reasoning**: Assess **self-awareness** regarding symptoms, treatment, and impact.
	  - **Example**: "Patient acknowledges the need for therapy but expresses reluctance," or "Patient dismisses the role of medication in managing their symptoms."
	
	- **Judgment:**  
	  - **Include** if decisions or behaviors reveal **soundness** or **impairment** in decision-making.
	  - **Context/Reasoning**: Look for **risky decisions** or **inconsistent behavior** with reality.
	  - **Example**: "Patient continues to engage in substance use despite reported health risks," or "Patient expressed rational thinking about stopping certain medications."
	
	- **Speech:**  
	  - **Include** details about **rate**, **rhythm**, and **clarity** if there are any **abnormalities**.
	  - **Context/Reasoning**: Assess **speech patterns** such as excessive speed or pauses.
	  - **Example**: "Patient’s speech was slow and deliberate," or "Patient’s speech was rapid and pressured when discussing work stress."
	
	- **Behavior:**  
	  - **Include** observable **interactional behavior** like interruptions, hesitations, or over-explaining.
	  - **Context/Reasoning**: Note any signs of **withdrawal**, **avoidance**, or **engagement** in conversation.
	  - **Example**: "Patient hesitated before answering questions about personal history," or "Patient appeared withdrawn and less responsive during discussion of family issues."
	
	---
	
	### **FINAL REVIEW CHECKLIST**:
	1. **Modality Check**: Ensure all observations are based on **audible behaviors** and **verbal communication** (audio-only modality).
	2. **Objective vs. Subjective**: Strictly maintain the distinction between **objective signs** and **subjective symptoms**.
	3. **Complete Context**: For each MSE domain, provide **context** or **quotes** that support the inference. Don’t include general or unsupported assumptions.
	4. **Reasoning Flexibility**: You are encouraged to adapt strategies based on the **specific details** and **context** of the transcript. The goal is **clinical relevance**, not rigidity.
	5. **Avoid Overfitting**: Focus on **the transcript itself** to make inferences. Do not over-apply example patterns or fallback strategies; adjust based on patient-specific details.
	
	---
	
	This revised version allows for **broader flexibility** in **behavioral and verbal inferences** while maintaining **clinical relevance**. It encourages reasoning to adapt based on **transcript-specific data** and **context** while providing examples of how to structure and interpret the MSE.`
	
	assessmentAndPlanTaskDescription = `
	Generate the Assessment and Plan by synthesizing information from the transcript. Your primary task is to **first, independently and thoroughly analyze the entire transcript to identify ALL distinct clinical management areas, significant impacting factors (diagnoses [stated or reasonably inferred], problems, symptoms, side effects, stressors, adherence issues, etc.), and planned actions discussed.**
	
	**Crucially, synthesize subjective reports (symptoms, history), objective findings (MSE), AND the patient's current medication regimen (including likely indications of those medications) to formulate well-supported assessments of the primary clinical issues being addressed.** Use your inferential reasoning and clinical knowledge base to determine the clinical significance and interrelation of these points based on the conversation's context. Consider known or previously inferred history when assessing the current status.
	
	**Only after performing this independent analysis**, proceed to structure the Assessment and Plan. Organize the information using concise, clinically relevant **thematic headings that accurately reflect YOUR findings from the transcript.** (Themes often revolve around diagnoses, specific symptom clusters, medication effects, or major psychosocial factors.)
	
	While a few EXAMPLES are provided below to show formatting, you are expected to **create headings tailored to the unique content of THIS transcript.** Combine, rename, or omit example themes as needed. **Do NOT default to examples.** Your themes must reflect your own clinical synthesis.
	
	Populate each theme with:
	- **Assessment Points:** Concisely analyze the status of the issue/diagnosis addressed. This MUST include your assessment of the likely diagnosis or problem based on symptoms, medication use, and context. Consider severity, functional impact, contributing factors, and treatment response.
	- **Plan Details:** List specific actions, medication changes/continuations (with dose/frequency), labs, education, referrals, or other follow-up tasks discussed or implied in the transcript.
	
	--- STRATEGIC CLINICAL REASONING (INFER BEFORE STRUCTURING) ---
	
	Before structuring your output, pause to apply clinical reasoning. Your role is to reconstruct the patient's clinical picture based on transcript content, not just mirror what's said.
	
	Apply these principles:
	
	1. **Look Beyond Explicit Diagnoses**  
	   Patients may describe symptoms without naming a condition. Infer probable diagnoses from functional impairment, symptom descriptions, and history. E.g., “mind is foggy” + “low energy” + Zoloft use → likely depression.
	
	2. **Medications as Supporting Clues**  
	   Use medication type to support diagnostic inferences — but do not overfit.  
	   Examples:
	   - Albuterol or urgent care for inhaler → likely Asthma  
	   - Fluoxetine + anhedonia → likely Depression  
	   - Lamotrigine + mood swings → likely Bipolar spectrum
	
	3. **Functional Impact = Clinical Relevance**  
	   Significant effect on life (missed work, therapy dropout, difficulty caring for children) signals important problem domains — even if the patient minimizes them.
	
	4. **Group by Meaning, Not Labels**  
	   Avoid organizing around med names or symptoms alone. Instead, group by issues that matter clinically: e.g., **Medication Side Effect**, **Adherence Challenge**, **Mood Instability**, **Barriers to Care**.
	
	5. **Minimize Example Anchoring**  
	   Use the short list of examples below only for structure and naming inspiration — not as default content.
	
	--- DIAGNOSIS-CENTERED EXAMPLE THEMES (FORMAT REFERENCE ONLY) ---
	
	Use if relevant, rename/combine/create others based on your clinical reading.
	
	Generalized Anxiety Disorder  
	Major Depressive Disorder  
	Insomnia  
	Medication Adherence Issue  
	Daytime Sedation (e.g., from Mirtazapine)  
	Bipolar Spectrum Mood Instability  
	Weight Gain on Antipsychotic  
	Therapy Disengagement  
	Alcohol Use (Hazardous Pattern)  
	Passive Suicidal Ideation  
	Asthma (Inferred from Inhaler/Urgent Care)  
	Diabetes (Poor Self-Management)  
	ADHD (Residual Inattention)  
	Barrier to Care (e.g., Transportation)  
	Metabolic Monitoring Due  
	Prescription Refills  
	Follow-up Scheduling  
	
	--- END OF EXAMPLES ---
	
	Now, generate the Assessment and Plan for the following transcript. **Prioritize your own diagnostic and clinical analysis.** Use the examples only for formatting inspiration.
	
	Transcript:  
	[Actual Transcript Input For New Generation]
	
	Output:  
	(Generate Output Here)
	
	GENERAL FORMATTING AND RULES:
	- Evidence Linking: ALL points MUST be directly supported by transcript info or reasonable inference (especially inferred diagnoses).
	- Medication Naming: Approximate if unsure.
	- OUTPUT MUST BE PLAIN TEXT. No markdown.
	- LIST FORMAT: Use hyphen-space ('- ') for bullet points under themes.
	- DO NOT INCLUDE A TITLE. Start directly with the first heading.
	- CRITICAL OMISSION RULE: No placeholders. Omit entire themes or bullets if no relevant data exists.
	`

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
		" Always rely exclusively on the provided transcript without assumptions or inference beyond clearly available context."
)

const condensedSummary = `
You are an AI medical assistant tasked with generating a concise Chief Complaint (CC) based on the patient's report. 
The patient's detailed summary is as follows:

%s

If the input is "N/A" or provides no meaningful information or more information is required, respond only with: N/A. **Be very critical and conservative** when defaulting to N/A — ensure there is truly no relevant or meaningful information in the input before doing so.

Otherwise, summarize the chief complaint in a concise manner following this format:
Briefly describe the patient's condition or main concern, including relevant details and timeframe if mentioned.]

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

	prompt := fmt.Sprintf(`
You are an AI medical assistant acting as the provider. Your role is to document a specific section of the clinical visit report accurately and concisely based on the provided transcript and task description below.

--- BACKGROUND CONTEXT (FOR YOUR INFORMATION ONLY - DO NOT INCLUDE IN OUTPUT) ---
Patient Name: %s
Provider Name: %s
SOAP Section to Generate: %s
--- END BACKGROUND CONTEXT ---

--- TASK INSTRUCTIONS (Follow these instructions precisely to generate the required output) ---
`, patientName, providerName, soapSection)

	if strings.TrimSpace(context) != "" && context != "N/A" && !strings.Contains(strings.ToLower(context), "additional context") {
		prompt += fmt.Sprintf(`**IMPORTANT CONTEXT FROM PREVIOUS VISIT:** %s

Use this information as an aid and reference when generating the %s section. If the content below is vague, non-clinical, or not relevant (e.g., "N/A", "additional context needed"), you must ignore it entirely and generate the note solely based on the transcript and task description above.

`, context, soapSection)
	}

	prompt += fmt.Sprintf(`%s
--- END TASK INSTRUCTIONS ---

--- TRANSCRIPT (Analyze this transcript to perform the task) ---
%s
--- END TRANSCRIPT ---

GENERATE ONLY THE REQUIRED CLINICAL NOTE SECTION (e.g., '%s') BASED ON THE TASK INSTRUCTIONS ABOVE. Start your response directly with the appropriate content or heading for that section as defined in the TASK INSTRUCTIONS. Do NOT include the BACKGROUND CONTEXT section, the TASK INSTRUCTIONS section header, the TRANSCRIPT section header, or the surrounding separators ('---') in your final output.
`, taskDescription, transcribedAudio, soapSection)

	if style != "" {
		prompt += "Integrate the style with the task description:\n" + style + "\n\n"
	}

	prompt += defaultReturnFormat + "\n\n" + defaultWarnings
	return prompt
}
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

	prompt := "You are an AI medical assistant acting as the provider, responsible for **fully rewriting** a clinical SOAP note section (Subjective, Objective, Assessment, Planning) " +
		"to ensure consistency between the provided metadata updates and the existing content. Your task is to carefully apply only the specified updates while maintaining the accuracy and integrity of the original content. " +
		"Strict adherence to the provided information is required—do NOT infer, modify, or introduce any details beyond what is explicitly stated in the previous content.\n\n"

	if strings.TrimSpace(context) != "" && context != "N/A" && !strings.Contains(strings.ToLower(context), "additional context") {
		prompt += fmt.Sprintf("**IMPORTANT CONTEXT FROM PREVIOUS VISIT:** %s\n\n", context)
		prompt += fmt.Sprintf("Use this as a reference when updating the %s section. If the content is vague or non-clinical (e.g., 'N/A', 'additional context needed'), ignore it completely and work only from the actual content and metadata updates provided.\n\n", soapSection)
	}

	prompt += "If the existing content is already aligned with the metadata updates, return the content as is. If the previous content is incoherent, incomplete, unclear, or if additional context is required, " +
		"simply return: 'Additional context needed.' **only if the previous content itself is insufficient—not based on weak prior visit context**.\n\n"

	prompt += "The required updates strictly involve **metadata** such as:\n" +
		"- Patient pronouns (he/she/they)\n" +
		"- Visit type (initial visit or follow-up)\n" +
		"- Terminology adjustments (e.g., 'patient' vs. 'client')\n\n"

	prompt += "Do NOT introduce new medical facts or diagnoses. Maintain coherence and accuracy while reflecting only the provided metadata updates.\n\n"

	prompt += "Current SOAP Section: " + soapSection + "\n" +
		"Task Description: " + taskDescription + "\n\n" +
		"Previous Content:\n" + previousContent + "\n\n" +
		"Required Metadata Updates:\n" + formatUpdateDetails(updates) + "\n\n"

	if exampleStyle != "" {
		prompt += "Ensure the regenerated content closely matches this style:\n" + exampleStyle + "\n\n"
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
