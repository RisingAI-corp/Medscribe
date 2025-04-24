import { useState } from 'react';
import ProfileSettings from '../../components/Profile/profileSettings';
import ProfileAffiliates from '../../components/Profile/profileAffiliates';
import ProfileSubscriptions from '../../components/Profile/profileSubscriptions';
import ProfileSidebar from '../../components/Profile/ProfileSidebar';
import { ProfileTabType } from '../../components/Profile/ProfileSidebar';
import { Header } from '../../components/Header/header';

const ProfileScreen = () => {
  const [activeTab, setActiveTab] = useState<ProfileTabType>('settings');

  const sidebarContentDisplay = () => {
    switch (activeTab) {
      case 'settings':
        return <ProfileSettings />;
      case 'affiliates':
        return <ProfileAffiliates />;
      case 'subscriptions':
        return <ProfileSubscriptions />;
      default:
        return <ProfileSettings />;
    }
  };

  return (
    <div className="flex flex-col h-screen w-full">
      <div>
        <Header />
      </div>
      <div className="flex flex-row flex-1">
        <ProfileSidebar onNavChange={setActiveTab} activeTab={activeTab} />
        <div className="flex-1 p-6 overflow-auto">
          {sidebarContentDisplay()}
        </div>
      </div>
    </div>
  );
};

export default ProfileScreen;
