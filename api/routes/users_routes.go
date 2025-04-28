package routes

import (
	userhandler "Medscribe/api/handlers/userHandler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(handler userhandler.UserHandler, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/signup", handler.SignUp)

	r.Post("/login", handler.Login)

	r.With(authMiddleware).Patch("/editProfileSettings", handler.UpdateProfileSettings)

	r.With(authMiddleware).Post("/logout", handler.Logout)

	return r
}
