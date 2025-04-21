package routes

import (
	"Medscribe/api/handlers/reportsHandler"
	userhandler "Medscribe/api/handlers/userHandler"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

// APIConfig is a struct that contains all the handlers, middleware and loggers for the API
type APIConfig struct {
	UserHandler        userhandler.UserHandler
	ReportsHandler     reportsHandler.ReportsHandler
	AuthMiddleware     func(http.Handler) http.Handler
	MetadataMiddleware func(http.Handler) http.Handler
}

func EntryRoutes(config APIConfig) *chi.Mux {
	userSubRoutes := UserRoutes(config.UserHandler)
	reportsSubRoutes := ReportRoutes(config.ReportsHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:6006", "http://localhost:8080", "https://medscribe.pro", "https://www.medscribe.pro","https://medscribe-dev-402133475168.us-central1.run.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r := chi.NewRouter()

	r.Use(config.MetadataMiddleware)
	r.Use(corsHandler.Handler)

	r.With(config.AuthMiddleware).Get("/checkAuth", config.UserHandler.GetMe)

	r.Route("/user", func(r chi.Router) {
		r.Mount("/", userSubRoutes)
	})

	r.Route("/report", func(r chi.Router) {
		r.Use(config.AuthMiddleware)
		r.Mount("/", reportsSubRoutes)
	})

	// Serve SPA frontend with static + fallback logic
	r.Handle("/*", spaHandler{
		staticPath: "./MedscribeUI/dist",
		indexPath:  "index.html",
	})

	return r
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := filepath.Join(h.staticPath, r.URL.Path)
	_, err := os.Stat(requestPath)

	if os.IsNotExist(err) {
		// Not a file? Serve index.html so React Router can handle it
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// Something went wrong (permissions, disk error, etc.)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// File exists, serve it
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
