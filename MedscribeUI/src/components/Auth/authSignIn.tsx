import { useState } from 'react';
import GoogleIcon from '../../assets/google-icon.webp';
import landingBg from '../../assets/landing-bg.png';
import logo from '../../assets/medscribe-logo.png';
import { useNavigate } from 'react-router-dom';
import { loginProvider } from '../../api/login';
import { useMutation } from '@tanstack/react-query';
import { isAuthenticatedAtom } from '../../states/userAtom';
import { userAtom } from '../../states/userAtom';
import { patientsAtom } from '../../states/patientsAtom';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { useAtom } from 'jotai';

interface SignInProps {
  onToggleAuth: () => void;
}

const SignIn = ({ onToggleAuth }: SignInProps) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loginFailed, setLoginFailed] = useState(false);
  const navigate = useNavigate();
  
  const [_, setProvider] = useAtom(userAtom);
  const [__, setIsAuthenticated] = useAtom(isAuthenticatedAtom);
  const [___, setPatients] = useAtom(patientsAtom);
  const [____, setCurrentlySelectedPatient] = useAtom(currentlySelectedPatientAtom);

  const loginMutation = useMutation({
    mutationFn: loginProvider,
    onSuccess: ({
      id,
      name,
      email,
      reports,
      subjectiveStyle,
      objectiveStyle,
      assessmentStyle,
      planningStyle,
      summaryStyle,
    }) => {
      setProvider({
        ID: id,
        name: name,
        email: email,
        subjectiveStyle,
        objectiveStyle,
        assessmentStyle,
        planningStyle,
        summaryStyle,
      });
      setIsAuthenticated(true);
      setPatients(reports ?? []);
      if (reports && reports.length > 0) {
        setCurrentlySelectedPatient(reports[0].id);
      } else {
        setCurrentlySelectedPatient('');
      }
    },
    onError: error => {
      setLoginFailed(true);
      console.error('Error logging in:', error);
    },
  });

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    loginMutation.mutate({ email, password });
  };

  return (
    <div className="flex h-screen">
      {/* Left side - Sign In Form */}
      <div className="w-1/3 px-[50px] flex flex-col min-w-[500px]">
        {/* Logo and Name */}
        <div className="mt-8 flex items-center gap-2">
          <img src={logo} alt="Medscribe Logo" className="w-16 h-16" />
          <span className="text-3xl font-semibold text-gray-800">Medscribe</span>
        </div>

        <div className="flex-1 flex items-center px-[50px]">
          <div className="flex flex-col gap-8 w-full">
            <h1 className="text-3xl font-semibold text-center text-gray-800">
              Welcome Back!
            </h1>

            {/* Form */}
            <form onSubmit={handleSubmit} className="w-full">
              <div className="flex flex-col gap-4">
                {loginFailed && (
                  <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                    <span className="block sm:inline">Invalid email or password. Please try again.</span>
                    <button
                      onClick={() => setLoginFailed(false)}
                      className="absolute top-0 bottom-0 right-0 px-4 py-3"
                    >
                      <span className="sr-only">Dismiss</span>
                      <svg className="fill-current h-4 w-4" role="button" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                        <title>Close</title>
                        <path d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z"/>
                      </svg>
                    </button>
                  </div>
                )}
                <div className="relative">
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Email"
                    required
                  />
                </div>

                <div className="relative">
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Password"
                    required
                  />
                </div>

                <button
                  type="submit"
                  className="w-full py-3 mt-4 text-white text-base rounded-lg bg-blue-600 hover:bg-blue-700 transition-colors duration-200"
                >
                  Sign In
                </button>
              </div>
            </form>

            {/* Divider */}
            <div className="flex items-center gap-4">
              <div className="flex-grow h-px bg-gray-300"></div>
              <span className="text-sm text-gray-500">OR</span>
              <div className="flex-grow h-px bg-gray-300"></div>
            </div>

            {/* Google Sign In Button */}
            <button
              onClick={() => console.log('Sign in with Google')}
              className="w-full py-3 px-4 flex items-center justify-center gap-3 border border-gray-300 rounded-lg hover:bg-gray-50 hover:border-blue-500 transition-colors duration-200"
            >
              <img src={GoogleIcon} alt="Google" className="w-6 h-6" />
              <span className="text-gray-700">Sign in with Google</span>
            </button>

            {/* Sign up link */}
            <div className="text-center">
              <p className="text-gray-700">
                Don't have an account?{' '}
                <button
                  onClick={onToggleAuth}
                  className="text-blue-600 hover:text-blue-700 transition-colors duration-200"
                >
                  Sign up now
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Right side - Background Image */}
      <div 
        className="w-2/3 h-full bg-cover bg-center bg-no-repeat"
        style={{ 
          backgroundImage: `url(${landingBg})`,
          backgroundSize: 'cover'
        }}
      />
    </div>
  );
};

export default SignIn;
