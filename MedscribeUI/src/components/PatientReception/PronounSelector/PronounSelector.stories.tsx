import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import PronounSelector from './PronounSelector';
import { useState } from 'react';

const PronounSelectorWrapper = (
  props: React.ComponentProps<typeof PronounSelector>,
) => {
  const [pronounState, setPronounState] = useState(props.pronoun || '');

  const handleSetPronoun = (value: string) => {
    setPronounState(value);
    props.setPronoun(value);
  };

  return (
    <MantineProvider>
      <PronounSelector pronoun={pronounState} setPronoun={handleSetPronoun} />
    </MantineProvider>
  );
};

const meta: Meta<typeof PronounSelectorWrapper> = {
  title: 'Components/PatientReception/PronounSelector',
  component: PronounSelectorWrapper,
  parameters: {
    layout: 'centered',
  },
  argTypes: {
    setPronoun: { action: 'setPronoun' },
  },
};

export default meta;
type Story = StoryObj<typeof PronounSelectorWrapper>;

export const Default: Story = {
  args: {
    pronoun: '',
    setPronoun: (value: string) => {
      console.log('Pronoun set to:', value);
    },
  },
  name: 'Default',
};

export const WithValue: Story = {
  args: {
    pronoun: 'they/them',
    setPronoun: (value: string) => {
      console.log('Pronoun set to:', value);
    },
  },
  name: 'With Selected Value',
};
