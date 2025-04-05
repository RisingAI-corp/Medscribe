import { atom } from 'jotai';
import { patientsAtom, ReportContent } from '../../states/patientsAtom';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { Report } from '../../states/patientsAtom';

export const NoteControlsAtom = atom(get => {
  const currentlySelectedPatient = get(currentlySelectedPatientAtom);
  const patients = get(patientsAtom);
  const patient = patients.find(p => p.id === currentlySelectedPatient);
  if (patient) {
    return {
      pronouns: patient.pronouns,
      visitType: patient.isFollowUp,
      patientOrClient: patient.patientOrClient,
    };
  }
  return {
    pronouns: '',
    visitType: false,
    patientOrClient: '',
  };
});

export interface noteControlEditsProps {
  id: string;
  changes: Partial<Pick<Report, 'patientOrClient' | 'pronouns' | 'isFollowUp'>>;
}

export const editPatientVisitAtom = atom(
  null,
  (get, set, updates: noteControlEditsProps) => {
    const currentVisits = get(patientsAtom);

    const updatedVisits = currentVisits.map(visit => {
      if (visit.id === updates.id) {
        return {
          ...visit,
          ...updates.changes,
        };
      }

      return visit;
    });
    set(patientsAtom, updatedVisits);
  },
);

const updateSectionLoading = (section: ReportContent): ReportContent => {
  return { ...section, loading: true };
};

export const toggleLoadingForReportSectionsAtom = atom(null, (get, set) => {
  const selectedReport = get(currentlySelectedPatientAtom);
  const updatedReports = get(patientsAtom).map(report => {
    if (report.id == selectedReport) {
      return {
        ...report,
        subjective: updateSectionLoading(report.subjective),
        objective: updateSectionLoading(report.objective),
        assessmentAndPlan: updateSectionLoading(report.assessmentAndPlan),
        patientInstructions: updateSectionLoading(report.patientInstructions),
        summary: updateSectionLoading(report.summary),
      };
    } else {
      return report;
    }
  });
  set(patientsAtom, updatedReports);
});
export const fetchContentDataAtom = atom(null, get => {
  const currentlySelectedPatient = get(currentlySelectedPatientAtom);
  const reports = get(patientsAtom);
  const report = reports.find(p => p.id === currentlySelectedPatient);
  if (!report) {
    return null;
  }

  return {
    subjective: report.subjective.data,
    objective: report.objective.data,
    assessmentAndPlan: report.assessmentAndPlan.data,
    patientInstructions: report.patientInstructions.data,
    summary: report.summary.data,
  };
});
