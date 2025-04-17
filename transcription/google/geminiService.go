package geminiTranscriber

import (
	transcriber "Medscribe/transcription"
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

const prompt = `Please provide a detailed diarized transcript of the following telehealth visit audio in JSON format, acting as a highly accurate medical transcriptionist.

**Context:** This is a two-way conversation between a healthcare provider and a patient during a telehealth visit. The primary language is English, and there may be variations in accents and medical terminology. Accuracy in medical terms, including medication names, dosages (e.g., milligrams, micrograms), medical conditions, anatomical terms, and proper names is crucial.

**Instructions:**

1.  **Highly Accurate Medical Transcription:** Transcribe the audio with the utmost accuracy, paying close attention to medical terminology, proper names (of people, medications, etc.), and units of measurement. Ensure correct spelling and pronunciation are reflected in the transcription.
2.  **Medical Terminology Focus:** Prioritize the correct transcription of all medical terms, including:
    * Generic and brand names of medications (e.g., lisinopril, Lipitor).
    * Units of measurement (e.g., milligrams as mg, micrograms as mcg, milliliters as mL). Spell out units when clarity is improved.
    * Medical conditions (e.g., hypertension, diabetes mellitus).
    * Anatomical terms (e.g., cardiovascular, gastrointestinal).
    * Medical procedures and tests (e.g., electrocardiogram, physical therapy).
3.  **Proper Name Accuracy:** Pay close attention to the pronunciation and spelling of patient and provider names. If a name is unclear, make your best educated guess based on phonetic sounds.
4.  **Speaker Diarization:** Clearly identify and label each speaker. Use "speaker" with values "provider" or "patient".
5.  **Timestamps:** Include timestamps for the beginning and end of each speaker's turn as floating-point numbers representing seconds under the keys "start_time" and "end_time".
6.  **Text:** The transcribed text for each turn should be under the key "text".
7.  **JSON Format:** Structure the output as a JSON array of objects, where each object represents a speaker turn.

**Audio:** (The audio data will be provided inline)
`

type geminiTranscriberStore struct {
	client *genai.Client
}

func NewGeminiTranscriberStore(client *genai.Client) transcriber.Transcription {
	return &geminiTranscriberStore{client: client}
}

func (i *geminiTranscriberStore) Transcribe(ctx context.Context, audioData []byte) (string, error) {
	if i.client == nil {
		return "", fmt.Errorf("gemini transcriber: gemini client is not initialized")
	}

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: audioData, MIMEType: "audio/wav"}},
	}

	contents := []*genai.Content{{Role: genai.RoleUser, Parts: parts}}
	resp, err := i.client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, nil)
	if err != nil {
		return "", fmt.Errorf("gemini transcriber:error generating content: %v", err)
	}
	return resp.Text(), nil
}



func (i *geminiTranscriberStore) TranscribeWithDiarization(ctx context.Context, audioData []byte) ([]transcriber.TranscriptTurn, error) {
	if i.client == nil {
		return nil, fmt.Errorf("gemini transcriber: gemini client is not initialized")
	}

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: audioData, MIMEType: "audio/wav"}},
	}
	contents := []*genai.Content{{Role: genai.RoleUser, Parts: parts}}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: "array",
			Items: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"speaker": {
						Type:        "string",
						Description: "is the person speaker patient or provider",
					},
					"start_time": {
						Type:        "number",
						Description: "The start time of the turn in seconds",
					},
					"end_time": {
						Type:        "number",
						Description: "The end time of the turn in seconds",
					},
					"text": {
						Type:        "string",
						Description: "The transcribed text of the turn",
					},
				},
				Required: []string{"speaker", "start_time", "end_time", "text"},
			},
		},
	}

	resp, err := i.client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, config)
	if err != nil {
		return nil, fmt.Errorf("gemini transcriber: error generating content: %v", err)
	}

	var transcript []transcriber.TranscriptTurn
	err = json.Unmarshal([]byte(resp.Text()), &transcript)
	if err != nil {
		return nil, fmt.Errorf("gemini transcriber: error unmarshalling response: %v", err)
	}

	return transcript, nil

	// responseText := ""
	// if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
	// 	responseText = resp.Candidates[0].Content.Parts[0].Text
	// } else {
	// 	return "", fmt.Errorf("no response text received from Gemini")
	// }

	// // Return the raw JSON response as a byte slice
	// return responseText, nil
}
