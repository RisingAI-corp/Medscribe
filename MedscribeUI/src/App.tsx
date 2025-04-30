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
import ProfileScreen from './pages/Profile/profileScreen';
import { useEffect, useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';

function App() {
  const [timerActive, setTimerActive] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useAtom(isAuthenticatedAtom);

  const [_ignoredProvider, setProvider] = useAtom(userAtom);
  const [_ignoredPatients, setPatients] = useAtom(patientsAtom);
  const [_ignoredSelectedPatient, setCurrentlySelectedPatient] = useAtom(
    currentlySelectedPatientAtom,
  );

  useEffect(() => {
    if (isAuthenticated) return;

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

        setPatients(reports);
        setIsAuthenticated(true);
        if (reports.length > 0) {
          setCurrentlySelectedPatient(reports[0].id);
        } else {
          setCurrentlySelectedPatient('');
        }
      },
    });
  }, []);

  const checkAuthMutation = useMutation({
    mutationFn: checkAuth,
    onSuccess: data => {
      console.log('Authenticated:', data);
    },
    onError: error => {
      console.error('Error adding todo:', error);
    },
  });

  const { isPending, isSuccess, isIdle } = checkAuthMutation;

  const renderAuthComponent = () => {
    if (isPending || isIdle || timerActive) {
      return <FallbackScreen />;
    } else if (isSuccess && isAuthenticated) {
      return <HomeScreen />;
    } else {
      return <LandingScreen />;
    }
  };

  const renderProfileComponent = () => {
    if (isAuthenticated || true) { // TODO: Remove this
      return <ProfileScreen />;
    } else {
      return <Navigate to="/" />;
    }
  };

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={renderAuthComponent()} />
        <Route path="/SignUp" element={<AuthScreen isSignUpRoute={true} />} />
        <Route path="/SignIn" element={<AuthScreen isSignUpRoute={false} />} />
        <Route path="/Profile" element={renderProfileComponent()} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
