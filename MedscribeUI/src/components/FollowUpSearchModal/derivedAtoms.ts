import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { format } from 'date-fns';
import { Report } from '../../states/patientsAtom';

export const searchVisitsAtom = atom((get) => {
  const allVisits = get(patientsAtom);
  return allVisits.map((visit: Report) => ({
    id: visit.id,
    patientName: visit.name,
    dateOfRecording: format(new Date(visit.timestamp), 'MM/dd/yy'),
    shortenedSummary: visit.oneLinerSummary,
  }));
});