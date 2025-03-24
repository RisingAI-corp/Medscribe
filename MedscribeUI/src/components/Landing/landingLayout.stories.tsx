import type { Meta, StoryObj } from '@storybook/react';
import LandingLayout from '../../pages/Landing/landingScreen';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const LandingLayoutWrapper = () => {
  return (
    <MantineProvider>
      <LandingLayout />
    </MantineProvider>
  );
};

const meta: Meta<typeof LandingLayoutWrapper> = {
  title: 'LandingLayout',
  component: LandingLayoutWrapper,
};
export default meta;

type Story = StoryObj<typeof LandingLayoutWrapper>;

export const Wrapped: Story = {};
