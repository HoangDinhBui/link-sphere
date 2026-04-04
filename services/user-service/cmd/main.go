package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes
	r.Post("/api/v1/users/register", userHandler.Register)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Post("/api/v1/users/follow", userHandler.Follow)
		r.Post("/api/v1/users/unfollow", userHandler.Unfollow)
		r.Get("/api/v1/users/profile", userHandler.GetProfile)
		r.Get("/api/v1/users/{id}", userHandler.GetUserByID)
		r.Get("/api/v1/users/{id}/following", userHandler.GetFollowing)
	})

	port := cfg.ServerPort
	if port == "" {
		port = "8001"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Info().Str("port", port).Msg("User Service starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exiting")
}
