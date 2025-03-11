import type { Meta, StoryObj } from '@storybook/react';
import HomeScreen from './homeScreen';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const HomeScreenWrapper = () => {
  const queryClient = new QueryClient();

  return (
    <div>
      <QueryClientProvider client={queryClient}>
        <MantineProvider>
          <HomeScreen />
        </MantineProvider>
      </QueryClientProvider>
    </div>
  );
};

const meta: Meta<typeof HomeScreenWrapper> = {
  title: 'HomeScreen',
  component: HomeScreenWrapper,
};
export default meta;

type Story = StoryObj<typeof HomeScreenWrapper>;

export const Wrapped: Story = {};
