import { helix } from 'ldrs';

helix.register();

const FallbackScreen = () => {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-gradient-to-b from-blue-100 to-blue-200">
      <l-helix size="200" speed="2.5" color="#1e44e1"></l-helix>
    </div>
  );
};

export default FallbackScreen;
