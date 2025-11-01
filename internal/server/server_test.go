package server

import (
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() returned an error: %v", err)
	}

	if s == nil {
		t.Fatal("NewServer() returned a nil server")
	}
}

func TestServerStart(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() returned an error: %v", err)
	}

	go func() {
		if err := s.Start("127.0.0.1:8081"); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server returned an error: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:8081/ping")
	if err != nil {
		t.Fatalf("Failed to make request to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}
