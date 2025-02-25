import axios from 'axios';

export interface DeleteReportPayload {
  ReportIDs: string[];
}

export async function deleteReport(payload: DeleteReportPayload) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.delete(`${baseURL}/report/delete`, {
    data: payload,
    withCredentials: true,
  });

  if (response.status == 401) {
    throw new Error('User is not authorized');
  }

  if (response.status != 200 && response.status != 201) {
    throw new Error(`'Error authenticated user' + ${String(response.status)}`);
  }
}
