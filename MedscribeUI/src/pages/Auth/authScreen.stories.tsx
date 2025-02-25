import type { Meta, StoryObj } from '@storybook/react';
import AuthScreen from './authScreen';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();
const AuthLayoutWrapper = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <AuthScreen />
      </MantineProvider>
    </QueryClientProvider>
  );
};

const meta: Meta<typeof AuthLayoutWrapper> = {
  title: 'AuthScreen',
  component: AuthLayoutWrapper,
};
export default meta;

type Story = StoryObj<typeof AuthLayoutWrapper>;

export const Wrapped: Story = {};
