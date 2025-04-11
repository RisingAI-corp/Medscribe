import { useState } from 'react';
import { Checkbox, Loader, Tooltip } from '@mantine/core';
import { IoMdMailOpen, IoMdMailUnread, IoMdTrash } from 'react-icons/io';

export interface PatientPreviewCardProps {
  id: string;
  patientName: string;
  dateOfRecording: string;
  timeOfRecording: string;
  durationOfRecording: string;
  sessionSummary: string;
  status: string;
  isSelected: boolean;
  isChecked: boolean;
  selectAllToggle: boolean;
  readStatus: boolean;
  loading: boolean;
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
  status,
  isSelected,
  isChecked,
  selectAllToggle,
  onClick,
  handleRemovePatient,
  handleToggleCheckbox,
  handleMarkRead,
  handleUnMarkRead,
  readStatus,
  loading,
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

      <div className="flex items-center" style={{ minWidth: '60px' }}>
        {' '}
        {/* Adjust minWidth as needed */}
        {loading && status !== 'failed' ? (
          <Loader size="sm" color="blue" />
        ) : status === 'failed' ? (
          <Tooltip label="Error generating report" position="right" withArrow>
            <div className="text-red-600 mr-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M8.257 3.099c.765-1.36 2.72-1.36 3.485 0l6.518 11.597c.75 1.335-.213 2.987-1.742 2.987H3.48c-1.53 0-2.492-1.652-1.742-2.987L8.257 3.1zM11 13a1 1 0 10-2 0 1 1 0 002 0zm-1-2a1 1 0 01-1-1V7a1 1 0 112 0v3a1 1 0 01-1 1z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
          </Tooltip>
        ) : (
          status !== 'failed' && (
            <div
              className="mr-2"
              style={{
                width: '24px',
                display: 'flex',
                justifyContent: 'center',
              }}
            >
              {readStatus ? (
                <span
                  onClick={event => {
                    event.stopPropagation();
                    handleUnMarkRead(id);
                  }}
                  className={`cursor-pointer ${isHovered ? 'block' : 'hidden'}`}
                  title="Open Mail Log"
                  role="button"
                  aria-label="open-mail-log"
                >
                  <IoMdMailOpen className="text-gray-500" size={18} />
                </span>
              ) : (
                <span
                  onClick={event => {
                    event.stopPropagation();
                    handleMarkRead(id);
                  }}
                  className={`cursor-pointer ${isHovered ? 'block' : 'hidden'}`}
                  title="Open Mail Log"
                  role="button"
                  aria-label="open-mail-log"
                >
                  <IoMdMailUnread className="text-gray-500" size={18} />
                </span>
              )}
            </div>
          )
        )}
        {/* Trash Icon (always render) */}
        <span
          onClick={event => {
            event.stopPropagation();
            handleRemovePatient(id);
          }}
          className={`cursor-pointer ml-4 ${isHovered ? 'block' : 'hidden'}`}
          title="Delete"
          role="button"
          aria-label="delete"
        >
          <IoMdTrash />
        </span>
      </div>
    </div>
  );
};

export default PatientPreviewCard;
