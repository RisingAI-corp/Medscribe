package routes

import (
	"Medscribe/api/handlers/reportsHandler"
	userhandler "Medscribe/api/handlers/userHandler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

// APIConfig is a struct that contains all the handlers, middleware and loggers for the API
type APIConfig struct {
	UserHandler      userhandler.UserHandler
	ReportsHandler   reportsHandler.ReportsHandler
	LoggerMiddleware func(http.Handler) http.Handler
	AuthMiddleware   func(http.Handler) http.Handler
	Logger           *zap.Logger
}

func EntryRoutes(config APIConfig) *chi.Mux {
	userSubRoutes := UserRoutes(config.UserHandler)
	reportsSubRoutes := ReportRoutes(config.ReportsHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with your frontend's URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Allows cookies & authentication headers
	})

	r := chi.NewRouter()

	r.Use(config.LoggerMiddleware)
	r.Use(corsHandler.Handler)

	r.With(config.AuthMiddleware).Get("/checkAuth", config.UserHandler.GetMe)

	r.Route("/user", func(r chi.Router) {
		r.Mount("/", userSubRoutes)
	})

	r.Route("/report", func(r chi.Router) {
		r.Use(config.AuthMiddleware)
		r.Mount("/", reportsSubRoutes)
	})

	return r
}
