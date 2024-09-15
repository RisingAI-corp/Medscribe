import { Button } from '@mantine/core';

interface BtnGroupSelector {
  buttonLabelOptions: string[];
  selectedBtn: string;
  onSelect: (value: string) => void;
}

function BtnGroupSelector({ buttonLabelOptions, selectedBtn, onSelect }: BtnGroupSelector) {
  return (
    <Button.Group variant="outline" aria-label="button group">
      {buttonLabelOptions.map(option => (
        <Button key={option} variant={selectedBtn === option ? 'filled' : 'outline'} onClick={() => onSelect(option)}>
          {option}
        </Button>
      ))}
    </Button.Group>
  );
}

export default BtnGroupSelector;
