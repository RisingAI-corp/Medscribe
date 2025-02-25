import type { Meta, StoryObj } from '@storybook/react';
import NoteControlsLayout from './noteControlsLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
const NoteControlsWrapper = () => {
  const wrapperStyles = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '90vh',
  };

  const innerStyles = {
    width: '300px',
    height: '500px',
    border: '1px dotted red',
    padding: '10px',
  };

  return (
    <div style={wrapperStyles}>
      <div style={innerStyles}>
        <MantineProvider>
          <NoteControlsLayout />
        </MantineProvider>
      </div>
    </div>
  );
};

const meta: Meta<typeof NoteControlsWrapper> = {
  title: 'NotecontrolsLayout',
  component: NoteControlsWrapper,
};
export default meta;

type Story = StoryObj<typeof NoteControlsWrapper>;

export const Wrapped: Story = {};
