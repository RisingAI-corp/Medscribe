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
  assessmentAndPlan: ReportContent;
  patientInstructions: ReportContent;
  summary: ReportContent;
  sessionSummary: string;
  condensedSummary: string;
  finishedGenerating: boolean;
  transcript: string;
}

// //TODO: remove when apis are built
// const reports: Report[] = [
//   {
//     id: 'report1',
//     providerID: 'provider123',
//     name: 'Emenike',
//     timestamp: '2024-07-26T10:00:00Z',
//     duration: 600000,
//     pronouns: She,
//     isFollowUp: false,
//     patientOrClient: Patient,
//     subjective: {
//       data: 'Patient presented with mild headache.',

//       loading: false,
//     },
//     objective: {
//       data: 'Blood pressure 120/80, no fever.',

//       loading: false,
//     },
//     assessmentAndPlan: {
//       data: 'Likely tension headache.',

//       loading: false,
//     },
//     patientInstructions: {
//       data: 'Rest and drink fluids.',

//       loading: false,
//     },
//     summary: {
//       data: 'Routine follow-up recommended.',

//       loading: false,
//     },
//     sessionSummary: 'Patient visit summary 1',
//     condensedSummary: 'condensed summary of visit 1',
//     finishedGenerating: true,
//     transcript: 'some random text',
//   },
//   {
//     id: 'report2',
//     providerID: 'provider456',
//     name: 'President',
//     timestamp: '2024-07-26T11:30:00Z',
//     duration: 1800000,
//     pronouns: He,
//     isFollowUp: true,
//     patientOrClient: Patient,
//     subjective: {
//       data: 'Patient reported chronic back pain.',

//       loading: false,
//     },
//     objective: {
//       data: 'Limited range of motion in lower back.',
//       loading: false,
//     },
//     assessmentAndPlan: {
//       data: 'Chronic lower back pain, further evaluation needed.',

//       loading: false,
//     },
//     patientInstructions: {
//       data: 'Referral to physical therapy.',
//       loading: false,
//     },
//     summary: {
//       data: 'Follow-up appointment scheduled.',
//       loading: false,
//     },
//     sessionSummary: 'Patient visit summary 2',
//     condensedSummary: 'condensed summary of visit 2',
//     finishedGenerating: true,
//     transcript: '',
//   },
//   {
//     id: 'report3',
//     providerID: 'provider123',
//     name: 'Cockroach',
//     timestamp: '2024-07-27T09:00:00Z',
//     duration: 200000,
//     pronouns: They, // Using imported constant
//     isFollowUp: false,
//     patientOrClient: Client, // Using imported constant
//     subjective: {
//       data: 'Patient presented with anxiety.',

//       loading: false,
//     },
//     objective: {
//       data: 'Restlessness observed.',

//       loading: false,
//     },
//     assessmentAndPlan: {
//       data: 'Anxiety disorder.',

//       loading: false,
//     },
//     patientInstructions: {
//       data: 'Prescribed medication and therapy.',

//       loading: false,
//     },
//     summary: {
//       data: 'Follow-up appointment in two weeks.',

//       loading: false,
//     },
//     sessionSummary: 'Patient visit summary 3',
//     condensedSummary: 'condensed summary of visit 3',
//     finishedGenerating: true,
//     transcript: '',
//   },
// ];

// const sampleCurrentlySelectedPatient = 'report1';

export const currentlySelectedPatientAtom = atom<string>('');

export const patientsAtom = atom<Report[]>([]);

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

export const UpdateSelectedPatientNameAtom = atom(
  null,
  (get, set, newName: string) => {
    const currentlySelectedPatient = get(currentlySelectedPatientAtom);
    const patients = get(patientsAtom);
    const patient = patients.find(p => p.id === currentlySelectedPatient);

    if (patient) {
      const updatedPatients = patients.map(p => {
        if (p.id === patient.id) {
          return {
            ...p,
            name: newName,
          };
        }
        return p;
      });
      set(patientsAtom, updatedPatients);
    }
  },
);
