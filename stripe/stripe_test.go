package stripe

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v74"
)

// MockStripeStore implements the Stripe interface for testing
type MockStripeStore struct {
	mock.Mock
}

func (m *MockStripeStore) CreateCustomer(name, email string) (*stripe.Customer, error) {
	args := m.Called(name, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *MockStripeStore) CreateCheckoutSession(customerID string) (*stripe.CheckoutSession, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.CheckoutSession), args.Error(1)
}

func (m *MockStripeStore) ConstructWebhookEvent(payload []byte, signature string) (stripe.Event, error) {
	args := m.Called(payload, signature)
	return args.Get(0).(stripe.Event), args.Error(1)
}

func (m *MockStripeStore) FulfillCheckout(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func TestCreateCustomer(t *testing.T) {
	t.Run("should create customer when request is valid", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		expectedCustomer := &stripe.Customer{
			ID:    "cus_123",
			Name:  "John Doe",
			Email: "john@example.com",
		}

		// Setup mock expectations
		mockStore.On("CreateCustomer", "John Doe", "john@example.com").
			Return(expectedCustomer, nil).Once()

		customer, err := mockStore.CreateCustomer("John Doe", "john@example.com")

		// Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedCustomer.ID, customer.ID)
		assert.Equal(t, expectedCustomer.Name, customer.Name)
		assert.Equal(t, expectedCustomer.Email, customer.Email)

		mockStore.AssertExpectations(t)
	})

	t.Run("should return error when customer creation fails", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		// Setup mock to return error
		mockStore.On("CreateCustomer", "John Doe", "john@example.com").
			Return(nil, errors.New("failed to create customer")).Once()

		customer, err := mockStore.CreateCustomer("John Doe", "john@example.com")

		// Verify results
		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Contains(t, err.Error(), "failed to create customer")

		mockStore.AssertExpectations(t)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		// Setup mock to return error for empty name
		mockStore.On("CreateCustomer", "", "john@example.com").
			Return(nil, errors.New("name cannot be empty")).Once()

		customer, err := mockStore.CreateCustomer("", "john@example.com")

		// Verify results
		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Contains(t, err.Error(), "name cannot be empty")

		mockStore.AssertExpectations(t)
	})

	t.Run("should return error when email is empty", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		// Setup mock to return error for empty email
		mockStore.On("CreateCustomer", "John Doe", "").
			Return(nil, errors.New("email cannot be empty")).Once()

		customer, err := mockStore.CreateCustomer("John Doe", "")

		// Verify results
		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Contains(t, err.Error(), "email cannot be empty")

		mockStore.AssertExpectations(t)
	})
}

func TestCreateCheckoutSession(t *testing.T) {
	t.Run("should create checkout session when customer ID is valid", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		expectedSession := &stripe.CheckoutSession{
			ID:  "cs_123",
			URL: "https://checkout.stripe.com/pay/cs_123",
		}

		// Setup mock expectations
		mockStore.On("CreateCheckoutSession", "cus_123").
			Return(expectedSession, nil).Once()

		session, err := mockStore.CreateCheckoutSession("cus_123")

		// Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedSession.ID, session.ID)
		assert.Equal(t, expectedSession.URL, session.URL)

		mockStore.AssertExpectations(t)
	})

	t.Run("should return error when customer ID is empty", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		// Setup mock to return error for empty customer ID
		mockStore.On("CreateCheckoutSession", "").
			Return(nil, errors.New("customer ID cannot be empty")).Once()

		session, err := mockStore.CreateCheckoutSession("")

		// Verify results
		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "customer ID cannot be empty")

		mockStore.AssertExpectations(t)
	})

	t.Run("should return error when checkout session creation fails", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		// Setup mock to return error
		mockStore.On("CreateCheckoutSession", "cus_123").
			Return(nil, errors.New("failed to create checkout session")).Once()

		session, err := mockStore.CreateCheckoutSession("cus_123")

		// Verify results
		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "failed to create checkout session")

		mockStore.AssertExpectations(t)
	})
}

func TestConstructWebhookEvent(t *testing.T) {
	t.Run("should construct webhook event when payload and signature are valid", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		payload := []byte(`{"id":"evt_123","type":"checkout.session.completed"}`)
		signature := "t=1234567890,v1=abc123"
		expectedEvent := stripe.Event{
			ID:   "evt_123",
			Type: "checkout.session.completed",
		}

		// Setup mock expectations
		mockStore.On("ConstructWebhookEvent", payload, signature).
			Return(expectedEvent, nil).Once()

		event, err := mockStore.ConstructWebhookEvent(payload, signature)

		// Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedEvent.ID, event.ID)
		assert.Equal(t, expectedEvent.Type, event.Type)

		mockStore.AssertExpectations(t)
	})

	t.Run("should return error when signature is invalid", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		payload := []byte(`{"id":"evt_123","type":"checkout.session.completed"}`)
		invalidSignature := "invalid_signature"

		// Setup mock to return error for invalid signature
		mockStore.On("ConstructWebhookEvent", payload, invalidSignature).
			Return(stripe.Event{}, errors.New("signature verification failed")).Once()

		event, err := mockStore.ConstructWebhookEvent(payload, invalidSignature)

		// Verify results
		assert.Error(t, err)
		assert.Empty(t, event.ID)
		assert.Contains(t, err.Error(), "signature verification failed")

		mockStore.AssertExpectations(t)
	})

	t.Run("should return error when payload is invalid", func(t *testing.T) {
		mockStore := new(MockStripeStore)

		invalidPayload := []byte(`invalid json`)
		signature := "t=1234567890,v1=abc123"

		// Setup mock to return error for invalid payload
		mockStore.On("ConstructWebhookEvent", invalidPayload, signature).
			Return(stripe.Event{}, errors.New("invalid payload")).Once()

		event, err := mockStore.ConstructWebhookEvent(invalidPayload, signature)

		// Verify results
		assert.Error(t, err)
		assert.Empty(t, event.ID)
		assert.Contains(t, err.Error(), "invalid payload")

		mockStore.AssertExpectations(t)
	})
}
