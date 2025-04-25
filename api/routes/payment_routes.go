package routes

import (
	"Medscribe/api/handlers/paymentHandler"

	"github.com/go-chi/chi/v5"
)

func PaymentRoutes(handler paymentHandler.PaymentHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/new-customer", handler.HandleCreateCustomer)

	r.Post("/new-checkout-session", handler.HandleCreateCheckoutSession)

	r.Post("/webhook", handler.HandleWebhook)

	return r
}
