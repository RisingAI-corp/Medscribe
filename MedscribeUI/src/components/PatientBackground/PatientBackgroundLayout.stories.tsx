import type { Meta, StoryObj } from '@storybook/react';
import PatientBackgroundLayout from './PatientBackgroundLayout';

const meta: Meta<typeof PatientBackgroundLayout> = {
  title: 'Components/PatientBackground/PatientBackgroundLayout',
  component: PatientBackgroundLayout,
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj<typeof PatientBackgroundLayout>;

export const Default: Story = {
  args: {},
};