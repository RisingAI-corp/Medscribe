import type { Meta, StoryObj } from '@storybook/react';
import FollowUpSearchModalLayout from './FollowUpSearchModalLayout';

const meta: Meta<typeof FollowUpSearchModalLayout> = {
  title: 'Components/FollowUpSearchModal/FollowUpSearchModalLayout',    
  component: FollowUpSearchModalLayout,
  parameters: {
    layout: 'centered',
  },    
};

export default meta;
type Story = StoryObj<typeof FollowUpSearchModalLayout>;

export const Default: Story = {
  args: {},
};

