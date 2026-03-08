package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/linksphere/pkg/config"
	"github.com/linksphere/user-service/internal/model"
)

// UserRepository handles database operations for users.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(cfg *config.Config) *UserRepository {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to database - running in degraded mode")
		return &UserRepository{}
	}
	log.Info().Msg("connected to PostgreSQL")
	return &UserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, phone, username, password) 
	           VALUES ($1, $2, $3, $4) 
	           RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		user.Email, user.Phone, user.Username, user.Password,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by their ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	return &user, err
}

// GetByEmail retrieves a user by their email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	return &user, err
}

// Follow creates a follow relationship between two users.
func (r *UserRepository) Follow(ctx context.Context, followerID, followeeID string) error {
	query := `INSERT INTO follows (follower_id, followee_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, followerID, followeeID)
	return err
}

// Unfollow removes a follow relationship.
func (r *UserRepository) Unfollow(ctx context.Context, followerID, followeeID string) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2`
	_, err := r.db.ExecContext(ctx, query, followerID, followeeID)
	return err
}

// GetFollowing returns the IDs of users that a given user follows.
func (r *UserRepository) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	var ids []string
	err := r.db.SelectContext(ctx, &ids, "SELECT followee_id FROM follows WHERE follower_id = $1", userID)
	return ids, err
}
