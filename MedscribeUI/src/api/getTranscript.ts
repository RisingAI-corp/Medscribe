import axios from 'axios';
import { TranscriptContainerSchema } from './serverResponseTypes';

export interface GetTranscriptPayload {
  reportID: string;
}

export async function getTranscript(payload: GetTranscriptPayload) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.post(
    `${baseURL}/report/getTranscript`,
    payload,
    {
      withCredentials: true,
    },
  );

  if (response.status == 401) {
    throw new Error('User is not authorized');
  }

  if (response.status != 200 && response.status != 201) {
    throw new Error(`'Error authenticated user' + ${String(response.status)}`);
  }

  console.log(response.data, ' here it is');
  const { success, data, error } = TranscriptContainerSchema.safeParse(
    response.data,
  );
  if (!success) {
    throw new Error('Error parsing Api request: ' + error.toString());
  }
  return data;
}
