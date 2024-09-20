import type { Meta, StoryObj } from '@storybook/react';
import NoteControlsLayout from '../components/NoteControls/NotecontrolsLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const WIDTH = '300px';
const HEIGHT = '500px';
const PADDING = '10px';

const NoteControlsWrapper = ({
  isStatus,
  defaultVisitType,
  defaultPronoun,
  defaultPatientClient,
}) => {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '90vh',
      }}
    >
      <div
        style={{
          width: WIDTH,
          height: HEIGHT,
          border: '1px dotted red',
          padding: PADDING,
        }}
      >
        <MantineProvider>
          <NoteControlsLayout
            isStatus={isStatus}
            defaultVisitType={defaultVisitType}
            defaultPronoun={defaultPronoun}
            defaultPatientClient={defaultPatientClient}
          />
        </MantineProvider>
      </div>
    </div>
  );
};

const Wrapper = () => {
  return (
    <NoteControlsWrapper
      isStatus={false}
      defaultVisitType="New Patient"
      defaultPronoun="HE"
      defaultPatientClient="Patient"
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
