import CaptureButton from './CaptureButton/CaptureButton';
import ControlButtonGroup from './ControlButtonGroup/ControlButtonGroup';
import { MutableRefObject } from 'react';

export interface AudioControlLayoutProps {
  isRecording: boolean;
  mediaRecorder: MediaRecorder | null;
  audioBlobRef: MutableRefObject<Blob | null>;
  recordingStartTime: MutableRefObject<number | null>;
  handleStartRecording: () => Promise<void>;
  handlePauseRecording: () => Promise<void>;
  handleResumeRecording: () => void;
  handleStopRecording: () => Promise<void>;
  handleResetMediaRecorder: () => void;
  onAudioCaptured?: (
    blob: Blob,
    duration: number,
    timestamp: number,
  ) => void | Promise<void>;
  handleEndVisit: () => void | Promise<void>;
}

const AudioControlLayout = ({
  isRecording,
  mediaRecorder,
  handleStartRecording,
  handlePauseRecording,
  handleResumeRecording,
  handleResetMediaRecorder,
  handleEndVisit,
}: AudioControlLayoutProps) => {
  return (
    <div className="flex flex-col items-center">
      {!isRecording && !mediaRecorder && (
        <CaptureButton
          onClick={() => {
            void handleStartRecording();
          }}
        />
      )}

      {mediaRecorder && (
        <div className="w-full flex items-center">
          <ControlButtonGroup
            isRecording={isRecording}
            onEndVisit={handleEndVisit}
            onPause={handlePauseRecording}
            onResume={handleResumeRecording}
            onReset={handleResetMediaRecorder}
            mediaRecorder={mediaRecorder}
          />
        </div>
      )}
    </div>
  );
};

export default AudioControlLayout;
