import { Button } from '@mantine/core';
import { useState, useEffect } from 'react';

interface ThreeBtnSelectorProps {
  options: string[]; // Array of button labels
  initialSelected?: string; // Optional initial selected value
  onSelect: (value: string) => void; // Callback when a button is selected
}

function ThreeBtnSelector({ options, initialSelected = '', onSelect }: ThreeBtnSelectorProps) {
  const [selected, setSelected] = useState(initialSelected);

  // Update the selected state when initialSelected prop changes
  useEffect(() => {
    setSelected(initialSelected);
  }, [initialSelected]);

  // Handler to set selected option and notify parent component
  const handleSelect = (value: string) => {
    setSelected(value);
    onSelect(value);
  };

  return (
    <Button.Group variant="outline" aria-label="button group">
      {options.map(option => (
        <Button key={option} variant={selected === option ? 'filled' : 'outline'} onClick={() => handleSelect(option)}>
          {option}
        </Button>
      ))}
    </Button.Group>
  );
}

export default ThreeBtnSelector;
