package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	h, err := NewHandlers()
	if err != nil {
		t.Fatalf("Failed to create handlers: %v", err)
	}

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	h.Ping(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "pong" {
		t.Errorf("Expected 'pong', got '%s'", w.Body.String())
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", contentType)
	}
}

func TestHealthz(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	h, err := NewHandlers()
	if err != nil {
		t.Fatalf("Failed to create handlers: %v", err)
	}

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	h.Healthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success=true, got %v", response["success"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object")
	}

	if status, ok := data["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status='healthy', got %v", data["status"])
	}
}

func TestVersion(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	h, err := NewHandlers()
	if err != nil {
		t.Fatalf("Failed to create handlers: %v", err)
	}

	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()

	h.Version(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object")
	}

	if name, ok := data["name"].(string); !ok || name != "learn-go" {
		t.Errorf("Expected name='learn-go', got %v", data["name"])
	}
}

func TestEcho(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	h, err := NewHandlers()
	if err != nil {
		t.Fatalf("Failed to create handlers: %v", err)
	}

	testData := map[string]interface{}{
		"message": "hello",
		"number":  42,
	}
	jsonData, _ := json.Marshal(testData)

	req := httptest.NewRequest("POST", "/echo", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Echo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object")
	}

	echo, ok := data["echo"]
	if !ok {
		t.Fatalf("Expected echo field")
	}

	// Check that echo contains the original data
	echoMap, ok := echo.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected echo to be map")
	}

	if msg, ok := echoMap["message"].(string); !ok || msg != "hello" {
		t.Errorf("Expected message='hello', got %v", echoMap["message"])
	}
}

func TestEchoInvalidJSON(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	h, err := NewHandlers()
	if err != nil {
		t.Fatalf("Failed to create handlers: %v", err)
	}

	req := httptest.NewRequest("POST", "/echo", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Echo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if errorFlag, ok := response["error"].(bool); !ok || !errorFlag {
		t.Errorf("Expected error=true, got %v", response["error"])
	}

	if message, ok := response["message"].(string); !ok || message != "Invalid JSON" {
		t.Errorf("Expected message='Invalid JSON', got %v", response["message"])
	}
}

func TestIndexEndpoint(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	h, err := NewHandlers()
	if err != nil {
		t.Fatalf("Failed to create handlers: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object")
	}

	if msg, ok := data["message"].(string); !ok || msg == "" {
		t.Errorf("Expected non-empty message, got %v", data["message"])
	}

	if endpoints, ok := data["endpoints"].([]interface{}); !ok || len(endpoints) == 0 {
		t.Errorf("Expected endpoints array, got %v", data["endpoints"])
	}
}

func TestInfoEndpoint(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	h, err := NewHandlers()
	if err != nil {
		t.Fatalf("Failed to create handlers: %v", err)
	}

	req := httptest.NewRequest("GET", "/info", nil)
	w := httptest.NewRecorder()

	h.Info(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data object")
	}

	if _, ok := data["application"]; !ok {
		t.Errorf("Expected application field")
	}

	if _, ok := data["system"]; !ok {
		t.Errorf("Expected system field")
	}
}

