import axios from 'axios';
import { AuthResponse } from './serverResponses';

interface SignUpProps {
  email: string;
  name: string;
  password: string;
}

export async function createProvider({ name, email, password }: SignUpProps) {
  console.log('sending');
  if (!email || !name || !password) {
    throw new Error('Email, name, and password are required.');
  }

  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);

  try {
    const response = await axios.post(
      `${baseURL}/user/signup`,
      { name, email, password },
      { withCredentials: true },
    );
    const { success, data, error } = AuthResponse.safeParse(response.data);
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
