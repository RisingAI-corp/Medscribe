import { useState } from 'react';
import { Checkbox, Loader } from '@mantine/core';

export interface PatientPreviewCardProps {
  id: string;
  patientName: string;
  dateOfRecording: string;
  timeOfRecording: string;
  durationOfRecording: string;
  shortenedSummary: string;
  loading: boolean;
  isSelected: boolean;
  isChecked: boolean;
  selectAllToggle: boolean;
  onClick: (id: string) => void;
  handleRemovePatient: (id: string) => void;
  handleToggleCheckbox: (id: string, checked: boolean) => void;
}

const PatientPreviewCard = ({
  id,
  patientName,
  dateOfRecording,
  timeOfRecording,
  durationOfRecording,
  shortenedSummary,
  loading = true,
  isSelected,
  isChecked,
  selectAllToggle,
  onClick,
  handleRemovePatient,
  handleToggleCheckbox,
}: PatientPreviewCardProps) => {
  const [isHovered, setIsHovered] = useState(false);

  const handleMouseEnter = () => {
    setIsHovered(true);
  };
  const handleMouseLeave = () => {
    setIsHovered(false);
  };

  const handleCheckboxChange = (checked: boolean) => {
    handleToggleCheckbox(id, checked);
  };

  return (
    <div
      className={`flex items-center justify-between px-3 py-2 border border-gray-300 rounded bg-white cursor-pointer transition-colors duration-300 relative 
        ${isSelected ? 'bg-gray-300' : isHovered ? 'bg-gray-100' : ''}`}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onClick={() => {
        onClick(id);
      }}
    >
      {selectAllToggle && (
        <Checkbox
          size="sm"
          checked={isChecked}
          onChange={e => {
            handleCheckboxChange(e.currentTarget.checked);
          }}
        />
      )}

      <div className="ml-3 flex-grow">
        <div className="font-bold text-sm text-gray-800 mb-1">
          {patientName}
        </div>
        <div className="text-xs text-gray-600">
          {dateOfRecording} {timeOfRecording} ({durationOfRecording})
        </div>
        <div className="text-xs text-gray-500 mt-1">{shortenedSummary}</div>
      </div>

      {loading ? (
        <Loader size="sm" color="blue" />
      ) : (
        <span
          onClick={event => {
            event.stopPropagation();
            handleRemovePatient(id);
          }}
          className={`cursor-pointer ${isHovered ? 'block' : 'hidden'}`}
          title="Delete"
          role="button"
          aria-label="delete"
        >
          🗑️
        </span>
      )}
    </div>
  );
};

export default PatientPreviewCard;
