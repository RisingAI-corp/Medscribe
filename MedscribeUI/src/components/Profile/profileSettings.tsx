import { useState } from 'react';
import { TextInput, PasswordInput, Button } from '@mantine/core';
import { useMutation } from '@tanstack/react-query';
import { editProfileSettings } from '../../api/editProfileSettings';
import { UpdateUserName } from './derviedAtoms';
import { useAtom } from 'jotai';
import { userAtom } from '../../states/userAtom';
import { AxiosError } from 'axios';

export default function Settings() {
  const [, updateUserName] = useAtom(UpdateUserName);
  const [user] = useAtom(userAtom);
  const [newName, setNewName] = useState(user.name);
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [successSettingsUpdate, setSuccessSettingsUpdate] = useState(false);
  const [invalidCurrentPassword, setInvalidCurrentPassword] = useState(false);
  const [formSubmitted, setFormSubmitted] = useState(false);

  const editProfileSettingsMutation = useMutation({
    mutationFn: editProfileSettings,
    onSuccess: () => {
      updateUserName(newName);
      setSuccessSettingsUpdate(true);
      setInvalidCurrentPassword(false);
      setFormSubmitted(false);
    },
    onError: error => {
      setFormSubmitted(true);
      if (error instanceof AxiosError && error.response?.status === 412) {
        console.log('boom boom');
        setInvalidCurrentPassword(true);
      }
    },
  });

  const validateSettingsForm = () => {
    return (
      newName.trim() &&
      currentPassword.trim() &&
      newPassword.trim() &&
      !(currentPassword && newPassword && currentPassword === newPassword) &&
      !(newPassword && confirmPassword && newPassword !== confirmPassword)
    );
  };

  const submitSettingsForm = () => {
    setFormSubmitted(true);
    if (!validateSettingsForm()) {
      return;
    }
    editProfileSettingsMutation.mutate({
      name: newName,
      currentPassword,
      newPassword,
    });
  };

  return (
    <div className="p-8 bg-white shadow-md rounded-lg">
      <h2 className="text-2xl font-semibold mb-6">Edit Profile</h2>

      <div className="space-y-6">
        <TextInput
          label="Name"
          placeholder="Your Name"
          value={newName}
          onChange={e => {
            setNewName(e.currentTarget.value);
          }}
          required
          error={formSubmitted && !newName.trim() ? 'Name is required' : ''}
        />

        <TextInput
          label="Email"
          placeholder="you@example.com"
          type="email"
          value={user.email}
          disabled
        />

        <div className="pt-6 border-t">
          <PasswordInput
            label="Current Password"
            placeholder="Current Password"
            value={currentPassword}
            onChange={e => {
              setCurrentPassword(e.currentTarget.value);
            }}
            required
            error={
              (formSubmitted && !currentPassword.trim()
                ? 'Current password is required'
                : '') ||
              (invalidCurrentPassword ? 'Current password is not a match' : '')
            }
          />

          <PasswordInput
            label="New Password"
            placeholder="New Password"
            value={newPassword}
            onChange={e => {
              setNewPassword(e.currentTarget.value);
            }}
            className="mt-8"
            required
            error={
              formSubmitted && !newPassword.trim()
                ? 'New password is required'
                : formSubmitted && newPassword === currentPassword
                  ? 'Current and new passwords cannot be the same'
                  : ''
            }
          />

          <PasswordInput
            label="Confirm New Password"
            placeholder="Confirm New Password"
            value={confirmPassword}
            onChange={e => {
              setConfirmPassword(e.currentTarget.value);
            }}
            className="mt-4"
            required
            error={
              formSubmitted && !confirmPassword.trim()
                ? 'Confirm password is required'
                : formSubmitted && newPassword !== confirmPassword
                  ? 'New password and confirm password do not match'
                  : ''
            }
          />
        </div>

        <Button
          onClick={submitSettingsForm}
          className="bg-blue-500 hover:bg-blue-600"
        >
          Submit
        </Button>

        {successSettingsUpdate && editProfileSettingsMutation.isSuccess && (
          <div className="mt-4 text-green-600">new changes saved!</div>
        )}
      </div>
    </div>
  );
}
