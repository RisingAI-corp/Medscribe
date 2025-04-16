import axios from 'axios';
import { ReportSchema } from './serverResponseTypes';

export interface GetReportPayload {
  reportID: string;
}

export async function getReport(payload: GetReportPayload) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.post(`${baseURL}/report/get`, payload, {
    withCredentials: true,
  });

  if (response.status == 401) {
    throw new Error('User is not authorized');
  }

  if (response.status != 200 && response.status != 201) {
    throw new Error(`'Error authenticated user' + ${String(response.status)}`);
  }

  const { success, data, error } = ReportSchema.safeParse(response.data);
  if (!success) {
    throw new Error('Error parsing Api request: ' + error.toString());
  }
  console.log(data, 'yo');
  return data;
}
