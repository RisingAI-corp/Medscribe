import { atom } from 'jotai';
import { patientsAtom } from '../../states/patientsAtom';
import { format } from 'date-fns';
import { Report } from '../../states/patientsAtom';
import { SearchResultItem } from './SearchResults/SearchResults';

export const searchVisitsAtom = atom<SearchResultItem[]>(get => {
  const allVisits = get(patientsAtom);
  return allVisits.map((visit: Report) => ({
    id: visit.id,
    patientName: visit.name,
    dateOfRecording: format(new Date(visit.timestamp), 'MM/dd/yy'),
    summary: visit.summary.data,
    condensedSummary: visit.summary.data,
    timeOfRecording: format(new Date(visit.timestamp), 'h:mm a'),
    durationOfRecording: visit.duration,
  }));
});
