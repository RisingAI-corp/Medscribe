import { Flex } from '@mantine/core';
import SoapSectionBox from './SoapSectionBox/soapSectionBox';
import { useAtom } from 'jotai';
import { SoapAtom } from './derivedAtoms';
import { useEffect, useState } from 'react';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { replaceReportAtom } from './derivedAtoms';
import { getReport, GetReportPayload } from '../../api/getReport';
import { useMutation } from '@tanstack/react-query';
import { reportStreamingAtom } from '../../states/patientsAtom';
import { learnStyle } from '../../api/learnStyle';
import { userAtom } from '../../states/userAtom';
import { updateContentSection } from '../../api/updateContentSection';

function VisitReportLayout() {
  const [soapData, updateSoapData] = useAtom(SoapAtom);

  const [selectedPatientID, _] = useAtom(currentlySelectedPatientAtom);
  const [__, replaceReport] = useAtom(replaceReportAtom);
  const [reportStreaming, ___] = useAtom(reportStreamingAtom);
  const [laggingSelectedPatient, setLaggingSelectedPatient] =
    useState(selectedPatientID);

  const getReportMutation = useMutation({
    mutationFn: async (props: GetReportPayload) => {
      console.log('fetching');
      const report = await getReport(props);
      console.log('this is report ', report);
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
      console.log('Report finished generating:', report);
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
    onSuccess: data => {
      console.log('success');
    },
    onError: error => {
      console.error('Error updating content section:', error);
    },
  });

  useEffect(() => {
    if (!reportStreaming.has(selectedPatientID) && !soapData?.loading) {
      console.log('ran', reportStreaming);
      const payload: GetReportPayload = {
        reportID: selectedPatientID,
      };
      getReportMutation.mutate(payload);
    }
  }, []);

  const handleSoapDataUpdate = (field: string, newData: string) => {
    console.log('hello', newData);
    // TODO: make api call later
    updateSoapData({ patientId: selectedPatientID, field, newData });
    updateContentSectionMutation.mutate({
      ReportID: selectedPatientID,
      ContentSection: field,
      Content: newData,
    });
  };

  const handleSoapDataAutoUpdate = (field: string, newData: string) => {
    // TODO: make api call later
    updateSoapData({ patientId: laggingSelectedPatient, field, newData });
    setLaggingSelectedPatient(selectedPatientID);
    updateContentSectionMutation.mutate({
      ReportID: selectedPatientID,
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
            title={section.type}
            text={section.content.data}
            isLoading={section.content.loading}
            handleSave={(newText: string) => {
              handleSoapDataUpdate(section.type, newText);
            }}
            handleAutoSave={(newText: string) => {
              handleSoapDataAutoUpdate(section.type, newText);
            }}
            handleLearnFormat={handleLearnFormat}
          />
        ))}
      </Flex>
    </>
  );
}

export default VisitReportLayout;
