import { useNavigate } from "react-router-dom";
import { Button } from "@mantine/core";

export const Header = () => {
  const navigate = useNavigate();
  return (
    <div className="flex items-center bg-blue-500 h-10 px-4">
      <button
        className="text-white text-lg font-semibold leading-none hover:opacity-80 transition-opacity"
        onClick={() => navigate('/')}
      >
        Medscribe
      </button>
      <div className="flex-1"></div>
      <Button
        variant="subtle"
        onClick={() => navigate('/profile')}
        className="text-white rounded-full w-8 h-8 p-0 flex items-center justify-center hover:bg-blue-600 transition-colors"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
          />
        </svg>
      </Button>
    </div>
  );
};
