# AI Coding Agent Instructions for learn-go

## Project Overview
This is a production-ready Go microservice demonstrating modern Go development practices, Kubernetes deployment, and comprehensive CI/CD pipelines. The application is a simple RESTful API with health checks, monitoring endpoints, and security middleware.

## Architecture & Code Organization

### Standard Go Project Layout (cmd/internal/pkg)
- **`cmd/api/main.go`**: Application entry point - minimal main() that delegates to `internal/server`
- **`internal/`**: Private application code (not importable by other projects)
  - `handlers/`: HTTP request handlers implementing business logic
  - `middleware/`: HTTP middleware chain (logging, CORS, security headers)
  - `router/`: Route definitions using gorilla/mux
  - `server/`: HTTP server configuration with timeouts
- **`pkg/models/`**: Shared data models (Response, AppInfo, HealthStatus, etc.)

### Key Architectural Patterns
1. **Dependency Injection via Constructors**: All components use `NewXxx()` constructors that return `(*Type, error)` - see `handlers.NewHandlers()`, `router.NewRouter()`, `server.NewServer()`
2. **Middleware Chain Pattern**: Router applies middleware in order: Logging → CORS → SecurityHeaders → Handlers
3. **Environment-Based Configuration**: All configuration via env vars (PORT, HOST, GO_ENV, APP_VERSION, CORS_ORIGIN)
4. **Test Mode Detection**: Middleware checks `GO_ENV=test` to suppress logging during tests

## Critical Developer Workflows

### Local Development
```bash
# Quick start (uses Makefile)
make dev                    # Runs: go run ./cmd/api

# Manual start with custom config
PORT=3000 GO_ENV=development APP_VERSION=0.1.0 go run ./cmd/api
```

### Testing Patterns
```bash
# Run all tests with GO_ENV set to avoid log spam
GO_ENV=test go test ./...

# Coverage (as used in CI)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Test structure: Each package has handlers_test.go, middleware_test.go, etc.
# Always set: os.Setenv("GO_ENV", "test") in test setup to disable middleware logging
```

### CI/CD Workflows

**GitHub Actions - Dual Pipeline Architecture:**
1. **`dockerimage.yml`** (primary build):
   - Builds ONCE: multi-arch image (amd64/arm64) → pushes to GHCR
   - Builds test artifact: single-arch (amd64) → exports as tar for K8s tests
   - Artifact shared between workflows (retention: 1 day)

2. **`k8s-deployment.yml`** (triggered by workflow_run):
   - Downloads Docker artifact from `dockerimage.yml`
   - NO rebuild - loads pre-built image into Kind cluster
   - Tests with Helm chart deployment

3. **`trigger-flux.yaml`**: Triggers FluxCD reconciliation after successful builds

**GitLab CI - 7 Stage Pipeline:**
Stages: lint → test → security → build → deploy-staging → deploy-production → cleanup
- Security scanning with Trivy (fails on HIGH/CRITICAL vulnerabilities)
- Kind-based staging deployment with 2 replicas
- Manual production deployment (tag-based only)

### Build & Deploy Commands
```bash
# Docker build (multi-stage, scratch-based final image)
docker build -t learn-go .

# Makefile shortcuts
make build          # Builds binary to bin/learn-go
make docker-build   # Builds Docker image
make test-coverage  # Generates coverage.html
```

## Project-Specific Conventions

### HTTP Response Structure
All JSON responses follow this model (see `pkg/models/responses.go`):
```go
type Response struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    Timestamp string      `json:"timestamp"`
}
```

### Handler Implementation Pattern
Every handler must:
1. Set `Content-Type: application/json` header
2. Use `json.NewEncoder(w).Encode()` for responses
3. Return `models.Response` struct with timestamp
4. Log errors before returning error responses

Example from `handlers.go`:
```go
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("pong"))
}
```

### Middleware Order Matters
Applied in `router.NewRouter()`:
```go
r.Use(middleware.LoggingMiddleware)       // First: log all requests
r.Use(middleware.CORSMiddleware)          // Second: handle CORS
r.Use(middleware.SecurityHeadersMiddleware) // Last: add security headers
```

### Environment Variables Pattern
Always provide defaults in code:
```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"  // Default in code, not in environment
}
```

## Integration Points & External Dependencies

### Third-Party Libraries
- **gorilla/mux**: HTTP router (not stdlib, chosen for better route matching)
- **gopsutil/v3**: System metrics for /healthz endpoint (CPU, memory, process stats)

### Kubernetes Integration
- **Helm chart**: `k8s/chart/` contains full deployment configuration
- **Health endpoints**: `/healthz` (detailed), `/ping` (simple) - used by K8s probes
- **Graceful shutdown**: Server configured with 15s read/write timeout, 60s idle

### CI/CD Artifact Sharing
GitHub Actions workflows share Docker images via artifacts:
```yaml
# In dockerimage.yml
- uses: actions/upload-artifact@v4
  with:
    name: docker-image
    path: /tmp/image.tar

# In k8s-deployment.yml
- uses: dawidd6/action-download-artifact@v3
  with:
    workflow: dockerimage.yml
    name: docker-image
```

## Common Pitfalls & Important Notes

1. **Don't rebuild Docker in K8s workflow** - Always download the artifact from `dockerimage.yml`
2. **Test environment isolation** - Set `GO_ENV=test` to disable logging during tests
3. **Handler initialization** - `NewHandlers()` reads env vars at startup, not per-request
4. **Middleware order** - Logging must be first to capture all requests
5. **Multi-stage Dockerfile** - Final image is `FROM scratch`, requires static binary (CGO_ENABLED=0)

## Quick Reference: Key Files

- **Entry point**: `cmd/api/main.go` (35 lines, minimal)
- **Route definitions**: `internal/router/router.go` (all 6 endpoints listed)
- **Handler logic**: `internal/handlers/handlers.go` (223 lines, all business logic)
- **CI/CD**: `.github/workflows/dockerimage.yml` (primary), `k8s-deployment.yml` (testing)
- **Deployment**: `k8s/chart/` (Helm), `.gitlab-ci.yml` (GitLab pipelines)

## Testing Quick Start

```bash
# Run specific package tests
go test ./internal/handlers -v

# Test with race detection (as in CI)
go test -race ./...

# Generate and view coverage
make test-coverage  # Opens coverage.html in browser
```
