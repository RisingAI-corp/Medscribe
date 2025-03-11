import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import PatientReception from './patientReceptionLayout';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const PatientReceptionWrapper = () => {
  const wrapperStyles = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '90vh',
  };

  const queryClient = new QueryClient();

  return (
    <div style={wrapperStyles}>
      <QueryClientProvider client={queryClient}>
        <MantineProvider>
          <PatientReception />
        </MantineProvider>
      </QueryClientProvider>
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
