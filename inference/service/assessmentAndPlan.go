package inferenceService

const assessmentAndPlanTaskDescription = `
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

