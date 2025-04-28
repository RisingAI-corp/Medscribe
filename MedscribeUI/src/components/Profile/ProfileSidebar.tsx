import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { IconArrowBackUp, IconLogout2 } from '@tabler/icons-react';
import { isAuthenticatedAtom } from '../../states/userAtom';
import { useAtom } from 'jotai';
import { logout } from '../../api/logout';
import { Modal, Button } from '@mantine/core';

export type ProfileTabType = 'settings' | 'affiliates' | 'subscriptions';

interface ProfileSidebarProps {
  onNavChange: (tab: ProfileTabType) => void;
  activeTab: ProfileTabType;
}

const ProfileSidebar = ({ onNavChange, activeTab }: ProfileSidebarProps) => {
  const [, setIsAuthenticated] = useAtom(isAuthenticatedAtom);
  const navigate = useNavigate();
  const [logoutModalOpen, setLogoutModalOpen] = useState(false);

  const navItems = [
    { id: 'settings', label: 'settings' },
    { id: 'affiliates', label: 'Affiliates' },
    { id: 'subscriptions', label: 'Subscriptions' },
  ] as const;

  const logoutMutation = {
    mutate: logout,
    onerror: (error: Error) => {
      console.error('Logout failed:', error);
      // Optionally show an error message to the user
    },
  };

  const handleLogoutConfirm = () => {
    void logoutMutation.mutate();
    setIsAuthenticated(false);
    void navigate('/SignUp', { replace: true });
    setLogoutModalOpen(false);
  };

  const handleLogoutClick = () => {
    setLogoutModalOpen(true);
  };

  const handleLogoutCancel = () => {
    setLogoutModalOpen(false);
  };

  return (
    <div className="w-64 border-r h-full p-4 flex flex-col">
      <button
        onClick={() => {
          void navigate('/');
        }}
        className="flex items-center text-gray-600 hover:text-gray-900 mb-6"
      >
        <IconArrowBackUp className="h-5 w-5 mr-2" />
        Back
      </button>
      <nav>
        <ul className="space-y-2">
          {navItems.map(item => (
            <li key={item.id}>
              <button
                onClick={() => {
                  onNavChange(item.id);
                }}
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
          onClick={handleLogoutClick}
          className="w-full flex items-center px-1 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-md transition"
        >
          <IconLogout2 className="h-5 w-5 mr-2" />
          Logout
        </button>
      </div>

      {/* Logout Confirmation Modal */}
      <Modal
        opened={logoutModalOpen}
        onClose={handleLogoutCancel}
        title="Confirm Logout"
        centered
      >
        <p className="text-center">Are you sure you want to logout?</p>
        <div className="mt-4 flex justify-center space-x-2">
          <Button variant="outline" onClick={handleLogoutCancel}>
            Cancel
          </Button>
          <Button color="red" onClick={handleLogoutConfirm}>
            Logout
          </Button>
        </div>
      </Modal>
    </div>
  );
};

export default ProfileSidebar;
