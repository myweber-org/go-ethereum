package main

import (
    "log"
    "net/http"
    "time"
)

var (
    requestCount    = 0
    totalLatency   = time.Duration(0)
)

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        latency := time.Since(start)

        requestCount++
        totalLatency += latency

        log.Printf("Request processed: %s, Latency: %v, Total Requests: %d, Avg Latency: %v",
            r.URL.Path, latency, requestCount, totalLatency/time.Duration(requestCount))
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/", metricsMiddleware(http.HandlerFunc(helloHandler)))

    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}