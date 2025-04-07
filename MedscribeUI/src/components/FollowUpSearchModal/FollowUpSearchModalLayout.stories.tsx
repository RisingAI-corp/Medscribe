import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import FollowUpSearchModalLayout from './FollowUpSearchModalLayout';
import { SearchResultItem } from './SearchResults/SearchResults';

// Wrapper component for state management in the story
const FollowUpSearchModalWithState = () => {
  const [selectedVisit, setSelectedVisit] = useState<SearchResultItem | null>(null);
  
  const handleSelectedVisit = (visit: SearchResultItem) => {
    setSelectedVisit(visit);
  };
  
  return (
    <div className="p-8">
      <FollowUpSearchModalLayout 
        handleSelectedVisit={handleSelectedVisit}
      >
        <button className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors">
          Search Follow-ups
        </button>
      </FollowUpSearchModalLayout>
      
      {/* Additional info for the story */}
      <div className="mt-8 p-4 border border-gray-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2 text-gray-800">Selected Visit:</h3>
        {selectedVisit ? (
          <div className="text-gray-700">
            <p><strong>Patient:</strong> {selectedVisit.patientName}</p>
            <p><strong>Date:</strong> {selectedVisit.dateOfRecording}</p>
            <p><strong>Summary:</strong> {selectedVisit.summary}</p>
          </div>
        ) : (
          <p className="text-gray-500">No visit selected</p>
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

