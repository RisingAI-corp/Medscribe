package routes

import (
	"Medscribe/api/handlers/paymentHandler"

	"github.com/go-chi/chi/v5"
)

func PaymentRoutes(handler paymentHandler.PaymentHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/new-customer", handler.CreateCustomer)

	r.Post("/new-checkout-session", handler.CreateCheckoutSession)

	r.Post("/webhook", handler.Webhook)

	return r
}
