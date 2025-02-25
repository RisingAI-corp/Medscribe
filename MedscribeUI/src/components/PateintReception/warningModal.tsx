import { Modal, Text, Button } from '@mantine/core';

interface WarningModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const WarningModal = ({ isOpen, onClose }: WarningModalProps) => {
  return (
    <Modal
      opened={isOpen}
      onClose={onClose}
      title="Warning: Recording Too Short"
    >
      <Text c="red" size="md">
        Recording must be at least 30 seconds long to generate meaningful data.
      </Text>
      <Button fullWidth mt="md" onClick={onClose}>
        Okay
      </Button>
    </Modal>
  );
};

export default WarningModal;
