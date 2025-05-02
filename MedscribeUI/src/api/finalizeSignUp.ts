import axios from 'axios';
import { AuthResponseSchema } from './serverResponseTypes';

export async function finalizeSignUp(token: string) {
  try {
    console.log('sending');
    if (!token) {
      throw new Error('Token is required.');
    }

    const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);

    const response = await axios.post(
      `${baseURL}/user/fianalizeSignup`,
      { token: token },
      { withCredentials: true },
    );
    const { success, data, error } = AuthResponseSchema.safeParse(
      response.data,
    );
    if (!success) {
      throw new Error('Error parsing API response: ' + error.toString());
    }
    return data;
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      if (err.response.status === 409) {
        throw new Error('status conflict: user already exists');
      }
    }
    throw err;
  }
}
