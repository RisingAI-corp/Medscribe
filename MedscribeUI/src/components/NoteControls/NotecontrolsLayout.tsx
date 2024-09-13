import { Select } from '@mantine/core';
import { Button } from '@mantine/core';

const visitSelectorLabel = 'Visit';

function NoteControlsLayout() {
  return (
    <>
      <Select label={visitSelectorLabel} placeholder={'Default'} data={['A', 'B']} />
      <Button.Group>
        <Button variant="default">First</Button>
        <Button variant="default">Second</Button>
        <Button variant="default">Third</Button>
      </Button.Group>
      {/* HIHIHIIi */}
    </>
  );
}

export default NoteControlsLayout;
