import type { Meta, StoryObj } from '@storybook/react';
import SearchInput from './SearchInput';
import { useState } from 'react';

const meta: Meta<typeof SearchInput> = {
  title: 'Components/FollowUpSearchModal/SearchInput',
  component: SearchInput,
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj<typeof SearchInput>;

export const Default: Story = {
  args: {
    query: '',
    setQuery: () => {
      console.log('nothing');
    },
  },
  decorators: [
    Story => {
      const [query, setQuery] = useState('');
      return <Story args={{ query, setQuery }} />;
    },
  ],
};
