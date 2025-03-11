import { useCallback } from 'react';
import {
  CreateReportProps,
  UpdateProps,
} from '../components/PateintReception/derivedAtoms';
import { UpdateResponse } from '../api/serverResponses';
import {
  setReportStreamStatusAtom,
  unsetReportStreamStatusAtom,
} from '../states/patientsAtom';

import { useAtom } from 'jotai';
import { currentlySelectedPatientAtom } from '../states/patientsAtom';

export interface UseStreamProcessorOptions {
  attemptCreateReport?: (report: CreateReportProps) => void;
  updateReports: (update: UpdateProps) => void;
  providerID: string;
  patientName?: string;
  timestamp?: string;
  duration?: number;
  reportID?: string;
}

export function useStreamProcessor({
  attemptCreateReport,
  updateReports,
  patientName = '',
  providerID = '',
  timestamp = '',
  duration = 0,
  reportID = '',
}: UseStreamProcessorOptions) {
  const [_, setReportStreamStatus] = useAtom(setReportStreamStatusAtom);
  const [__, unsetReportStreamStatus] = useAtom(unsetReportStreamStatusAtom);
  const [currentlySelectedPatient, setCurrentlySelectedPatient] = useAtom(
    currentlySelectedPatientAtom,
  );

  const processStream = useCallback(
    async (eventReader: ReadableStreamDefaultReader<Uint8Array>) => {
      const decoder = new TextDecoder();
      let buffer = '';

      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      while (true) {
        const { value, done } = await eventReader.read();
        if (done) {
          // Process any remaining data in the buffer
          if (buffer.trim()) {
            try {
              const { success, data, error } = UpdateResponse.safeParse(buffer);
              if (!success) {
                throw new Error(
                  'Error parsing API request: ' + error.toString(),
                );
              }
              updateReports({ id: reportID, ...data });
              unsetReportStreamStatus(reportID);
            } catch (err) {
              console.error(
                'Error parsing final JSON:',
                err,
                'Buffer:',
                buffer,
              );
            }
          }
          break;
        }
        buffer += decoder.decode(value, { stream: true });

        // Split into NDJSON lines
        const lines = buffer.split('\n');
        // The last element may be incomplete; keep it in the buffer
        buffer = lines.pop() ?? '';

        for (const line of lines) {
          const trimmedLine = line.trim();
          if (!trimmedLine) continue;
          try {
            const { success, data, error } = UpdateResponse.safeParse(
              JSON.parse(trimmedLine),
            );
            if (!success) {
              throw new Error('Error parsing API request: ' + error.toString());
            }
            if (data.Key === '_id') {
              if (typeof data.Value === 'string' && attemptCreateReport) {
                // eslint-disable-next-line react-hooks/exhaustive-deps
                reportID = data.Value;
                setReportStreamStatus(reportID);
                // For new report creation (or when regenerating an existing report)
                attemptCreateReport({
                  id: data.Value,
                  providerID,
                  name: patientName,
                  timestamp,
                  duration,
                });
              } else {
                throw new Error(
                  'Error parsing update from SSE: id needs to be string',
                );
              }
            } else {
              if (currentlySelectedPatient === '') {
                setCurrentlySelectedPatient(reportID);
              }
              updateReports({ id: reportID, ...data });
            }
          } catch (err) {
            console.error('Error parsing JSON:', err, 'Line:', trimmedLine);
          }
        }
      }
    },
    [
      updateReports,
      unsetReportStreamStatus,
      setReportStreamStatus,
      attemptCreateReport,
      providerID,
      patientName,
      timestamp,
      duration,
    ],
  );

  return processStream;
}
