import type { Meta, StoryObj } from '@storybook/react';
import NoteControlsLayout from '../components/NoteControls/NotecontrolsLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const Wrapper = () => {
  return (
    <MantineProvider>
      <NoteControlsLayout />
    </MantineProvider>
  );
};
const meta: Meta<typeof Wrapper> = {
  title: 'NotecontrolsLayout',
  component: Wrapper,
};
export default meta;

type Story = StoryObj<typeof Wrapper>;

export const Wrapped: Story = {};
