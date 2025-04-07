import { Select } from '@mantine/core';

interface PronounSelectorProps {
  pronoun: string;
  setPronoun: (pronoun: string) => void;
}

const PronounSelector = ({ pronoun, setPronoun }: PronounSelectorProps) => {
  const pronounOptions = [
    { value: 'he/him', label: 'He/Him' },
    { value: 'she/her', label: 'She/Her' },
    { value: 'they/them', label: 'They/Them' },
    { value: 'other', label: 'Other' },
  ];

  return (
    <Select
      placeholder="Select pronouns"
      data={pronounOptions}
      variant="unstyled"
      value={pronoun}
      onChange={value => {
        setPronoun(value ?? '');
      }}
    />
  );
};

export default PronounSelector;
