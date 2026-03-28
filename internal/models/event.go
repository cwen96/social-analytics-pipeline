package models

import "time"

// EventType represents the kind of social media engagement.
type EventType string

const (
	EventLike    EventType = "like"
	EventShare   EventType = "share"
	EventClick   EventType = "click"
	EventComment EventType = "comment"
	EventRepost  EventType = "repost"
)

// EngagementEvent represents a single social media engagement action.
type EngagementEvent struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Platform  string    `json:"platform"` // e.g. "twitter", "instagram", "facebook"
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}
