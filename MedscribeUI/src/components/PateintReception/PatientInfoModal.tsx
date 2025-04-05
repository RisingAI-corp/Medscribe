import { Modal, TextInput, Button, Group } from '@mantine/core';
import { useRef } from 'react';

interface PatientInfoModalProps {
  isOpen: boolean;
  patientName: string;
  onClose: () => void;
  onSubmit: (name: string) => void;
}

const PatientInfoModal = ({
  isOpen,
  onClose,
  onSubmit,
}: PatientInfoModalProps) => {
  const patientNameRef = useRef<HTMLInputElement>(null);
  return (
    <Modal
      centered
      opened={isOpen}
      onClose={onClose}
      title="Patient Information"
    >
      <TextInput
        label="Patient Name"
        ref={patientNameRef}
        placeholder="Enter patient's name"
        required
      />
      <Group className="flex justify-end mt-4">
        <Button
          color="blue"
          onClick={() => {
            if (patientNameRef.current != null) {
              onSubmit(patientNameRef.current.value);
            }
          }}
        >
          Start Recording
        </Button>
      </Group>
    </Modal>
  );
};

export default PatientInfoModal;
