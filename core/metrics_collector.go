
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()
		status := http.StatusText(recorder.statusCode)

		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
		httpRequestCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.Handle("/metrics", promhttp.Handler())

	wrappedMux := metricsMiddleware(mux)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		log.Fatal(err)
	}
}