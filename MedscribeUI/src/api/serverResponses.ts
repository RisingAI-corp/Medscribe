import { z } from 'zod';

export const ReportContentSchema = z.object({
  data: z.string(),
  loading: z.boolean(),
});

export const ReportSchema = z.object({
  id: z.string(),
  providerID: z.string(),
  name: z.string(),
  timeStamp: z.string(),
  duration: z.number(),
  pronouns: z.string(),
  isFollowUp: z.boolean(),
  patientOrClient: z.string(),
  subjective: ReportContentSchema,
  objective: ReportContentSchema,
  assessment: ReportContentSchema,
  planning: ReportContentSchema,
  summary: ReportContentSchema,
  oneLinerSummary: z.string(),
  shortSummary: z.string(),
  finishedGenerating: z.boolean(),
});

export const AuthResponse = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string(),
  reports: z.array(ReportSchema).nullable().default([]),
  subjectiveStyle: z.string(),
  objectiveStyle: z.string(),
  assessmentStyle: z.string(),
  planningStyle: z.string(),
  summaryStyle: z.string(),
});

export const UpdateResponse = z.object({
  Key: z.string(),
  Value: z.unknown(),
});
