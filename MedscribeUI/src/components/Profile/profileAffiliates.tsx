import { useState } from 'react';
import { Button, Card, Text, Group, CopyButton, Badge, Stack, Divider } from '@mantine/core';
import { IconCopy, IconUsers, IconGift, IconChartBar } from '@tabler/icons-react';

const ProfileAffiliates = () => {
  const [affiliateLink] = useState('https://medscribe.com/signup?ref=USER123'); // This would come from your backend
  const [referralStats] = useState({
    totalReferrals: 12,
    activeUsers: 8,
    earnedCredits: 240
  });

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-4">MedScribe Affiliate Program</h1>
        <Text size="lg" className="text-gray-600 mb-6">
          Share MedScribe with your colleagues and earn credits for every new user you bring to our platform.
          Help others streamline their medical documentation while getting rewarded for your referrals.
        </Text>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Group justify="space-between" mb="md">
            <IconUsers size={24} className="text-blue-500" />
            <Badge color="blue" variant="light">Referral Bonus</Badge>
          </Group>
          <Text size="sm" c="dimmed">
            Earn $20 in credits for each new user who signs up through your link and subscribes to a paid plan.
          </Text>
        </Card>

        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Group justify="space-between" mb="md">
            <IconGift size={24} className="text-green-500" />
            <Badge color="green" variant="light">No Limit</Badge>
          </Group>
          <Text size="sm" c="dimmed">
            There's no limit to how many people you can refer. The more you share, the more you earn!
          </Text>
        </Card>

        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Group justify="space-between" mb="md">
            <IconChartBar size={24} className="text-purple-500" />
            <Badge color="purple" variant="light">Track Progress</Badge>
          </Group>
          <Text size="sm" c="dimmed">
            Monitor your referral success with real-time statistics and earnings tracking.
          </Text>
        </Card>
      </div>

      <Card shadow="sm" padding="lg" radius="md" withBorder className="mb-8">
        <Text size="xl" fw={500} mb="md">Your Affiliate Link</Text>
        <Group justify="space-between" className="bg-gray-50 p-4 rounded-md">
          <Text size="sm" className="font-mono">{affiliateLink}</Text>
          <CopyButton value={affiliateLink} timeout={2000}>
            {({ copied, copy }) => (
              <Button
                color={copied ? 'teal' : 'blue'}
                onClick={copy}
                leftSection={<IconCopy size={16} />}
              >
                {copied ? 'Copied' : 'Copy Link'}
              </Button>
            )}
          </CopyButton>
        </Group>
      </Card>

      <Card shadow="sm" padding="lg" radius="md" withBorder>
        <Text size="xl" fw={500} mb="md">Your Referral Statistics</Text>
        <Stack gap="md">
          <Group justify="space-between">
            <Text>Total Referrals</Text>
            <Badge size="lg" variant="light">{referralStats.totalReferrals}</Badge>
          </Group>
          <Divider />
          <Group justify="space-between">
            <Text>Active Users</Text>
            <Badge size="lg" variant="light" color="green">{referralStats.activeUsers}</Badge>
          </Group>
          <Divider />
          <Group justify="space-between">
            <Text>Earned Credits</Text>
            <Badge size="lg" variant="light" color="blue">${referralStats.earnedCredits}</Badge>
          </Group>
        </Stack>
      </Card>
    </div>
  );
};

export default ProfileAffiliates;
