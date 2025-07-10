package router

import (
	"net/http"

	"github.com/1sh-repalto/url-monitoring-api/internal/handler"
	"github.com/1sh-repalto/url-monitoring-api/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRoutes(authHandler *handler.AuthHandler, urlHandler *handler.URLHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},   // optional
		AllowCredentials: true,
		MaxAge:           300,                // 5 minutes cache for pre‑flight
	}))


	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Handle("/metrics", metrics.Handler())

	r.Mount("/auth", AuthRoutes(authHandler))
	r.Mount("/urls", URLRoutes(urlHandler))

	return r
}
