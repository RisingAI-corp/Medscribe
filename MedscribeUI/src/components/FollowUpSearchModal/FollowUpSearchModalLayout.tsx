import { useState } from 'react';
import { useAtom } from 'jotai';
import { SearchResultItem } from './SearchResults/SearchResults';
import SearchInput from '../SearchInput/SearchInput';
import SearchResults from './SearchResults/SearchResults';
import useSearch from '../../hooks/useSearch';
import { searchVisitsAtom } from './derivedAtoms';

interface FollowUpSearchModalLayoutProps {
  handleSelectedVisit: (visit: SearchResultItem) => void;
  children: React.ReactNode;
}

const FollowUpSearchModalLayout = ({
  handleSelectedVisit,
  children,
}: FollowUpSearchModalLayoutProps) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchVisits] = useAtom(searchVisitsAtom);
  const [selectedVisitID, setSelectedVisitID] = useState<string>('');

  // Use the search hook to manage the query and filtered results
  const [filteredResults, query, setQuery] = useSearch<SearchResultItem>(
    searchVisits,
    item => `${item.patientName} ${item.dateOfRecording} ${item.summary}`,
    300, // debounce time
  );

  const handleToggleModal = () => {
    setIsModalOpen(!isModalOpen);
  };

  const handleSelectItem = (item: SearchResultItem) => {
    console.log('hit');
    setQuery(item.patientName);
    handleSelectedVisit(item);
    setSelectedVisitID(item.id);
    console.log('this is item', item);
    setIsModalOpen(false); // Close modal after selection
  };

  return (
    <div className="relative">
      {/* Main button to open the modal */}
      <div onClick={handleToggleModal}>{children}</div>

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
            onClick={e => {
              e.stopPropagation();
            }}
          >
            <SearchInput query={query} setQuery={setQuery} />

            <SearchResults
              filteredResults={filteredResults}
              onSelectItem={handleSelectItem}
              selectedVisitID={selectedVisitID}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default FollowUpSearchModalLayout;
