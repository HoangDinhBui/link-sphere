package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/linksphere/feed-service/internal/handler"
	"github.com/linksphere/feed-service/internal/service"
	"github.com/linksphere/pkg/config"
	"github.com/linksphere/pkg/middleware"
	pkgredis "github.com/linksphere/pkg/redis"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()

	// Initialize Redis
	redisClient := pkgredis.NewClient(cfg.RedisAddr(), cfg.RedisPassword)

	feedSvc := service.NewFeedService(cfg, redisClient)
	feedHandler := handler.NewFeedHandler(feedSvc)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Post("/api/v1/feed/get", feedHandler.GetFeed)
	})

	port := cfg.ServerPort
	if port == "" {
		port = "8005"
	}

	log.Info().Str("port", port).Msg("Feed Service starting")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
