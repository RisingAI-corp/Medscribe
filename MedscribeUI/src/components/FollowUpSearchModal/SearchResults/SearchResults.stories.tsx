import type { Meta, StoryObj } from '@storybook/react';
import SearchResults from './SearchResults';

const meta: Meta<typeof SearchResults> = {
  title: 'Components/FollowUpSearchModal/SearchResults',
  component: SearchResults,
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj<typeof SearchResults>;

export const Default: Story = {
  args: {},
};

