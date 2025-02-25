import React from 'react';
import { Button, Group } from '@mantine/core';
import { IconPlayerPause, IconPlayerPlay } from '@tabler/icons-react';

interface ControlButtonsProps {
  isRecording: boolean;
  mediaRecorder: MediaRecorder | null;
  onEndVisit: () => void | Promise<void>;
  onPause: () => void;
  onResume: () => void;
}

const ControlButtons: React.FC<ControlButtonsProps> = ({
  isRecording,
  mediaRecorder,
  onEndVisit,
  onPause,
  onResume,
}) => {
  return (
    <Group className="justify-center mb-4">
      {mediaRecorder && (
        <Button
          // eslint-disable-next-line @typescript-eslint/no-misused-promises
          onClick={onEndVisit}
          className="bg-gray-300 rounded-2xl shadow-md text-black font-bold hover:bg-gray-400"
        >
          End Visit
        </Button>
      )}

      {isRecording && mediaRecorder && (
        <Button
          onClick={onPause}
          className="bg-gray-300 rounded-full shadow-md p-2 hover:bg-gray-400"
        >
          <IconPlayerPause size={24} className="text-black" />
        </Button>
      )}

      {!isRecording && mediaRecorder && (
        <Button
          onClick={onResume}
          className="bg-gray-300 rounded-full shadow-md p-2 hover:bg-gray-400"
        >
          <IconPlayerPlay size={24} className="text-black" />
        </Button>
      )}
    </Group>
  );
};

export default ControlButtons;
