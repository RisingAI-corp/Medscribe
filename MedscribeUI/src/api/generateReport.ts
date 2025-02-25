export interface GenerateReportMetadata {
  providerID: string;
  patientName: string;
  timestamp: string;
  duration: number;
  subjectiveStyle: string;
  objectiveStyle: string;
  assessmentStyle: string;
  planningStyle: string;
  summaryStyle: string;
}

export async function generateReport(
  formData: FormData,
  metadataData: GenerateReportMetadata,
) {
  formData.append('metadata', JSON.stringify(metadataData));
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);

  // Create an AbortController to manage a custom timeout
  const controller = new AbortController();
  const signal = controller.signal;

  // Set a custom timeout (for example, 5 minutes = 300000 ms)
  const timeout = 300000;
  const timeoutId = setTimeout(() => {
    console.warn('Fetch timeout reached. Aborting request...');
    controller.abort();
  }, timeout);

  // Start the fetch with the abort signal
  const response = await fetch(`${baseURL}/report/generate`, {
    method: 'POST',
    body: formData,
    credentials: 'include',
    headers: {
      Accept: 'application/x-ndjson',
    },
    signal,
  });

  // Clear the timeout once the response is received
  clearTimeout(timeoutId);

  if (response.status === 401) {
    throw new Error('User is not authorized');
  }
  if (!response.ok) {
    throw new Error(`Error with status code: ${String(response.status)}`);
  }
  if (!response.body) {
    throw new Error('ReadableStream not supported.');
  }

  console.log('✅ Starting to receive streaming data...');
  const reader = response.body.getReader();
  return reader;
}
