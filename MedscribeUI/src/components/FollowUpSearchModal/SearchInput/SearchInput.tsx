import React from 'react';

interface SearchInputProps {
  query: string;
  setQuery?: (query: string) => void;
  className?: string;
}

const SearchInput: React.FC<SearchInputProps> = ({
  query,
  setQuery,
  className,
}) => {
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!setQuery) return;
    setQuery(e.target.value);
  };

  return (
    <div
      className={` h-[60px] w-full px-4 py-2 bg-white border border-gray-300 shadow-sm hover:border-gray-400 rounded-t-lg ${className ?? ''}`}
    >
      <div className="flex items-center h-full">
        <svg
          className="w-6 h-6 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
        <input
          type="text"
          value={query}
          onChange={handleSearchChange}
          className="ml-2 w-full h-full bg-transparent text-gray-700 text-lg outline-none"
          placeholder="Search Previous Visits"
        />
      </div>
    </div>
  );
};

export default SearchInput;
