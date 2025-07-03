package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalChecks = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "url_checks_total",
			Help: "Total number of URL checks performed",
		},
	)

	FailedChecks = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "url_checks_failed_total",
			Help: "Total number of failed URL checks",
		},
	)

	ResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "url_response_time_ms",
			Help: "Histogram of response time for monitored URLs",
			Buckets: prometheus.ExponentialBuckets(10, 2 ,8),
		},
		[]string{"url_id", "hostname"},
	)

	CheckStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_check_status_code_total",
			Help: "Count of HTTP status codes returned",
		},
		[]string{"status_code"},
	)
)

func Init() {
	prometheus.MustRegister(TotalChecks)
	prometheus.MustRegister(FailedChecks)
	prometheus.MustRegister(ResponseTime)
	prometheus.MustRegister(CheckStatus)
}