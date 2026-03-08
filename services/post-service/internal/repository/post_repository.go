package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/linksphere/pkg/config"
	"github.com/linksphere/post-service/internal/model"
)

// PostRepository handles database operations for posts.
type PostRepository struct {
	db *sqlx.DB
}

// NewPostRepository creates a new PostRepository.
func NewPostRepository(cfg *config.Config) *PostRepository {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to database - running in degraded mode")
		return &PostRepository{}
	}
	log.Info().Msg("Post Service connected to PostgreSQL")
	return &PostRepository{db: db}
}

// Create inserts a new post.
func (r *PostRepository) Create(ctx context.Context, post *model.Post) error {
	query := `INSERT INTO posts (user_id, content, images, hashtags) 
	           VALUES ($1, $2, $3, $4) 
	           RETURNING id, like_count, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		post.UserID, post.Content,
		pq.Array(post.Images), pq.Array(post.Hashtags),
	).Scan(&post.ID, &post.LikeCount, &post.CreatedAt, &post.UpdatedAt)
}

// GetByID retrieves a post by ID.
func (r *PostRepository) GetByID(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post
	err := r.db.GetContext(ctx, &post, "SELECT * FROM posts WHERE id = $1", id)
	return &post, err
}

// List returns a paginated list of posts ordered by newest first.
func (r *PostRepository) List(ctx context.Context, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.SelectContext(ctx, &posts,
		"SELECT * FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	return posts, err
}

// Like adds a like to a post.
func (r *PostRepository) Like(ctx context.Context, postID, userID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert like record
	_, err = tx.ExecContext(ctx,
		"INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		postID, userID,
	)
	if err != nil {
		return err
	}

	// Increment like count
	_, err = tx.ExecContext(ctx,
		"UPDATE posts SET like_count = like_count + 1 WHERE id = $1",
		postID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetByUserIDs retrieves posts by multiple user IDs.
func (r *PostRepository) GetByUserIDs(ctx context.Context, userIDs []string, limit, offset int) ([]model.Post, error) {
	query, args, err := sqlx.In(
		"SELECT * FROM posts WHERE user_id IN (?) ORDER BY created_at DESC LIMIT ? OFFSET ?",
		userIDs, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var posts []model.Post
	err = r.db.SelectContext(ctx, &posts, query, args...)
	return posts, err
}
