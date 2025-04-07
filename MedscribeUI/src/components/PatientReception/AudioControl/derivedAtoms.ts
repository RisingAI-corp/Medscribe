import { atom } from 'jotai';
import { Report, patientsAtom } from '../../../states/patientsAtom';
import { REPORT_CONTENT_SECTIONS } from '../../../constants';

export interface UpdateProps {
  id: string;
  Key: string;
  Value?: unknown;
}

export interface CreateReportProps {
  id: string;
  providerID: string;
  name: string;
  timestamp: string;
  duration: number;
}

// Atom to update an existing report in patientsAtom
export const UpdateReportsAtom = atom(
  null,
  (get, set, { id, Key, Value }: UpdateProps) => {
    const reports = get(patientsAtom);

    const updatedReports = reports.map(report => {
      if (report.id !== id) return report;

      const updatedReport = { ...report };

      if (REPORT_CONTENT_SECTIONS.includes(Key)) {
        updatedReport[Key] = {
          data: Value,
          loading: false,
        };
      } else {
        updatedReport[Key] = Value;
      }

      return updatedReport;
    });

    set(patientsAtom, updatedReports);
  },
);

// Atom to create a new report and insert it into patientsAtom
export const createReportAtom = atom(
  null,
  (
    get,
    set,
    { id, providerID, name, timestamp, duration }: CreateReportProps,
  ) => {
    const reports = get(patientsAtom);

    if (reports.some(report => report.id === id)) {
      console.warn(`Report with ID '${id}' already exists. Skipping creation.`);
      return;
    }

    const newReport: Report = {
      id,
      providerID,
      name,
      timestamp,
      duration,
      pronouns: '',
      isFollowUp: false,
      patientOrClient: '',
      subjective: { data: '', loading: true },
      objective: { data: '', loading: true },
      assessmentAndPlan: { data: '', loading: true },
      patientInstructions: { data: '', loading: true },
      summary: { data: '', loading: true },
      sessionSummary: '',
      condensedSummary: '',
      finishedGenerating: false,
      transcript: '',
      readStatus: false,
    };

    set(patientsAtom, [newReport, ...reports]);
  },
);
