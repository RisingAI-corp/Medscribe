import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { Report } from '../../api/serverResponseTypes';
import { REPORT_CONTENT_SECTIONS } from '../../constants';

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

export const UpdateReportsAtom = atom(
  null,
  (get, set, { id, Key, Value }: UpdateProps) => {
    const reports = get(patientsAtom);

    const newReports = reports.map(report => {
      if (report.id === id) {
        // Create a new report object with updated values
        const updatedReport = { ...report }; // Shallow copy
        if (REPORT_CONTENT_SECTIONS.find(member => member === Key)) {
          updatedReport[Key] = {
            data: Value,
            loading: false,
          };
        } else {
          updatedReport[Key] = Value;
        }

        return updatedReport;
      }
      return report; // Return unchanged reports
    });

    set(patientsAtom, newReports);
  },
);

export const createReportAtom = atom(
  null,
  (
    get,
    set,
    { id, providerID, name, timestamp, duration }: CreateReportProps,
  ) => {
    const reports = get(patientsAtom);
    if (reports.find(report => report.id === id)) {
      console.error('report already exists');
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
      transcriptContainer: {
        transcript: '',
        diarizedTranscript: [],
        providerID: '',
        usedDiarization: false,
      },
      readStatus: false,
      status: 'pending',
      lastVisitID: '',
    };
    set(patientsAtom, [newReport, ...reports]);
  },
);
