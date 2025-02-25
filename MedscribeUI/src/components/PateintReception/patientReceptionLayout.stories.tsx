import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import PatientReception from './patientReceptionLayout';

const PatientReceptionWrapper = () => {
  const wrapperStyles = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '90vh',
  };

  return (
    <div style={wrapperStyles}>
      <MantineProvider>
        <PatientReception />
      </MantineProvider>
    </div>
  );
};

const meta: Meta<typeof PatientReceptionWrapper> = {
  title: 'PatientReception',
  component: PatientReceptionWrapper,
};
export default meta;

type Story = StoryObj<typeof PatientReceptionWrapper>;

export const Wrapped: Story = {};
