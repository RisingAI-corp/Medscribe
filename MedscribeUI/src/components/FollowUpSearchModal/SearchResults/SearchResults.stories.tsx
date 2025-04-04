import type { Meta, StoryObj } from '@storybook/react';
import SearchResults, { SearchResultItem } from './SearchResults';
import { useState } from 'react';

const meta: Meta<typeof SearchResults> = {
  title: 'Components/FollowUpSearchModal/SearchResults',
  component: SearchResults,
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj<typeof SearchResults>;

// Sample data
const sampleResults: SearchResultItem[] = [
  {
    id: '1',
    patientName: 'John Doe',
    dateOfRecording: '2023-05-15',
    shortenedSummary: 'Patient presented with symptoms of seasonal allergies and requested refill of antihistamine.'
  },
  {
    id: '2',
    patientName: 'Jane Smith',
    dateOfRecording: '2023-05-14',
    shortenedSummary: 'Follow-up for hypertension management. Blood pressure readings reviewed and medication adjusted.'
  },
  {
    id: '3',
    patientName: 'Robert Johnson',
    dateOfRecording: '2023-05-12',
    shortenedSummary: 'Annual wellness visit. Patient reports feeling well. Routine labs ordered.'
  },
  {
    id: '4',
    patientName: 'Emily Williams',
    dateOfRecording: '2023-05-10',
    shortenedSummary: 'Patient complains of persistent cough for 2 weeks. Prescribed antibiotics and cough suppressant.'
  },
  {
    id: '5',
    patientName: 'Michael Brown',
    dateOfRecording: '2023-05-08',
    shortenedSummary: 'Post-surgical follow-up. Incision healing well. Continue physical therapy.'
  },
];

// Interactive example with state
const InteractiveSearchResults = () => {
  const [selectedItems, setSelectedItems] = useState<SearchResultItem[]>([]);
  
  const handleSelectItem = (item: SearchResultItem) => {
    setSelectedItems(prev => {
      const isAlreadySelected = prev.some(i => i.id === item.id);
      if (isAlreadySelected) {
        return prev.filter(i => i.id !== item.id);
      } else {
        return [...prev, item];
      }
    });
  };
  
  return (
    <div style={{ width: '600px' }}>
      <SearchResults 
        filteredResults={sampleResults} 
        selectedItems={selectedItems}
        onSelectItem={handleSelectItem}
      />
    </div>
  );
};

export const WithResults: Story = {
  render: () => <InteractiveSearchResults />
};

export const WithPreselectedItems: Story = {
  args: {
    filteredResults: sampleResults,
    selectedItems: [sampleResults[0], sampleResults[2]],
    onSelectItem: () => {},
  },
  decorators: [
    (Story) => (
      <div style={{ width: '600px' }}>
        <Story />
      </div>
    ),
  ],
};

export const Empty: Story = {
  args: {
    filteredResults: [],
    selectedItems: [],
    onSelectItem: () => {},
  },
  decorators: [
    (Story) => (
      <div style={{ width: '600px' }}>
        <Story />
      </div>
    ),
  ],
};

export const ManyResults: Story = {
  args: {
    filteredResults: Array(20).fill(null).map((_, index) => ({
      id: `${index + 6}`,
      patientName: `Patient ${index + 6}`,
      dateOfRecording: `2023-04-${30 - index}`,
      shortenedSummary: `This is a sample summary for patient ${index + 6}. Contains medical information and follow-up details.`
    })) as SearchResultItem[],
    selectedItems: [],
    onSelectItem: () => {},
  },
  decorators: [
    (Story) => (
      <div style={{ width: '600px' }}>
        <Story />
      </div>
    ),
  ],
};

