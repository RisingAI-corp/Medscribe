import type { Meta, StoryObj } from '@storybook/react';
import SearchButton from './SearchButton';

const meta: Meta<typeof SearchButton> = {
  title: 'Components/FollowUpSearchModal/SearchButton',
  component: SearchButton,
  parameters: {
    layout: 'centered',
  },    
};

export default meta;
type Story = StoryObj<typeof SearchButton>;

export const Default: Story = {
  args: {},
};

