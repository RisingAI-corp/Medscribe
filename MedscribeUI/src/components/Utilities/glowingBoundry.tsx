import { Button, keyframes } from '@mantine/core/styles'; // Corrected import
import { IconExternalLink } from '@tabler/icons-react';

const pulseGlowPurple = keyframes({
  '0%': {
    boxShadow: '0 0 0 0 rgba(128, 0, 128, 0)', // Transparent purple
    borderColor: 'transparent',
  },
  '50%': {
    boxShadow: '0 0 8px 4px rgba(128, 0, 128, 0.4)', // Glowing purple
    borderColor: 'purple',
  },
  '100%': {
    boxShadow: '0 0 0 0 rgba(128, 0, 128, 0)',
    borderColor: 'transparent',
  },
});

interface Props {
  truthyValue: any; // Or a more specific boolean type
  visitSearchValue?: string;
}

export function PulsingGlowButton({ truthyValue, visitSearchValue }: Props) {
  const styles = {
    root: {
      '--pulse-animation': truthyValue
        ? `${pulseGlowPurple} 1.5s ease-in-out infinite`
        : 'none',
      animation: 'var(--pulse-animation)',
      border: '1px solid transparent', // Initial transparent border
      transition: 'border-color 0.3s ease', // Smooth transition for when pulsing stops
    },
  };

  return (
    <Button
      rightSection={<IconExternalLink size={16} />}
      variant="light"
      color="blue"
      className="h-[36px] text-sm"
      style={styles.root}
    >
      {visitSearchValue || 'Link Visit'}
    </Button>
  );
}

// Example usage in a component:
import { useState, useEffect } from 'react';

function MyComponent() {
  const [isActive, setIsActive] = useState(false);

  useEffect(() => {
    const intervalId = setInterval(() => {
      setIsActive(prev => !prev);
    }, 3000); // Example: Pulse every 3 seconds

    return () => clearInterval(intervalId);
  }, []);

  return (
    <div>
      <PulsingGlowButton
        truthyValue={isActive}
        visitSearchValue="My Custom Link"
      />
    </div>
  );
}

export default MyComponent;
