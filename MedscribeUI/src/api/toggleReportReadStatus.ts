import axios from 'axios';

export interface UpdateReadStatusPayload {
  ReportID: string;
  Opened: boolean;
}

export async function markRead(payload: UpdateReadStatusPayload) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.patch(`${baseURL}/report/markRead`, payload, {
    withCredentials: true,
  });

  if (response.status == 401) {
    throw new Error('User is not authorized');
  }

  if (response.status != 200 && response.status != 201) {
    throw new Error(`'Error authenticated user' + ${String(response.status)}`);
  }
}

export async function markUnRead(payload: UpdateReadStatusPayload) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.patch(`${baseURL}/report/markUnRead`, payload, {
    withCredentials: true,
  });

  if (response.status == 401) {
    throw new Error('User is not authorized');
  }

  if (response.status != 200 && response.status != 201) {
    throw new Error(`'Error authenticated user' + ${String(response.status)}`);
  }
}
