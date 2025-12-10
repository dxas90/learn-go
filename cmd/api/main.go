package main

import (
	"log"
	"os"
	"time"

	"github.com/dxas90/learn-go/internal/server"
	"github.com/dxas90/learn-go/internal/telemetry"
)

func main() {
	// Initialize OpenTelemetry tracing
	shutdown, err := telemetry.InitTracer()
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize tracer: %v", err)
	}
	defer shutdown()

	// Create and initialize the server
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("[ERROR] Failed to create server: %v", err)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get host from environment or use default
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	// Print startup information
	log.Printf("[INFO] 🚀 Server starting at http://%s:%s/", host, port)
	log.Printf("[INFO] 📊 Environment: %s", os.Getenv("GO_ENV"))
	log.Printf("[INFO] 📦 Version: %s", os.Getenv("APP_VERSION"))
	log.Printf("[INFO] 🕐 Started at: %s", time.Now().UTC().Format(time.RFC3339))

	// Start the server (blocks until error or shutdown)
	if err := srv.Start(host + ":" + port); err != nil {
		log.Fatalf("[ERROR] Server failed to start: %v", err)
	}
}
