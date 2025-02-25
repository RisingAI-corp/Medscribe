import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import ProfileSummaryCard, {
  ProfileSummaryCardProps,
} from './profileSummaryCard';

const ProfileSummaryCardWrapper = (props: ProfileSummaryCardProps) => {
  const wrapperStyles = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '90vh',
  };

  const innerStyles = {
    width: '600px',
    height: '500px',
    border: '1px dotted red',
    padding: '10px',
  };

  return (
    <div style={wrapperStyles}>
      <div style={innerStyles}>
        <MantineProvider>
          <ProfileSummaryCard {...props} />
        </MantineProvider>
      </div>
    </div>
  );
};

const meta: Meta<typeof ProfileSummaryCardWrapper> = {
  title: 'ProfileSummaryCard',
  component: ProfileSummaryCardWrapper,
};
export default meta;

type Story = StoryObj<typeof ProfileSummaryCardWrapper>;

export const Wrapped: Story = {};
