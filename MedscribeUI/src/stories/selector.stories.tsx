// SampleComponent.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
// import { WrappedVisitSelector } from './wrappedVisitSelector';
import NoteControlsLayout from '../components/NoteControls/NotecontrolsLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

// Define Meta for the SampleComponent
const Wrapper = () => {
  return (
    <MantineProvider>
      <NoteControlsLayout />
    </MantineProvider>
  );
};
const meta: Meta<typeof Wrapper> = {
  title: 'NotecontrolsLayout',
  component: Wrapper,
};
export default meta;

// Create a type for the StoryObj to ensure type safety in stories
type Story = StoryObj<typeof Wrapper>;

// Bind a copy of the template to create different stories
export const Wrapped: Story = {};
