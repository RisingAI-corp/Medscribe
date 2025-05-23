package user

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockUserStore is a mock implementation of the UserStore interface.
type MockUserStore struct {
	mock.Mock
}

// Put mocks the Put method.
func (m *MockUserStore) Put(ctx context.Context, name, email, password string) (string, error) {
	args := m.Called(ctx, name, email, password)
	return args.String(0), args.Error(1)
}

// Get mocks the Get method.
func (m *MockUserStore) Get(ctx context.Context, id string) (User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1) // Type assertion
}

// GetByAuth mocks the GetByAuth method.
func (m *MockUserStore) GetByAuth(ctx context.Context, email, password string) (User, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(User), args.Error(1) // Type assertion
}

// GetStyleField mocks the GetStyleField method.
func (m *MockUserStore) GetStyleField(ctx context.Context, userID, styleField string) (string, error) {
	args := m.Called(ctx, userID, styleField)
	return args.String(0), args.Error(1)
}

// UpdateStyle mocks the UpdateStyle method.
func (m *MockUserStore) UpdateStyle(ctx context.Context, providerID, contentType, newStyle string) error {
	args := m.Called(ctx, providerID, contentType, newStyle)
	return args.Error(0)
}

func (s *MockUserStore) UpdateProfileSettings(ctx context.Context, userID string, name string, currentPassword string, newPassword string) error {
	args := s.Called(ctx, userID, name, currentPassword, newPassword)
	return args.Error(0)
}

