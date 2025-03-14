import { Flex } from '@mantine/core';
import SoapSectionBox from './SoapSectionBox/soapSectionBox';
import { useAtom } from 'jotai';
import { SoapAtom } from './derivedAtoms';
import { useEffect } from 'react';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { replaceReportAtom } from './derivedAtoms';
import { getReport, GetReportPayload } from '../../api/getReport';
import { useMutation } from '@tanstack/react-query';
import { reportStreamingAtom } from '../../states/patientsAtom';
import { learnStyle } from '../../api/learnStyle';
import { updateContentSection } from '../../api/updateContentSection';

function VisitReportLayout() {
  const [soapData, updateSoapData] = useAtom(SoapAtom);

  const [selectedPatientID, _] = useAtom(currentlySelectedPatientAtom);
  const [__, replaceReport] = useAtom(replaceReportAtom);
  const [reportStreaming, ___] = useAtom(reportStreamingAtom);

  const getReportMutation = useMutation({
    mutationFn: async (props: GetReportPayload) => {
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

  useEffect(() => {
    if (!reportStreaming.has(selectedPatientID) && !soapData?.loading) {
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

  const handleLearnFormat = (contentSection: string, content: string) => {
    learnStyleMutation.mutate({
      ReportID: selectedPatientID,
      ContentSection: contentSection,
      Content: content,
    });
  };

  return (
    <>
      <Flex direction="column" gap="xl">
        {soapData?.content.map(section => (
          <SoapSectionBox
            key={`${selectedPatientID}-${section.type}`}
            reportID={selectedPatientID}
            title={section.type}
            text={section.content.data}
            isLoading={section.content.loading}
            handleSave={handleSoapDataUpdate}
            handleLearnFormat={handleLearnFormat}
          />
        ))}
      </Flex>
    </>
  );
}

export default VisitReportLayout;
