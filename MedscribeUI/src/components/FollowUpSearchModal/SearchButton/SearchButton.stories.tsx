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

export const WithOnePatient: Story = {
  args: {
    selectedItems: ['John Smith'],
  },
};

export const WithMultiplePatients: Story = {
  args: {
    selectedItems: ['John Smith', 'Maria Garcia', 'David Chen'],
  },
};

export const WithLongPatientNames: Story = {
  args: {
    selectedItems: ['Elizabeth Williamson-Thompson', 'Christopher Rodriguez-Martinez'],
  },
};

