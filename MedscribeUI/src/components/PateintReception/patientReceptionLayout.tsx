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

  const [timeStamp, setTimeStamp] = useState<string>('');
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
    timeStamp,
    duration,
  });

  const convertBlobToFormData = (blob: Blob) => {
    console.log('before conversion ');
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
      await processStream(reader);
      setPatientName('');
    },
    onError: error => {
      console.error('Error generating report:', error);
    },
  });

  const handleGenerateReport = (duration: number) => {
    console.log('before conversion ', audioBlobRef.current);
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
      providerID: provider.ID,
      patientName: patientName,
      timestamp: reportTime,
      duration: duration,
      subjectiveStyle: provider.subjectiveStyle,
      objectiveStyle: provider.objectiveStyle,
      assessmentStyle: provider.assessmentStyle,
      planningStyle: provider.planningStyle,
      summaryStyle: provider.summaryStyle,
    };

    generateReportMutation.mutate({ formData, metadata });
  };

  const handleEndVisit = async () => {
    console.log('this is the current audio blob ', audioBlobRef.current);
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

      <ControlButtons
        isRecording={isRecording}
        mediaRecorder={mediaRecorder}
        onEndVisit={handleEndVisit}
        onPause={handlePauseRecording}
        onResume={() => {
          handleResumeRecording();
        }}
      />

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
        onChange={value => {
          setPatientName(value);
        }}
        onSubmit={() => {
          if (patientName.trim()) {
            setCaptureModalOpen(false);
            void handleStartRecording();
          } else {
            alert('Please enter a valid patient name.');
          }
        }}
      />
    </div>
  );
};

export default PatientReception;
