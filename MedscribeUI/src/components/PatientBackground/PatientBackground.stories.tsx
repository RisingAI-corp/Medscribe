import type { Meta, StoryObj } from '@storybook/react';
import PatientBackgroundDetails from './PatientBackground';

const meta: Meta<typeof PatientBackgroundDetails> = {
  title: 'Components/PatientBackground/PatientBackgroundDetails',
  component: PatientBackgroundDetails,
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj<typeof PatientBackgroundDetails>;

export const Default: Story = {
  args: {
    condensedSummary:
      'Patient has a history of hypertension and type 2 diabetes. Currently on medication for both conditions.',
    dateOfRecording: '2024-03-15T10:30:00Z',
    durationOfRecording: 1800000, // 30 minutes in milliseconds
    summary:
      'Patient reported improved blood pressure readings. Medication dosage adjusted. Recommended continued diet and exercise regimen.',
  },
};

export const MinimalData: Story = {
  args: {
    condensedSummary: 'New patient with no prior medical history.',
    dateOfRecording: '2024-04-01T09:00:00Z',
    durationOfRecording: 900000, // 15 minutes in milliseconds
    summary: 'Initial consultation. Baseline measurements taken.',
  },
};
