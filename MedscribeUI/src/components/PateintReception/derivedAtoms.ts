import { atom } from 'jotai';
import { patientsAtom, Report } from '../../states/patientsAtom';
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
  timeStamp: string;
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
          console.log('attempting udpate of ', Key);
          updatedReport[Key] = {
            data: Value,
            loading: false,
          };
        } else {
          updatedReport[Key] = Value;
        }

        console.log('update successful', Key, Value);

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
    { id, providerID, name, timeStamp, duration }: CreateReportProps,
  ) => {
    const reports = get(patientsAtom);
    if (reports.find(report => report.id === id)) {
      console.error('report already exists');
      return;
    }
    console.log('creating report');
    const newReport: Report = {
      id,
      providerID,
      name,
      timeStamp,
      duration,
      pronouns: '',
      isFollowUp: false,
      patientOrClient: '',
      subjective: { data: '', loading: true },
      objective: { data: '', loading: true },
      assessment: { data: '', loading: true },
      planning: { data: '', loading: true },
      summary: { data: '', loading: true },
      oneLinerSummary: '',
      shortSummary: '',
      finishedGenerating: false,
    };
    set(patientsAtom, [newReport, ...reports]);
  },
);
