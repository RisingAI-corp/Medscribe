import { useState } from 'react';
import { notifications } from '@mantine/notifications';
import WarningModal from './warningModal';
import PronounSelector from './PronounSelector/PronounSelector';
import AudioControlLayout from './AudioControl/audioControlLayout';
import { useMutation } from '@tanstack/react-query';
import {
  generateReport,
  GenerateReportMetadata,
} from '../../api/generateReport';
import { userAtom } from '../../states/userAtom';
import { UpdateReportsAtom, createReportAtom } from './derivedAtoms';

import { useAtom } from 'jotai';
import { useStreamProcessor } from '../../hooks/useStreamProcessor';
import { Button } from '@mantine/core';
import FollowUpSearchModalLayout from '../FollowUpSearchModal/FollowUpSearchModalLayout';
import { useDebouncedNameChange } from '../../hooks/useDebounceNameChange';
import { SearchResultItem } from '../FollowUpSearchModal/SearchResults/SearchResults';
import { IconExternalLink, IconX } from '@tabler/icons-react';
import PatientBackgroundDetails from '../PatientBackground/PatientBackground';
import useAudioRecorder from '../../hooks/useAudioRecorder';
import { TOP_CENTER } from '../../constants';

const PatientReception = () => {
  const [__, updateReports] = useAtom(UpdateReportsAtom);
  const [___, attemptCreateReport] = useAtom(createReportAtom);
  const [provider, _____] = useAtom(userAtom);

  const [timestamp, setTimeStamp] = useState<string>('');
  const [_, setPatientInfoModalOpen] = useState(false);
  const [warningModalOpen, setWarningModalOpen] = useState(false);
  const [patientName, setPatientName] = useState('');
  const [duration, setDuration] = useState(0);
  const [pronoun, setPronoun] = useState('');
  const [lastVisitID, setLastVisitID] = useState('');
  const [lastVisitContext, setLastVisitContext] = useState<SearchResultItem>();
  const [visitSearchValue, setVisitSearchValue] = useState('');
  const [visibleToolTip, setVisibleToolTip] = useState(false);

  const processStream = useStreamProcessor({
    attemptCreateReport,
    updateReports,
    providerID: provider.ID,
    patientName,
    timestamp,
    duration,
  });

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

  const { nameRef, nameValue, setNameValue, debouncedNameChange } =
    useDebouncedNameChange({
      name: patientName,
      onChange: setPatientName,
      handleUpdateName: () => {
        console.log('Wha');
      },
    });

  const isEmpty = !nameValue.trim();

  const convertBlobToFormData = (blob: Blob) => {
    const file = new File([blob], 'audio', { type: 'audio/webm' });

    const formData = new FormData();
    formData.append('audio', file);

    return formData;
  };

  const handleVisitContextSelect = (visitContext: SearchResultItem) => {
    setNameValue(visitContext.patientName);
    setVisitSearchValue(visitContext.patientName);
    setLastVisitID(visitContext.id);
    setLastVisitContext(visitContext);
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
      lastVisitID: lastVisitID,
      visitContext: lastVisitContext?.summary ?? '',
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
    console.log(blob, 'here is the blob');
    setDuration(duration);
    handleResetMediaRecorder();
    handleGenerateReport(duration, timestamp, blob);
    setLastVisitContext(undefined);
  };

  const handleEndVisit = async () => {
    if (!recordingStartTime.current) return;
    if (patientName === '') {
      console.log('Patient name is empty');
      setVisibleToolTip(true);
      await handlePauseRecording();
      notifications.show({
        title: `Cannot Generate Report Yet`,
        message: `Please include patient name`,
        position: TOP_CENTER,
        color: 'red',
        icon: <IconX />,
        autoClose: 1000,
      });
      return;
    }

    const elapsedSeconds = (Date.now() - recordingStartTime.current) / 1000;
    if (elapsedSeconds <= 0) {
      await handlePauseRecording();
      return;
    }

    await handlePauseRecording();
    if (elapsedSeconds < 0) {
      setWarningModalOpen(true);
      return;
    }
    await handleStopRecording();
    if (audioBlobRef.current) {
      console.log(audioBlobRef.current, 'here is the blob');
      handleAudioCaptured(
        audioBlobRef.current,
        elapsedSeconds,
        recordingStartTime.current,
      );
    }
    handleResetMediaRecorder();
  };

  console.log('tool', visibleToolTip);

  return (
    <>
      <div className="bg-white border-b border-gray-200 shadow-sm">
        <div className="flex items-center justify-between px-6 py-4">
          {/* Left side controls */}
          <div className="flex items-center gap-20">
            {/* Group: input + or + link visit */}
            <div className="flex items-center gap-2">
              <div className="inline-block">
                <input
                  type="text"
                  ref={nameRef}
                  value={nameValue}
                  onChange={e => {
                    setNameValue(e.target.value);
                    debouncedNameChange(e.target.value);
                  }}
                  placeholder="Add a New Patient"
                  required
                  className={`rounded-md px-3 py-2 text-sm font-medium border ${
                    isEmpty && visibleToolTip
                      ? 'border-red-500'
                      : 'border-gray-300'
                  } focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition`}
                />
              </div>

              <span className="text-sm text-gray-500">or</span>

              <FollowUpSearchModalLayout
                handleSelectedVisit={handleVisitContextSelect}
              >
                <Button
                  rightSection={<IconExternalLink size={16} />}
                  variant="light"
                  color="blue"
                  className="h-[36px] text-sm"
                >
                  {visitSearchValue || 'Link Visit'}
                </Button>
              </FollowUpSearchModalLayout>
            </div>

            {/* Spacer */}
            <div className="ml-6">
              <PronounSelector pronoun={pronoun} setPronoun={setPronoun} />
            </div>
          </div>

          {/* Right side - Capture button */}
          <AudioControlLayout
            isRecording={isRecording}
            mediaRecorder={mediaRecorder}
            audioBlobRef={audioBlobRef}
            handleStartRecording={handleStartRecording}
            handlePauseRecording={handlePauseRecording}
            handleResumeRecording={handleResumeRecording}
            handleStopRecording={handleStopRecording}
            recordingStartTime={recordingStartTime}
            handleResetMediaRecorder={handleResetMediaRecorder}
            onAudioCaptured={handleAudioCaptured}
            handleEndVisit={handleEndVisit}
          />
        </div>
      </div>

      {lastVisitContext?.summary && (
        <div className="flex justify-center mt-8 px-4">
          <div
            className="w-full max-w-6xl sm:max-w-4xl md:max-w-3xl lg:max-w-2xl 
                 transform transition-all duration-500 ease-out 
                 animate-fade-in-slide-up"
          >
            <PatientBackgroundDetails {...lastVisitContext} />
          </div>
        </div>
      )}

      {/* Warning Modal */}
      <WarningModal
        isOpen={warningModalOpen}
        onClose={() => {
          setWarningModalOpen(false);
        }}
      />
    </>
  );
};

export default PatientReception;
