import type { Meta, StoryObj } from '@storybook/react';
import AuthLayout from './authLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const AuthLayoutWrapper = () => {
  return (
    <MantineProvider>
        <AuthLayout />
    </MantineProvider>
  );
};

const meta: Meta<typeof AuthLayoutWrapper> = {
  title: 'AuthLayout',
  component: AuthLayoutWrapper,
};
export default meta;

type Story = StoryObj<typeof AuthLayoutWrapper>;

export const Wrapped: Story = {};
