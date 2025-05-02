import { useState, useRef } from 'react';
import { Button, Modal, Loader } from '@mantine/core';
import { IconCheck, IconAlertTriangle } from '@tabler/icons-react'; // Fixed icons

interface VerificationCodeModalProps {
  isOpen: boolean;
  onClose: () => void;
  email: string;
  verificationStatus: 'idle' | 'pending' | 'success' | 'error';
  handleSubmitVerificationCode: (code: string) => void;
}

const MAX_CODE_LENGTH = 6;

const VerificationCodeModal = ({
  isOpen,
  email,
  verificationStatus,
  handleSubmitVerificationCode,
}: VerificationCodeModalProps) => {
  console.log('verification status', verificationStatus);
  const [code, setCode] = useState('');
  const [prevCode, setPrevCode] = useState(code);
  const inputRef = useRef<HTMLInputElement>(null);
  const isVerifying = verificationStatus === 'pending';
  const isVerified = verificationStatus === 'success';
  const isError = verificationStatus === 'error';
  console.log('isError', isError);

  const handleVerification = () => {
    if (code.length !== MAX_CODE_LENGTH && prevCode !== code) return;
    setPrevCode(code);
    handleSubmitVerificationCode(code);
  };

  const handleCodeChange = () => {
    if (inputRef.current) {
      const value = inputRef.current.value.replace(/\D+/g, ''); // Remove non-digit characters
      if (value.length <= MAX_CODE_LENGTH) {
        setCode(value);
      }
    }
  };

  return (
    <Modal
      opened={isOpen}
      onClose={() => {
        return;
      }}
      title="Please Enter Your Email Verification Code"
      overlayProps={{ opacity: 0.55, blur: 3 }}
      transitionProps={{ transition: 'fade', duration: 300 }}
      className="rounded-xl"
      centered
      withCloseButton={false}
    >
      <div className="p-4">
        <p className="text-gray-600 mb-4">
          A verification code has been sent to{' '}
          <span className="font-medium">{email}</span>.
        </p>
        <div
          className="flex justify-between gap-2 mb-4"
          onClick={() => inputRef.current?.focus()}
        >
          {Array.from({ length: MAX_CODE_LENGTH }).map((_, i) => (
            <div
              key={i}
              className={`w-10 h-12 flex items-center justify-center text-xl font-mono rounded-md border-2 transition-colors duration-200 ${isVerifying ? 'bg-gray-100 text-gray-500' : ''} ${code[i] ? 'border-blue-500' : 'border-gray-300'} ${isError ? 'border-red-500' : ''}`}
            >
              {code[i] || '\u2000'}{' '}
              {/* Use a zero-width space as a placeholder */}
            </div>
          ))}
        </div>

        <input
          ref={inputRef}
          type="text" // Use type="text" for custom validation
          value={code}
          onChange={handleCodeChange}
          className="absolute inset-0 opacity-0 cursor-none pointer-events-none"
          disabled={verificationStatus === 'pending'}
          aria-hidden="true"
        />

        {isVerifying && (
          <div className="flex items-center justify-center mt-4 text-gray-500">
            <Loader size="sm" className="mr-2" />
            Verifying...
          </div>
        )}

        {isError && (
          <div className="mt-4 text-red-500 flex items-center">
            <IconAlertTriangle size={20} className="mr-2" />
            Incorrect code. Please try again.
          </div>
        )}

        {isVerified && (
          <div className="mt-4 text-green-500 flex items-center">
            <IconCheck size={20} className="mr-2" />
            Email verified!
          </div>
        )}

        <div className="mt-6 flex justify-end">
          <Button
            onClick={handleVerification}
            disabled={isVerifying || prevCode === code}
            className={
              (isVerifying ? 'opacity-70 cursor-not-allowed' : '') +
              (code.length !== MAX_CODE_LENGTH
                ? ' opacity-50 cursor-not-allowed'
                : '')
            }
          >
            Verify
          </Button>
        </div>
      </div>
    </Modal>
  );
};

export default VerificationCodeModal;
