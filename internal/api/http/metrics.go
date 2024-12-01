package http

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	healthStatusMetric = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "request",
			Subsystem: "http",
			Name:      "payment_service_health",
			Help:      "Payment service components health",
		}, []string{"component"},
	)

	performEmissionClearanceHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "request",
			Subsystem: "http",
			Name:      "emission_clearance",
			Help:      "Perform emission clearance",
			Buckets:   prometheus.LinearBuckets(0.01, 0.05, 10),
		}, []string{"status"})
)
