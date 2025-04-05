import React from 'react';

interface SearchButtonProps {
  selectedItems?: string[];
  maxDisplayLength?: number;
}

const SearchButton: React.FC<SearchButtonProps> = ({ 
  selectedItems = [], 
  maxDisplayLength = 35
}) => {
  const getDisplayText = () => {
    if (selectedItems.length === 0) {
      return "Search follow-ups...";
    }
    
    if (selectedItems.length === 1) {
      return selectedItems[0];
    }
    
    const text = selectedItems.join(', ');
    if (text.length <= maxDisplayLength) {
      return text;
    }
    
    // If we have multiple items and the text is too long
    return `${selectedItems.length} patients selected`;
  };
  
  const displayText = getDisplayText();
    
  return (
    <button className="w-full px-4 py-2 text-left bg-white border border-gray-300 rounded-lg shadow-sm hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
      <div className="flex items-center">
        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <span className="ml-2 text-gray-500">{displayText}</span>
      </div>
    </button>
  );
};

export default SearchButton;
