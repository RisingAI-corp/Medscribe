import { Button, Text } from '@mantine/core';

interface DeleteAllNotesModalContent {
  closeModal: () => void;
  onDeleteSelectedPatients: () => void;
}
function DeleteAllNotesModalContent({
  closeModal,
  onDeleteSelectedPatients,
}: DeleteAllNotesModalContent) {
  return (
    <>
      <Text className="text-red-500 text-lg font-bold text-center">
        Are you sure you want to delete all your notes?
      </Text>
      <div className="flex justify-center mt-5 gap-2">
        <Button
          color="red"
          onClick={() => {
            onDeleteSelectedPatients();
            closeModal();
          }}
        >
          Yes
        </Button>
        <Button
          variant="default"
          onClick={() => {
            closeModal();
          }}
        >
          No
        </Button>
      </div>
    </>
  );
}

export default DeleteAllNotesModalContent;
