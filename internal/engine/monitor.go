package engine

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/1sh-repalto/url-monitoring-api/internal/metrics"
	"github.com/1sh-repalto/url-monitoring-api/internal/model"
	"github.com/1sh-repalto/url-monitoring-api/internal/service"
	"github.com/google/uuid"
)

type MonitorEngine struct {
	urlService *service.URLService
	client     *http.Client
}

func NewMonitorEngine(s *service.URLService) *MonitorEngine {
	return &MonitorEngine{
		urlService: s,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (e *MonitorEngine) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := e.CheckURLs()
		if err != nil {
			log.Printf("monitoring error: failed to check URLs: %v", err)
		}
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

	urlLog := &model.URLLog{
		ID:             uuid.NewString(),
		URLID:          u.ID,
		ResponseTimeMs: int(duration.Milliseconds()),
		CheckedAt:      time.Now().UTC(),
	}

	hostname := extractHostname(u.URL)

	metrics.TotalChecks.Inc()

	metrics.ResponseTime.WithLabelValues(u.ID, hostname).Observe(float64(duration.Milliseconds()))

	if err != nil {
		urlLog.StatusCode = 0
		urlLog.IsUp = false

		metrics.FailedChecks.Inc()
	} else {
		defer resp.Body.Close()
		urlLog.StatusCode = resp.StatusCode
		urlLog.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

		metrics.CheckStatus.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()

		if !urlLog.IsUp {
			metrics.FailedChecks.Inc()
		}
	}

	if err := e.urlService.LogURLCheck(urlLog); err != nil {
		log.Printf("Failed to log URL check for %s: %v", u.URL, err)
	}
}

func extractHostname(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "invalid"
	}
	return u.Hostname()
}
