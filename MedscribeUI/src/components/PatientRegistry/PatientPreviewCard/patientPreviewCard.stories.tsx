import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import PatientPreviewCard, {
  PatientPreviewCardProps,
} from './patientPreviewCard';

const samplePatientData: PatientPreviewCardProps = {
  id: '1',
  patientName: 'John Doe',
  dateOfRecording: '2025-01-18',
  timeOfRecording: '10:30 AM',
  durationOfRecording: '15 min',
  sessionSummary: 'Follow-up consultation for back pain.',
  loading: false,
  isChecked: false,
  selectAllToggle: false,
  readStatus: false,
  handleMarkRead: (id: string) => {
    console.log(`Mark patient with ID: ${id} as read`);
  },
  handleUnMarkRead: (id: string) => {
    console.log(`Mark patient with ID: ${id} as unread`);
  },
  handleRemovePatient: (id: string) => {
    console.log(`this is the patient you want to remove: ${id}`);
  },
  handleToggleCheckbox: (id: string, checked: boolean) => {
    console.log(
      `Remove patient with ID: ${id} and is checked: ${String(checked)}`,
    );
  },
  onClick: (id: string) => {
    console.log(`Clicked on patient with ID: ${id}`);
  },
  isSelected: false,
};

const PatientPreviewCardWrapper = (props: PatientPreviewCardProps) => {
  const wrapperStyles = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '90vh',
  };

  const innerStyles = {
    width: '600px',
    height: '500px',
    border: '1px dotted red',
    padding: '10px',
  };

  return (
    <div style={wrapperStyles}>
      <div style={innerStyles}>
        <MantineProvider>
          <PatientPreviewCard {...props} />
        </MantineProvider>
      </div>
    </div>
  );
};

const meta: Meta<typeof PatientPreviewCardWrapper> = {
  title: 'PatientPreviewCard',
  component: PatientPreviewCardWrapper,
};
export default meta;

type Story = StoryObj<typeof PatientPreviewCardWrapper>;

export const Wrapped: Story = {
  args: samplePatientData,
};
