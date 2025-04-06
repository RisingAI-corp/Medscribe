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
    name: 'John Doe',
    condensedSummary: 'Patient has a history of hypertension and type 2 diabetes. Currently on medication for both conditions.',
    lastVisitDate: '2024-03-15T10:30:00Z',
    duration: 1800000, // 30 minutes in milliseconds
    lastVisitSummary: 'Patient reported improved blood pressure readings. Medication dosage adjusted. Recommended continued diet and exercise regimen.',
  },
};

export const MinimalData: Story = {
  args: {
    name: 'Jane Smith',
    condensedSummary: 'New patient with no prior medical history.',
    lastVisitDate: '2024-04-01T09:00:00Z',
    duration: 900000, // 15 minutes in milliseconds
    lastVisitSummary: 'Initial consultation. Baseline measurements taken.',
  },
}; 