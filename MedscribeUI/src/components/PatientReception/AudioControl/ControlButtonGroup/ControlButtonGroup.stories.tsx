import type { Meta, StoryObj } from '@storybook/react';
import ControlButtonGroup from './ControlButtonGroup';

const meta: Meta<typeof ControlButtonGroup> = {
  title: 'Components/AudioControl/ControlButtonGroup',
  component: ControlButtonGroup,
  parameters: {
    layout: 'centered',
  },
  argTypes: {
    onEndVisit: { action: 'End Visit clicked' },
    onPause: { action: 'Pause clicked' },
    onResume: { action: 'Resume clicked' },
    onReset: { action: 'Reset clicked' },
  },
};

export default meta;
type Story = StoryObj<typeof ControlButtonGroup>;

export const Paused: Story = {
  args: {
    isRecording: false,
  },
  name: 'Paused State',
};

export const Recording: Story = {
  args: {
    isRecording: true,
  },
  name: 'Recording State',
};
