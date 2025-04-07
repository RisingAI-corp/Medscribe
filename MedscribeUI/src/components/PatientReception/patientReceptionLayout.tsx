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
import { Tooltip, Button } from '@mantine/core';
import FollowUpSearchModalLayout from '../FollowUpSearchModal/FollowUpSearchModalLayout';
import { useDebouncedNameChange } from '../../hooks/useDebounceNameChange';
import { SearchResultItem } from '../FollowUpSearchModal/SearchResults/SearchResults';
import { IconExternalLink } from '@tabler/icons-react';
import PatientBackgroundDetails from '../PatientBackground/PatientBackground';

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

  const processStream = useStreamProcessor({
    attemptCreateReport,
    updateReports,
    providerID: provider.ID,
    patientName,
    timestamp,
    duration,
  });

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

  console.log(visitSearchValue);

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
    <>
      <div className="bg-white border-b border-gray-200 shadow-sm">
        <div className="flex items-center justify-between px-6 py-4">
          {/* Left side controls */}
          <div className="flex items-center gap-4">
            {/* Patient name input with tooltip */}
            <Tooltip
              label="Name field cannot be empty"
              opened={isEmpty}
              position="bottom"
              withArrow
            >
              <input
                type="text"
                ref={nameRef}
                value={nameValue}
                onChange={e => {
                  setNameValue(e.target.value);
                  debouncedNameChange(e.target.value);
                }}
                placeholder="Add or select patient"
                required
                className={`rounded-md px-3 py-2 text-sm font-medium border ${isEmpty ? 'border-red-500' : 'border-gray-300'} 
                          focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition`}
              />
            </Tooltip>

            {/* Link previous visit */}
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

            {/* Pronoun Selector */}
            <PronounSelector pronoun={pronoun} setPronoun={setPronoun} />
          </div>

          {/* Right side - Capture button */}
          <AudioControlLayout onAudioCaptured={handleAudioCaptured} />
        </div>
      </div>

      {lastVisitContext && <PatientBackgroundDetails {...lastVisitContext} />}

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
