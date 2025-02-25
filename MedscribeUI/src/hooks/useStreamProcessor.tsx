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

export interface UseStreamProcessorOptions {
  attemptCreateReport?: (report: CreateReportProps) => void;
  updateReports: (update: UpdateProps) => void;
  providerID: string;
  patientName?: string;
  timeStamp?: string;
  duration?: number;
  reportID?: string;
}

export function useStreamProcessor({
  attemptCreateReport,
  updateReports,
  patientName = '',
  providerID = '',
  timeStamp = '',
  duration = 0,
  reportID = '',
}: UseStreamProcessorOptions) {
  const [_, setReportStreamStatus] = useAtom(setReportStreamStatusAtom);
  const [__, unsetReportStreamStatus] = useAtom(unsetReportStreamStatusAtom);

  const processStream = useCallback(
    async (eventReader: ReadableStreamDefaultReader<Uint8Array>) => {
      const decoder = new TextDecoder();
      let buffer = '';
      console.log('Starting to receive streaming data...');

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
              console.log('Parsed final chunk:', data);
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
          console.log('✅ Server sent EOF (End of Stream).');
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
                reportID = data.Value;
                setReportStreamStatus(reportID);
                // For new report creation (or when regenerating an existing report)
                attemptCreateReport({
                  id: data.Value,
                  providerID,
                  name: patientName,
                  timeStamp,
                  duration,
                });
              } else {
                throw new Error(
                  'Error parsing update from SSE: id needs to be string',
                );
              }
            } else {
              updateReports({ id: reportID, ...data });
            }
            console.log('Parsed chunk:', data);
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
      timeStamp,
      duration,
    ],
  );

  return processStream;
}
