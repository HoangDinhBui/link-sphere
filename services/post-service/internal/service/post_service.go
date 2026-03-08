package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/linksphere/pkg/config"
	"github.com/linksphere/pkg/kafka"
	"github.com/linksphere/post-service/internal/model"
	"github.com/linksphere/post-service/internal/repository"
)

// PostService contains business logic for post operations.
type PostService struct {
	repo     *repository.PostRepository
	producer *kafka.Producer
}

// NewPostService creates a new PostService.
func NewPostService(repo *repository.PostRepository, cfg *config.Config) *PostService {
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	producer := kafka.NewProducer(brokers, "post-events")

	return &PostService{
		repo:     repo,
		producer: producer,
	}
}

// Create creates a new post.
func (s *PostService) Create(ctx context.Context, userID string, req *model.CreatePostRequest) (*model.Post, error) {
	if req.Content == "" {
		return nil, errors.New("content is required")
	}

	post := &model.Post{
		UserID:   userID,
		Content:  req.Content,
		Images:   req.Images,
		Hashtags: req.Hashtags,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	// Publish post created event to Kafka
	s.publishEvent(ctx, "post.created", userID, post.ID)

	return post, nil
}

// GetByID retrieves a post by ID.
func (s *PostService) GetByID(ctx context.Context, id string) (*model.Post, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns a paginated list of posts.
func (s *PostService) List(ctx context.Context, page, limit int) ([]model.Post, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, limit, offset)
}

// Like adds a like to a post and publishes a Kafka event.
func (s *PostService) Like(ctx context.Context, postID, userID string) error {
	if err := s.repo.Like(ctx, postID, userID); err != nil {
		return err
	}

	// Publish like event to Kafka
	s.publishEvent(ctx, "post.liked", userID, postID)

	return nil
}

func (s *PostService) publishEvent(ctx context.Context, event, userID, postID string) {
	evt := model.KafkaEvent{
		Event:     event,
		UserID:    userID,
		PostID:    postID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(evt)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal kafka event")
		return
	}

	if err := s.producer.Publish(ctx, []byte(postID), data); err != nil {
		log.Error().Err(err).Msg("failed to publish kafka event")
	}
}
