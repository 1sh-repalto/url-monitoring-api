package engine

import (
	"log"
	"net/http"
	"sync"
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
			Timeout: 10 * time.Second, // Timeout per request
		},
	}
}

func (e *MonitorEngine) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		metrics.ResetCycleMetrics() // 🔁 Reset metrics at start of each cycle

		if err := e.CheckURLs(); err != nil {
			log.Printf("monitoring error: failed to check URLs: %v", err)
		}
	}
}


func (e *MonitorEngine) CheckURLs() error {
	urls, err := e.urlService.GetAllActiveURLs()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(u *model.MonitoredURL) {
			defer wg.Done()
			e.checkAndLog(u)
		}(u)
	}
	wg.Wait()
	return nil
}

func (e *MonitorEngine) checkAndLog(u *model.MonitoredURL) {
	start := time.Now()
	resp, err := e.client.Get(u.URL)
	duration := time.Since(start)

	// Count every check
	metrics.TotalChecksPerCycle.Inc()

	// Report exact response time for this cycle
	metrics.ResponseTimeGaugeVec.WithLabelValues(u.ID, u.URL).Set(float64(duration.Milliseconds()))

	urlLog := &model.URLLog{
		ID:             uuid.NewString(),
		URLID:          u.ID,
		ResponseTimeMs: int(duration.Milliseconds()),
		CheckedAt:      time.Now().UTC(),
	}

	if err != nil {
		urlLog.StatusCode = 0
		urlLog.IsUp = false

		metrics.FailedChecksPerCycle.Inc()
	} else {
		defer resp.Body.Close()
		urlLog.StatusCode = resp.StatusCode
		urlLog.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

		if !urlLog.IsUp {
			metrics.FailedChecksPerCycle.Inc()
		}
	}

	if err := e.urlService.LogURLCheck(urlLog); err != nil {
		log.Printf("failed to log URL check for %s: %v", u.URL, err)
	}
}


