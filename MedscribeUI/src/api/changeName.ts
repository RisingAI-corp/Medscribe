import axios from 'axios';

export interface changeNameRequest {
  ReportID: string;
  NewName: string;
}

export async function changeName(payload: changeNameRequest) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.patch(`${baseURL}/report/changeName`, payload, {
    withCredentials: true,
  });

  if (response.status == 401) {
    throw new Error('User is not authorized');
  }

  if (response.status != 200 && response.status != 201) {
    throw new Error(`'Error authenticated user' + ${String(response.status)}`);
  }
}
