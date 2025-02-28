import { atom } from 'jotai';
import { Client, He, Patient, She, They } from '../constants';

export interface ReportContent {
  data: string;
  loading: boolean;
}
export interface Report {
  id: string;
  providerID: string;
  name: string;
  timestamp: string;
  duration: number;
  pronouns: string;
  isFollowUp: boolean;
  patientOrClient: string;
  subjective: ReportContent;
  objective: ReportContent;
  assessment: ReportContent;
  planning: ReportContent;
  summary: ReportContent;
  oneLinerSummary: string;
  shortSummary: string;
  finishedGenerating: boolean;
}

//TODO: remove when apis are built
const reports: Report[] = [
  {
    id: 'report1',
    providerID: 'provider123',
    name: 'Emenike',
    timestamp: '2024-07-26T10:00:00Z',
    duration: 600000,
    pronouns: She,
    isFollowUp: false,
    patientOrClient: Patient,
    subjective: {
      data: 'Patient presented with mild headache.',

      loading: false,
    },
    objective: {
      data: 'Blood pressure 120/80, no fever.',

      loading: false,
    },
    assessment: {
      data: 'Likely tension headache.',

      loading: false,
    },
    planning: {
      data: 'Recommend rest and hydration.',

      loading: false,
    },
    summary: {
      data: 'Routine follow-up recommended.',

      loading: false,
    },
    oneLinerSummary: 'Patient visit summary 1',
    shortSummary: 'Short summary of visit 1',
    finishedGenerating: true,
  },
  {
    id: 'report2',
    providerID: 'provider456',
    name: 'President',
    timestamp: '2024-07-26T11:30:00Z',
    duration: 1800000,
    pronouns: He,
    isFollowUp: true,
    patientOrClient: Patient,
    subjective: {
      data: 'Patient reported chronic back pain.',

      loading: false,
    },
    objective: {
      data: 'Limited range of motion in lower back.',
      loading: false,
    },
    assessment: {
      data: 'Chronic lower back pain, further evaluation needed.',

      loading: false,
    },
    planning: {
      data: 'Referral to physical therapy.',
      loading: false,
    },
    summary: {
      data: 'Follow-up appointment scheduled.',
      loading: false,
    },
    oneLinerSummary: 'Patient visit summary 2',
    shortSummary: 'Short summary of visit 2',
    finishedGenerating: true,
  },
  {
    id: 'report3',
    providerID: 'provider123',
    name: 'Cockroach',
    timestamp: '2024-07-27T09:00:00Z',
    duration: 200000,
    pronouns: They, // Using imported constant
    isFollowUp: false,
    patientOrClient: Client, // Using imported constant
    subjective: {
      data: 'Patient presented with anxiety.',

      loading: false,
    },
    objective: {
      data: 'Restlessness observed.',

      loading: false,
    },
    assessment: {
      data: 'Anxiety disorder.',

      loading: false,
    },
    planning: {
      data: 'Prescribed medication and therapy.',

      loading: false,
    },
    summary: {
      data: 'Follow-up appointment in two weeks.',

      loading: false,
    },
    oneLinerSummary: 'Patient visit summary 3',
    shortSummary: 'Short summary of visit 3',
    finishedGenerating: true,
  },
];

//TODO: remove when e2e with the backend is built
const sampleCurrentlySelectedPatient = 'report1';

export const currentlySelectedPatientAtom = atom<string>(
  sampleCurrentlySelectedPatient,
);

export const patientsAtom = atom<Report[]>(reports);

export const reportStreamingAtom = atom<Set<string>>(new Set<string>());

export const setReportStreamStatusAtom = atom(null, (get, set, id: string) => {
  const currentSet: Set<string> = get(reportStreamingAtom);
  const newSet = new Set(currentSet);
  newSet.add(id);
  set(reportStreamingAtom, newSet);
});

export const unsetReportStreamStatusAtom = atom(
  null,
  (get, set, id: string) => {
    const currentSet: Set<string> = get(reportStreamingAtom);
    const newSet = new Set(currentSet);
    newSet.delete(id);
    set(reportStreamingAtom, newSet);
  },
);
