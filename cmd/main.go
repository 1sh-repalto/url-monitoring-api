package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/1sh-repalto/url-monitoring-api/internal/engine"
	"github.com/1sh-repalto/url-monitoring-api/internal/handler"
	"github.com/1sh-repalto/url-monitoring-api/internal/metrics"
	"github.com/1sh-repalto/url-monitoring-api/internal/repository"
	"github.com/1sh-repalto/url-monitoring-api/internal/router"
	"github.com/1sh-repalto/url-monitoring-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env file found")
		}
	}

	// initialize prometheus metric collection
	metrics.Init()

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	defer dbpool.Close()
	fmt.Println("Connected to:", os.Getenv("DB_URL"))


	// Initialize repositories
	userRepo := repository.NewUserRepository(dbpool)
	urlRepo := repository.NewPgxURLRepository(dbpool)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	urlService := service.NewURLService(urlRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	urlHandler := handler.NewURLHandler(urlService)

	// start monitoring engine
	monitorEngine := engine.NewMonitorEngine(urlService)
	go monitorEngine.Start()

	// Setup router
	r := router.SetupRoutes(authHandler, urlHandler)

	log.Println("Server started on :3000")
	http.ListenAndServe(":3000", r)
}
