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

func mountRoutes(config APIConfig) *chi.Mux {
	r := chi.NewRouter()

	// Core middlewares
	r.Use(config.MetadataMiddleware)

	// Public + auth health check
	r.With(config.AuthMiddleware).Get("/checkAuth", config.UserHandler.GetMe)


	// Mount user and report subroutes
	r.Route("/user", func(r chi.Router) {
		r.Mount("/", UserRoutes(config.UserHandler, config.AuthMiddleware))
	})

	r.Route("/report", func(r chi.Router) {
		r.Use(config.AuthMiddleware)
		r.Mount("/", ReportRoutes(config.ReportsHandler))
	})

	// Static frontend fallback
	r.Handle("/*", spaHandler{
		staticPath: "./MedscribeUI/dist",
		indexPath:  "index.html",
	})

	return r
}

func getCORSHandler() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{
			"http://localhost:3000",
			"http://localhost:6006",
			"http://localhost:8080",
			"https://medscribe.pro",
			"https://www.medscribe.pro",
			"https://dev.medscribe.pro",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler
}

func EntryRoutes(config APIConfig) http.Handler {
	router := mountRoutes(config)
	corsMiddleware := getCORSHandler()
	return corsMiddleware(router)
}

type responseLogger struct {
	http.ResponseWriter
}

func (w *responseLogger) WriteHeader(statusCode int) {
	// Let the headers be written but intercept for inspection
	w.ResponseWriter.WriteHeader(statusCode)
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
