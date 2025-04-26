package stripe

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/webhook"

	"go.uber.org/zap"
)

// Stripe defines the interface for Stripe payment processing operations
type Stripe interface {
	CreateCustomer(name, email string) (*stripe.Customer, error)
	CreateCheckoutSession(customerID string) (*stripe.CheckoutSession, error)
	ConstructWebhookEvent(payload []byte, signature string) (stripe.Event, error)
	FulfillCheckout(sessionID string) error
}

// stripeStore handles Stripe API interactions
type stripeStore struct {
	apiKey        string
	webhookSecret string
	logger        *zap.Logger
}

// NewStripeStore creates a new Stripe client with the provided credentials
func NewStripeStore(apiKey, webhookSecret string, logger *zap.Logger) Stripe {
	// Set the global Stripe API key
	stripe.Key = apiKey

	return &stripeStore{
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
		logger:        logger,
	}
}

// CreateCustomer creates a new Stripe customer with the given name and email
func (s *stripeStore) CreateCustomer(name, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}

	return customer.New(params)
}

// CreateCheckoutSession creates a new Stripe checkout session for a customer
func (s *stripeStore) CreateCheckoutSession(customerID string) (*stripe.CheckoutSession, error) {
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

// ConstructWebhookEvent validates and constructs a webhook event from a request payload
func (s *stripeStore) ConstructWebhookEvent(payload []byte, signature string) (stripe.Event, error) {
	return webhook.ConstructEventWithOptions(
		payload,
		signature,
		s.webhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
}

// FulfillCheckout fulfills a Stripe checkout session
func (s *stripeStore) FulfillCheckout(sessionID string) error {
	s.logger.Info("Fulfilling Checkout Session", zap.String("sessionID", sessionID))

	// Retrieve the Checkout Session from the API with line_items expanded
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")

	cs, err := session.Get(sessionID, params)
	if err != nil {
		s.logger.Error("Failed to retrieve checkout session", zap.String("sessionID", sessionID), zap.Error(err))
		return err
	}

	// Check the Checkout Session's payment_status property
	// to determine if fulfillment should be performed
	if cs.PaymentStatus != stripe.CheckoutSessionPaymentStatusUnpaid {
		s.logger.Info("Processing paid checkout session",
			zap.String("sessionID", sessionID),
			zap.String("customerID", cs.Customer.ID),
		)
		// TODO: Perform fulfillment of the line items
		// TODO: Record/save fulfillment status for this Checkout Session
	}

	return nil
}
