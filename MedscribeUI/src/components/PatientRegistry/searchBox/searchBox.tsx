import '@mantine/core/styles.css';

import { IconSearch } from '@tabler/icons-react';
interface SearchBoxProps {
  value: string;
  onChange?: (value: string) => void;
  classname?: string;
}

const SearchBox = ({ value = '', onChange, classname }: SearchBoxProps) => {
  return (
    <div
      className={`flex items-center border border-gray-300 rounded px-2 py-1 w-full max-w-sm bg-white ${classname ?? ''}`}
    >
      <IconSearch size={20} className="text-gray-400 mr-2" />

      <input
        type="text"
        value={value}
        onChange={e => {
          if (!onChange) return;
          onChange(e.target.value);
        }}
        placeholder="Search for notes by name"
        className="border-none outline-none flex-grow text-sm text-gray-800 placeholder-gray-500"
      />
    </div>
  );
};

export default SearchBox;
