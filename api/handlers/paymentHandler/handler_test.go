package paymentHandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	stripeGo "github.com/stripe/stripe-go/v74"
	"go.uber.org/zap"
)

// MockStripeClient implements the Stripe interface for testing
type MockStripeClient struct {
	mock.Mock
}

func (m *MockStripeClient) CreateCustomer(name, email string) (*stripeGo.Customer, error) {
	args := m.Called(name, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripeGo.Customer), args.Error(1)
}

func (m *MockStripeClient) CreateCheckoutSession(customerID string) (*stripeGo.CheckoutSession, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripeGo.CheckoutSession), args.Error(1)
}

func (m *MockStripeClient) ConstructWebhookEvent(payload []byte, signature string) (stripeGo.Event, error) {
	args := m.Called(payload, signature)
	return args.Get(0).(stripeGo.Event), args.Error(1)
}

func (m *MockStripeClient) FulfillCheckout(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func TestCreateCustomer(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockStripe := new(MockStripeClient)
	handler := NewPaymentHandler(mockStripe, logger)

	t.Run("should create customer successfully", func(t *testing.T) {
		// Setup test data
		reqBody := CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		jsonData, _ := json.Marshal(reqBody)

		// Setup mock expectations
		expectedCustomer := &stripeGo.Customer{
			ID:    "cus_123",
			Name:  "John Doe",
			Email: "john@example.com",
		}
		mockStripe.On("CreateCustomer", "John Doe", "john@example.com").
			Return(expectedCustomer, nil).Once()

		// Create request and response recorder
		req := httptest.NewRequest("POST", "/customers", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateCustomer(w, req)

		// Verify response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response CustomerResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "cus_123", response.CustomerID)
		assert.Equal(t, "John Doe", response.Name)
		assert.Equal(t, "john@example.com", response.Email)

		mockStripe.AssertExpectations(t)
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		// Create request with invalid JSON
		req := httptest.NewRequest("POST", "/customers", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateCustomer(w, req)

		// Verify response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when Stripe fails", func(t *testing.T) {
		// Setup test data
		reqBody := CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		jsonData, _ := json.Marshal(reqBody)

		// Setup mock to return error
		mockStripe.On("CreateCustomer", "John Doe", "john@example.com").
			Return(nil, assert.AnError).Once()

		// Create request and response recorder
		req := httptest.NewRequest("POST", "/customers", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateCustomer(w, req)

		// Verify response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockStripe.AssertExpectations(t)
	})
}

func TestCreateCheckoutSession(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockStripe := new(MockStripeClient)
	handler := NewPaymentHandler(mockStripe, logger)

	t.Run("should create checkout session successfully", func(t *testing.T) {
		// Setup test data
		reqBody := CheckoutSessionRequest{
			CustomerID: "cus_123",
		}
		jsonData, _ := json.Marshal(reqBody)

		// Setup mock expectations
		expectedSession := &stripeGo.CheckoutSession{
			ID:  "cs_123",
			URL: "https://checkout.stripe.com/pay/cs_123",
		}
		mockStripe.On("CreateCheckoutSession", "cus_123").
			Return(expectedSession, nil).Once()

		// Create request and response recorder
		req := httptest.NewRequest("POST", "/checkout-sessions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateCheckoutSession(w, req)

		// Verify response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response CheckoutSessionResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "cs_123", response.SessionID)
		assert.Equal(t, "https://checkout.stripe.com/pay/cs_123", response.SessionURL)

		mockStripe.AssertExpectations(t)
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		// Create request with invalid JSON
		req := httptest.NewRequest("POST", "/checkout-sessions", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateCheckoutSession(w, req)

		// Verify response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when Stripe fails", func(t *testing.T) {
		// Setup test data
		reqBody := CheckoutSessionRequest{
			CustomerID: "cus_123",
		}
		jsonData, _ := json.Marshal(reqBody)

		// Setup mock to return error
		mockStripe.On("CreateCheckoutSession", "cus_123").
			Return(nil, assert.AnError).Once()

		// Create request and response recorder
		req := httptest.NewRequest("POST", "/checkout-sessions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateCheckoutSession(w, req)

		// Verify response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockStripe.AssertExpectations(t)
	})
}

func TestWebhook(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockStripe := new(MockStripeClient)
	handler := NewPaymentHandler(mockStripe, logger)

	// Test cases
	tests := []struct {
		name           string
		eventType      string
		payload        []byte
		signature      string
		expectedStatus int
		setupMock      func()
	}{
		{
			name:      "Valid checkout.session.completed event",
			eventType: checkoutSessionCompleted,
			payload: []byte(`{
				"type": "checkout.session.completed",
				"data": {
					"object": {
						"id": "cs_test_123",
						"customer": {
							"id": "cus_123"
						},
						"status": "complete"
					}
				}
			}`),
			signature:      "valid_signature",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				// Mock ConstructWebhookEvent
				mockStripe.On("ConstructWebhookEvent", mock.Anything, "valid_signature").
					Return(stripeGo.Event{
						Type: checkoutSessionCompleted,
						Data: &stripeGo.EventData{
							Raw: []byte(`{
								"id": "cs_test_123",
								"customer": {
									"id": "cus_123"
								},
								"status": "complete"
							}`),
						},
					}, nil)

				// Mock FulfillCheckout
				mockStripe.On("FulfillCheckout", "cs_test_123").
					Return(nil)
			},
		},
		{
			name:           "Invalid signature",
			eventType:      checkoutSessionCompleted,
			payload:        []byte(`{}`),
			signature:      "invalid_signature",
			expectedStatus: http.StatusBadRequest,
			setupMock: func() {
				mockStripe.On("ConstructWebhookEvent", mock.Anything, "invalid_signature").
					Return(stripeGo.Event{}, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock()

			// Create request
			req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(tt.payload))
			req.Header.Set("Stripe-Signature", tt.signature)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Webhook(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify mock expectations
			mockStripe.AssertExpectations(t)
		})
	}
}
