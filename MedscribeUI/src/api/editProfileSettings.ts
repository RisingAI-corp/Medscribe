import axios from 'axios';

export interface EditProfileSettingsPayload {
  name: string;
  currentPassword: string;
  newPassword: string;
}

export class NotAuthorizedError extends Error {
  constructor() {
    super('User is not authorized');
    this.name = 'NotAuthorizedError';
  }
}

export class InvalidCurrentPasswordError extends Error {
  constructor() {
    super('Invalid current password');
    this.name = 'InvalidCurrentPasswordError';
  }
}

export class UnknownError extends Error {
  constructor(status: number) {
    super(`'Error authenticated user' + ${String(status)}`);
    this.name = 'UnknownError';
  }
}

export class ApiError extends Error {
  data: unknown;
  status: number;
  constructor(message: string, data: unknown, status: number) {
    super(message);
    this.name = 'ApiError';
    this.data = data;
    this.status = status;
  }
}

export async function editProfileSettings(payload: EditProfileSettingsPayload) {
  const baseURL = String(import.meta.env.VITE_MEDSCRIBE_BASE_URL);

  await axios.patch(`${baseURL}/user/editProfileSettings`, payload, {
    withCredentials: true,
  });
}
