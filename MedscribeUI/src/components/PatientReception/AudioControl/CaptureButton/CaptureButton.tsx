import React from 'react';
import { IconMicrophone } from '@tabler/icons-react';

interface CaptureButtonProps {
  onClick?: () => void;
  disabled?: boolean;
}

const CaptureButton: React.FC<CaptureButtonProps> = ({ onClick }) => {
  return (
    <button
      className={`flex flex-row items-center justify-center bg-white border-[1px] border-blue-500 rounded-md px-4 py-2`}
      onClick={onClick}
    >
      <IconMicrophone size={24} className="text-blue-500 mr-2" />
      <span className="text-sm text-blue-700 ">Capture</span>
    </button>
  );
};

export default CaptureButton; 