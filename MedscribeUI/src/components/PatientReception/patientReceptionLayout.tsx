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

const PatientReception = () => {
  const [__, updateReports] = useAtom(UpdateReportsAtom);
  const [___, attemptCreateReport] = useAtom(createReportAtom);
  const [provider, _____] = useAtom(userAtom);

  const [timestamp, setTimeStamp] = useState<string>('');
  const [_, setPatientInfoModalOpen] = useState(false);
  const [warningModalOpen, setWarningModalOpen] = useState(false);
  const [patientName, setPatientName] = useState('TEMP'); // TODO: Remove 'TEMP' to ''
  const [duration, setDuration] = useState(0);
  const [pronoun, setPronoun] = useState('');
  const [lastVisitID, setLastVisitID] = useState('');
  const [lastVisitContext, setLastVisitContext] = useState('');
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
      handleUpdateName: () => {},
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
    setLastVisitContext(visitContext.summary);
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
      visitContext: lastVisitContext,
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
      <div className="flex flex-row gap-2 justify-between w-full px-8">
        <div className="flex flex-row gap-2">
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
              placeholder="Enter patient's name"
              required
              className={`border-b-2 ${isEmpty ? 'border-red-500' : 'border-gray-400'} 
                          focus:outline-none hover:border-blue-700 focus:border-blue-500 pl-0 pb-1 pt-1 text-sm bg-transparent font-bold`}
              style={{ width: '10rem' }}
            />
          </Tooltip>

          <FollowUpSearchModalLayout handleSelectedVisit={handleVisitContextSelect}>
            <Button 
              rightSection={<IconExternalLink size={16} />}
              variant="light"
              color="blue"
              fullWidth
              className="h-[60px]"
            >
              {visitSearchValue ? visitSearchValue : 'Link Visit'}
            </Button>
          </FollowUpSearchModalLayout>
          
          <PronounSelector pronoun={pronoun} setPronoun={setPronoun} />
        </div>

        
        <AudioControlLayout onAudioCaptured={handleAudioCaptured} />

        
      </div>

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
