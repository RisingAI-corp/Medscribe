import type { Meta, StoryObj } from '@storybook/react';
import NoteControlsLayout from '../components/NoteControls/NotecontrolsLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const NoteControlsWrapper = ({ isStatus, defaultHover }) => {
  return (
    <MantineProvider>
      <NoteControlsLayout isStatus={isStatus} defaultHover={defaultHover} />
    </MantineProvider>
  );
};

const Wrapper = () => {
  return (
    <NoteControlsWrapper
      isStatus={false}
      defaultHover={{
        visitType: 'New Patient',
        pronoun: 'HE',
        patientClient: 'Patient',
      }}
    />
  );
};

const meta: Meta<typeof Wrapper> = {
  title: 'NotecontrolsLayout',
  component: Wrapper,
};
export default meta;

type Story = StoryObj<typeof Wrapper>;

export const Wrapped: Story = {};
