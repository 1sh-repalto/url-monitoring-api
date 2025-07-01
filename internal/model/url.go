package model

import "time"

type MonitoredURL struct {
	ID			string
	URL			string
	UserID		int
	IsActive	bool
	CreatedAt	time.Time
}