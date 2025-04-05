import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import FollowUpSearchModalLayout from './FollowUpSearchModalLayout';
import { SearchResultItem } from './SearchResults/SearchResults';

// Wrapper component for state management in the story
const FollowUpSearchModalWithState = () => {
  const [selectedItems, setSelectedItems] = useState<SearchResultItem[]>([]);
  
  return (
    <div className="p-8">
      <FollowUpSearchModalLayout 
        selectedItems={selectedItems} 
        setSelectedItems={setSelectedItems}
      />
      
      {/* Additional info for the story */}
      <div className="mt-8 p-4 border border-gray-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2 text-gray-800">Selected Items:</h3>
        {selectedItems.length > 0 ? (
          <ul className="list-disc pl-5">
            {selectedItems.map(item => (
              <li key={item.id} className="text-gray-700">{item.patientName} - {item.dateOfRecording}</li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-500">No items selected</p>
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

