package paymentHandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/webhook"

	"go.uber.org/zap"
)

// PaymentHandler defines the interface for payment processing handlers
type PaymentHandler interface {
	HandleCreateCustomer(w http.ResponseWriter, r *http.Request)
	HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request)
	HandleWebhook(w http.ResponseWriter, r *http.Request)
}

// CustomerRequest represents the incoming request to create a customer
type CustomerRequest struct {
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
	stripeAPIKey        string
	stripeWebhookSecret string
	logger              *zap.Logger
}

func NewPaymentHandler(stripeAPIKey string, stripeWebhookSecret string, logger *zap.Logger) PaymentHandler {
	// Set the global Stripe API key
	stripe.Key = stripeAPIKey

	return &paymentHandler{
		stripeAPIKey:        stripeAPIKey,
		stripeWebhookSecret: stripeWebhookSecret,
		logger:              logger,
	}
}

// HandleCreateCustomer handles HTTP requests to create a Stripe customer
func (h *paymentHandler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call Stripe API to create customer
	customer, err := h.createCustomer(req.Name, req.Email)
	if err != nil {
		http.Error(w, "Failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	resp := CustomerResponse{
		CustomerID: customer.ID,
		Name:       customer.Name,
		Email:      customer.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// createCustomer creates a new Stripe customer with the given name and email
func (h *paymentHandler) createCustomer(name string, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}

	return customer.New(params)
}

// HandleCreateCheckoutSession handles HTTP requests to create a Stripe checkout session
func (h *paymentHandler) HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call Stripe API to create checkout session
	sess, err := h.createCheckoutSession(req.CustomerID)
	if err != nil {
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

// createCheckoutSession creates a new Stripe checkout session for a customer
func (h *paymentHandler) createCheckoutSession(customerID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer:   stripe.String(customerID),
		SuccessURL: stripe.String("https://example.com/success"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1RHckAPMOCbWba8pYYUTOdgG"),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
	}

	return session.New(params)
}

func (h *paymentHandler) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	endpointSecret := h.stripeWebhookSecret
	event, err := webhook.ConstructEventWithOptions(body, req.Header.Get("Stripe-Signature"), endpointSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	if event.Type == "checkout.session.completed" ||
		event.Type == "checkout.session.async_payment_succeeded" {
		var cs stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.FulfillCheckout(cs.ID)

	}

	w.WriteHeader(http.StatusOK)
}

// FulfillCheckout fulfills a Stripe checkout session
func (h *paymentHandler) FulfillCheckout(sessionId string) {
	fmt.Println("Fulfilling Checkout Session " + sessionId)

	// TODO: Make this function safe to run multiple times,
	// even concurrently, with the same session ID

	// TODO: Make sure fulfillment hasn't already been
	// performed for this Checkout Session

	// Retrieve the Checkout Session from the API with line_items expanded
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")

	cs, _ := session.Get(sessionId, params)

	// Check the Checkout Session's payment_status property
	// to determine if fulfillment should be performed
	if cs.PaymentStatus != stripe.CheckoutSessionPaymentStatusUnpaid {
		// TODO: Perform fulfillment of the line items
		fmt.Println("Fulfilling line items for Checkout Session " + sessionId)
		fmt.Println(cs.Customer.ID)
		// TODO: Record/save fulfillment status for this
		// Checkout Session
	}
}
