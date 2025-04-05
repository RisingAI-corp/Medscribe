import { useState } from 'react';
import { Modal, Checkbox, ScrollArea } from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { useAtom } from 'jotai';
import SearchBox from './searchBox/searchBox';
import PatientPreviewCard from './PatientPreviewCard/patientPreviewCard';
import DeleteAllNotesModalContent from './deleteAllNotes';
import useSearch from '../../hooks/useSearch';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import {
  removePatientsByIdsAtom,
  PatientRegistryAtom,
  setReadStatusAtom,
} from './derivedAtoms';
import { useMutation } from '@tanstack/react-query';
import { deleteReport } from '../../api/deleteReport';
import { markRead, markUnRead } from '../../api/toggleReportReadStatus';

export interface PatientPreviewRecord {
  id: string;
  patientName: string;
  dateOfRecording: string;
  timeOfRecording: string;
  durationOfRecording: string;
  sessionSummary: string;
  finishedGenerating: boolean;
  loading: boolean;
  readStatus: boolean;
}

const PatientRegistryLayout = () => {
  const [checked, setChecked] = useState(false);
  const [selectedPatientsToRemove, setSelectedPatientsToRemove] = useState<
    string[]
  >([]);
  const [opened, { close: closeModal, open: openModal }] = useDisclosure(false);
  const [currentlySelectedPatient, setCurrentlySelectedPatient] = useAtom(
    currentlySelectedPatientAtom,
  );

  const [registryList] = useAtom(PatientRegistryAtom);
  const removePatientById = useAtom(removePatientsByIdsAtom);
  const [, setReadStatus] = useAtom(setReadStatusAtom);

  const [filteredResults, query, setQuery] = useSearch(
    registryList,
    (patient: PatientPreviewRecord) => patient.patientName,
  );

  const deleteReportsMutation = useMutation({
    mutationFn: deleteReport,
    onError: error => {
      console.error('Error deleting reports:', error);
    },
  });

  const markReadMutation = useMutation({
    mutationFn: markRead,
    onError: error => {
      console.error('Error marking read:', error);
    },
  });

  const markUnReadMutation = useMutation({
    mutationFn: markUnRead,
    onError: error => {
      console.error('Error marking unread:', error);
    },
  });

  const handleAddPatientToRemove = (id: string) => {
    setSelectedPatientsToRemove(prev =>
      prev.includes(id) ? prev : [...prev, id],
    );
  };

  const handleRemovePatientToRemove = (id: string) => {
    setSelectedPatientsToRemove(prev =>
      prev.filter(patientId => patientId !== id),
    );
  };

  const handleToggleCheckbox = (id: string, isChecked: boolean) => {
    if (isChecked) {
      handleAddPatientToRemove(id);
    } else {
      handleRemovePatientToRemove(id);
    }
  };

  const handleSelectAll = (status: boolean) => {
    if (status) {
      setSelectedPatientsToRemove(registryList.map(patient => patient.id));
    } else {
      setSelectedPatientsToRemove([]);
    }
  };

  const handleDeleteSelectedPatients = () => {
    deleteReportsMutation.mutate({ ReportIDs: selectedPatientsToRemove });
    removePatientById[1](selectedPatientsToRemove);
    setChecked(false);
    setSelectedPatientsToRemove([]);
  };

  const handlePatientClick = (id: string) => {
    setCurrentlySelectedPatient(id);
  };

  const handleRemovePatient = (id: string) => {
    deleteReportsMutation.mutate({ ReportIDs: [id] });
    removePatientById[1]([id]);
  };

  const handleMarkRead = (id: string) => {
    setReadStatus({ reportId: id, readStatus: true });
    markReadMutation.mutate({ ReportID: id, Opened: true });
  };

  const handleUnMarkRead = (id: string) => {
    setReadStatus({ reportId: id, readStatus: false });
    markUnReadMutation.mutate({ ReportID: id, Opened: false });
  };

  return (
    <div className="flex flex-col gap-4 h-screen relative">
      <SearchBox value={query} onChange={setQuery} />

      <Checkbox
        checked={checked}
        onChange={e => {
          const isChecked = e.currentTarget.checked;
          setChecked(isChecked);
          handleSelectAll(isChecked);
        }}
        label={
          <span
            onClick={e => {
              if (checked) {
                e.preventDefault();
                e.stopPropagation();
                openModal();
              }
            }}
            className={`ml-2 ${checked ? 'text-red-500 cursor-pointer' : 'cursor-default'}`}
          >
            {checked ? 'Delete all notes' : 'Select all notes'}
          </span>
        }
        styles={{ label: { color: 'inherit' } }}
      />

      <ScrollArea className="flex-grow h-full mb-4">
        <div className="pb-40">
          {filteredResults.map(patient => (
            <div className="mt-2" key={patient.id}>
              <PatientPreviewCard
                {...patient}
                isChecked={selectedPatientsToRemove.includes(patient.id)}
                selectAllToggle={checked}
                isSelected={currentlySelectedPatient === patient.id}
                handleToggleCheckbox={handleToggleCheckbox}
                onClick={handlePatientClick}
                handleRemovePatient={handleRemovePatient}
                handleMarkRead={handleMarkRead}
                handleUnMarkRead={handleUnMarkRead}
                readStatus={patient.readStatus}
              />
            </div>
          ))}
        </div>
      </ScrollArea>

      <Modal opened={opened} onClose={closeModal} size="auto">
        <DeleteAllNotesModalContent
          closeModal={closeModal}
          onDeleteSelectedPatients={handleDeleteSelectedPatients}
        />
      </Modal>
    </div>
  );
};

export default PatientRegistryLayout;
