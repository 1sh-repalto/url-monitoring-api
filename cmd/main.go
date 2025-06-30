package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/1sh-repalto/url-monitoring-api/internal/handler"
	"github.com/1sh-repalto/url-monitoring-api/internal/repository"
	"github.com/1sh-repalto/url-monitoring-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main(){
	_ = godotenv.Load()

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	defer dbpool.Close()

	urlRepo := repository.NewPgxURLRepository(dbpool)

	urlService := service.NewURLService(urlRepo)

	urlHandler := handler.NewURLHandler(urlService)
	
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// routes
	r.Post("/urls", urlHandler.Register)
	r.Get("/urls", urlHandler.GetURL)
	r.Delete("/urls/{urlID}", urlHandler.DeleteURL)

	log.Println("Server started on :3000")
	http.ListenAndServe(":3000", r)
}