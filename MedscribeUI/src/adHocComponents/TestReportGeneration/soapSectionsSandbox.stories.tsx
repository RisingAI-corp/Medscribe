import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import GenerateReportTest from './reportGeneration';
import '@mantine/core/styles.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const GenerateReportTestWrapper = () => {
  const queryClient = new QueryClient();

  return (
    <>
      <QueryClientProvider client={queryClient}>
        <MantineProvider>
          <GenerateReportTest />
        </MantineProvider>
      </QueryClientProvider>
    </>
  );
};

const meta: Meta<typeof GenerateReportTestWrapper> = {
  title: 'GenerateReportTest',
  component: GenerateReportTestWrapper,
};
export default meta;

type Story = StoryObj<typeof GenerateReportTestWrapper>;

export const Wrapped: Story = {};
