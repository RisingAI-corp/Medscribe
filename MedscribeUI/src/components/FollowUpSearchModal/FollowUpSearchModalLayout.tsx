import React, { useState } from 'react';
import { SearchResultItem } from './SearchResults/SearchResults';
import SearchButton from './SearchButton/SearchButton';
import SearchInput from './SearchInput/SearchInput';
import SearchResults from './SearchResults/SearchResults';
import useSearch from '../../hooks/useSearch';

interface FollowUpSearchModalLayoutProps {
  selectedItems: SearchResultItem[];
  setSelectedItems: React.Dispatch<React.SetStateAction<SearchResultItem[]>>;
  mockData?: SearchResultItem[]; // For demo/story purposes
}

const FollowUpSearchModalLayout: React.FC<FollowUpSearchModalLayoutProps> = ({ 
  selectedItems, 
  setSelectedItems,
  mockData = [] 
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  
  // Use the search hook to manage the query and filtered results
  const [filteredResults, query, setQuery] = useSearch<SearchResultItem>(
    mockData,
    (item) => `${item.patientName} ${item.dateOfRecording} ${item.shortenedSummary || ''}`,
    300 // debounce time
  );

  const handleToggleModal = () => {
    setIsModalOpen(!isModalOpen);
  };

  const handleSelectItem = (item: SearchResultItem) => {
    setSelectedItems((prev) => {
      const isAlreadySelected = prev.some((selected) => selected.id === item.id);
      
      if (isAlreadySelected) {
        return prev.filter((selected) => selected.id !== item.id);
      } else {
        return [...prev, item];
      }
    });
  };

  return (
    <div className="relative">
      {/* Main button to open the modal */}
      <div onClick={handleToggleModal}>
        <SearchButton selectedItems={selectedItems.map(item => item.patientName)} />
      </div>

      {/* Modal overlay */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50">
          {/* Backdrop with blur effect */}
          <div 
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            onClick={handleToggleModal}
          />
          
          {/* Modal content */}
          <div 
            className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg shadow-xl p-6 w-[50%]"
            onClick={(e) => e.stopPropagation()}
          >
            <SearchInput query={query} setQuery={setQuery} />
            
            <SearchResults 
              filteredResults={filteredResults} 
              onSelectItem={handleSelectItem} 
              selectedItems={selectedItems} 
            />

            <div className="flex justify-end mt-4">
              <button 
                onClick={handleToggleModal}
                className="text-green-500 hover:text-green-700"
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default FollowUpSearchModalLayout;
