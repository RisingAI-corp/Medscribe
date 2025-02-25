import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import PatientRegistryLayout from './patientRegistryLayout';
const PatientRegistryLayoutWrapper = () => {
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
    overflow: 'hidden',
  };

  return (
    <div style={wrapperStyles}>
      <div style={innerStyles}>
        <MantineProvider>
          <PatientRegistryLayout />
        </MantineProvider>
      </div>
    </div>
  );
};

const meta: Meta<typeof PatientRegistryLayoutWrapper> = {
  title: 'PatientRegistryLayout',
  component: PatientRegistryLayoutWrapper,
};
export default meta;

type Story = StoryObj<typeof PatientRegistryLayoutWrapper>;

export const Wrapped: Story = {};
