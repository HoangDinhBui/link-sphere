package model

import "time"

// Post represents a social media post.
type Post struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	Images    []string  `json:"images,omitempty"`
	Hashtags  []string  `json:"hashtags,omitempty"`
	LikeCount int       `json:"like_count" db:"like_count"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePostRequest is the request body for creating a post.
type CreatePostRequest struct {
	Content  string   `json:"content"`
	Images   []string `json:"images,omitempty"`
	Hashtags []string `json:"hashtags,omitempty"`
}

// PostLike represents a like on a post.
type PostLike struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"post_id" db:"post_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// KafkaEvent represents an event published to Kafka.
type KafkaEvent struct {
	Event     string `json:"event"`
	UserID    string `json:"userId"`
	PostID    string `json:"postId"`
	OwnerID   string `json:"ownerId,omitempty"`
	Timestamp string `json:"timestamp"`
}
