// Package models defines the data structures used for API requests and responses.
// It includes standard response formats, error responses, and endpoint-specific data models.
package models

// AppInfo holds application metadata including name, version, environment, and timestamp
type AppInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Timestamp   string `json:"timestamp"`
}

// Response represents a standard API response
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error      bool   `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Timestamp  string `json:"timestamp"`
}

// WelcomeData for the index endpoint
type WelcomeData struct {
	Message       string        `json:"message"`
	Description   string        `json:"description"`
	Documentation Documentation `json:"documentation"`
	Links         Links         `json:"links"`
	Endpoints     []Endpoint    `json:"endpoints"`
}

// Documentation links
type Documentation struct {
	Swagger *string `json:"swagger"`
	Postman *string `json:"postman"`
}

// Links for repository and issues
type Links struct {
	Repository string `json:"repository"`
	Issues     string `json:"issues"`
}

// Endpoint description
type Endpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

// HealthData for health check
type HealthData struct {
	Status      string     `json:"status"`
	Uptime      float64    `json:"uptime"`
	Timestamp   string     `json:"timestamp"`
	Memory      MemoryInfo `json:"memory"`
	Version     string     `json:"version"`
	Environment string     `json:"environment"`
}

// MemoryInfo for memory statistics
type MemoryInfo struct {
	RSS       uint64 `json:"rss"`
	VMS       uint64 `json:"vms"`
	Percent   uint64 `json:"percent"`
	Available uint64 `json:"available"`
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used,omitempty"`
}

// InfoData for system information
type InfoData struct {
	Application AppInfo         `json:"application"`
	System      SystemInfo      `json:"system"`
	Environment EnvironmentInfo `json:"environment"`
}

// SystemInfo for system details
type SystemInfo struct {
	Platform        string     `json:"platform"`
	PlatformRelease string     `json:"platform_release"`
	PlatformVersion string     `json:"platform_version"`
	Architecture    string     `json:"architecture"`
	Processor       string     `json:"processor"`
	GoVersion       string     `json:"go_version"`
	Uptime          float64    `json:"uptime"`
	Memory          MemoryInfo `json:"memory"`
	CPU             CPUInfo    `json:"cpu"`
}

// CPUInfo for CPU details
type CPUInfo struct {
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

// EnvironmentInfo for environment variables
type EnvironmentInfo struct {
	GoEnv string `json:"go_env"`
	Port  string `json:"port"`
	Host  string `json:"host"`
}

// VersionData for version endpoint
type VersionData struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
}

// EchoData for echo endpoint
type EchoData struct {
	Echo    interface{}       `json:"echo"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
}
