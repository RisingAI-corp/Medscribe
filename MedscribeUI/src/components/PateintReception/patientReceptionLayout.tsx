import { useState } from 'react';
import { Text } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import WarningModal from './warningModal';
import PatientInfoModal from './PatientInfoModal';
import CaptureConversationButton from './captureConversation';
import AudioControlLayout from '../AudioControl/audioControlLayout';
import { useMutation } from '@tanstack/react-query';
import {
  generateReport,
  GenerateReportMetadata,
} from '../../api/generateReport';
import { userAtom } from '../../states/userAtom';
import { UpdateReportsAtom, createReportAtom } from './derivedAtoms';

import { useAtom } from 'jotai';
import { useStreamProcessor } from '../../hooks/useStreamProcessor';

const PatientReception = () => {
  const [__, updateReports] = useAtom(UpdateReportsAtom);
  const [___, attemptCreateReport] = useAtom(createReportAtom);
  const [provider, _____] = useAtom(userAtom);

  const [timestamp, setTimeStamp] = useState<string>('');
  const [patientInfoModalOpen, setPatientInfoModalOpen] = useState(false);
  const [warningModalOpen, setWarningModalOpen] = useState(false);
  const [patientName, setPatientName] = useState('');
  const [duration, setDuration] = useState(0);

  const processStream = useStreamProcessor({
    attemptCreateReport,
    updateReports,
    providerID: provider.ID,
    patientName,
    timestamp,
    duration,
  });

  const convertBlobToFormData = (blob: Blob) => {
    const file = new File([blob], 'audio', { type: 'audio/webm' });

    const formData = new FormData();
    formData.append('audio', file);

    return formData;
  };

  const generateReportMutation = useMutation({
    mutationFn: ({
      formData,
      metadata,
    }: {
      formData: FormData;
      metadata: GenerateReportMetadata;
    }) => generateReport(formData, metadata),
    onMutate: () => {
      notifications.show({
        loading: true,
        title: 'Generating Report',
        message: 'This make take anywhere from 10 seconds to 2 minutes',
        autoClose: 5000,
        withCloseButton: false,
      });
    },
    onSuccess: async reader => {
      setPatientName('');
      await processStream(reader);
    },
    onError: error => {
      console.error('Error generating report:', error);
    },
  });

  const handleGenerateReport = (
    duration: number,
    recordingTime: number,
    audioBlob: Blob,
  ) => {
    if (!patientName) {
      setPatientInfoModalOpen(true);
      return;
    }

    const reportTime = new Date(recordingTime).toISOString();
    setTimeStamp(reportTime);

    const metadata: GenerateReportMetadata = {
      providerName: provider.name,
      patientName: patientName,
      timestamp: reportTime,
      duration: duration,
      subjectiveStyle: provider.subjectiveStyle,
      objectiveStyle: provider.objectiveStyle,
      assessmentAndPlanStyle: provider.assessmentAndPlanStyle,
      patientInstructionsStyle: provider.patientInstructionsStyle,
      summaryStyle: provider.summaryStyle,
    };

    generateReportMutation.mutate({
      formData: convertBlobToFormData(audioBlob),
      metadata,
    });
  };

  const handleAudioCaptured = (
    blob: Blob,
    duration: number,
    timestamp: number,
  ) => {
    if (duration < 0) {
      setWarningModalOpen(true);
      return;
    }

    setDuration(duration);
    handleGenerateReport(duration, timestamp, blob);
  };

  return (
    <div>
      {!patientName && (
        <CaptureConversationButton
          onClick={() => {
            setPatientInfoModalOpen(true);
          }}
        />
      )}

      {patientName && (
        <div>
          <Text size="sm" className="mb-2 text-center">
            Patient: {patientName}
          </Text>
          <AudioControlLayout onAudioCaptured={handleAudioCaptured} />
        </div>
      )}

      <WarningModal
        isOpen={warningModalOpen}
        onClose={() => {
          setWarningModalOpen(false);
        }}
      />

      <PatientInfoModal
        isOpen={patientInfoModalOpen}
        patientName={patientName}
        onClose={() => {
          setPatientInfoModalOpen(false);
        }}
        onSubmit={(name: string) => {
          if (name.trim()) {
            setPatientInfoModalOpen(false);
            setPatientName(name);
          } else {
            alert('Please enter a valid patient name.');
          }
        }}
      />
    </div>
  );
};

export default PatientReception;
