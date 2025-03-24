import { atom } from 'jotai';

export interface User {
  ID: string;
  name: string;
  email: string;
  subjectiveStyle: string; // Added style fields
  objectiveStyle: string;
  assessmentAndPlanStyle: string;
  summaryStyle: string;
  patientInstructionsStyle: string;
}

//TODO:delete once api are created
const sampleUser: User = {
  ID: '1',
  name: 'Emenike',
  email: 'emenikeani3@gmail.com',
  subjectiveStyle: '',
  objectiveStyle: '',
  summaryStyle: '',
  assessmentAndPlanStyle: '',
  patientInstructionsStyle: '',
};

export const isAuthenticatedAtom = atom(false);
export const userAtom = atom<User>(sampleUser);
