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
    setQuery: () => {},
  },
  decorators: [
    (Story) => {
      const [query, setQuery] = useState('');
      return (
        <div className="w-[300px]">
          <Story args={{ query, setQuery }} />
        </div>
      );
    },
  ],
};

