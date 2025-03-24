import LandingNav from '../../components/Landing/landingSections/landingNav';
import LandingHero from '../../components/Landing/landingSections/landingHero';
import LandingPricing from '../../components/Landing/landingSections/landingPricing';
import LandingFAQ from '../../components/Landing/landingSections/landingFAQ';
import LandingFooter from '../../components/Landing/landingSections/landingFooter';

function LandingLayout() {
  return (
    <div className="w-full flex flex-col min-h-screen">
      <LandingNav />
      <div className="flex flex-col w-full snap-y snap-mandatory h-screen overflow-y-auto scroll-smooth">
        <section
          id="hero"
          className="min-h-screen w-full flex items-center justify-center snap-start"
        >
          <LandingHero />
        </section>
        <section
          id="pricing"
          className="min-h-screen w-full flex items-center justify-center snap-start"
        >
          <LandingPricing />
        </section>
        <section
          id="faq"
          className="min-h-screen w-full flex items-center justify-center snap-start"
        >
          <LandingFAQ />
        </section>
        <section
          id="footer"
          className="h-[100px] w-full flex items-center justify-center snap-start"
        >
          <LandingFooter />
        </section>
      </div>
    </div>
  );
}

export default LandingLayout;
