import React from 'react';
import { Button, Group } from '@mantine/core';
import { IconPlayerPause, IconPlayerPlay } from '@tabler/icons-react';

interface ControlButtonsProps {
  isRecording: boolean;
  mediaRecorder: MediaRecorder | null;
  onEndVisit: () => void | Promise<void>;
  onPause: () => void;
  onResume: () => void;
  onReset: () => void;
}

const ControlButtons: React.FC<ControlButtonsProps> = ({
  isRecording,
  mediaRecorder,
  onEndVisit,
  onPause,
  onResume,
  onReset,
}) => {
  return (
    <Group className="justify-center mb-4">
      <Button
        // eslint-disable-next-line @typescript-eslint/no-misused-promises
        onClick={onEndVisit}
        className="bg-gray-300 rounded-2xl shadow-md text-black font-bold hover:bg-gray-400"
      >
        End Visit
      </Button>
      {isRecording && (
        <Button
          onClick={onPause}
          className="bg-gray-300 rounded-full shadow-md p-2 hover:bg-gray-400"
        >
          <IconPlayerPause size={24} className="text-black" />
        </Button>
      )}

      {!isRecording && (
        <>
          <Button
            onClick={onResume}
            className="bg-gray-300 rounded-full shadow-md p-2 hover:bg-gray-400"
          >
            <IconPlayerPlay size={24} className="text-black" />
          </Button>
          <Button
            onClick={onReset}
            className="bg-gray-500 rounded-full shadow-md p-2 ml-2 hover:bg-red-600"
          >
            Clear
          </Button>
        </>
      )}
    </Group>
  );
};

export default ControlButtons;
