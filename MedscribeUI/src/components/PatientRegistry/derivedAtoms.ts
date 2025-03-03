import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { format } from 'date-fns';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { PatientPreviewRecord } from './patientRegistryLayout';

export const PatientRegistryAtom = atom(get =>
  get(patientsAtom).map<PatientPreviewRecord>(patient => {
    const dateOfRecording = new Date(patient.timestamp);
    const date = format(dateOfRecording, 'MM/dd/yy').toLowerCase();
    const time = format(dateOfRecording, 'h:mmaaa').toLowerCase();
    const durationInMinutes = Math.floor(patient.duration / 60);
    const duration = `${String(durationInMinutes)} min`;

    return {
      id: patient.id,
      patientName: patient.name,
      dateOfRecording: date,
      timeOfRecording: time,
      durationOfRecording: duration,
      shortenedSummary: patient.oneLinerSummary,
      finishedGenerating: patient.finishedGenerating,
      loading: !patient.finishedGenerating,
    };
  }),
);

export const removePatientsByIdsAtom = atom(
  null,
  (get, set, patientIdentifiers: string[] | string) => {
    const patients = get(patientsAtom);
    const currentlySelectedPatient = get(currentlySelectedPatientAtom);

    const patientIdsToRemove = Array.isArray(patientIdentifiers)
      ? patientIdentifiers
      : [patientIdentifiers];

    const updatedPatients = patients.filter(
      patient => !patientIdsToRemove.includes(patient.id),
    );

    set(patientsAtom, updatedPatients);

    if (patientIdsToRemove.includes(currentlySelectedPatient)) {
      set(currentlySelectedPatientAtom, '');
    }
  },
);
