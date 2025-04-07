import type { Meta, StoryObj } from '@storybook/react';
import AudioControlLayout from './audioControlLayout';

const meta: Meta<typeof AudioControlLayout> = {
  title: 'Components/PatientReception/AudioControl/AudioControlLayout',
  component: AudioControlLayout,
  parameters: {
    layout: 'centered',
  },
  argTypes: {
    onAudioCaptured: { action: 'audioCaptured' },
  },
};

export default meta;
type Story = StoryObj<typeof AudioControlLayout>;

export const Default: Story = {
  args: {},
};

export const WithAudioCaptureHandler: Story = {
  args: {
    onAudioCaptured: (blob, duration, timestamp) => {
      console.log('Audio captured:', {
        blobSize: blob.size,
        duration,
        timestamp: new Date(timestamp).toISOString(),
      });
    },
  },
};

export const SimpleRecordingDemo: Story = {
  args: {},
  parameters: {
    docs: {
      description: 'A simple recording control for capturing patient audio',
    },
  },
};
