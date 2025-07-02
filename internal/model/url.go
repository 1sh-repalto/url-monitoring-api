package model

import "time"

type MonitoredURL struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	UserID    int       `json:"user_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
