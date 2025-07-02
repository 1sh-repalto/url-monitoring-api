package model

import "time"

type URLLog struct {
	ID             string    `json:"id"`
	URLID          string    `json:"url_id"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int       `json:"response_time_ms"`
	CheckedAt      time.Time `json:"checked_at"`
	IsUp           bool      `json:"is_up"`
}
