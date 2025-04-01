import { z } from 'zod';
import { Report } from '../states/patientsAtom';

export interface AuthResponse {
  id: string;
  name: string;
  email: string;
  subjectiveStyle: string;
  objectiveStyle: string;
  assessmentAndPlanStyle: string;
  summaryStyle: string;
  patientInstructionsStyle: string;
  reports: Report[];
}

export const ReportContentSchema = z.object({
  data: z.string(),
  loading: z.boolean(),
});

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
  finishedGenerating: z.boolean(),
  transcript: z.string().default(''),
  readStatus: z.boolean(),
});

export const AuthResponse = z.object({
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

export const UpdateResponse = z.object({
  Key: z.string(),
  Value: z.unknown(),
});
