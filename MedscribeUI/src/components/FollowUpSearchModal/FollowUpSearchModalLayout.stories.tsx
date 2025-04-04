import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import FollowUpSearchModalLayout from './FollowUpSearchModalLayout';
import { SearchResultItem } from './SearchResults/SearchResults';

// Mock data for the story
const mockSearchResults: SearchResultItem[] = [
  {
    id: '1',
    patientName: 'John Smith',
    dateOfRecording: '2023-05-10',
    shortenedSummary: 'Patient complains of persistent headaches for the past week.'
  },
  {
    id: '2',
    patientName: 'Maria Garcia',
    dateOfRecording: '2023-05-12',
    shortenedSummary: 'Follow-up for hypertension management. Blood pressure readings stable.'
  },
  {
    id: '3',
    patientName: 'David Chen',
    dateOfRecording: '2023-05-15',
    shortenedSummary: 'Post-surgical follow-up. Incision healing well, no signs of infection.'
  },
  {
    id: '4',
    patientName: 'Sarah Johnson',
    dateOfRecording: '2023-05-18',
    shortenedSummary: 'Annual wellness check. All vitals within normal ranges.'
  },
  {
    id: '5',
    patientName: 'Robert Williams',
    dateOfRecording: '2023-05-20',
    shortenedSummary: 'Diabetic check-up. HbA1c shows improvement from previous visit.'
  }
];

// Wrapper component for state management in the story
const FollowUpSearchModalWithState = () => {
  const [selectedItems, setSelectedItems] = useState<SearchResultItem[]>([]);
  
  return (
    <div className="p-8">
      <FollowUpSearchModalLayout 
        selectedItems={selectedItems} 
        setSelectedItems={setSelectedItems}
        mockData={mockSearchResults}
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

