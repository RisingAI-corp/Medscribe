package inferenceService

const summaryTaskDescription = `
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
