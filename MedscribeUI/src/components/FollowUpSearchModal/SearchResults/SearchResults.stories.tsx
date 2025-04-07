import type { Meta, StoryObj } from '@storybook/react';
import SearchResults, { SearchResultItem } from './SearchResults';

const meta: Meta<typeof SearchResults> = {
  title: 'Components/FollowUpSearchModal/SearchResults',
  component: SearchResults,
  parameters: {
    layout: 'centered',
  },
};

export default meta;
type Story = StoryObj<typeof SearchResults>;

const sampleResults: SearchResultItem[] = [
  {
    id: '1',
    patientName: 'John Doe',
    dateOfRecording: '2023-05-15',
    condensedSummary:
      'Patient presented with symptoms of seasonal allergies and requested refill of antihistamine.',
    summary:
      'Patient presented with sneezing, runny nose, and itchy eyes consistent with seasonal allergies. Requested a refill of loratadine. No other concerns reported.',
    timeOfRecording: '09:00',
    durationOfRecording: 600000,
  },
  {
    id: '2',
    patientName: 'Jane Smith',
    dateOfRecording: '2023-05-14',
    condensedSummary:
      'Follow-up for hypertension management. Blood pressure readings reviewed and medication adjusted.',
    summary:
      'Follow-up for hypertension. Home BP readings showed elevated systolic pressure. Lisinopril increased to 20mg daily. Patient educated on dietary changes.',
    timeOfRecording: '10:30',
    durationOfRecording: 720000,
  },
  {
    id: '3',
    patientName: 'Robert Johnson',
    dateOfRecording: '2023-05-12',
    condensedSummary:
      'Annual wellness visit. Patient reports feeling well. Routine labs ordered.',
    summary:
      'Annual physical. No current complaints. Reviewed health maintenance items. Ordered CBC, CMP, and lipid panel. Follow-up in one year unless symptoms arise.',
    timeOfRecording: '14:15',
    durationOfRecording: 900000,
  },
];

export const Default: Story = {
  args: {
    filteredResults: sampleResults,
    onSelectItem: () => {
      console.log('yo');
    },
  },
  decorators: [
    Story => (
      <div style={{ width: '600px' }}>
        <Story />
      </div>
    ),
  ],
};
