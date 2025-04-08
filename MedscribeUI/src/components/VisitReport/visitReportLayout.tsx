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

function VisitReportLayout() {
  const [soapData, updateSoapData] = useAtom(SoapAtom);

  const [selectedPatientID, _] = useAtom(currentlySelectedPatientAtom);
  const [__, replaceReport] = useAtom(replaceReportAtom);
  const [reportStreaming, ___] = useAtom(reportStreamingAtom);
  const [____, updateTranscript] = useAtom(updateTranscriptAtom);

  const getReportMutation = useMutation({
    mutationFn: async (props: GetReportPayload) => {
      console.log('yoski', props);
      const report = await getReport(props);
      replaceReport(report);
      if (!report.finishedGenerating) {
        throw new Error('Report not finished generating yet.');
      }
      return report;
    },
    retry: 20,
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

    onSuccess: transcript => {
      updateTranscript({ id: selectedPatientID, transcript: transcript });
    },

    onError: error => {
      console.error('Error getting transcripts:', error);
    },
  });

  useEffect(() => {
    if (
      selectedPatientID &&
      !reportStreaming.has(selectedPatientID) &&
      !soapData?.loading
    ) {
      const payload: GetReportPayload = {
        reportID: selectedPatientID,
      };

      getReportMutation.mutate(payload);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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
  console.log(variables?.ContentSection, 'here here');
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
          {soapData.loading && (
            <SoapSectionBox
              key={`${selectedPatientID}-${soapData.transcript}`}
              reportID={selectedPatientID}
              title={'Full Transcript'}
              text={soapData.transcript}
              isLoading={isPending}
              onExpand={() => {
                handleGetTranscript();
              }}
              isExpanded={isPending || isSuccess}
              readonly={true}
            />
          )}
        </Flex>
      )}
    </>
  );
}

export default VisitReportLayout;
