package engine

import (
	"log"
	"net/http"
	"strconv"
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
			Timeout: 10 * time.Second,
		},
	}
}

func (e *MonitorEngine) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		metrics.ResetCycleMetrics()

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

	var (
		wg            sync.WaitGroup
		userCheckMap  = make(map[string]int)
		userFailMap   = make(map[string]int)
		mapMutex      sync.Mutex
	)

	for _, u := range urls {
		wg.Add(1)
		go func(u *model.MonitoredURL) {
			defer wg.Done()

			start := time.Now()
			resp, err := e.client.Get(u.URL)
			duration := time.Since(start)
			userID := strconv.Itoa(u.UserID)

			mapMutex.Lock()
			userCheckMap[userID]++
			mapMutex.Unlock()

			metrics.ResponseTimeGaugeVec.WithLabelValues(userID, u.ID, u.URL).Set(float64(duration.Milliseconds()))

			urlLog := &model.URLLog{
				ID:             uuid.NewString(),
				URLID:          u.ID,
				ResponseTimeMs: int(duration.Milliseconds()),
				CheckedAt:      time.Now().UTC(),
			}

			if err != nil {
				urlLog.StatusCode = 0
				urlLog.IsUp = false

				mapMutex.Lock()
				userFailMap[userID]++
				mapMutex.Unlock()

				metrics.StatusCodeGaugeVec.WithLabelValues(userID, u.ID, u.URL).Set(0)
			} else {
				defer resp.Body.Close()
				urlLog.StatusCode = resp.StatusCode
				urlLog.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

				metrics.StatusCodeGaugeVec.WithLabelValues(userID, u.ID, u.URL).Set(float64(resp.StatusCode))

				if !urlLog.IsUp {
					mapMutex.Lock()
					userFailMap[userID]++
					mapMutex.Unlock()
				}
			}

			if err := e.urlService.LogURLCheck(urlLog); err != nil {
				log.Printf("failed to log URL check for %s: %v", u.URL, err)
			}
		}(u)
	}

	wg.Wait()

	for userID, count := range userCheckMap {
		metrics.TotalChecksPerCycle.WithLabelValues(userID).Set(float64(count))
	}
	for userID, count := range userFailMap {
		metrics.FailedChecksPerCycle.WithLabelValues(userID).Set(float64(count))
	}

	for userID := range userCheckMap {
    if _, ok := userFailMap[userID]; !ok {
        metrics.FailedChecksPerCycle.WithLabelValues(userID).Set(0)
    }
}

	return nil
}
