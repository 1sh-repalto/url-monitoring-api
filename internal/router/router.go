package router

import (
	"net/http"

	"github.com/1sh-repalto/url-monitoring-api/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(authHandler *handler.AuthHandler, urlHandler *handler.URLHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Mount("/auth", AuthRoutes(authHandler))
	r.Mount("/urls", URLRoutes(urlHandler))

	return r
}
