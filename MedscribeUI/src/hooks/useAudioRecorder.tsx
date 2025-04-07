import { useState, useRef } from 'react';
import { MIME_TYPE } from '../constants';

const useAudioRecorder = () => {
  const [isRecording, setIsRecording] = useState(false);
  const [mediaRecorder, setMediaRecorder] = useState<MediaRecorder | null>(
    null,
  );
  const recordingStartTime = useRef<number | null>(null);

  const audioChunksRef = useRef<BlobPart[]>([]);
  const audioBlobRef = useRef<Blob | null>(null);

  const handleStartRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const recorder = new MediaRecorder(stream, { mimeType: MIME_TYPE });

      recorder.ondataavailable = event => {
        audioChunksRef.current.push(event.data);
      };

      recorder.onstop = () => {
        audioBlobRef.current = new Blob(audioChunksRef.current, {
          type: MIME_TYPE,
        });
      };

      setMediaRecorder(recorder);
      recorder.start();
      recordingStartTime.current = Date.now();
      setIsRecording(true);
    } catch (error) {
      console.error('Failed to access microphone:', error);
    }
  };

  const handlePauseRecording = (): Promise<void> => {
    return new Promise<void>(resolve => {
      if (!mediaRecorder || mediaRecorder.state !== 'recording') {
        resolve();
        return;
      }

      mediaRecorder.onpause = () => {
        audioBlobRef.current = new Blob(audioChunksRef.current, {
          type: MIME_TYPE,
        });
        setIsRecording(false);
        resolve();
      };

      mediaRecorder.pause();
    });
  };

  const handleResumeRecording = () => {
    mediaRecorder?.resume();
    setIsRecording(true);
  };

  const handleStopRecording = () => {
    return new Promise<void>(resolve => {
      if (!mediaRecorder) {
        resolve();
        return;
      }
      mediaRecorder.onstop = () => {
        audioBlobRef.current = new Blob(audioChunksRef.current, {
          type: MIME_TYPE,
        });
        setIsRecording(false);
        mediaRecorder.stream.getTracks().forEach(track => {
          track.stop();
        });
        resolve();
      };
      mediaRecorder.stop();
    });
  };

  const handleResetMediaRecorder = () => {
    audioChunksRef.current = [];
    audioBlobRef.current = null;
    setMediaRecorder(null);
  };

  return {
    isRecording,
    mediaRecorder,
    audioBlobRef,
    handleStartRecording,
    handlePauseRecording,
    handleResumeRecording,
    handleStopRecording,
    recordingStartTime,
    handleResetMediaRecorder,
  };
};

export default useAudioRecorder;
