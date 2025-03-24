import React from 'react';
import { Tooltip } from '@mantine/core';
import GoogleIcon from '../../assets/google-icon.webp';

interface SocialAuthProps {
  isSignUp: boolean;
}

const SocialAuth: React.FC<SocialAuthProps> = ({ isSignUp }) => {
  return (
    <>
      <div className="flex items-center gap-4">
        <div className="flex-grow h-px bg-gray-300"></div>
        <span className="text-sm text-gray-500">OR</span>
        <div className="flex-grow h-px bg-gray-300"></div>
      </div>

      <Tooltip label="In beta - coming soon!" position="bottom" withArrow>
        <button
          disabled
          className="w-full py-3 px-4 flex items-center justify-center gap-3 border border-gray-300 rounded-lg opacity-50 cursor-not-allowed"
        >
          <img src={GoogleIcon} alt="Google" className="w-6 h-6" />
          <span className="text-gray-700">
            {isSignUp ? 'Sign up with Google' : 'Sign in with Google'}
          </span>
        </button>
      </Tooltip>
    </>
  );
};

export default SocialAuth; 