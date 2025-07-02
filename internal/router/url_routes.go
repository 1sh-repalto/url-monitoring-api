package router

import (
	"github.com/1sh-repalto/url-monitoring-api/internal/handler"
	appmiddleware "github.com/1sh-repalto/url-monitoring-api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func URLRoutes(urlHandler *handler.URLHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(appmiddleware.JWTMiddleware)

	r.Post("/", urlHandler.Register)
	r.Get("/", urlHandler.GetURL)
	r.Delete("/{urlID}", urlHandler.DeleteURL)

	return r
}
