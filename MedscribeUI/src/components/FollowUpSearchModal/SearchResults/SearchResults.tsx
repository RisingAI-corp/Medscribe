export interface SearchResultItem {
  id: string;
  patientName: string;
  dateOfRecording: string;
  shortenedSummary?: string;
}

interface SearchResultsProps {
  filteredResults: SearchResultItem[];
  onSelectItem: (item: SearchResultItem) => void;
  selectedItems: SearchResultItem[];
}

const SearchResults: React.FC<SearchResultsProps> = ({ 
  filteredResults, 
  onSelectItem, 
  selectedItems 
}) => {
  const isSelected = (id: string) => {
    return selectedItems.some(item => item.id === id);
  };

  return (
    <div className="mt-4">
      <div className="h-[300px] overflow-y-auto border border-gray-200 rounded-lg">
        {filteredResults.length > 0 ? (
          <ul className="divide-y divide-gray-200">
            {filteredResults.map((item) => (
              <li 
                key={item.id}
                className={`px-4 py-3 cursor-pointer hover:bg-gray-50 transition-colors ${
                  isSelected(item.id) ? 'bg-blue-50' : ''
                }`}
                onClick={() => onSelectItem(item)}
              >
                <div className="flex items-center">
                  <div className={`w-5 h-5 rounded-full flex items-center justify-center mr-3 ${
                    isSelected(item.id) ? 'bg-blue-500 text-white' : 'border border-gray-300'
                  }`}>
                    {isSelected(item.id) && (
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    )}
                  </div>
                  <div className="flex-1">
                    <div className="font-medium text-gray-800">{item.patientName}</div>
                    <div className="text-sm text-gray-500">{item.dateOfRecording}</div>
                    {item.shortenedSummary && (
                      <div className="text-sm text-gray-600 mt-1 line-clamp-2">{item.shortenedSummary}</div>
                    )}
                  </div>
                </div>
              </li>
            ))}
          </ul>
        ) : (
          <div className="flex items-center justify-center h-full">
            <p className="text-gray-500">No results found</p>
          </div>
        )}
      </div>
      <div className="text-sm text-gray-600 mt-2">
        {filteredResults.length} results found
      </div>
    </div>
  );
};

export default SearchResults;
