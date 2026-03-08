package service

import (
	"context"
	"errors"
	"net/mail"

	"golang.org/x/crypto/bcrypt"

	"github.com/linksphere/user-service/internal/model"
	"github.com/linksphere/user-service/internal/repository"
)

// UserService contains business logic for user operations.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register creates a new user account.
func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// Validate email
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, errors.New("invalid email address")
	}

	// Validate username
	if len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}

	// Validate password
	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Email:    req.Email,
		Phone:    req.Phone,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByEmail retrieves a user by email.
func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// Follow creates a follow relationship.
func (s *UserService) Follow(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	return s.repo.Follow(ctx, followerID, followeeID)
}

// Unfollow removes a follow relationship.
func (s *UserService) Unfollow(ctx context.Context, followerID, followeeID string) error {
	return s.repo.Unfollow(ctx, followerID, followeeID)
}

// GetFollowing returns the list of user IDs that a user follows.
func (s *UserService) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	return s.repo.GetFollowing(ctx, userID)
}
