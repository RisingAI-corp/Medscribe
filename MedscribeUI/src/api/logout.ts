import axios from 'axios';

export async function logout() {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);

  await axios.post(`${baseURL}/user/logout`, null, {
    withCredentials: true,
  });
}
