import type { Meta, StoryObj } from '@storybook/react';
import LandingScreen from './landingScreen';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const LandingLayoutWrapper = () => {
  return (
    <MantineProvider>
      <LandingScreen />
    </MantineProvider>
  );
};

const meta: Meta<typeof LandingLayoutWrapper> = {
  title: 'LandingScreen',
  component: LandingLayoutWrapper,
};
export default meta;

type Story = StoryObj<typeof LandingLayoutWrapper>;

export const Wrapped: Story = {};
