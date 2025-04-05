import { Tooltip } from '@mantine/core';
import { useDebouncedNameChange } from '../../hooks/useDebounceNameChange';

export interface ProfileSummaryCardProps {
  name: string;
  description: string;
  onChange: (newName: string) => void;
  handleUpdateName: (newName: string) => void;
}

function ProfileSummaryCard({
  name,
  description,
  onChange,
  handleUpdateName,
}: ProfileSummaryCardProps) {
  const { nameRef, nameValue, setNameValue, debouncedNameChange } =
    useDebouncedNameChange({
      name,
      onChange,
      handleUpdateName,
    });

  const isEmpty = !nameValue.trim();

  return (
    <div className="bg-gray-100 p-5 border border-gray-300 shadow-md">
      <div className="flex items-center">
        <div className="relative w-full">
          <Tooltip
            label="Name field cannot be empty"
            opened={isEmpty}
            position="top"
            withArrow
          >
            <input
              type="text"
              ref={nameRef}
              value={nameValue}
              onChange={e => {
                setNameValue(e.target.value);
                debouncedNameChange(e.target.value);
              }}
              placeholder="Enter patient's name"
              required
              className={`border-b-2 ${isEmpty ? 'border-red-500' : 'border-gray-400'} 
                          focus:outline-none hover:border-blue-700 focus:border-blue-500 pl-0 pb-1 pt-1 text-sm bg-transparent font-bold`}
              style={{ width: '15rem' }}
            />
          </Tooltip>
        </div>
      </div>
      <div className="text-sm text-gray-600 mb-5">{description}</div>
    </div>
  );
}

export default ProfileSummaryCard;
