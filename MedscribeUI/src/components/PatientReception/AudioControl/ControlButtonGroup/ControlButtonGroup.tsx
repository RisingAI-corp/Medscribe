import React from 'react';
import { IconPlayerPause, IconPlayerPlay } from '@tabler/icons-react';
import { LiveAudioVisualizer } from 'react-audio-visualize';

interface ControlButtonGroupProps {
  isRecording: boolean;
  onEndVisit: () => void | Promise<void>;
  onPause: () => Promise<void>;
  onResume: () => void;
  onReset: () => void;
  mediaRecorder?: MediaRecorder;
}

const ControlButtonGroup: React.FC<ControlButtonGroupProps> = ({
  isRecording,
  mediaRecorder = null,
  onEndVisit,
  onPause,
  onResume,
  onReset,
}) => {
  return (
    <div className="flex justify-center gap-2">
      <button
        // eslint-disable-next-line @typescript-eslint/no-misused-promises
        onClick={onEndVisit}
        className="bg-red-200 rounded-2xl shadow-md text-black font-bold hover:bg-red-300 px-4 py-2"
      >
        End Visit
      </button>
      {isRecording && (
        <>
          <button
            onClick={() => {
              void onPause();
            }}
            className="bg-blue-200 rounded-full shadow-md p-2 hover:bg-blue-300 text-black"
          >
            <IconPlayerPause size={24} />
          </button>
          {mediaRecorder && (
            <LiveAudioVisualizer
              mediaRecorder={mediaRecorder}
              width={100}
              height={20}
            />
          )}
        </>
      )}

      {!isRecording && (
        <>
          <button
            onClick={onResume}
            className="bg-blue-200 rounded-full shadow-md p-2 hover:bg-blue-300 text-black"
          >
            <IconPlayerPlay size={24} className="text-black" />
          </button>
          <button
            onClick={onReset}
            className="bg-gray-200 rounded-2xl shadow-md p-2 hover:bg-gray-300 text-black"
          >
            Clear
          </button>
        </>
      )}
    </div>
  );
};

export default ControlButtonGroup;
