package router

import (
	"net/http"

	"github.com/dxas90/learn-go/internal/handlers"
	"github.com/dxas90/learn-go/internal/middleware"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Router wraps the mux router with application-specific configuration
type Router struct {
	mux *mux.Router
}

// NewRouter creates and configures a new Router instance.
// It sets up all application routes and applies middleware.
// Returns an error if handler initialization fails.
func NewRouter() (*Router, error) {
	r := mux.NewRouter()

	// Create handlers
	h, err := handlers.NewHandlers()
	if err != nil {
		return nil, err
	}

	// Apply middleware (order matters!)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.SecurityHeadersMiddleware)
	r.Use(middleware.MetricsMiddleware)
	// OpenTelemetry tracing middleware
	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "http-server")
	})

	// Routes
	r.HandleFunc("/", h.Index).Methods("GET")
	r.HandleFunc("/ping", h.Ping).Methods("GET")
	r.HandleFunc("/healthz", h.Healthz).Methods("GET")
	r.HandleFunc("/info", h.Info).Methods("GET")
	r.HandleFunc("/version", h.Version).Methods("GET")
	r.HandleFunc("/echo", h.Echo).Methods("POST")
	r.HandleFunc("/metrics", h.Metrics).Methods("GET")

	return &Router{
		mux: r,
	}, nil
}

// Mux returns the underlying mux.Router instance
func (r *Router) Mux() *mux.Router {
	return r.mux
}

