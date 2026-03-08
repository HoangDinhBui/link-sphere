package model

import "time"

// User represents a user in the system.
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Follow represents a follow relationship.
type Follow struct {
	ID         string    `json:"id" db:"id"`
	FollowerID string    `json:"follower_id" db:"follower_id"`
	FolloweeID string    `json:"followee_id" db:"followee_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// FollowRequest is the request body for follow/unfollow.
type FollowRequest struct {
	TargetUserID string `json:"targetUserId"`
}
