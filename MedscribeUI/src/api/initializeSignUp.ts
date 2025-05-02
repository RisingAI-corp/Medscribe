import axios from 'axios';

interface initializeSignUpProps {
  email: string;
  name: string;
  password: string;
}

export async function initializeSignUp({
  name,
  email,
  password,
}: initializeSignUpProps) {
  console.log('sending');
  if (!email || !name || !password) {
    throw new Error('Email, name, and password are required.');
  }

  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);

  try {
    const response = await axios.post(
      `${baseURL}/user/initializeSignup`,
      { name, email, password },
      { withCredentials: true },
    );
    console.log('response', response);
  } catch (err) {
    console.error('Error initializing sign up:', err);
    if (axios.isAxiosError(err) && err.response) {
      console.log('response', err.response);
      if (typeof err.response.data === 'string') {
        if (err.response.status === 409) {
          throw new Error(err.response.data);
        }
        throw new Error(err.response.data);
      }
    }
  }
}
