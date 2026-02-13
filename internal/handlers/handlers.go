package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/dxas90/learn-go/internal/apispec"
	"github.com/dxas90/learn-go/pkg/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"gopkg.in/yaml.v3"
)

// Handlers contains all HTTP request handlers for the application
type Handlers struct {
	appInfo   models.AppInfo
	startTime time.Time
}

// NewHandlers creates a new Handlers instance with application metadata
// It reads configuration from environment variables and initializes the start time
func NewHandlers() (*Handlers, error) {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "0.0.1"
	}

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	log.Printf("Creating handlers with version=%s, env=%s", version, env)

	return &Handlers{
		appInfo: models.AppInfo{
			Name:        "learn-go",
			Version:     version,
			Environment: env,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		},
		startTime: time.Now(),
	}, nil
}

// Index handles the root endpoint (/)
// Returns a welcome message with application information
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Data: models.WelcomeData{
			Message:     "Welcome to learn-go API",
			Description: "A simple Go microservice for learning and demonstration",
			Documentation: models.Documentation{
				Swagger: nil,
				Postman: nil,
			},
			Links: models.Links{
				Repository: "https://github.com/dxas90/learn-go",
				Issues:     "https://github.com/dxas90/learn-go/issues",
			},
			Endpoints: []models.Endpoint{
				{Path: "/", Method: "GET", Description: "API welcome and documentation"},
				{Path: "/ping", Method: "GET", Description: "Simple ping-pong response"},
				{Path: "/healthz", Method: "GET", Description: "Health check endpoint"},
				{Path: "/info", Method: "GET", Description: "Application and system information"},
				{Path: "/version", Method: "GET", Description: "Application version information"},
				{Path: "/echo", Method: "POST", Description: "Echo back the request body"},
				{Path: "/openapi.json", Method: "GET", Description: "OpenAPI specification (JSON)"},
				{Path: "/openapi.yaml", Method: "GET", Description: "OpenAPI specification (YAML)"},
				{Path: "/metrics", Method: "GET", Description: "Prometheus metrics"},
			},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Ping handles the /ping endpoint
// Returns a simple "pong" text response to verify the service is responsive
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// Healthz handles the /healthz endpoint
// Returns detailed health information including memory usage and uptime
func (h *Handlers) Healthz(w http.ResponseWriter, r *http.Request) {
	p, _ := process.NewProcess(int32(os.Getpid()))
	memInfo, _ := p.MemoryInfo()
	virtualMem, _ := mem.VirtualMemory()

	uptime := time.Since(h.startTime).Seconds()

	response := models.Response{
		Success: true,
		Data: models.HealthData{
			Status:    "healthy",
			Uptime:    uptime,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Memory: models.MemoryInfo{
				RSS:       memInfo.RSS,
				VMS:       memInfo.VMS,
				Percent:   memInfo.RSS * 100 / virtualMem.Total,
				Available: virtualMem.Available,
				Total:     virtualMem.Total,
			},
			Version:     h.appInfo.Version,
			Environment: h.appInfo.Environment,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Info handles the /info endpoint
// Returns comprehensive system and runtime information including CPU, memory, and process details
func (h *Handlers) Info(w http.ResponseWriter, r *http.Request) {
	p, _ := process.NewProcess(int32(os.Getpid()))
	memInfo, _ := p.MemoryInfo()
	virtualMem, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(time.Millisecond*100, false)
	cpuCount, _ := cpu.Counts(true)

	response := models.Response{
		Success: true,
		Data: models.InfoData{
			Application: h.appInfo,
			System: models.SystemInfo{
				Platform:        runtime.GOOS,
				PlatformRelease: "", // Not easily available in Go
				PlatformVersion: "", // Not easily available in Go
				Architecture:    runtime.GOARCH,
				Processor:       "", // Not easily available in Go
				GoVersion:       runtime.Version(),
				Uptime:          time.Since(h.startTime).Seconds(),
				Memory: models.MemoryInfo{
					RSS:       memInfo.RSS,
					VMS:       memInfo.VMS,
					Percent:   memInfo.RSS * 100 / virtualMem.Total,
					Available: virtualMem.Available,
					Total:     virtualMem.Total,
					Used:      virtualMem.Used,
				},
				CPU: models.CPUInfo{
					Count:   cpuCount,
					Percent: cpuPercent[0],
				},
			},
			Environment: models.EnvironmentInfo{
				GoEnv: os.Getenv("GO_ENV"),
				Port:  os.Getenv("PORT"),
				Host:  os.Getenv("HOST"),
			},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Version handles the /version endpoint
// Returns application version and environment information
func (h *Handlers) Version(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Success: true,
		Data: models.VersionData{
			Version:     h.appInfo.Version,
			Name:        h.appInfo.Name,
			Environment: h.appInfo.Environment,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Echo handles the /echo endpoint
// Accepts JSON in the request body and echoes it back along with request metadata
// Returns a 400 Bad Request if the JSON payload is invalid
func (h *Handlers) Echo(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response := models.ErrorResponse{
			Error:      true,
			Message:    "Invalid JSON",
			StatusCode: 400,
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	response := models.Response{
		Success: true,
		Data: models.EchoData{
			Echo:    data,
			Headers: headers,
			Method:  r.Method,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OpenAPISpec handles the /openapi.json endpoint
// Returns the embedded OpenAPI YAML spec converted to JSON
func (h *Handlers) OpenAPISpec(w http.ResponseWriter, r *http.Request) {
	// Convert embedded YAML to JSON
	var yamlData interface{}
	if err := yaml.Unmarshal(apispec.OpenAPISpec, &yamlData); err != nil {
		log.Printf("Error parsing OpenAPI spec: %v", err)
		http.Error(w, "Failed to parse OpenAPI spec", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(yamlData)
	if err != nil {
		log.Printf("Error converting OpenAPI spec to JSON: %v", err)
		http.Error(w, "Failed to convert OpenAPI spec to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// OpenAPISpecYAML handles the /openapi.yaml endpoint
// Returns the embedded OpenAPI spec in YAML format
func (h *Handlers) OpenAPISpecYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(apispec.OpenAPISpec)
}
