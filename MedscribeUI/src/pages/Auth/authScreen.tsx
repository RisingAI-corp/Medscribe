import { useState } from 'react';
import landingBg from '../../assets/hero-image.png';
import { loginProvider } from '../../api/login';
import { createProvider } from '../../api/signUp';
import { AuthResponse } from '../../api/serverResponses';
import { useMutation } from '@tanstack/react-query';
import { isAuthenticatedAtom, userAtom } from '../../states/userAtom';
import { useNavigate } from 'react-router-dom';
import {
  patientsAtom,
  currentlySelectedPatientAtom,
} from '../../states/patientsAtom';
import { useAtom } from 'jotai';
import AuthForm from '../../components/Auth/AuthForm';
import SocialAuth from '../../components/Auth/SocialAuth';
import AuthToggle from '../../components/Auth/AuthToggle';
import logo from '../../assets/medscribe-logo.png';

const AuthScreen = ({ isSignUpRoute }: { isSignUpRoute: boolean }) => {
  const navigate = useNavigate();
  const [isSignUp, setIsSignUp] = useState(isSignUpRoute);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<{
    message: string;
    type: 'login' | 'signup';
  } | null>(null);

  const [_, setProvider] = useAtom(userAtom);
  const [__, setIsAuthenticated] = useAtom(isAuthenticatedAtom);
  const [___, setPatients] = useAtom(patientsAtom);
  const [____, setCurrentlySelectedPatient] = useAtom(
    currentlySelectedPatientAtom,
  );

  const handleSuccess = (data: AuthResponse) => {
    setProvider({
      ID: data.id,
      name: data.name,
      email: data.email,
      subjectiveStyle: data.subjectiveStyle,
      objectiveStyle: data.objectiveStyle,
      assessmentAndPlanStyle: data.assessmentAndPlanStyle,
      summaryStyle: data.summaryStyle,
      patientInstructionsStyle: data.patientInstructionsStyle,
    });
    setIsAuthenticated(true);
    if (data.reports.length > 0) {
      setPatients(data.reports);
      setCurrentlySelectedPatient(data.reports[0].id);
    } else {
      setCurrentlySelectedPatient('');
    }
    void navigate('/');
  };

  const loginMutation = useMutation({
    mutationFn: loginProvider,
    onSuccess: handleSuccess,
    onError: error => {
      setError({
        message: 'Invalid email or password. Please try again.',
        type: 'login',
      });
      console.error('Error logging in:', error);
    },
  });

  const signUpMutation = useMutation({
    mutationFn: createProvider,
    onSuccess: handleSuccess,
    onError: error => {
      console.log('failed to sign up');
      if (error.message === 'status conflict: user already exists') {
        setError({
          message:
            'This email is already registered. Please use a different email or sign in.',
          type: 'signup',
        });
      }
      console.error('Error signing up:', error);
    },
  });

  const handleSubmit = () => {
    setError(null);

    if (isSignUp && password !== confirmPassword) {
      setError({ message: 'Passwords do not match', type: 'signup' });
      return;
    }
    if (isSignUp) {
      signUpMutation.mutate({ name, email, password });
    } else {
      loginMutation.mutate({ email, password });
    }
  };

  return (
    <div className="flex min-h-[700px] h-screen">
      {/* Left side - Auth Form */}
      <div className="w-1/3 px-[50px] flex flex-col min-w-[500px]">
        <div className="mt-8 flex items-center gap-2" onClick={() => void navigate('/')}>
          <img src={logo} alt="Medscribe Logo" className="w-16 h-16"/>
          <span className="text-3xl font-semibold text-gray-800">
            Medscribe
          </span>
        </div>

        <div className="flex-1 flex items-center px-[50px]">
          <div className="flex flex-col gap-8 w-full">
            <AuthForm
              isSignUp={isSignUp}
              name={name}
              email={email}
              password={password}
              confirmPassword={confirmPassword}
              error={error}
              onNameChange={e => {
                setName(e.target.value);
              }}
              onEmailChange={e => {
                setEmail(e.target.value);
              }}
              onPasswordChange={e => {
                setPassword(e.target.value);
              }}
              onConfirmPasswordChange={e => {
                setConfirmPassword(e.target.value);
              }}
              onErrorDismiss={() => {
                setError(null);
              }}
              onSubmit={e => {
                e.preventDefault();
                handleSubmit();
              }}
            />

            <SocialAuth isSignUp={isSignUp} />

            <AuthToggle
              isSignUp={isSignUp}
              onToggle={() => {
                setIsSignUp(!isSignUp);
                setError(null);
              }}
            />
          </div>
        </div>
      </div>

      {/* Right side - Background Image */}
      <div
        className="w-2/3 h-full bg-cover bg-center bg-no-repeat"
        style={{
          backgroundImage: `url(${landingBg})`,
          backgroundSize: 'cover',
        }}
      />
    </div>
  );
};

export default AuthScreen;
