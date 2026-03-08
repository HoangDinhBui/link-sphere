package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/linksphere/comment-service/internal/model"
	"github.com/linksphere/pkg/config"
)

// CommentRepository handles database operations for comments.
type CommentRepository struct {
	db *sqlx.DB
}

// NewCommentRepository creates a new CommentRepository.
func NewCommentRepository(cfg *config.Config) *CommentRepository {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to database - running in degraded mode")
		return &CommentRepository{}
	}
	log.Info().Msg("Comment Service connected to PostgreSQL")
	return &CommentRepository{db: db}
}

// Create inserts a new comment.
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) 
	           VALUES ($1, $2, $3) 
	           RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		comment.PostID, comment.UserID, comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt)
}

// ListByPostID returns comments for a post with pagination.
func (r *CommentRepository) ListByPostID(ctx context.Context, postID string, limit, offset int) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.SelectContext(ctx, &comments,
		"SELECT * FROM comments WHERE post_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3",
		postID, limit, offset,
	)
	return comments, err
}
