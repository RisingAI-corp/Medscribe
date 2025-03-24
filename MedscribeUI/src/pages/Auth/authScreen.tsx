import { useState } from 'react';
import SignIn from '../../components/Auth/authSignIn';
import SignUp from '../../components/Auth/authSignUp';

export function AuthScreen() {
  const [isSignUp, setIsSignUp] = useState(true);

  return (
    <div>
      {isSignUp ? (
        <SignUp onToggleAuth={() => setIsSignUp(false)} />
      ) : (
        <SignIn onToggleAuth={() => setIsSignUp(true)} />
      )}
    </div>
  );
}

export default AuthScreen;


