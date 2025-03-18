import { Paper, Title } from '@mantine/core';

function LandingLayout() {
  return (
    <div className="h-screen w-full flex">
      <div className="flex-1 bg-contain bg-center bg-[url('./assets/authBackgroundImage.png')]" />

      <div className="flex-1 max-w-[450px] flex items-center justify-center shadow-md">
        <Paper radius={0} p={30} className="shadow-3d">
          <Title order={2} className="text-black text-center mb-6 ">
            Landing Page... Coming Soon!
          </Title>
        </Paper>
      </div>
    </div>
  );
}

export default LandingLayout;
