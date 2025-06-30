package repository

import "github.com/1sh-repalto/url-monitoring-api/internal/model"

type URLRepository interface {
	SaveURL	(u *model.MonitoredURL)	error
	GetURLByUserID (userID string) ([]* model.MonitoredURL, error)
	DeleteURL (id string) error
}