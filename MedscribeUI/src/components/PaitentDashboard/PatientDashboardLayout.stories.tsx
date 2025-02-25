import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import PatientDashBoardLayout from './patientDashBoardLayout';
import '@mantine/core/styles.css';
const PatientDashBoardLayoutWrapper = () => {
  return (
    <>
      <MantineProvider>
        <PatientDashBoardLayout />
      </MantineProvider>
    </>
  );
};

const meta: Meta<typeof PatientDashBoardLayoutWrapper> = {
  title: 'PatientDashBoardLayout',
  component: PatientDashBoardLayoutWrapper,
};
export default meta;

type Story = StoryObj<typeof PatientDashBoardLayoutWrapper>;

export const Wrapped: Story = {};
