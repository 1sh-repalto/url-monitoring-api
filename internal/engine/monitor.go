package engine

import (
	"net/http"
	"time"

	"github.com/1sh-repalto/url-monitoring-api/internal/model"
	"github.com/1sh-repalto/url-monitoring-api/internal/service"
	"github.com/google/uuid"
)

type MonitorEngine struct {
	urlService 	*service.URLService
	client		*http.Client
}

func NewMonitorEngine(s *service.URLService) *MonitorEngine {
	return &MonitorEngine{
		urlService: s,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (e *MonitorEngine) CheckURLs() error {
	urls, err := e.urlService.GetAllActiveURLs()
	if err != nil {
		return err
	}

	for _, u := range urls {
		go e.checkAndLog(u)
	}

	return nil
}

func (e *MonitorEngine) checkAndLog(u *model.MonitoredURL) {
	start := time.Now()
	resp, err := e.client.Get(u.URL)
	duration := time.Since(start)

	log := &model.URLLog{
		ID:				uuid.NewString(),
		URLID:			u.ID,
		ResponseTimeMs:	int(duration.Milliseconds()),
		CheckedAt: 		time.Now(),
	}

	if err != nil {
		log.StatusCode = 0
		log.IsUp = false
	} else {
		defer resp.Body.Close()
		log.StatusCode = resp.StatusCode
		log.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400
	}

	if err := e.urlService.LogURLCheck(log); err != nil {
		
	}

}