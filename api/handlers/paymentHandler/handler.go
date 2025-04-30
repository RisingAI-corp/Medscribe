package paymentHandler

import (
	"Medscribe/stripe"
	"encoding/json"
	"io/ioutil"
	"net/http"

	stripeGo "github.com/stripe/stripe-go/v74"

	"go.uber.org/zap"
)

const (
	checkoutSessionCompleted = "checkout.session.completed"
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
	stripeClient stripe.Stripe
	logger       *zap.Logger
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(stripeClient stripe.Stripe, logger *zap.Logger) PaymentHandler {
	return &paymentHandler{
		stripeClient: stripeClient,
		logger:       logger,
	}
}

// CreateCustomer handles HTTP requests to create a Stripe customer
func (h *paymentHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to create customer")

	// Parse request body
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Creating customer with provided details",
		zap.String("name", req.Name),
		zap.String("email", req.Email),
	)

	// Call Stripe API to create customer
	customer, err := h.stripeClient.CreateCustomer(req.Name, req.Email)
	if err != nil {
		h.logger.Error("Failed to create customer via Stripe",
			zap.String("name", req.Name),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		http.Error(w, "Failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully created customer",
		zap.String("customerID", customer.ID),
		zap.String("name", customer.Name),
		zap.String("email", customer.Email),
	)

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
		h.logger.Error("Failed to decode request body", zap.Error(err))
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
	h.logger.Info("Received Stripe webhook request")

	req.Body = http.MaxBytesReader(w, req.Body, 65536)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("Error reading request body", zap.Error(err))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	h.logger.Debug("Received webhook payload",
		zap.String("signature", req.Header.Get("Stripe-Signature")),
		zap.String("body", string(body)),
	)

	// Pass the request body and Stripe-Signature header to ConstructEvent
	event, err := h.stripeClient.ConstructWebhookEvent(body, req.Header.Get("Stripe-Signature"))
	if err != nil {
		h.logger.Error("Error verifying webhook signature",
			zap.Error(err),
			zap.String("signature", req.Header.Get("Stripe-Signature")),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Info("Processing Stripe webhook event",
		zap.String("eventType", event.Type),
		zap.String("eventID", event.ID),
	)

	if event.Type == checkoutSessionCompleted {
		var cs *stripeGo.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			h.logger.Error("Error parsing webhook JSON",
				zap.String("eventType", event.Type),
				zap.String("eventID", event.ID),
				zap.Error(err),
				zap.String("rawData", string(event.Data.Raw)),
			)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.logger.Info("Processing checkout session",
			zap.String("sessionID", cs.ID),
			zap.String("customerID", cs.Customer.ID),
			zap.String("paymentStatus", string(cs.Status)),
		)

		if err := h.stripeClient.FulfillCheckout(cs.ID); err != nil {
			h.logger.Error("Error fulfilling checkout",
				zap.String("sessionID", cs.ID),
				zap.String("customerID", cs.Customer.ID),
				zap.Error(err),
			)
			// We still return 200 to Stripe to avoid repeated webhook calls
		} else {
			h.logger.Info("Successfully fulfilled checkout",
				zap.String("sessionID", cs.ID),
				zap.String("customerID", cs.Customer.ID),
			)
		}
	}

	w.WriteHeader(http.StatusOK)
}
