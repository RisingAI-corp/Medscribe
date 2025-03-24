import './App.css';
import HomeScreen from './pages/Home/homeScreen';
import { useMutation } from '@tanstack/react-query';
import { isAuthenticatedAtom } from './states/userAtom';
import { userAtom } from './states/userAtom';
import {
  currentlySelectedPatientAtom,
  patientsAtom,
} from './states/patientsAtom';
import { useAtom } from 'jotai';
import { checkAuth } from './api/checkAuth';
import AuthScreen from './pages/Auth/authScreen';
import FallbackScreen from './pages/Fallback/fallbackScreen';
import LandingScreen from './pages/Landing/landingScreen';
import { useEffect, useState } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

function App() {
  const [timerActive, setTimerActive] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useAtom(isAuthenticatedAtom);

  const [_ignoredProvider, setProvider] = useAtom(userAtom);
  const [_ignoredPatients, setPatients] = useAtom(patientsAtom);
  const [_ignoredSelectedPatient, setCurrentlySelectedPatient] = useAtom(
    currentlySelectedPatientAtom,
  );

  useEffect(() => {
    setTimeout(() => {
      setTimerActive(false);
    }, 800);
    checkAuthMutation.mutate(undefined, {
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
        setPatients(reports ?? []);
        setIsAuthenticated(true);
        if (reports && reports.length > 0) {
          setCurrentlySelectedPatient(reports[0].id);
          return;
        }
        setCurrentlySelectedPatient('');
      },
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const checkAuthMutation = useMutation({
    mutationFn: checkAuth,
    onError: error => {
      console.error('Error adding todo:', error);
    },
  });

  const { isPending, isSuccess, isIdle } = checkAuthMutation;

  const renderAuthComponent = () => {
    if (isPending || isIdle || timerActive) {
      return <FallbackScreen />;
    } else if (isSuccess || isAuthenticated) {
      return <HomeScreen />;
    } else {
      return <AuthScreen />;
    }
  };

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/landing" element={<LandingScreen />} />
        <Route path="/" element={renderAuthComponent()} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
