package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 1. Total checks per user per cycle
	TotalChecksPerCycle = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "url_checks_total",
		Help: "Total number of URLs checked in the current monitoring cycle per user",
	}, []string{"user_id"})

	// 2. Failed checks per user per cycle
	FailedChecksPerCycle = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "url_checks_failed_total",
		Help: "Number of failed URL checks in the current monitoring cycle per user",
	}, []string{"user_id"})

	// 3. Response time per URL per user
	ResponseTimeGaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "url_response_time_ms",
		Help: "Response time in milliseconds per URL in the current monitoring cycle per user",
	}, []string{"user_id", "url_id", "host"})
	
	// Add this new metric
	StatusCodeGaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "url_last_status_code",
		Help: "HTTP status code of the last check per URL per user",
	}, []string{"user_id", "url_id", "host"})

	metricsMu sync.Mutex
)

func ResetCycleMetrics() {
	metricsMu.Lock()
	defer metricsMu.Unlock()

	TotalChecksPerCycle.Reset()
	FailedChecksPerCycle.Reset()
	ResponseTimeGaugeVec.Reset()
	StatusCodeGaugeVec.Reset()
}

func Handler() http.Handler {
	return promhttp.Handler()
}
