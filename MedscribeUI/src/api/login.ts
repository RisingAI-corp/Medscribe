import axios from 'axios';
import { AuthResponseSchema } from './serverResponseTypes';

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

  console.log(response.data, ' here it is');
  const { success, data, error } = AuthResponseSchema.safeParse(response.data);
  if (!success) {
    throw new Error('Error parsing Api request: ' + error.toString());
  }

  return data;
}
