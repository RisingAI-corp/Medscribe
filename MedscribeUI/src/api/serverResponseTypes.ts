import { z } from 'zod';

export const ReportContentSchema = z.object({
  data: z.string(),
  loading: z.boolean(),
});

export type ReportContent = z.infer<typeof ReportContentSchema>;

export const TranscriptTurnSchema = z.object({
  speaker: z.string(),
  startTime: z.number(),
  endTime: z.number(),
  text: z.string(),
});

export type TranscriptTurn = z.infer<typeof TranscriptTurnSchema>;

export const DiarizedTranscriptSchema = z.array(TranscriptTurnSchema);

export type DiarizedTranscript = z.infer<typeof DiarizedTranscriptSchema>;

export const TranscriptContainerSchema = z.object({
  transcript: z.string(),
  diarizedTranscript: z.nullable(z.array(TranscriptTurnSchema)).default([]),
  providerID: z.string(),
  usedDiarization: z.boolean(),
});

export type TranscriptContainer = z.infer<typeof TranscriptContainerSchema>;

export const ReportSchema = z.object({
  id: z.string(),
  providerID: z.string(),
  name: z.string(),
  timestamp: z.string(),
  duration: z.number(),
  pronouns: z.string(),
  isFollowUp: z.boolean(),
  patientOrClient: z.string(),
  subjective: ReportContentSchema,
  objective: ReportContentSchema,
  assessmentAndPlan: ReportContentSchema,
  patientInstructions: ReportContentSchema,
  summary: ReportContentSchema,
  condensedSummary: z.string(),
  sessionSummary: z.string(),
  readStatus: z.boolean(),
  status: z.string(),
  transcriptContainer: TranscriptContainerSchema.default({
    transcript: '',
    diarizedTranscript: [],
    providerID: '',
    usedDiarization: false,
  }),
  lastVisitID: z.string(),
});

export type Report = z.infer<typeof ReportSchema>;

export const AuthResponseSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string(),
  reports: z.array(ReportSchema).default([]),
  subjectiveStyle: z.string(),
  objectiveStyle: z.string(),
  assessmentAndPlanStyle: z.string(),
  patientInstructionsStyle: z.string(),
  summaryStyle: z.string(),
});

export type AuthResponse = z.infer<typeof AuthResponseSchema>;

export const UpdateResponseSchema = z.object({
  Key: z.string(),
  Value: z.unknown(),
});

export type UpdateResponse = z.infer<typeof UpdateResponseSchema>;
