package routes

import (
	userhandler "Medscribe/api/handlers/userHandler"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(handler userhandler.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/signup", handler.SignUp)

	r.Post("/login", handler.Login)

	return r
}
