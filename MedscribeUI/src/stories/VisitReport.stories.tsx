import type { Meta, StoryObj } from '@storybook/react';
import VisitReportLayout from '../components/VisitReport/visitReportLayout';
import { MantineProvider } from '@mantine/core';
import '@mantine/core/styles.css';

const WIDTH = '1000px';
const HEIGHT = '800px';
const PADDING = '10px';

const NoteControlsWrapper = ({}) => {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '90vh' }}>
      <div style={{ width: WIDTH, height: HEIGHT, border: '1px dotted red', padding: PADDING }}>
        <MantineProvider>
          <VisitReportLayout />
        </MantineProvider>
      </div>
    </div>
  );
};

const Wrapper = () => {
  return <NoteControlsWrapper />;
};

const meta: Meta<typeof Wrapper> = {
  title: 'VisitReportLayout',
  component: Wrapper,
};
export default meta;

type Story = StoryObj<typeof Wrapper>;

export const Wrapped: Story = {};
