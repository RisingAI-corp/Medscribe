import type { Meta, StoryObj } from '@storybook/react';
import CaptureButton from './CaptureButton';

const meta: Meta<typeof CaptureButton> = {
  title: 'Components/PatientReception/AudioControl/CaptureButton',
  component: CaptureButton,
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj<typeof CaptureButton>;

export const Default: Story = {
  args: {},
};
