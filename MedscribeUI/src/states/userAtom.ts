import { atom } from 'jotai';

export interface User {
  ID: string;
  name: string;
  email: string;
  subjectiveStyle: string; // Added style fields
  objectiveStyle: string;
  assessmentStyle: string;
  planningStyle: string;
  summaryStyle: string;
}

//TODO:delete once api are created
const sampleUser: User = {
  ID: '1',
  name: 'Emenike',
  email: 'emenikeani3@gmail.com',
  subjectiveStyle: 'Creative', // Added sample style values
  objectiveStyle: 'Formal',
  assessmentStyle: 'Analytical',
  planningStyle: 'Detailed',
  summaryStyle: 'Concise',
};

export const isAuthenticatedAtom = atom(false);
export const userAtom = atom<User>(sampleUser);
