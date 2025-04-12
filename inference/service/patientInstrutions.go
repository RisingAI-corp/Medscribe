package inferenceService

const patientInstruction = `
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