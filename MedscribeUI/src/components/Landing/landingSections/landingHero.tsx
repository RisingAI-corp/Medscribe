import landingBackground from '../../../assets/landing-bg.png';
import heroImage from '../../../assets/LandingImage.png';
import hippaImage from '../../../assets/hippa-logo.svg';
import { landingContent } from '../landingContent';
import { Button } from '@mantine/core';
import LandingCarousel from './landingCarousel';
import { useNavigate } from 'react-router-dom';

function LandingHero() {
  const navigate = useNavigate();
  return (
    <div className="h-full w-full flex flex-col relative gap-10">
      <img
        src={landingBackground}
        alt="Landing Background"
        className="absolute inset-0 w-full h-full object-contain -z-10"
      />
      <div className="flex flex-row w-full px-[200px] pt-[150px]">
        <div className="flex flex-col w-1/2 gap-[90px]">
          <h1 className="text-7xl font-bold bg-gradient-to-r from-[#0772BA] to-[#0493B3] bg-clip-text text-transparent">
            {landingContent.hero.title}
          </h1>
          <p className="text-3xl">{landingContent.hero.subtitle}</p>
          <div className="flex items-center gap-4">
            <Button
              variant="filled"
              size="lg"
              onClick={() => {
                void navigate('/Auth');
              }}
            >
              Start for Free
            </Button>
          </div>
          <img src={hippaImage} alt="HIPPA" className="w-1/5" />
        </div>
        <div className="flex flex-col items-center w-1/2">
          <img src={heroImage} alt="Hero Image" className="w-3/4" />
        </div>
      </div>
      <LandingCarousel />
    </div>
  );
}

export default LandingHero;
