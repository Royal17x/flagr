package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flagr_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "flagr_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5},
		},
		[]string{"method", "status"},
	)

	FlagEvaluationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flagr_flag_evaluations_total",
			Help: "Total number of flag evaluations",
		},
		[]string{"flag_key", "result", "source"},
	)

	FlagEvaluationDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "flagr_flag_evaluation_duration_seconds",
			Help:    "Flag evaluation duration",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05},
		},
	)

	CacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "flagr_cache_hits_total",
			Help: "Total Redis cache hits",
		},
	)

	CacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "flagr_cache_misses_total",
			Help: "Total Redis cache misses",
		},
	)

	KafkaPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flagr_kafka_published_total",
			Help: "Total Kafka messages published",
		},
		[]string{"topic", "status"},
	)

	ActiveGRPCConnection = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "flagr_grpc_active_connections",
			Help: "Number of active gRPC connections",
		},
	)
)
