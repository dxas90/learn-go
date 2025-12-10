package middleware

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dxas90/learn-go/internal/handlers"
	"github.com/gorilla/mux"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs incoming HTTP requests with timestamp and user agent.
// Logging is disabled when GO_ENV is set to "test" to avoid cluttering test output.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("GO_ENV") != "test" {
			timestamp := time.Now().UTC().Format(time.RFC3339)
			userAgent := r.Header.Get("User-Agent")
			if userAgent == "" {
				userAgent = "Unknown"
			}
			log.Printf("[INFO] %s %s %s - User-Agent: %s", timestamp, r.Method, r.URL.Path, userAgent)
		}
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds Cross-Origin Resource Sharing (CORS) headers to responses.
// The CORS_ORIGIN environment variable can be used to configure allowed origins.
// Defaults to "*" (allow all origins) if not set.
// Handles OPTIONS preflight requests automatically.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security-related HTTP headers to all responses.
// Headers include:
// - X-Content-Type-Options: nosniff
// - X-Frame-Options: DENY
// - X-XSS-Protection: 1; mode=block
// - Referrer-Policy: strict-origin-when-cross-origin
// - Content-Security-Policy: default-src 'self'
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware tracks Prometheus metrics for HTTP requests
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip metrics endpoint to avoid recursion
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		route := mux.CurrentRoute(r)
		path := r.URL.Path
		if route != nil {
			if template, err := route.GetPathTemplate(); err == nil {
				path = template
			}
		}

		handlers.HTTPRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		handlers.HTTPRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rw.statusCode)).Inc()
	})
}

