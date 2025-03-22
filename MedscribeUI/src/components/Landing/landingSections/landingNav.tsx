import { Button } from '@mantine/core';
import logo from '../../../assets/medscribe-logo.png';
import { useNavigate } from 'react-router-dom';

function LandingNav() {
  const navigate = useNavigate();
  return (
    <div className="h-20 w-full flex bg-white justify-between items-center px-8 fixed top-0 z-10">
        <div className="flex items-center">
            <img src={logo} alt="MedScribe Logo" className="w-10 h-10" />
            <span className="ml-2 text-2xl text-gray-800">MedScribe</span>
        </div>
        <div className="flex items-center gap-8">
            <a href="#hero" className="text-gray-600 hover:text-gray-900">Home</a>
            <a href="#pricing" className="text-gray-600 hover:text-gray-900">Pricing</a>
            <a href="#faq" className="text-gray-600 hover:text-gray-900">FAQ</a>
        </div>

        <div className="flex items-center gap-4">
            <Button variant="outline" onClick={() => navigate('/signin')}>Sign Up</Button>
            <Button variant="filled" onClick={() => navigate('/signin')}>Sign In</Button>
        </div>
    </div>
  );
}

export default LandingNav;
