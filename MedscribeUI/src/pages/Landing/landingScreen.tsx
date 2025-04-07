import LandingNav from '../../components/Landing/landingSections/landingNav';
import LandingHero from '../../components/Landing/landingSections/landingHero';
import LandingPricing from '../../components/Landing/landingSections/landingPricing';
import LandingFAQ from '../../components/Landing/landingSections/landingFAQ';
import LandingFooter from '../../components/Landing/landingSections/landingFooter';
import landingBackground from '../../assets/landing-bg.png';

function LandingLayout() {
  return (
    <div className="w-full min-w-[750px] overflow-x-auto">
      <img
        src={landingBackground}
        alt="Landing Background"
        className="fixed inset-0 w-full h-full object-cover -z-50"
      />
      
      {/* Navbar */}
      <div className="fixed top-0 left-0 w-full min-w-[750px] h-[60px] shadow-md z-50 flex items-center justify-between">
        <LandingNav />
      </div>

      {/* Spacing to account for fixed navbar */}
      <div className="h-[60px]"/>

      {/* Sections */}
      <section id="hero" className="min-h-[700px] h-[calc(100vh-70px)] flex items-center justify-center">
        <LandingHero />
      </section>
      <section id="pricing" className="min-h-[700px] h-[calc(100vh-60px)] flex items-center justify-center">
        <LandingPricing />
      </section>
      <section id="faq" className="min-h-[700px] h-[calc(100vh-60px)] flex items-center justify-center">
        <LandingFAQ />
      </section>
      <section id="footer" className=" w-full flex items-center justify-center">
        <LandingFooter />
      </section>
    </div>

  );
}

export default LandingLayout;
