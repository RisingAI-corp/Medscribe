export interface Updates {
  Key: string;
  Value: unknown;
}
export interface ReportContentSection {
  ContentType: string;
  Content: string;
}

export interface RegenerateReportMetadata {
  ID: string;
  subjectiveStyle: string;
  objectiveStyle: string;
  summaryStyle: string;
  assessmentAndPlanStyle: string;
  patientInstructionsStyle: string;
  updates: Updates[];
  subjectiveContent: string;
  objectiveContent: string;
  assessmentAndPlanContent: string;
  patientInstructionsContent: string;
  summaryContent: string;
  lastVisitID?: string;
  visitContext?: string;
}

export async function regenerateReport(payload: RegenerateReportMetadata) {
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
  const response = await fetch(`${baseURL}/report/regenerate`, {
    method: 'PATCH',
    body: JSON.stringify(payload),
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
  if (response.status === 500) {
    throw new Error('There was an internal server error ');
  }
  if (!response.ok) {
    throw new Error(`Error with status code: ${String(response.status)}`);
  }
  if (!response.body) {
    throw new Error('ReadableStream not supported.');
  }

  const reader = response.body.getReader();
  return reader;
}
