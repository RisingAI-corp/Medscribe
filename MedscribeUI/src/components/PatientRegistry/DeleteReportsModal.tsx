import { Modal, Button, Text } from '@mantine/core';

interface DeleteReportsModalProps {
  opened: boolean;
  handleClose: () => void;
  handleDeleteReports: () => void;
}
function DeleteReportsModal({
  opened,
  handleClose,
  handleDeleteReports,
}: DeleteReportsModalProps) {
  return (
    <>
      <Modal
        opened={opened}
        onClose={() => {
          handleClose();
        }}
        title="Confirm Deletion"
      >
        <Text>Are you sure you want to delete reports?</Text>
        <div
          style={{ marginTop: 20, display: 'flex', justifyContent: 'flex-end' }}
        >
          <Button
            variant="default"
            onClick={() => {
              handleClose();
            }}
            style={{ marginRight: 10 }}
          >
            Cancel
          </Button>
          <Button
            color="red"
            onClick={() => {
              handleDeleteReports();
              handleClose();
            }}
          >
            Delete
          </Button>
        </div>
      </Modal>
    </>
  );
}

export default DeleteReportsModal;
