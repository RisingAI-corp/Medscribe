import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import {
  REPORT_CONTENT_TYPE_ASSESSMENT_AND_PLAN,
  REPORT_CONTENT_TYPE_CONDENSED_SUMMARY,
  REPORT_CONTENT_TYPE_OBJECTIVE,
  REPORT_CONTENT_TYPE_PATIENT_INSTRUCTIONS,
  REPORT_CONTENT_TYPE_SESSION_SUMMARY,
  REPORT_CONTENT_TYPE_SUBJECTIVE,
  REPORT_CONTENT_TYPE_SUMMARY,
} from '../../constants';
import { Report, TranscriptContainer } from '../../api/serverResponseTypes';

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

export const updateTranscriptAtom = atom(
  null,
  (
    get,
    set,
    {
      id,
      transcriptContainer,
    }: { id: string; transcriptContainer: TranscriptContainer },
  ) => {
    const reports = get(patientsAtom);
    set(
      patientsAtom,
      reports.map(report => {
        if (report.id == id) {
          return {
            ...report,
            transcriptContainer: transcriptContainer,
          };
        }
        return report;
      }),
    );
  },
);

export const SoapAtom = atom(
  get => {
    const currentlySelectedPatient = get(currentlySelectedPatientAtom);
    const patients = get(patientsAtom);
    const patient = patients.find(p => p.id === currentlySelectedPatient);
    if (patient) {
      return {
        soapContent: [
          {
            title: 'Visit Summary',
            sectionType: REPORT_CONTENT_TYPE_SUMMARY,
            content: {
              data: patient.summary.data,
              loading: patient.summary.loading,
            },
          },
          {
            title: 'Subjective',
            sectionType: REPORT_CONTENT_TYPE_SUBJECTIVE,
            content: {
              data: patient.subjective.data,
              loading: patient.subjective.loading,
            },
          },
          {
            title: 'Objective',
            sectionType: REPORT_CONTENT_TYPE_OBJECTIVE,
            content: {
              data: patient.objective.data,
              loading: patient.objective.loading,
            },
          },
          {
            title: 'Assessment & Plan',
            sectionType: REPORT_CONTENT_TYPE_ASSESSMENT_AND_PLAN,
            content: {
              data: patient.assessmentAndPlan.data,
              loading: patient.assessmentAndPlan.loading,
            },
          },
          {
            title: 'Patient Instructions',
            sectionType: REPORT_CONTENT_TYPE_PATIENT_INSTRUCTIONS,
            content: {
              data: patient.patientInstructions.data,
              loading: patient.patientInstructions.loading,
            },
          },
        ],
        status: patient.status,
        transcriptContainer: patient.transcriptContainer,
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
  try {
    switch (field) {
      case REPORT_CONTENT_TYPE_SUBJECTIVE:
        return {
          ...report,
          subjective: { ...report.subjective, data: newData },
        };
      case REPORT_CONTENT_TYPE_OBJECTIVE:
        return { ...report, objective: { ...report.objective, data: newData } };
      case REPORT_CONTENT_TYPE_ASSESSMENT_AND_PLAN:
        return {
          ...report,
          assessmentAndPlan: { ...report.assessmentAndPlan, data: newData },
        };
      case REPORT_CONTENT_TYPE_PATIENT_INSTRUCTIONS:
        return {
          ...report,
          patientInstructions: { ...report.patientInstructions, data: newData },
        };
      case REPORT_CONTENT_TYPE_SUMMARY:
        return { ...report, summary: { ...report.summary, data: newData } };
      case REPORT_CONTENT_TYPE_CONDENSED_SUMMARY:
        return { ...report, condensedSummary: newData };
      case REPORT_CONTENT_TYPE_SESSION_SUMMARY:
        return { ...report, sessionSummary: newData };
      default:
        console.error(`Unknown report content type: ${field}`);
        return report; // Return the original report if the field is unknown
    }
  } catch (error) {
    // Handle the error here.  What you do depends on the needs of your application.
    console.error(`Error updating report content for field "${field}":`, error);
    //  Possible error handling strategies:
    //  1.  Return a specific error value:
    //      return { ...report, error: 'Failed to update content' }; // If your Report type has an error field
    //  2.  Throw the error to be caught by the caller:
    //      throw error;
    //  3.  Return the original report (as you were doing):
    return report;
    //  4.  Return a report with the newData and an error flag
    //       return {...report, [field]: {data: newData, error: true}
  }
}
