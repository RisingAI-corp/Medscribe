import { Select, Button } from '@mantine/core';
import { useState } from 'react';
import ThreeBtnSelector from '../Utilities/ThreeBtnSelector';

const visitSelectorLabel = 'Visit';
const defaultPronoun = 'He';

function NoteControlsLayout() {
  const [selectedPronoun, setSelectedPronoun] = useState(defaultPronoun);
  const [selectedVisitType, setSelectedVisitType] = useState('New Patient');

  const handleVisitTypeSelect = (value: string | null) => {
    setSelectedVisitType(value || '');
  };

  const handlePronounSelect = (value: string) => {
    setSelectedPronoun(value);
  };

  const handleRegenerate = () => {};

  return (
    <>
      <Select
        label={visitSelectorLabel}
        placeholder="Select visit type"
        data={['New Patient', 'Returning Patient']}
        value={selectedVisitType}
        onChange={handleVisitTypeSelect}
      />

      <hr />

      <span>Pronoun Selector</span>
      <ThreeBtnSelector
        buttonLabelOptions={['She', 'He', 'They']}
        initialSelected={selectedPronoun}
        onSelect={handlePronounSelect}
      />

      <hr />

      <Button onClick={handleRegenerate} style={{ width: '100%' }}>
        {' '}
        Regenerate Report{' '}
      </Button>
    </>
  );
}

export default NoteControlsLayout;
