import type { Meta, StoryObj } from '@storybook/react';
import VisitReportLayout from './visitReportLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const WIDTH = '500px';
const HEIGHT = '600px';
const PADDING = '10px';

const VisitReportWrapper = () => {
  return (
    <div
      style={{
        width: WIDTH,
        height: HEIGHT,
        border: '1px dotted red',
        padding: PADDING,
      }}
    >
      <MantineProvider>
        <VisitReportLayout />
      </MantineProvider>
    </div>
  );
};

const Wrapper = () => {
  return <VisitReportWrapper />;
};

const meta: Meta<typeof Wrapper> = {
  title: 'VisitReportLayout',
  component: Wrapper,
};
export default meta;

type Story = StoryObj<typeof Wrapper>;

export const Wrapped: Story = {};
