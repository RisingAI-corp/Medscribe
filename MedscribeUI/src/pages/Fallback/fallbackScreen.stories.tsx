import type { Meta, StoryObj } from '@storybook/react';
import FallbackScreen from './fallbackScreen';
import '@mantine/core/styles.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();
const FallBackPage = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <FallbackScreen />
    </QueryClientProvider>
  );
};

const meta: Meta<typeof FallBackPage> = {
  title: 'FallBackScreen',
  component: FallBackPage,
};
export default meta;

type Story = StoryObj<typeof FallBackPage>;

export const Wrapped: Story = {};
