import axios from 'axios';
import { AuthResponse } from './serverResponses';

interface SignUpProps {
  email: string;
  name: string;
  password: string;
}

export async function createProvider({ name, email, password }: SignUpProps) {
  if (!email || !name || !password) {
    throw new Error('Email, name, and password are required.');
  }

  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.post(
    `${baseURL}/user/signup`,
    {
      name: name,
      email: email,
      password: password,
    },
    { withCredentials: true },
  );

  const { success, data, error } = AuthResponse.safeParse(response.data);
  if (!success) {
    throw new Error('Error parsing Api request: ' + error.toString());
  }

  return data;
}
