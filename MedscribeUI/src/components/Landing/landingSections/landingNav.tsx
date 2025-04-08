import { Button } from '@mantine/core';
import logo from '../../../assets/medscribe-logo.png';
import { useNavigate } from 'react-router-dom';

function LandingNav() {
  const navigate = useNavigate();

  const scrollToSection = (sectionId: string) => {
    const element = document.getElementById(sectionId);
    if (element) {
      const offset = 60;
      const elementPosition = element.getBoundingClientRect().top;
      const offsetPosition = elementPosition + window.pageYOffset - offset;

      window.scrollTo({
        top: offsetPosition,
        behavior: 'smooth',
      });
    }
  };

  return (
    <div className="h-full w-full flex bg-white justify-between items-center px-8">
      <div className="flex items-center">
        <img src={logo} alt="MedScribe Logo" className="w-10 h-10" />
        <span className="ml-2 text-2xl text-gray-800">MedScribe</span>
      </div>
      <div className="flex items-center gap-8">
        <button
          onClick={() => {
            scrollToSection('hero');
          }}
          className="text-gray-600 hover:text-gray-900"
        >
          Home
        </button>
        <button
          onClick={() => {
            scrollToSection('pricing');
          }}
          className="text-gray-600 hover:text-gray-900"
        >
          Pricing
        </button>
        <button
          onClick={() => {
            scrollToSection('faq');
          }}
          className="text-gray-600 hover:text-gray-900"
        >
          FAQ
        </button>
      </div>

      <div className="flex items-center gap-4">
        <Button variant="outline" onClick={() => void navigate('/SignUp')}>
          Sign Up
        </Button>
        <Button variant="filled" onClick={() => void navigate('/SignIn')}>
          Sign In
        </Button>
      </div>
    </div>
  );
}

export default LandingNav;
