import { Select, LoadingOverlay, Button } from '@mantine/core';
import { useState } from 'react';
import BtnGroupSelector from '../Utilities/btnGroupSelector';
import { useDisclosure } from '@mantine/hooks';
import {
  fetchContentDataAtom,
  NoteControlsAtom,
  toggleLoadingForReportSectionsAtom,
} from './derivedAtoms';
import { useAtom } from 'jotai';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { editPatientVisitAtom, noteControlEditsProps } from './derivedAtoms';
import {
  Client,
  FollowUp,
  He,
  NewVisit,
  Patient,
  She,
  They,
} from '../../constants';
import { useMutation } from '@tanstack/react-query';
import { regenerateReport, Updates } from '../../api/regenerateReport';
import { userAtom } from '../../states/userAtom';
import { useStreamProcessor } from '../../hooks/useStreamProcessor';
import { UpdateReportsAtom } from '../PateintReception/derivedAtoms';

function NoteControlsLayout() {
  const [noteControls] = useAtom(NoteControlsAtom);
  const editPatientVisit = useAtom(editPatientVisitAtom);
  const [______, fetchContentData] = useAtom(fetchContentDataAtom);
  const [_____, updateReports] = useAtom(UpdateReportsAtom);

  const [currentlySelectedPatient, _] = useAtom(currentlySelectedPatientAtom);
  const [provider, __] = useAtom(userAtom);
  const [___, toggleContentLoading] = useAtom(
    toggleLoadingForReportSectionsAtom,
  );

  const [selectedPronoun, setSelectedPronoun] = useState(noteControls.pronouns);
  const [selectedVisitType, setSelectedVisitType] = useState(
    noteControls.visitType,
  );
  const [changes, setChanges] = useState<noteControlEditsProps>({
    id: currentlySelectedPatient,
    changes: {},
  });

  const [selectedPatientClient, setSelectedPatientClient] = useState(
    noteControls.patientOrClient,
  );
  const [visible, { toggle }] = useDisclosure(false);

  const regenerateReportMutation = useMutation({
    mutationFn: regenerateReport,
    onSuccess: async reader => {
      await processStream(reader);
      toggle();
    },
    onError: error => {
      if (error.message.includes('500')) {
        toggle();
        toggleContentLoading();
      }
    },
  });

  const processStream = useStreamProcessor({
    updateReports,
    providerID: provider.ID,
    reportID: currentlySelectedPatient,
  });

  const isDirty =
    selectedVisitType !== noteControls.visitType ||
    selectedPronoun !== noteControls.pronouns ||
    selectedPatientClient !== noteControls.patientOrClient;

  const handleVisitTypeSelect = (value: boolean) => {
    if (value !== noteControls.visitType) {
      setSelectedVisitType(value);
      setChanges(prevChanges => ({
        ...prevChanges,
        changes: {
          ...prevChanges.changes,
          isFollowUp: value,
        },
      }));
      return;
    }
  };

  const handlePronounSelect = (value: string) => {
    if (value === She || value === He || value === They) {
      setSelectedPronoun(value);
      setChanges(prevChanges => ({
        ...prevChanges,
        changes: {
          ...prevChanges.changes,
          pronouns: value,
        },
      }));

      return;
    }
  };

  const handlePatientClientSelect = (value: string) => {
    if (value === Patient || value === Client) {
      setSelectedPatientClient(value);
      setChanges(prevChanges => ({
        ...prevChanges,
        changes: {
          ...prevChanges.changes,
          patientOrClient: value,
        },
      }));
      return;
    }
  };

  const handleRegenerate = () => {
    const updates: Updates[] = [];
    //Todo: send changes upstream to endPoint
    if (changes.changes.pronouns !== undefined) {
      updates.push({ Key: 'pronouns', Value: changes.changes.pronouns });
      setSelectedPronoun(changes.changes.pronouns);
    }
    if (changes.changes.isFollowUp !== undefined) {
      updates.push({ Key: 'isFollowUp', Value: changes.changes.isFollowUp });
      setSelectedVisitType(changes.changes.isFollowUp);
    }
    if (changes.changes.patientOrClient !== undefined) {
      updates.push({
        Key: 'patientOrClient',
        Value: changes.changes.patientOrClient,
      });
      setSelectedPatientClient(changes.changes.patientOrClient);
    }
    const contentData = fetchContentData();
    if (!contentData) {
      return;
    }
    toggle();
    toggleContentLoading();
    regenerateReportMutation.mutate({
      providerID: provider.ID,
      ID: currentlySelectedPatient,
      subjectiveStyle: provider.subjectiveStyle,
      objectiveStyle: provider.objectiveStyle,
      assessmentStyle: provider.assessmentStyle,
      planningStyle: provider.planningStyle,
      summaryStyle: provider.summaryStyle,
      updates: updates,
      SubjectiveContent: contentData.subjective,
      ObjectiveContent: contentData.objective,
      AssessmentContent: contentData.assessment,
      PlanningContent: contentData.planning,
      SummaryContent: contentData.summary,
    });

    editPatientVisit[1](changes);
  };

  return (
    <div className="relative">
      <LoadingOverlay
        visible={visible}
        zIndex={1000}
        overlayProps={{ radius: 'sm', blur: 2 }}
        loaderProps={{ color: 'blue', type: 'bars' }}
      />
      <span className="block mb-2">Visit Type</span>
      <Select
        defaultValue={''}
        data={[NewVisit, FollowUp]}
        value={selectedVisitType ? FollowUp : NewVisit}
        onChange={e => {
          if (e === FollowUp) {
            handleVisitTypeSelect(true);
          } else {
            handleVisitTypeSelect(false);
          }
        }}
      />

      <hr className="my-4" />

      <span className="block mb-2">Pronoun Selector</span>
      <BtnGroupSelector
        buttonLabelOptions={[He, She, They]}
        selectedBtn={selectedPronoun}
        onSelect={handlePronounSelect}
      />

      <span className="block my-4 mb-2">Patient/Client</span>
      <BtnGroupSelector
        buttonLabelOptions={[Patient, Client]}
        selectedBtn={selectedPatientClient}
        onSelect={handlePatientClientSelect}
      />

      <hr className="my-4" />
      <Button onClick={handleRegenerate} className="w-full" disabled={!isDirty}>
        Regenerate Report
      </Button>
    </div>
  );
}

export default NoteControlsLayout;
