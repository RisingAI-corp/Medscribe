import { IconX } from '@tabler/icons-react';

export interface SearchResultItem {
  id: string;
  patientName: string;
  dateOfRecording: string;
  summary: string;
  condensedSummary: string;
  timeOfRecording: string;
  durationOfRecording: number;
}

interface SearchResultsProps {
  filteredResults: SearchResultItem[];
  onSelectItem: (item: SearchResultItem) => void;
  selectedVisitID: string;
}

const SearchResults: React.FC<SearchResultsProps> = ({
  filteredResults,
  onSelectItem,
  selectedVisitID,
}) => {
  const isSelected = (visitID: string) => {
    return visitID === selectedVisitID;
  };

  const getDuration = (duration: number) => {
    duration = Math.floor(duration / 60);
    const durationToString = `${String(duration)} min`;
    return duration > 1 ? durationToString : '< 1 min';
  };

  const handleUnselectVisit = () => {
    onSelectItem({} as SearchResultItem);
  };

  return (
    <div className="h-[400px] overflow-y-auto border border-gray-200 rounded-b-lg">
      {filteredResults.length > 0 ? (
        <ul className="divide-y divide-gray-200">
          {filteredResults.map(item => (
            <li
              key={item.id}
              className={`px-4 py-3 cursor-pointer hover:bg-[#e8f1ff] transition-colors ${
                isSelected(item.id) ? 'bg-blue-100' : ''
              }`}
              onClick={() => {
                onSelectItem(item);
              }}
            >
              <div className="flex items-center">
                <div className="flex-1">
                  <div className="font-bold text-gray-800">
                    {item.patientName}
                  </div>
                  <div className="text-sm text-gray-500">
                    {item.dateOfRecording} {item.timeOfRecording} (
                    {getDuration(item.durationOfRecording)})
                  </div>
                  {item.summary && (
                    <div className="text-sm text-gray-600 mt-1 line-clamp-2">
                      {item.summary}
                    </div>
                  )}
                </div>
                {isSelected(item.id) && (
                  <IconX
                    // Modify the onClick handler for the icon
                    onClick={event => {
                      event.stopPropagation();
                      console.log('IconX onClick triggered for item:', item.id); // For debugging
                      handleUnselectVisit();
                    }}
                    className="text-red-500 hover:text-red-700 transition-colors cursor-pointer" // Added cursor-pointer here too
                  />
                )}
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
