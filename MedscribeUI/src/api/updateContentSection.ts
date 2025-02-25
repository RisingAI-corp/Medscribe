import axios from 'axios';

export interface updateContentSectionRequest {
  ReportID: string;
  ContentSection: string;
  Content: string;
}

export async function updateContentSection(
  payload: updateContentSectionRequest,
) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);
  const response = await axios.patch(
    `${baseURL}/report/updateContentSection`,
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
}
