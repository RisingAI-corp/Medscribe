import { useEffect, useRef, useState } from 'react';
import { useDebouncedCallback } from 'use-debounce';

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
  const localNameRef = useRef(name);
  const [currentName, setName] = useState(name);
  localNameRef.current = currentName;

  const debouncedNameChange = useDebouncedCallback((value: string) => {
    onChange(value);
  }, 200);

  useEffect(() => {
    return () => {
      debouncedNameChange.cancel();
      if (localNameRef.current !== name) {
        console.log('updated Name');
        handleUpdateName(localNameRef.current);
      }
    };
  }, []);

  return (
    <div className="bg-gray-100 p-5 border border-gray-300 shadow-md">
      <input
        type="text"
        value={currentName}
        onChange={e => {
          setName(e.target.value);
          debouncedNameChange(e.target.value);
        }}
        className="text-lg font-bold border-none bg-transparent outline-none w-full"
      />
      <div className="text-sm text-gray-600">One Liner</div>
      <div className="text-sm text-gray-600 mb-5">{description}</div>
    </div>
  );
}

export default ProfileSummaryCard;
