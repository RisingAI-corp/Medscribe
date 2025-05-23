import { atom } from 'jotai';
import {
  currentlySelectedPatientAtom,
  patientsAtom,
} from '../../states/patientsAtom';
import { format } from 'date-fns';
import { SearchResultItem } from './SearchResults/SearchResults';
import { Report } from '../../api/serverResponseTypes';

export const searchVisitsAtom = atom<SearchResultItem[]>(get => {
  const allVisits = get(patientsAtom);
  const currentlySelectedReport = get(currentlySelectedPatientAtom);
  return allVisits
    .filter(visit => visit.id !== currentlySelectedReport)
    .map((visit: Report) => ({
      id: visit.id,
      patientName: visit.name,
      dateOfRecording: format(new Date(visit.timestamp), 'MM/dd/yy'),
      summary: visit.summary.data,
      condensedSummary: visit.condensedSummary,
      timeOfRecording: format(new Date(visit.timestamp), 'h:mm a'),
      durationOfRecording: visit.duration,
    }));
});
