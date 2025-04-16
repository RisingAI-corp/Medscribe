import { Flex } from '@mantine/core';
import SoapSectionBox from './SoapSectionBox/soapSectionBox';
import { useAtom } from 'jotai';
import { SoapAtom, updateTranscriptAtom } from './derivedAtoms';
import { useEffect } from 'react';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { replaceReportAtom } from './derivedAtoms';
import { getReport, GetReportPayload } from '../../api/getReport';
import { useMutation } from '@tanstack/react-query';
import { reportStreamingAtom } from '../../states/patientsAtom';
import { learnStyle } from '../../api/learnStyle';
import { updateContentSection } from '../../api/updateContentSection';
import { getTranscript } from '../../api/getTranscript';
import TranscriptAccordion from './SoapSectionBox/transcriptAccordian';

function VisitReportLayout() {
  const [soapData, updateSoapData] = useAtom(SoapAtom);
  const [selectedPatientID, _] = useAtom(currentlySelectedPatientAtom);
  const [__, replaceReport] = useAtom(replaceReportAtom);
  const [reportStreaming, ___] = useAtom(reportStreamingAtom);
  const [, updateTranscript] = useAtom(updateTranscriptAtom); // Renamed for clarity

  const getReportMutation = useMutation({
    mutationFn: async (props: GetReportPayload) => {
      console.log('yoski', props);
      const report = await getReport(props);
      replaceReport(report);
      if (report.status === 'pending') {
        throw new Error('Report not finished generating yet.');
      }
      return report;
    },
    retry: 10,
    retryDelay: 4000,
    onSuccess: report => {
      replaceReport(report);
    },
    onError: error => {
      console.error('Error generating report (will retry):', error);
    },
  });

  const learnStyleMutation = useMutation({
    mutationFn: learnStyle,
    onError: error => {
      console.error('Error learning style:', error);
    },
  });

  const updateContentSectionMutation = useMutation({
    mutationFn: updateContentSection,
    onError: error => {
      console.error('Error updating content section:', error);
    },
  });

  const fetchTranscriptMutation = useMutation({
    mutationFn: getTranscript,
    onSuccess: transcriptContainer => {
      console.log('fetched transcript', transcriptContainer);
      updateTranscript({
        id: selectedPatientID,
        transcriptContainer: transcriptContainer, // Assuming the API returns the DiarizedTranscript array directly
      });
      console.log(transcriptContainer, 'success here is transcript');
    },
    onError: error => {
      console.error('Error getting transcripts:', error);
    },
  });

  useEffect(() => {
    console.log(soapData?.status, 'here is status');
    console.log(
      selectedPatientID,
      !reportStreaming.has(selectedPatientID), // if its not streaming
      soapData?.status === 'pending',
      'checking status',
    );
    if (
      selectedPatientID &&
      !reportStreaming.has(selectedPatientID) &&
      soapData?.status === 'pending'
    ) {
      console.log('getting report');
      const payload: GetReportPayload = {
        reportID: selectedPatientID,
      };
      getReportMutation.mutate(payload);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  console.log(soapData, 'here is soap data');
  const handleSoapDataUpdate = (
    field: string,
    newData: string,
    reportID: string,
  ) => {
    updateSoapData({ patientId: reportID, field, newData });
    updateContentSectionMutation.mutate({
      ReportID: reportID,
      ContentSection: field,
      Content: newData,
    });
  };

  const handleLearnFormat = (
    contentSection: string,
    previous: string,
    current: string,
  ) => {
    learnStyleMutation.mutate({
      ReportID: selectedPatientID,
      ContentSection: contentSection,
      Previous: previous,
      Current: current,
    });
  };

  const handleGetTranscript = () => {
    fetchTranscriptMutation.mutate({ reportID: selectedPatientID });
  };

  const { isPending, isSuccess } = fetchTranscriptMutation;
  const { variables } = updateContentSectionMutation;

  return (
    <>
      {soapData && (
        <Flex direction="column" gap="xl">
          {soapData.soapContent.map(section => (
            <SoapSectionBox
              key={`${selectedPatientID}-${section.title}-${section.content.data}`}
              reportID={selectedPatientID}
              title={section.title}
              sectionType={section.sectionType}
              text={section.content.data}
              isLoading={section.content.loading}
              handleSave={handleSoapDataUpdate}
              handleLearnFormat={handleLearnFormat}
              onExpand={() => {
                return;
              }}
              isExpanded={true}
              readonly={false}
              isContentSaveLoading={
                updateContentSectionMutation.isPending &&
                updateContentSectionMutation.variables.ContentSection ===
                  section.sectionType
              }
              isLearnStyleLoading={
                learnStyleMutation.isPending &&
                learnStyleMutation.variables.ContentSection ===
                  section.sectionType
              }
            />
          ))}
          {soapData.status === 'success' &&
            (soapData.transcriptContainer.usedDiarization ? (
              <TranscriptAccordion
                reportID={selectedPatientID}
                title={'Full Transcript'}
                transcriptTurns={
                  soapData.transcriptContainer.diarizedTranscript
                }
                isLoading={isPending}
                onExpand={handleGetTranscript}
                isExpanded={isPending || isSuccess}
                readonly={true}
              />
            ) : (
              <SoapSectionBox
                key={`${selectedPatientID}-Full Transcript`}
                reportID={selectedPatientID}
                title={'Full Transcript'}
                text={soapData.transcriptContainer.transcript}
                isLoading={isPending}
                onExpand={() => {
                  handleGetTranscript();
                }}
                isExpanded={isPending || isSuccess}
                readonly={true}
              />
            ))}
        </Flex>
      )}
    </>
  );
}

export default VisitReportLayout;
