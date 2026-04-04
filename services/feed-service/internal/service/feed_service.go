package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/linksphere/pkg/config"
	pkgredis "github.com/linksphere/pkg/redis"
)

// FeedItem represents a single item in the news feed.
type FeedItem struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Images    []string  `json:"images,omitempty"`
	Hashtags  []string  `json:"hashtags,omitempty"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
}

// FeedService handles feed generation logic.
type FeedService struct {
	cfg         *config.Config
	redisClient *pkgredis.Client
}

// NewFeedService creates a new FeedService.
func NewFeedService(cfg *config.Config, redisClient *pkgredis.Client) *FeedService {
	return &FeedService{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

// GetFeed generates the news feed for a user.
// Strategy: Fetch posts from followed users, sort by time, cache in Redis.
func (s *FeedService) GetFeed(ctx context.Context, userID string, page, limit int) ([]FeedItem, error) {
	cacheKey := fmt.Sprintf("feed:%s:page:%d:limit:%d", userID, page, limit)

	// Try cache first
	cached, err := s.redisClient.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var feed []FeedItem
		if err := json.Unmarshal([]byte(cached), &feed); err == nil {
			log.Debug().Str("user_id", userID).Msg("feed served from cache")
			return feed, nil
		}
	}

	// Cache miss — fetch from user service (get following list)
	// In production, this would call user-service and post-service internally
	feed, err := s.fetchFeedFromServices(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	// Cache the result for 5 minutes
	if data, err := json.Marshal(feed); err == nil {
		if err := s.redisClient.Set(ctx, cacheKey, string(data), 5*time.Minute); err != nil {
			log.Warn().Err(err).Msg("failed to cache feed")
		}
	}

	return feed, nil
}

// fetchFeedFromServices calls other microservices to build the feed.
func (s *FeedService) fetchFeedFromServices(ctx context.Context, userID string, page, limit int) ([]FeedItem, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Step 1: Get following list from user-service
	userServiceURL := fmt.Sprintf("http://user-service:8001/api/v1/users/%s/following", userID)
	req1, err := http.NewRequestWithContext(ctx, http.MethodGet, userServiceURL, nil)
	if err != nil {
		return []FeedItem{}, nil
	}

	resp, err := client.Do(req1)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch following list from user-service")
		return []FeedItem{}, nil
	}
	defer resp.Body.Close()

	var followingIDs []string
	if err := json.NewDecoder(resp.Body).Decode(&followingIDs); err != nil {
		log.Warn().Err(err).Msg("failed to decode following list")
		return []FeedItem{}, nil
	}

	// If the user isn't following anyone, return empty feed (or global feed)
	if len(followingIDs) == 0 {
		return []FeedItem{}, nil
	}

	// Step 2: Get posts from post-service for those users
	postServiceURL := "http://post-service:8003/api/v1/posts/by-users"
	postReq := map[string]interface{}{
		"userIds": followingIDs,
		"page":    page,
		"limit":   limit,
	}
	postData, _ := json.Marshal(postReq)

	req2, err := http.NewRequestWithContext(ctx, http.MethodPost, postServiceURL, bytes.NewBuffer(postData))
	if err != nil {
		return []FeedItem{}, nil
	}
	req2.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req2)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch posts from post-service")
		return []FeedItem{}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []FeedItem{}, nil
	}

	var feed []FeedItem
	if err := json.Unmarshal(body, &feed); err != nil {
		log.Warn().Err(err).Msg("failed to unmarshal feed items")
		return []FeedItem{}, nil
	}

	return feed, nil
}
