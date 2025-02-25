import type { Meta, StoryObj } from '@storybook/react';
import HomeScreen from './HomeScreen';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const HomeScreenWrapper = () => {
  return (
    <MantineProvider>
      <HomeScreen />
    </MantineProvider>
  );
};

const meta: Meta<typeof HomeScreenWrapper> = {
  title: 'HomeScreen',
  component: HomeScreenWrapper,
};
export default meta;

type Story = StoryObj<typeof HomeScreenWrapper>;

export const Wrapped: Story = {};
