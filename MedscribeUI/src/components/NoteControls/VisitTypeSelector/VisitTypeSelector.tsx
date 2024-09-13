import { Select } from '@mantine/core';

interface VisitTypeSelecterProps {
  label: string;
  placeholder: string;
  data: string[];
}

function VisitTypeSelecter({ label, placeholder, data }: VisitTypeSelecterProps) {
  return <Select label={label} placeholder={placeholder} data={data} />;
}

export default VisitTypeSelecter;
