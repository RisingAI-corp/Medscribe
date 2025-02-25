import type { Meta, StoryObj } from '@storybook/react';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';
import SearchBox from './searchBox';
import { useState } from 'react';

const Wrapper = () => {
  const [value, setValue] = useState('');
  const wrapperStyles = {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '90vh',
  };

  const innerStyles = {
    width: '600px',
    height: '500px',
    border: '1px dotted red',
    padding: '10px',
  };

  return (
    <div style={wrapperStyles}>
      <div style={innerStyles}>
        <MantineProvider>
          <SearchBox
            value={value}
            onChange={newValue => {
              setValue(newValue);
            }}
          />
        </MantineProvider>
      </div>
    </div>
  );
};

const meta: Meta<typeof Wrapper> = {
  title: 'searchBox',
  component: Wrapper,
};
export default meta;

type Story = StoryObj<typeof Wrapper>;

export const Wrapped: Story = {};
