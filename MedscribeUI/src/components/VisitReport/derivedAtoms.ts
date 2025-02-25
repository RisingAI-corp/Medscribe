import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import {
  REPORT_CONTENT_TYPE_ASSESSMENT,
  REPORT_CONTENT_TYPE_OBJECTIVE,
  REPORT_CONTENT_TYPE_PLANNING,
  REPORT_CONTENT_TYPE_SUBJECTIVE,
  REPORT_CONTENT_TYPE_SUMMARY,
} from '../../constants';
import { Report } from '../../states/patientsAtom';

export const replaceReportAtom = atom(null, (get, set, newReport: Report) => {
  const reports = get(patientsAtom);
  set(
    patientsAtom,
    reports.map(report => {
      if (report.id == newReport.id) {
        return newReport;
      }
      return report;
    }),
  );
});

interface SoapSection {
  type: string;
  content: {
    data: string;
    loading: boolean;
  };
}

export const SoapAtom = atom(
  get => {
    const currentlySelectedPatient = get(currentlySelectedPatientAtom);
    const patients = get(patientsAtom);
    const patient = patients.find(p => p.id === currentlySelectedPatient);
    if (patient) {
      return {
        content: [
          {
            type: REPORT_CONTENT_TYPE_SUBJECTIVE,
            content: {
              data: patient.subjective.data,
              loading: patient.subjective.loading,
            },
          },
          {
            type: REPORT_CONTENT_TYPE_OBJECTIVE,
            content: {
              data: patient.objective.data,
              loading: patient.objective.loading,
            },
          },
          {
            type: REPORT_CONTENT_TYPE_ASSESSMENT,
            content: {
              data: patient.assessment.data,
              loading: patient.assessment.loading,
            },
          },
          {
            type: REPORT_CONTENT_TYPE_PLANNING,
            content: {
              data: patient.planning.data,
              loading: patient.planning.loading,
            },
          },
        ],
        loading: patient.finishedGenerating,
      };
    }

    console.error(`Patient ${currentlySelectedPatient} not found `);
    return null;
  },
  (
    get,
    set,
    {
      patientId,
      field,
      newData,
    }: { patientId: string; field: string; newData: string },
  ) => {
    const reports = get(patientsAtom);
    const updatedReports = reports.map(report => {
      if (report.id === patientId) {
        return updateReportContent(report, field, newData);
      }
      return report;
    });
    set(patientsAtom, updatedReports);
  },
);

function updateReportContent(
  report: Report,
  field: string,
  newData: string,
): Report {
  switch (field) {
    case REPORT_CONTENT_TYPE_SUBJECTIVE:
      return { ...report, subjective: { ...report.subjective, data: newData } };
    case REPORT_CONTENT_TYPE_OBJECTIVE:
      return { ...report, objective: { ...report.objective, data: newData } };
    case REPORT_CONTENT_TYPE_ASSESSMENT:
      return { ...report, assessment: { ...report.assessment, data: newData } };
    case REPORT_CONTENT_TYPE_PLANNING:
      return { ...report, planning: { ...report.planning, data: newData } };
    case REPORT_CONTENT_TYPE_SUMMARY:
      return { ...report, summary: { ...report.summary, data: newData } };
    default:
      console.error(`Unknown report content type: ${field}`);
      return report; // Return the original report if the field is unknown
  }
}
