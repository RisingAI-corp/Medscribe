import React, { useState } from 'react';
import { useAtom } from 'jotai';
import { SearchResultItem } from './SearchResults/SearchResults';
import SearchButton from './SearchButton/SearchButton';
import SearchInput from './SearchInput/SearchInput';
import SearchResults from './SearchResults/SearchResults';
import useSearch from '../../hooks/useSearch';
import { searchVisitsAtom } from './derivedAtoms';

interface FollowUpSearchModalLayoutProps {
  selectedItem: string;
  setSelectedItem: React.Dispatch<React.SetStateAction<string>>;
}

const FollowUpSearchModalLayout: React.FC<FollowUpSearchModalLayoutProps> = ({ 
  selectedItem, 
  setSelectedItem,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchVisits] = useAtom(searchVisitsAtom);
  
  // Use the search hook to manage the query and filtered results
  const [filteredResults, query, setQuery] = useSearch<SearchResultItem>(
    searchVisits,
    (item) => `${item.patientName} ${item.dateOfRecording} ${item.shortenedSummary || ''}`,
    300 // debounce time
  );

  const handleToggleModal = () => {
    setIsModalOpen(!isModalOpen);
  };

  const handleSelectItem = (item: SearchResultItem) => {
    setSelectedItem(item.patientName);
    setIsModalOpen(false); // Close modal after selection
  };

  return (
    <div className="relative">
      {/* Main button to open the modal */}
      <div onClick={handleToggleModal}>
        <SearchButton selectedItem={selectedItem} />
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
            className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg shadow-xl w-[50%]"
            onClick={(e) => e.stopPropagation()}
          >
            <SearchInput query={query} setQuery={setQuery} />
            
            <SearchResults 
              filteredResults={filteredResults} 
              onSelectItem={handleSelectItem} 
              selectedItemName={selectedItem} 
            />

          </div>
        </div>
      )}
    </div>
  );
};

export default FollowUpSearchModalLayout;
