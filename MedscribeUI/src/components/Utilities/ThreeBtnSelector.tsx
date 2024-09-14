import { Button } from '@mantine/core';
import { useState, useEffect } from 'react';

interface ThreeBtnSelectorProps {
  buttonLabelOptions: string[];
  initialSelected?: string;
  onSelect: (value: string) => void;
}

function ThreeBtnSelector({ buttonLabelOptions, initialSelected = '', onSelect }: ThreeBtnSelectorProps) {
  const [selected, setSelected] = useState(initialSelected);

  useEffect(() => {
    setSelected(initialSelected);
  }, [initialSelected]);

  const handleSelect = (value: string) => {
    setSelected(value);
    onSelect(value);
  };

  return (
    <Button.Group variant="outline" aria-label="button group">
      {buttonLabelOptions.map(option => (
        <Button key={option} variant={selected === option ? 'filled' : 'outline'} onClick={() => handleSelect(option)}>
          {option}
        </Button>
      ))}
    </Button.Group>
  );
}

export default ThreeBtnSelector;
