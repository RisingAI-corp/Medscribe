import { useState } from 'react';
import { Checkbox, Loader } from '@mantine/core';
import { IoMdMailOpen, IoMdMailUnread, IoMdTrash } from 'react-icons/io';

export interface PatientPreviewCardProps {
  id: string;
  patientName: string;
  dateOfRecording: string;
  timeOfRecording: string;
  durationOfRecording: string;
  sessionSummary: string;
  loading: boolean;
  isSelected: boolean;
  isChecked: boolean;
  selectAllToggle: boolean;
  readStatus: boolean;
  onClick: (id: string) => void;
  handleRemovePatient: (id: string) => void;
  handleToggleCheckbox: (id: string, checked: boolean) => void;
  handleMarkRead: (id: string) => void;
  handleUnMarkRead: (id: string) => void;
}

const PatientPreviewCard = ({
  id,
  patientName,
  dateOfRecording,
  timeOfRecording,
  durationOfRecording,
  sessionSummary,
  loading = true,
  isSelected,
  isChecked,
  selectAllToggle,
  onClick,
  handleRemovePatient,
  handleToggleCheckbox,
  handleMarkRead,
  handleUnMarkRead,
  readStatus,
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

  console.log('new read status ', readStatus);

  return (
    <div
      className={`flex items-center justify-between px-3 py-2 border border-gray-300 rounded bg-white cursor-pointer transition-colors duration-300 relative 
        ${isSelected ? 'bg-blue-200' : isHovered ? 'bg-blue-100' : ''}`}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onClick={() => {
        if (!readStatus) {
          handleMarkRead(id);
        }
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
        <div
          className={` ${
            readStatus ? 'font-normal' : 'font-bold'
          } font-bold text-sm text-gray-800 mb-1 
          }`}
        >
          {patientName}
        </div>
        <div className="text-xs text-gray-600">
          {dateOfRecording} {timeOfRecording} ({durationOfRecording})
        </div>
        <div className="text-xs text-gray-500 mt-1">{sessionSummary}</div>
      </div>

      {loading ? (
        <Loader size="sm" color="blue" />
      ) : (
        <div className="flex items-center">
          {/* Conditional Mail Log Icon (Open or Unread) */}
          {readStatus ? (
            <span
              onClick={event => {
                event.stopPropagation();
                handleUnMarkRead(id); // Handle the mark as unread action
              }}
              className={`cursor-pointer ml-4 ${isHovered ? 'block' : 'hidden'}`} // Increased margin and space
              title="Open Mail Log"
              role="button"
              aria-label="open-mail-log"
            >
              <IoMdMailOpen className="text-gray-500" size={18} />{' '}
              {/* Grey color for open mail */}
            </span>
          ) : (
            <span
              onClick={event => {
                event.stopPropagation();
                handleMarkRead(id); // Handle the mark as unread action
              }}
              className={`cursor-pointer ml-4 ${isHovered ? 'block' : 'hidden'}`} // Increased margin and space
              title="Open Mail Log"
              role="button"
              aria-label="open-mail-log"
            >
              <IoMdMailUnread className="text-gray-500" size={18} />{' '}
              {/* Grey color for unread mail */}
            </span>
          )}

          {/* Trash Icon */}
          <span
            onClick={event => {
              event.stopPropagation();
              handleRemovePatient(id);
            }}
            className={`cursor-pointer ml-4 ${isHovered ? 'block' : 'hidden'}`} // Adjust margin for spacing
            title="Delete"
            role="button"
            aria-label="delete"
          >
            <IoMdTrash />
          </span>
        </div>
      )}
    </div>
  );
};

export default PatientPreviewCard;
