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
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:6006", "http://localhost:8080", "https://medscribe.pro", "https://www.medscribe.pro"}, // dev server, storybook, backend server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, 
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

	fs := http.FileServer(http.Dir("./MedscribeUI/dist"))
	r.Handle("/*", fs) // Serves all static files correctly

	// Ensure `index.html` is served for the root
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./MedscribeUI/dist/index.html")
	})

	// Fallback: Serve `index.html` for unknown frontend routes (for SPAs)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./MedscribeUI/dist/index.html")
	})

	return r
}
