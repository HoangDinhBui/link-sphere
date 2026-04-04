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
	"github.com/linksphere/post-service/internal/handler"
	"github.com/linksphere/post-service/internal/repository"
	"github.com/linksphere/post-service/internal/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()
	postRepo := repository.NewPostRepository(cfg)
	postSvc := service.NewPostService(postRepo, cfg)
	postHandler := handler.NewPostHandler(postSvc)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Post("/api/v1/posts", postHandler.Create)
		r.Post("/api/v1/posts/list", postHandler.List)
		r.Post("/api/v1/posts/detail", postHandler.GetByID)
		r.Post("/api/v1/posts/{id}/like", postHandler.Like)
		r.Post("/api/v1/posts/by-users", postHandler.GetByUserIDs)
	})

	port := cfg.ServerPort
	if port == "" {
		port = "8003"
	}

	log.Info().Str("port", port).Msg("Post Service starting")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
