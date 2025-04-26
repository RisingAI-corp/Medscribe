package paymentHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/stripe/stripe-go/v74"

	"go.uber.org/zap"
)

const (
	checkoutSessionCompleted             = "checkout.session.completed"
	checkoutSessionAsyncPaymentSucceeded = "checkout.session.async_payment_succeeded"
)

// PaymentHandler defines the interface for payment processing handlers
type PaymentHandler interface {
	CreateCustomer(w http.ResponseWriter, r *http.Request)
	CreateCheckoutSession(w http.ResponseWriter, r *http.Request)
	Webhook(w http.ResponseWriter, r *http.Request)
}

// CustomerRequest represents the incoming request to create a customer
type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CustomerResponse represents the response with customer info
type CustomerResponse struct {
	CustomerID string `json:"customer_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

// CheckoutSessionRequest represents the incoming request to create a checkout session
type CheckoutSessionRequest struct {
	CustomerID string `json:"customer_id"`
}

// CheckoutSessionResponse represents the response with session info
type CheckoutSessionResponse struct {
	SessionID  string `json:"session_id"`
	SessionURL string `json:"session_url"`
}

type paymentHandler struct {
	stripeClient Stripe
	logger       *zap.Logger
}

// Stripe defines the interface for stripe operations used by the handler
type Stripe interface {
	CreateCustomer(name, email string) (*stripe.Customer, error)
	CreateCheckoutSession(customerID string) (*stripe.CheckoutSession, error)
	ConstructWebhookEvent(payload []byte, signature string) (stripe.Event, error)
	FulfillCheckout(sessionID string) error
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(stripeClient Stripe, logger *zap.Logger) PaymentHandler {
	return &paymentHandler{
		stripeClient: stripeClient,
		logger:       logger,
	}
}

// CreateCustomer handles HTTP requests to create a Stripe customer
func (h *paymentHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call Stripe API to create customer
	customer, err := h.stripeClient.CreateCustomer(req.Name, req.Email)
	if err != nil {
		http.Error(w, "Failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		CustomerResponse{
			CustomerID: customer.ID,
			Name:       customer.Name,
			Email:      customer.Email,
		},
	)
}

// CreateCheckoutSession handles HTTP requests to create a Stripe checkout session
func (h *paymentHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call Stripe API to create checkout session
	sess, err := h.stripeClient.CreateCheckoutSession(req.CustomerID)
	if err != nil {
		h.logger.Error("Failed to create checkout session", zap.Error(err))
		http.Error(w, "Failed to create checkout session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	resp := CheckoutSessionResponse{
		SessionID:  sess.ID,
		SessionURL: sess.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *paymentHandler) Webhook(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, 65536)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("Error reading request body", zap.Error(err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent
	event, err := h.stripeClient.ConstructWebhookEvent(body, req.Header.Get("Stripe-Signature"))
	if err != nil {
		h.logger.Error("Error verifying webhook signature", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	if event.Type == checkoutSessionCompleted ||
		event.Type == checkoutSessionAsyncPaymentSucceeded {
		var cs stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			h.logger.Error("Error parsing webhook JSON", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h.stripeClient.FulfillCheckout(cs.ID); err != nil {
			h.logger.Error("Error fulfilling checkout", zap.String("sessionID", cs.ID), zap.Error(err))
			// We still return 200 to Stripe to avoid repeated webhook calls
		}
	}

	w.WriteHeader(http.StatusOK)
}
