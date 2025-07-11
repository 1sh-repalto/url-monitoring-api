package router

import (
	"github.com/1sh-repalto/url-monitoring-api/internal/handler"
	appmiddleware "github.com/1sh-repalto/url-monitoring-api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes(authHandler *handler.AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", authHandler.Signup)
	r.Post("/login", authHandler.Login)
	r.Get("/refresh", authHandler.Refresh)
	r.Post("/logout", authHandler.Logout)
	r.With(appmiddleware.JWTMiddleware).Get("/me", authHandler.Me)

	return r
}
