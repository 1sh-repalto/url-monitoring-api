package service

import (
	"errors"

	"github.com/1sh-repalto/url-monitoring-api/internal/model"
	"github.com/1sh-repalto/url-monitoring-api/internal/repository"
	"github.com/google/uuid"
)

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(r repository.URLRepository) *URLService {
	return &URLService{repo: r}
}

func (s *URLService) RegisterURL(rawUrl string, userId int) error {
	url := &model.MonitoredURL{
		ID:       uuid.NewString(),
		URL:      rawUrl,
		UserID:   userId,
		IsActive: true,
	}

	return s.repo.SaveURL(url)
}

func (s *URLService) GetURLByUser(userID int) ([]*model.MonitoredURL, error) {
	return s.repo.GetURLByUserID(userID)
}

func (s *URLService) DeleteURL(urlID string, userID int) error {
	url, err := s.repo.GetURLByID(urlID)
	if err != nil {
		return err
	}

	if url.UserID != userID {
		return errors.New("unauthorized: not the owner of this URL")
	}

	return s.repo.DeleteURL(urlID)
}

func (s *URLService) LogURLCheck(log *model.URLLog) error {
	return s.repo.SaveURLLog(log)
}

func (s *URLService) GetAllActiveURLs() ([]*model.MonitoredURL, error) {
	return s.repo.GetAllActiveURLs()
}
