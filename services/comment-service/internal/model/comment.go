package model

import "time"

// Comment represents a comment on a post.
type Comment struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"post_id" db:"post_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateCommentRequest is the request body for creating a comment.
type CreateCommentRequest struct {
	PostID  string `json:"postId"`
	Content string `json:"content"`
}

// ListCommentsRequest is the request body for listing comments.
type ListCommentsRequest struct {
	PostID string `json:"postId"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}
