import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import FollowUpSearchModalLayout from './FollowUpSearchModalLayout';

// Wrapper component for state management in the story
const FollowUpSearchModalWithState = () => {
  const [selectedItem, setSelectedItem] = useState<string>("");
  
  return (
    <div className="p-8">
      <FollowUpSearchModalLayout 
        selectedItem={selectedItem} 
        setSelectedItem={setSelectedItem}
      />
      
      {/* Additional info for the story */}
      <div className="mt-8 p-4 border border-gray-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2 text-gray-800">Selected Item:</h3>
        {selectedItem ? (
          <p className="text-gray-700">{selectedItem}</p>
        ) : (
          <p className="text-gray-500">No item selected</p>
        )}
      </div>
    </div>
  );
};

const meta: Meta<typeof FollowUpSearchModalLayout> = {
  title: 'Components/FollowUpSearchModal/FollowUpSearchModalLayout',
  component: FollowUpSearchModalLayout,
  parameters: {
    layout: 'centered',
  },
  decorators: [
    (Story) => <FollowUpSearchModalWithState />
  ],
};

export default meta;
type Story = StoryObj<typeof FollowUpSearchModalLayout>;

export const Default: Story = {};

