import type { Meta, StoryObj } from '@storybook/react';
import NoteControlsLayout from '../components/NoteControls/NotecontrolsLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const WIDTH = '300px';
const HEIGHT = '500px';

const NoteControlsWrapper = ({ isStatus, defaultHover }) => {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '90vh' }}>
      <div style={{ width: WIDTH, height: HEIGHT, border: '1px dotted red' }}>
        <MantineProvider>
          <NoteControlsLayout isStatus={isStatus} defaultHover={defaultHover} />
        </MantineProvider>
      </div>
    </div>
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
