import { Select, Button } from '@mantine/core';
import { useState } from 'react';
import ThreeBtnSelector from '../Utilities/ThreeBtnSelector';

const visitSelectorLabel = 'Visit';
const defaultPronoun = 'He';

function NoteControlsLayout() {
  const [selectedPronoun, setSelectedPronoun] = useState(defaultPronoun);
  const [selectedNoteLength, setSelectedNoteLength] = useState('Standard');
  const [selectedVisitType, setSelectedVisitType] = useState('New Patient');

  const handleVisitTypeSelect = (value: string | null) => {
    setSelectedVisitType(value || ''); // Set default value to empty string if value is null
    console.log('Selected Visit Type:', value); // !!!DEBUG ONLY!!! - Remove before production
  };

  const handlePronounSelect = (value: string) => {
    setSelectedPronoun(value);
    console.log('Selected Pronoun:', value); // !!!DEBUG ONLY!!! - Remove before production
  };

  const handleNoteLengthSelect = (value: string) => {
    setSelectedNoteLength(value);
    console.log('Selected Note Length:', value); // !!!DEBUG ONLY!!! - Remove before production
  };

  const handleMagicEdit = () => {
    console.log('Magic Edit Clicked'); // !!!DEBUG ONLY!!! - Remove before production
  };

  return (
    <>
      {/* Visit Selector */}
      <Select
        label={visitSelectorLabel}
        placeholder="Select visit type"
        data={['New Patient', 'Returning Patient']}
        value={selectedVisitType}
        onChange={handleVisitTypeSelect}
      />
      <hr />

      {/* Pronoun Selector */}
      <span>Pronoun Selector</span>
      <ThreeBtnSelector
        options={['She', 'He', 'They']}
        initialSelected={selectedPronoun}
        onSelect={handlePronounSelect}
      />
      {/* Note Length Selector */}
      <span>Note Length Selector</span>
      <ThreeBtnSelector
        options={['Concise', 'Standard', 'Detailed']}
        initialSelected={selectedNoteLength}
        onSelect={handleNoteLengthSelect}
      />
      <hr />

      <Button onClick={handleMagicEdit} style={{ width: '100%' }}>
        {' '}
        Magic Edit{' '}
      </Button>

      {/* Display Selected [FOR DEBUGGING] */}
      {/* 
      <br />
      <span>Selected Pronoun: {selectedPronoun}</span>
      <br />
      <span>Selected Note Length: {selectedNoteLength}</span> 
      <br />
      <span>Selected Visit Type: {selectedVisitType}</span>
      */}
    </>
  );
}

export default NoteControlsLayout;
