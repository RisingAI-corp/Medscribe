import { useNavigate } from 'react-router-dom';
import { IconArrowBackUp, IconLogout2} from '@tabler/icons-react';


interface ProfileSidebarProps {
  onNavChange: (tab: 'settings' | 'affiliates' | 'subscriptions') => void;
  activeTab: 'settings' | 'affiliates' | 'subscriptions';
}

const ProfileSidebar = ({ onNavChange, activeTab }: ProfileSidebarProps) => {
  const navigate = useNavigate();
  const navItems = [
    { id: 'settings', label: 'Settings' },
    { id: 'affiliates', label: 'Affiliates' },
    { id: 'subscriptions', label: 'Subscriptions' },
  ];

  const handleLogout = () => {
    // Implement logout functionality here
    // For example, clear local storage, reset auth state, etc.
    navigate('/login');
  };

  return (
    <div className="w-64 border-r h-full p-4 flex flex-col">
      <button onClick={() => navigate('/')} className="flex items-center text-gray-600 hover:text-gray-900 mb-6">
        <IconArrowBackUp className="h-5 w-5 mr-2" />
        Back
      </button>
      <nav>
        <ul className="space-y-2">
          {navItems.map((item) => (
            <li key={item.id}>
              <button
                onClick={() => onNavChange(item.id as any)}
                className={`w-full text-left px-4 py-2 rounded-md transition ${
                  activeTab === item.id
                    ? 'bg-blue-100 text-blue-700 font-medium'
                    : 'hover:bg-gray-100'
                }`}
              >
                {item.label}
              </button>
            </li>
          ))}
        </ul>
      </nav>
      <div className="mt-auto pt-4">
        <button
          onClick={handleLogout}
          className="w-full flex items-center px-1 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-md transition"
        >
          <IconLogout2 className="h-5 w-5 mr-2" />
          Logout
        </button>
      </div>
    </div>
  );
};

export default ProfileSidebar; 