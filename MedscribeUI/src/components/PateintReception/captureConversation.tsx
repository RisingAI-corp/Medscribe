import { Button, Group } from '@mantine/core';

interface CaptureConversationButtonProps {
  onClick: () => void;
}

const CaptureConversationButton = ({
  onClick,
}: CaptureConversationButtonProps) => {
  return (
    <Group className="flex justify-center items-center mb-4">
      <Button color="blue" onClick={onClick} size="xl">
        Capture Conversation
      </Button>
    </Group>
  );
};

export default CaptureConversationButton;
