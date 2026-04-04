package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/linksphere/comment-service/internal/handler"
	"github.com/linksphere/comment-service/internal/repository"
	"github.com/linksphere/comment-service/internal/service"
	"github.com/linksphere/pkg/config"
	"github.com/linksphere/pkg/middleware"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()
	commentRepo := repository.NewCommentRepository(cfg)
	commentSvc := service.NewCommentService(commentRepo, cfg)
	commentHandler := handler.NewCommentHandler(commentSvc)

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
		r.Post("/api/v1/posts/comment", commentHandler.Create)
		r.Post("/api/v1/posts/comments", commentHandler.List)
	})

	port := cfg.ServerPort
	if port == "" {
		port = "8004"
	}

	log.Info().Str("port", port).Msg("Comment Service starting")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
