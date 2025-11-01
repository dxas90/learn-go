package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dxas90/learn-go/internal/server"
)

func main() {
	// Create and initialize the server
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get host from environment or use default
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// Print startup information
	fmt.Printf("🚀 Server starting at http://%s:%s/\n", host, port)
	fmt.Printf("📊 Environment: %s\n", os.Getenv("GO_ENV"))
	fmt.Printf("📦 Version: %s\n", os.Getenv("APP_VERSION"))
	fmt.Printf("🕐 Started at: %s\n", time.Now().UTC().Format(time.RFC3339))

	// Start the server (blocks until error or shutdown)
	if err := srv.Start(host + ":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
