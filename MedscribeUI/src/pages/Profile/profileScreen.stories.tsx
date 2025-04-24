import type { Meta, StoryObj } from '@storybook/react';
import ProfileScreen from './profileScreen';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const ProfileScreenWrapper = () => {
  const queryClient = new QueryClient();

  return (
    <div>
      <QueryClientProvider client={queryClient}>
        <MantineProvider>
          <ProfileScreen />
        </MantineProvider>
      </QueryClientProvider>
    </div>
  );
};

const meta: Meta<typeof ProfileScreenWrapper> = {
  title: 'ProfileScreen',
  component: ProfileScreenWrapper,
  parameters: {
    layout: 'fullscreen',
  }
};
export default meta;

type Story = StoryObj<typeof ProfileScreenWrapper>;

export const Wrapped: Story = {};
