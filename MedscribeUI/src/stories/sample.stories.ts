// SampleComponent.stories.tsx
import React from 'react';
import type { Meta, StoryObj } from '@storybook/react';
import SampleComponent from './SampleComponent';

// Define Meta for the SampleComponent
const meta: Meta<typeof SampleComponent> = {
  title: 'SampleComponent',
  component: SampleComponent,
};
export default meta;

// Create a type for the StoryObj to ensure type safety in stories
type Story = StoryObj<typeof SampleComponent>;


// Bind a copy of the template to create different stories
export const Default: Story = {};
Default.args = {
  title: 'Welcome to my app',
  subtitle: 'This is a sample subtitle',
  backgroundColor: '#f0f0f0',
  textColor: '#333',
  isDisabled: false,
};

export const Disabled: Story = {};
Disabled.args = {
  title: 'Disabled button',
  subtitle: 'This button is disabled',
  backgroundColor: '#f0f0f0',
  textColor: '#333',
  isDisabled: true,
};

export const DarkMode: Story = {};
DarkMode.args = {
  title: 'Dark mode',
  subtitle: 'This is a sample subtitle',
  backgroundColor: '#333',
  textColor: '#fff',
  isDisabled: false,
};
