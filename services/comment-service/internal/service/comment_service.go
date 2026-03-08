package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/linksphere/comment-service/internal/model"
	"github.com/linksphere/comment-service/internal/repository"
	"github.com/linksphere/pkg/config"
	"github.com/linksphere/pkg/kafka"
)

// CommentService contains business logic for comment operations.
type CommentService struct {
	repo     *repository.CommentRepository
	producer *kafka.Producer
}

// NewCommentService creates a new CommentService.
func NewCommentService(repo *repository.CommentRepository, cfg *config.Config) *CommentService {
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	producer := kafka.NewProducer(brokers, "comment-events")

	return &CommentService{
		repo:     repo,
		producer: producer,
	}
}

// Create creates a new comment.
func (s *CommentService) Create(ctx context.Context, userID string, req *model.CreateCommentRequest) (*model.Comment, error) {
	if req.Content == "" {
		return nil, errors.New("content is required")
	}

	comment := &model.Comment{
		PostID:  req.PostID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Publish comment event to Kafka
	evt := map[string]string{
		"event":     "post.commented",
		"userId":    userID,
		"postId":    req.PostID,
		"commentId": comment.ID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.Marshal(evt)

	if err := s.producer.Publish(ctx, []byte(req.PostID), data); err != nil {
		log.Error().Err(err).Msg("failed to publish comment event")
	}

	return comment, nil
}

// ListByPostID returns comments for a given post.
func (s *CommentService) ListByPostID(ctx context.Context, postID string, page, limit int) ([]model.Comment, error) {
	offset := (page - 1) * limit
	return s.repo.ListByPostID(ctx, postID, limit, offset)
}
