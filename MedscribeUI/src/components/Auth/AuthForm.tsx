import React from 'react';
import ErrorAlert from './ErrorAlert';

interface AuthFormProps {
  isSignUp: boolean;
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
  error: { message: string; type: 'login' | 'signup' } | null;
  onNameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onEmailChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onPasswordChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onConfirmPasswordChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onErrorDismiss: () => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
}

const AuthForm: React.FC<AuthFormProps> = ({
  isSignUp,
  name,
  email,
  password,
  confirmPassword,
  error,
  onNameChange,
  onEmailChange,
  onPasswordChange,
  onConfirmPasswordChange,
  onErrorDismiss,
  onSubmit,
}) => {
  return (
    <form onSubmit={onSubmit} className="w-full">
      <h1 className="text-2xl font-semibold text-center text-gray-800 mb-6">
        {isSignUp ? 'Create Your Account' : 'Welcome Back!'}
      </h1>
      <div className="flex flex-col gap-4">
        {error && (
          <ErrorAlert message={error.message} onDismiss={onErrorDismiss} />
        )}

        {isSignUp && (
          <div className="relative">
            <input
              type="text"
              value={name}
              onChange={onNameChange}
              className="w-full px-4 py-2.5 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Full Name"
              required
            />
          </div>
        )}

        <div className="relative">
          <input
            type="email"
            value={email}
            onChange={onEmailChange}
            className="w-full px-4 py-2.5 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Email"
            required
          />
        </div>

        <div className="relative">
          <input
            type="password"
            value={password}
            onChange={onPasswordChange}
            className="w-full px-4 py-2.5 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Password"
            required
          />
        </div>

        {isSignUp && (
          <div className="relative">
            <input
              type="password"
              value={confirmPassword}
              onChange={onConfirmPasswordChange}
              className="w-full px-4 py-2.5 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Confirm Password"
              required
            />
          </div>
        )}

        <button
          type="submit"
          className="w-full py-2.5 mt-4 text-white text-sm rounded-lg bg-blue-600 hover:bg-blue-700 transition-colors duration-200"
        >
          {isSignUp ? 'Sign Up' : 'Sign In'}
        </button>
      </div>
    </form>
  );
};

export default AuthForm;
