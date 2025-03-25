import { useState } from 'react';
import { Text } from '@mantine/core';
import { LiveAudioVisualizer } from 'react-audio-visualize';
import WarningModal from './warningModal';
import PatientInfoModal from './PatientInfoModal';
import CaptureConversationButton from './captureConversation';
import ControlButtons from './audioControls';
import useAudioRecorder from '../../hooks/useAudioRecorder';
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
  const [captureModalOpen, setCaptureModalOpen] = useState(false);
  const [warningModalOpen, setWarningModalOpen] = useState(false);
  const [patientName, setPatientName] = useState('');
  const [duration, setDuration] = useState(0);

  const {
    isRecording,
    mediaRecorder,
    handleStartRecording,
    handlePauseRecording,
    handleResumeRecording,
    handleStopRecording,
    recordingStartTime,
    audioBlobRef,
    handleResetMediaRecorder,
  } = useAudioRecorder();

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
    onSuccess: async reader => {
      setPatientName('');
      await processStream(reader);
    },
    onError: error => {
      console.error('Error generating report:', error);
    },
  });

  const handleGenerateReport = (duration: number) => {
    if (audioBlobRef.current == null) {
      throw Error('cannot generate Report when Blob file is null');
    }
    if (recordingStartTime.current == null) {
      throw Error('timestamp cannot be null');
    }
    const formData = convertBlobToFormData(audioBlobRef.current);
    const reportTime = new Date(recordingStartTime.current).toISOString();
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

    generateReportMutation.mutate({ formData, metadata });
  };

  const handleEndVisit = async () => {
    if (!recordingStartTime.current) return;

    const elapsedSeconds = (Date.now() - recordingStartTime.current) / 1000;
    if (elapsedSeconds < 0) {
      setWarningModalOpen(true);
      handlePauseRecording();
    } else {
      setDuration(elapsedSeconds);
      await handleStopRecording();
      handleGenerateReport(elapsedSeconds);
      handleResetMediaRecorder();
    }
  };

  return (
    <div>
      {!isRecording && !mediaRecorder && (
        <CaptureConversationButton
          onClick={() => {
            setCaptureModalOpen(true);
          }}
        />
      )}

      {mediaRecorder && (
        <div>
          <Text size="lg" className="mb-2 text-center">
            {isRecording
              ? 'Listening to patient conversation. Keep this screen open, please.'
              : 'Press "End Visit" to generate the report or "Resume" to continue the visit.'}
          </Text>
          <LiveAudioVisualizer
            mediaRecorder={mediaRecorder}
            width={500}
            height={75}
          />
        </div>
      )}
      {mediaRecorder && (
        <ControlButtons
          isRecording={isRecording}
          onEndVisit={handleEndVisit}
          onPause={handlePauseRecording}
          onReset={handleResetMediaRecorder}
          onResume={() => {
            handleResumeRecording();
          }}
        />
      )}

      <WarningModal
        isOpen={warningModalOpen}
        onClose={() => {
          setWarningModalOpen(false);
        }}
      />

      <PatientInfoModal
        isOpen={captureModalOpen}
        patientName={patientName}
        onClose={() => {
          setCaptureModalOpen(false);
        }}
        onSubmit={(name: string) => {
          if (name.trim()) {
            setCaptureModalOpen(false);
            void handleStartRecording();
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
