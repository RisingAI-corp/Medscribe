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
func NewStripeStore(apiKey string, webhookSecret string, logger *zap.Logger) Stripe {
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
	s.logger.Info("Creating new Stripe customer",
		zap.String("name", name),
		zap.String("email", email),
	)

	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}

	customer, err := customer.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe customer",
			zap.String("name", name),
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("Successfully created Stripe customer",
		zap.String("customerID", customer.ID),
		zap.String("name", name),
		zap.String("email", email),
	)

	return customer, nil
}

// CreateCheckoutSession creates a new Stripe checkout session for a customer
func (s *stripeStore) CreateCheckoutSession(customerID string) (*stripe.CheckoutSession, error) {
	s.logger.Info("Creating new checkout session",
		zap.String("customerID", customerID),
	)

	params := &stripe.CheckoutSessionParams{
		Customer:   stripe.String(customerID),
		SuccessURL: stripe.String("https://localhost:3000/Profile"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1RGQ9NFmFJe4gD9zFGwaqhMU"),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
	}

	session, err := session.New(params)
	if err != nil {
		s.logger.Error("Failed to create checkout session",
			zap.String("customerID", customerID),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("Successfully created checkout session",
		zap.String("sessionID", session.ID),
		zap.String("customerID", customerID),
		zap.String("url", session.URL),
	)

	return session, nil
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
	// TODO: Make this function safe to run multiple times,
	// even concurrently, with the same session ID

	// TODO: Make sure fulfillment hasn't already been
	// performed for this Checkout Session

	// Retrieve the Checkout Session from the API with line_items expanded
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")

	cs, err := session.Get(sessionID, params)
	if err != nil {
		s.logger.Error("Failed to retrieve checkout session", zap.Error(err))
		return err
	}

	// TODO: Change the Plan Type to Pro for the Customer in the DB
	s.logger.Info("Checkout Session",
		zap.String("customerID", cs.Customer.ID),
	)

	return nil
}
