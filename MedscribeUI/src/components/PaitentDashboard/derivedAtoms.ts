import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';

export const HeaderInformationAtom = atom(
  get => {
    const currentlySelectedPatient = get(currentlySelectedPatientAtom);
    const patients = get(patientsAtom);
    const patient = patients.find(p => p.id === currentlySelectedPatient);
    if (patient) {
      return {
        name: patient.name,
        oneLiner: patient.oneLinerSummary,
        id: patient.id,
      };
    } else {
      console.error(`Patient not found: ${JSON.stringify(patient)}`);
      return {
        name: '',
        oneLiner: '',
        id: '',
      };
    }
  },
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
      return;
    }
    console.error(`Patient not found: ${JSON.stringify(patient)}`);
  },
);
