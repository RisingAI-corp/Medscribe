import axios from 'axios';
import { AuthResponse } from './serverResponses';

interface LoginProps {
  email: string;
  password: string;
}

export async function loginProvider({ email, password }: LoginProps) {
  if (!email || !password) {
    throw new Error('Email, name, and password are required.');
  }
  const baseUrl = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.post(
    `${baseUrl}/user/login`,
    { email, password },
    { withCredentials: true },
  );

  console.log(response.data);
  const { success, data, error } = AuthResponse.safeParse(response.data);
  if (!success) {
    throw new Error('Error parsing Api request: ' + error.toString());
  }

  return data;
}
