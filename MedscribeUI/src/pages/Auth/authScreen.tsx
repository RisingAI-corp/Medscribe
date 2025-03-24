import { Paper, Title } from '@mantine/core';
import AuthenticationForm from '../../components/Auth/authFrom';
import { loginProvider } from '../../api/login';
import { createProvider } from '../../api/signUp';
import { useMutation } from '@tanstack/react-query';
import { isAuthenticatedAtom } from '../../states/userAtom';
import { userAtom } from '../../states/userAtom';
import { patientsAtom } from '../../states/patientsAtom';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import { useAtom } from 'jotai';
import { useState } from 'react';

export function AuthScreen() {
  const [_, setProvider] = useAtom(userAtom);
  const [__, setIsAuthenticated] = useAtom(isAuthenticatedAtom);
  const [___, setPatients] = useAtom(patientsAtom);
  const [____, setCurrentlySelectedPatient] = useAtom(
    currentlySelectedPatientAtom,
  );
  const [emailInUse, setUseEmail] = useState(false);
  const [loginFailed, setLoginFailed] = useState(false);

  const loginMutation = useMutation({
    mutationFn: loginProvider,
    onSuccess: ({
      id,
      name,
      email,
      reports,
      subjectiveStyle,
      objectiveStyle,
      assessmentAndPlanStyle,
      patientInstructionsStyle,
      summaryStyle,
    }) => {
      setProvider({
        ID: id,
        name: name,
        email: email,
        subjectiveStyle,
        objectiveStyle,
        assessmentAndPlanStyle,
        patientInstructionsStyle,
        summaryStyle,
      });
      setIsAuthenticated(true);
      setPatients(reports ?? []);
      if (reports && reports.length > 0) {
        setCurrentlySelectedPatient(reports[0].id);
        return;
      }
      setCurrentlySelectedPatient('');
    },
    onError: error => {
      setLoginFailed(true);
      console.error('Error adding todo:', error);
    },
  });

  const signUpMutation = useMutation({
    mutationFn: createProvider,
    onSuccess: ({ id, name, email }) => {
      setProvider({
        ID: id,
        name: name,
        email: email,
        subjectiveStyle: '',
        objectiveStyle: '',
        assessmentAndPlanStyle: '',
        patientInstructionsStyle: '',
        summaryStyle: '',
      });
      setIsAuthenticated(true);
      setCurrentlySelectedPatient('');
    },
    onError: error => {
      if (error.message === 'status conflict: user already exists') {
        setUseEmail(true);
      }
      console.error('Error signingUp User:', error);
    },
  });
  const handleSignUp = (name: string, email: string, password: string) => {
    signUpMutation.mutate({ name, email, password });
  };

  const handleLogin = (email: string, password: string) => {
    loginMutation.mutate({ email, password });
  };

  return (
    <div className="h-screen w-full flex">
      <div className="flex-1 bg-contain bg-center bg-[url('./assets/authBackgroundImage.png')]" />

      <div className="flex-1 max-w-[450px] flex items-center justify-center shadow-md">
        <Paper radius={0} p={30} className="shadow-3d">
          <Title order={2} className="text-black text-center mb-6 ">
            Welcome To Medscribe!
          </Title>
          <AuthenticationForm
            handleRegister={handleSignUp}
            handleLogin={handleLogin}
            emailInUse={emailInUse}
            showLoginFailedNotification={loginFailed}
            handleCloseNotification={() => {
              setLoginFailed(false);
            }}
          />
        </Paper>
      </div>
    </div>
  );
}

export default AuthScreen;
