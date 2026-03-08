package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/linksphere/pkg/config"
	"github.com/linksphere/pkg/middleware"
	"github.com/linksphere/user-service/internal/handler"
	"github.com/linksphere/user-service/internal/repository"
	"github.com/linksphere/user-service/internal/service"
)

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load config
	cfg := config.Load()

	// Initialize repository
	userRepo := repository.NewUserRepository(cfg)

	// Initialize service
	userSvc := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handler.NewUserHandler(userSvc)

	// Setup router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Public routes
	r.Post("/api/v1/users/register", userHandler.Register)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Post("/api/v1/users/follow", userHandler.Follow)
		r.Post("/api/v1/users/unfollow", userHandler.Unfollow)
		r.Get("/api/v1/users/profile", userHandler.GetProfile)
		r.Get("/api/v1/users/{id}", userHandler.GetUserByID)
	})

	port := cfg.ServerPort
	if port == "" {
		port = "8001"
	}

	log.Info().Str("port", port).Msg("User Service starting")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
