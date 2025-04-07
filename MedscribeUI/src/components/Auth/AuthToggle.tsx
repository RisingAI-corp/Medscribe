import React from 'react';

interface AuthToggleProps {
  isSignUp: boolean;
  onToggle: () => void;
}

const AuthToggle: React.FC<AuthToggleProps> = ({ isSignUp, onToggle }) => {
  return (
    <div className="text-center">
      <p className="text-gray-700">
        {isSignUp ? 'Already have an account?' : "Don't have an account?"}{' '}
        <button
          onClick={onToggle}
          className="text-blue-600 hover:text-blue-700 transition-colors duration-200"
        >
          {isSignUp ? 'Sign in here' : 'Sign up now'}
        </button>
      </p>
    </div>
  );
};

export default AuthToggle;
