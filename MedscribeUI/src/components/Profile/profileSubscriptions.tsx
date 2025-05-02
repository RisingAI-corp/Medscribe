import { useEffect, useState } from 'react';
import { Button, Card, Text, Title } from '@mantine/core';
import { useAtom } from 'jotai';
import { userAtom } from '../../states/userAtom';
import { set } from 'zod';

const ProfileSubscriptions = () => {
  const [user] = useAtom(userAtom);
  const [currentPlan, setCurrentPlan] = useState('Free Tier');
  const [isLoading, setIsLoading] = useState(false);

  const handleManageSubscription = async () => {
    if (currentPlan === 'Free Tier') {
      setCurrentPlan('Pro Tier');
    } else {
      setCurrentPlan('Free Tier');
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <Card 
        shadow="sm" 
        padding="xl" 
        radius="xl" 
        withBorder
        className="bg-white/95 backdrop-blur-sm"
      >
        <Title 
          order={2} 
          className="mb-6 font-light text-gray-700"
        >
          Subscription Details
        </Title>
        
        <div className="mb-8">
          <Text fw={400} size="lg" className="text-gray-600">Current Plan</Text>
          <Text size="xl" className="text-blue-500 font-light mt-1">{currentPlan}</Text>
        </div>

        <div className="mb-8">
          <Text fw={400} size="lg" className="text-gray-600">Features</Text>
          <ul className="list-disc pl-6 mt-3 space-y-2 text-gray-600">
            {currentPlan === 'Free Tier' ? (
              <>
                <li className="font-light">10 Free Lightning-Fast Transcriptions</li>
                <li className="font-light">1GB of Secure Notes Storage (HIPAA Compliant)</li>
                <li className="font-light">10 Free Sessions</li>
                <li className="font-light">24/7 Email Support</li>
              </>
            ) : (
              <>
                <li className="font-light">Unlimited Lightning-Fast Transcriptions</li>
                <li className="font-light">10GB of Secure Notes Storage (HIPAA Compliant)</li>
                <li className="font-light">Unlimited Sessions</li>
                <li className="font-light">24/7 Email Support</li>
              </>
            )}
          </ul>
        </div>

        <Button
          fullWidth
          size="lg"
          loading={isLoading}
          onClick={handleManageSubscription}
          className="bg-blue-500 hover:bg-blue-600 transition-colors duration-200 rounded-xl font-light"
        >
          Manage Subscription
        </Button>
      </Card>
    </div>
  );
};

export default ProfileSubscriptions;
