import landingBackground from '../../../assets/landing-bg.png';
import { landingContent } from '../landingContent';

interface PricingCardProps {
  title: string;
  price: string;
  features: string[];
  isHighlighted?: boolean;
}

const PricingCard = ({
  title,
  price,
  features,
  isHighlighted,
}: PricingCardProps) => {
  return (
    <div
      className={`relative rounded-lg border-[1px] border-black bg-white flex flex-col h-full`}
    >
      <div
        className={`${isHighlighted ? 'bg-blue-600 text-white' : 'bg-blue-100 text-black'} py-4 rounded-md mb-4`}
      >
        <h3 className="text-3xl font-bold text-center mb-2">{title}</h3>
        <div className="text-5xl font-bold text-center">{price}</div>
      </div>

      <ul className="list-disc pt-[20px] pb-[50px] px-[50px] space-y-4 text-xl text-black">
        {features.map((feature, index) => (
          <li key={index}>{feature}</li>
        ))}
      </ul>
    </div>
  );
};

const LandingPricing = () => {
  return (
    <div className="w-full h-full flex items-center justify-center relative bg-cover bg-center">
      <div className="max-w-6xl mx-auto px-4 py-12 relative">
        {/* Centered Pricing grid */}
        <div className="flex flex-wrap justify-center gap-10 max-w-7xl mx-auto">
          {landingContent.pricing.plans.map((plan, index) => (
            <div
              className="w-full md:w-[calc(33.33%-5px)] min-w-[400px]"
              key={index}
            >
              <PricingCard
                title={plan.title}
                price={plan.price}
                features={plan.features}
                isHighlighted={plan.isHighlighted}
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default LandingPricing;
