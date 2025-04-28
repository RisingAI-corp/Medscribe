import { atom } from 'jotai';

import { userAtom } from '../../states/userAtom';

export const UpdateUserName = atom(null, (get, set, newName: string) => {
  const user = get(userAtom);
  const updatedUser = { ...user, name: newName };
  set(userAtom, updatedUser);
});

export const getNameAtom = atom(get => get(userAtom).name);
