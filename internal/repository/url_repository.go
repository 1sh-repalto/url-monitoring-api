package repository

import (
	"context"

	"github.com/1sh-repalto/url-monitoring-api/internal/model"
)

type URLRepository interface {
	SaveURL(u *model.MonitoredURL) error
	ExistsByUserAndURL(ctx context.Context, url string, userID int) (bool, error) 
	GetURLByUserID(userID int) ([]*model.MonitoredURL, error)
	GetURLByID(id string) (*model.MonitoredURL, error)
	DeleteURL(id string) error
	GetAllActiveURLs() ([]*model.MonitoredURL, error)
	SaveURLLog(log *model.URLLog) error
	GetLogsByURLID(urlID string) ([]*model.URLLog, error)
}
