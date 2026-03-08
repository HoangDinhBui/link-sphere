package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/linksphere/pkg/config"
)

// UserRecord represents a user row from the database.
type UserRecord struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
}

// AuthService handles authentication logic.
type AuthService struct {
	db     *sqlx.DB
	secret string
	expiry time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(cfg *config.Config) *AuthService {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to database - running in degraded mode")
	}
	return &AuthService{
		db:     db,
		secret: cfg.JWTSecret,
		expiry: cfg.JWTExpiry,
	}
}

// Login validates credentials and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	if s.db == nil {
		return "", errors.New("database not available")
	}

	var user UserRecord
	err := s.db.GetContext(ctx, &user, "SELECT id, email, password, username FROM users WHERE email = $1", email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(s.expiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenStr, nil
}
