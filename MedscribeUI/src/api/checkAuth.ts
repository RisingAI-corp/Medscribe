import axios from 'axios';
import { AuthResponseSchema } from './serverResponseTypes';

export async function checkAuth() {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.get(`${baseURL}/checkAuth`, {
    withCredentials: true,
  });

  if (response.status == 401) {
    throw new Error('User is not authorized');
  }

  if (response.status != 200 && response.status != 201) {
    throw new Error(`'Error authenticated user' + ${String(response.status)}`);
  }

  console.log(response.data);
  const { success, data, error } = AuthResponseSchema.safeParse(response.data);
  if (!success) {
    throw new Error('Error parsing Api request: ' + error.toString());
  }
  return data;
}
