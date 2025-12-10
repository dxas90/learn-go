package server

import (
	"log"
	"net/http"
	"time"

	"github.com/dxas90/learn-go/internal/router"
)

// Server represents the HTTP server with its router
type Server struct {
	router *router.Router
}

// NewServer creates a new Server instance with an initialized router.
// Returns an error if router initialization fails.
func NewServer() (*Server, error) {
	r, err := router.NewRouter()
	if err != nil {
		return nil, err
	}

	return &Server{
		router: r,
	}, nil
}

// Start starts the HTTP server on the specified address.
// It configures timeouts and logs any errors that occur.
// The server will block until it encounters an error or is shut down.
func (s *Server) Start(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router.Mux(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting HTTP server on %s", addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Printf("HTTP server error: %v", err)
	}
	return err
}

