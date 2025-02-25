import { Modal, TextInput, Button, Group } from '@mantine/core';

interface PatientInfoModalProps {
  isOpen: boolean;
  patientName: string;
  onClose: () => void;
  onChange: (value: string) => void;
  onSubmit: () => void;
}

const PatientInfoModal = ({
  isOpen,
  patientName,
  onClose,
  onChange,
  onSubmit,
}: PatientInfoModalProps) => {
  return (
    <Modal
      centered
      opened={isOpen}
      onClose={onClose}
      title="Patient Information"
    >
      <TextInput
        label="Patient Name"
        placeholder="Enter patient's name"
        value={patientName}
        onChange={event => {
          onChange(event.currentTarget.value);
        }}
        required
      />
      <Group className="flex justify-end mt-4">
        <Button color="blue" onClick={onSubmit}>
          Start Recording
        </Button>
      </Group>
    </Modal>
  );
};

export default PatientInfoModal;
