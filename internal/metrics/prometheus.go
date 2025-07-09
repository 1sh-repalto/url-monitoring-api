package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 1. Total checks per cycle
	TotalChecksPerCycle = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "url_checks_total",
		Help: "Total number of URLs checked in the current monitoring cycle",
	})

	// 2. Failed checks per cycle
	FailedChecksPerCycle = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "url_checks_failed_total",
		Help: "Number of failed URL checks in the current monitoring cycle",
	})

	// 3. Actual response time per URL (non-averaged, per cycle)
	ResponseTimeGaugeVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "url_response_time_ms",
		Help: "Response time in milliseconds per URL in the current monitoring cycle",
	}, []string{"url_id", "host"})

	// For reset protection (optional, Prometheus is single-threaded for metrics but safe practice)
	metricsMu sync.Mutex
)

// Reset all per-cycle gauges (called at start of each monitoring cycle)
func ResetCycleMetrics() {
	metricsMu.Lock()
	defer metricsMu.Unlock()

	TotalChecksPerCycle.Set(0)
	FailedChecksPerCycle.Set(0)
	ResponseTimeGaugeVec.Reset()
}

// Register the Prometheus HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}
