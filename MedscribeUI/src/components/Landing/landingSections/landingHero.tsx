import heroImage from '../../../assets/LandingImage.png';
import hippaImage from '../../../assets/hippa-logo.svg';
import { landingContent } from '../landingContent';
import { Button } from '@mantine/core';
import LandingCarousel from './landingCarousel';
import { useNavigate } from 'react-router-dom';


function LandingHero() {
  const navigate = useNavigate();
  return (
    <div className="h-full w-full flex flex-col relative">
      <div className="flex h-full w-full">
        <div className="flex flex-col h-full w-1/2 justify-center">
          <div className="flex flex-col gap-[30px] pl-[50px] pr-[100px]">
          
            <span className="text-5xl font-bold bg-gradient-to-r from-[#0772BA] to-[#0493B3] bg-clip-text text-transparent leading-relaxed">
              {landingContent.hero.title}
            </span>
            <span className="text-2xl">{landingContent.hero.subtitle}</span>
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
            <img src={hippaImage} alt="HIPPA" className="w-1/5"/>
          </div>
        </div>
        <div className="flex h-full w-1/2 justify-center items-center">
          <img src={heroImage} alt="Hero Image" className="w-3/5" />
        </div>
      </div>
      <div className="flex h-[100px] w-full">
        <LandingCarousel />
      </div>
    </div>
  );
}

export default LandingHero;