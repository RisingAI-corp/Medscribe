import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';

export const SelectedPatientHeaderInformationAtom = atom(get => {
  const currentlySelectedPatient = get(currentlySelectedPatientAtom);
  const patients = get(patientsAtom);
  const patient = patients.find(p => p.id === currentlySelectedPatient);
  if (patient) {
    return {
      name: patient.name,
      condensedSummary: patient.condensedSummary,
      id: patient.id,
    };
  } else {
    console.error(`Patient not found: ${JSON.stringify(patient)}`);
    return {
      name: '',
      condensedSummary: '',
      id: '',
    };
  }
});
