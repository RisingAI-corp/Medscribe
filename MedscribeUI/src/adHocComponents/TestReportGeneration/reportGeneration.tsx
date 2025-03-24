/* 

the purpose of this sandbox is to allow for quick visible back of
what the generateReport pipeline feeds us given that the UI loads to slowly currently

I suggest removing this, once the frontend starts to load quickly

*/

import { Button } from '@mantine/core';
import { useAtom } from 'jotai';
import { useMutation } from '@tanstack/react-query';
import {
  generateReport,
  GenerateReportMetadata,
} from '../../api/generateReport';
import ProfileSummaryCard from '../../components/ProfileSummaryCard/profileSummaryCard';
import {
  UpdateReportsAtom,
  createReportAtom,
} from '../../components/PateintReception/derivedAtoms';
import { useStreamProcessor } from '../../hooks/useStreamProcessor';
import { Input } from '@mantine/core';

import { SelectedPatientHeaderInformationAtom } from '../../components/PaitentDashboard/derivedAtoms';
import { UpdateSelectedPatientNameAtom } from '../../states/patientsAtom';
import { userAtom } from '../../states/userAtom';
import VisitReportLayout from '../../components/VisitReport/visitReportLayout';
import NoteControlsLayout from '../../components/NoteControls/noteControlsLayout';
import { MIME_TYPE } from '../../constants';
import { useEffect } from 'react';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';

const GenerateReportTest = () => {
  const [headerInformation] = useAtom(SelectedPatientHeaderInformationAtom);
  const [, updateHeaderInformation] = useAtom(UpdateSelectedPatientNameAtom);

  const [__, updateReports] = useAtom(UpdateReportsAtom);
  const [___, attemptCreateReport] = useAtom(createReportAtom);
  const [_____, setSelectedReport] = useAtom(currentlySelectedPatientAtom);
  const [provider, setProvider] = useAtom(userAtom);

  useEffect(() => {
    setProvider({
      ...provider,
      ID: '67bb97c737c2f469d7ac7225',
    });
  }, []);

  const processStream = useStreamProcessor({
    attemptCreateReport,
    updateReports,
    providerID: provider.ID,
    patientName: 'test name',
    timestamp: new Date().toISOString(),
    duration: 10,
  });

  const createSampleFormData = () => {
    const formData = new FormData();

    const randomData = new Uint8Array(1024);
    window.crypto.getRandomValues(randomData);
    const randomBlob = new Blob([randomData], {
      type: MIME_TYPE,
    });
    const file = new File([randomBlob], 'audio', { type: 'audio/webm' });
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
    },
    onError: error => {
      console.error('Error generating report:', error);
    },
  });

  const handleGenerateReport = () => {
    setSelectedReport(''); // need this because of the stream processor only switches when no patient is currently selected
    const formData = createSampleFormData();
    const metadata: GenerateReportMetadata = {
      patientName: 'test name',
      providerName: provider.name,
      timestamp: new Date().toISOString(),
      duration: 10,
      subjectiveStyle: provider.subjectiveStyle,
      objectiveStyle: provider.objectiveStyle,
      assessmentAndPlanStyle: provider.assessmentAndPlanStyle,
      patientInstructionsStyle: provider.patientInstructionsStyle,
      summaryStyle: provider.summaryStyle,
    };

    generateReportMutation.mutate({ formData, metadata });
  };

  return (
    <div className="flex flex-col h-full max-h-screen overflow-hidden">
      <div className="flex gap-4 items-center">
        <Button
          onClick={handleGenerateReport}
          color="blue"
          style={{ maxWidth: '200px', minHeight: '30px' }}
        >
          Generate Sample Report
        </Button>
        <div
          className="flex-shrink-0 border-b border-gray-300 p-4"
          style={{ maxWidth: '300px' }}
        >
          <label
            htmlFor="patientID"
            className="block text-sm font-medium text-gray-700"
          >
            Patient ID
          </label>
          <Input
            placeholder={provider.ID}
            id="patientID"
            name="patientID"
            className="mt-1 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
            onChange={e => {
              setProvider({
                ...provider,
                ID: e.target.value,
              });
            }}
          />
        </div>
      </div>
      <div className="flex-shrink-0 border-b border-gray-300">
        <ProfileSummaryCard
          name={headerInformation.name}
          description={headerInformation.condensedSummary}
          onChange={updateHeaderInformation}
          handleUpdateName={() => {
            console.log('updating name');
          }}
        />
      </div>

      <div className="flex-1 overflow-y-auto p-4">
        <div className="flex gap-16">
          <div className="flex-1 pb-24">
            <VisitReportLayout />
          </div>

          <div className="flex-2 pr-3 flex-shrink-0 min-w-[300px]">
            <NoteControlsLayout />
          </div>
        </div>
      </div>
    </div>
  );
};

export default GenerateReportTest;
