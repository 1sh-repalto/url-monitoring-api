package repository

import "github.com/1sh-repalto/url-monitoring-api/internal/model"

type URLRepository interface {
	SaveURL	(u *model.MonitoredURL)	error
	GetURLByUserID (userID int) ([]* model.MonitoredURL, error)
	GetURLByID (id string) (*model.MonitoredURL, error)
	DeleteURL (id string) error
}