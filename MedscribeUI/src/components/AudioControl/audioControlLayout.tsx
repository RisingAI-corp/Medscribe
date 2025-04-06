import CaptureButton from './CaptureButton/CaptureButton';
import ControlButtonGroup from './ControlButtonGroup/ControlButtonGroup';
import useAudioRecorder from '../../hooks/useAudioRecorder';

interface AudioControlLayoutProps {
  onAudioCaptured?: (blob: Blob, duration: number, timestamp: number) => void | Promise<void>;
}

const AudioControlLayout = ({ onAudioCaptured }: AudioControlLayoutProps) => {
  const {
    isRecording,
    mediaRecorder,
    handleStartRecording,
    handlePauseRecording,
    handleResumeRecording,
    handleStopRecording,
    recordingStartTime,
    audioBlobRef,
    handleResetMediaRecorder,
  } = useAudioRecorder();

  const handleEndVisit = async () => {
    if (!recordingStartTime.current) return;

    const elapsedSeconds = (Date.now() - recordingStartTime.current) / 1000;
    if (elapsedSeconds <= 0) {
      handlePauseRecording();
      return;
    }
    
    await handleStopRecording();
    
    if (onAudioCaptured && audioBlobRef.current) {
      await onAudioCaptured(
        audioBlobRef.current, 
        elapsedSeconds, 
        recordingStartTime.current
      );
    }
    
    handleResetMediaRecorder();
  };

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
