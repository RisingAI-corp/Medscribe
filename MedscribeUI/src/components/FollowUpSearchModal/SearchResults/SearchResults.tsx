export interface SearchResultItem {
  id: string;
  patientName: string;
  dateOfRecording: string;
  shortenedSummary?: string;
}

interface SearchResultsProps {
  filteredResults: SearchResultItem[];
  onSelectItem: (item: SearchResultItem) => void;
  selectedItemName: string;
}

const SearchResults: React.FC<SearchResultsProps> = ({ 
  filteredResults, 
  onSelectItem, 
  selectedItemName 
}) => {
  const isSelected = (name: string) => {
    return name === selectedItemName;
  };

  return (
    <div className="h-[300px] overflow-y-auto border border-gray-200 rounded-b-lg">
      {filteredResults.length > 0 ? (
        <ul className="divide-y divide-gray-200">
          {filteredResults.map((item) => (
            <li 
              key={item.id}
              className={`px-4 py-3 cursor-pointer hover:bg-gray-50 transition-colors ${
                isSelected(item.patientName) ? 'bg-blue-50' : ''
              }`}
              onClick={() => onSelectItem(item)}
            >
              <div className="flex items-center">
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
  );
};

export default SearchResults;
