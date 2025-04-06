import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';

export interface PatientBackgroundData {
  id: string;
  name: string;
  condensedSummary: string;
  lastVisitDate: string;
  duration: number;
  lastVisitSummary: string;
}

// Atom to store the currently selected patient ID for the background view
export const selectedPatientIdForBackgroundAtom = atom<string | null>(null);

// Derived atom that gets the patient background data based on the selected ID
export const patientBackgroundDataAtom = atom<PatientBackgroundData | null>((get) => {
  const patientId = get(selectedPatientIdForBackgroundAtom);
  if (!patientId) return null;
  
  const patients = get(patientsAtom);
  const patient = patients.find(p => p.id === patientId);
  
  if (!patient) return null;
  
  return {
    id: patient.id,
    name: patient.name,
    condensedSummary: patient.condensedSummary,
    lastVisitDate: patient.timestamp,
    duration: patient.duration,
    lastVisitSummary: patient.summary.data,
  };
}); 