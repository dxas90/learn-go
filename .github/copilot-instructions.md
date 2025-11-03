````instructions
# AI Coding Agent Instructions for learn-go

## Project Overview
Production-ready Go microservice demonstrating modern development practices, Kubernetes deployment, and comprehensive CI/CD. RESTful API with 6 endpoints (/, /ping, /healthz, /info, /version, /echo), security middleware, and multi-environment testing (local, Docker, Kubernetes).

## Architecture & Code Organization

### Standard Go Project Layout (cmd/internal/pkg)
- **`cmd/api/main.go`**: Minimal entry point (~35 lines) - delegates to `internal/server`
- **`internal/`**: Private application code (NOT importable by other projects)
  - `handlers/`: HTTP request handlers with business logic (all 6 endpoints)
  - `middleware/`: Logging → CORS → SecurityHeaders chain
  - `router/`: Route definitions using gorilla/mux
  - `server/`: HTTP server with 15s read/write timeout, 60s idle
- **`pkg/models/`**: Shared models (Response, AppInfo, HealthStatus)
- **`scripts/`**: Testing scripts (smoke-test.sh, e2e-test.sh, integration-test.sh)

### Key Architectural Patterns
1. **Constructor-based DI**: ALL components use `NewXxx() (*Type, error)` pattern - see `handlers.NewHandlers()`, `router.NewRouter()`, `server.NewServer()`
2. **Middleware Chain Order** (CRITICAL): `r.Use(LoggingMiddleware)` → `r.Use(CORSMiddleware)` → `r.Use(SecurityHeadersMiddleware)` - order matters!
3. **Environment-Based Config**: PORT, HOST, GO_ENV, APP_VERSION, CORS_ORIGIN - defaults in code, not env
4. **Test Mode Detection**: `GO_ENV=test` suppresses logging - MUST set in all test files: `os.Setenv("GO_ENV", "test")`

## Critical Developer Workflows

### Local Development
```bash
make dev                    # Quickest: go run ./cmd/api (port 8080)
PORT=3000 go run ./cmd/api  # Custom port

# Test endpoints
curl http://localhost:8080/healthz  # Detailed health with system metrics
curl http://localhost:8080/ping     # Simple pong response
curl -X POST -H "Content-Type: application/json" -d '{"message":"test"}' http://localhost:8080/echo
```

### Testing (3-Layer Strategy)
```bash
# 1. Unit tests (GO_ENV=test required to suppress logs)
GO_ENV=test go test ./...
go test ./internal/handlers -v  # Test specific package

# 2. Integration tests (local + Docker) - runs before Docker build in CI
./scripts/integration-test.sh    # Auto-detects local vs k8s environment

# 3. E2E tests (Kubernetes) - runs in Kind cluster after deployment
./scripts/smoke-test.sh          # Quick validation (cluster health, basic endpoint)
./scripts/e2e-test.sh            # Comprehensive (all endpoints, pod resilience, probes)

# Coverage (as in CI)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### CI/CD Workflow (Unified full-workflow.yml)
**Single pipeline replaces old dockerimage.yml + k8s-deployment.yml + trigger-flux.yaml**

```
lint → test (matrix: 1.24.x, 1.25.x) → integration-test → build → helm-test
                                                              ↓
                                                       test-deployment (Kind)
                                                              ├─ smoke-test.sh
                                                              ├─ e2e-test.sh
                                                              └─ health checks
                                                              ↓
                                          ┌─────────────────────────────────┐
                                          ↓                                 ↓
                                  deploy-staging                  deploy-production
                                  (main branch)                   (tags only)
                                  FluxCD webhook                  FluxCD webhook
```

**Key Implementation Details:**
- Docker build: Single-platform (amd64) with `load: true` for ALL branches/PRs
- Artifact sharing: Image exported to tar → uploaded → downloaded by test-deployment
- Kind cluster: Uses `.github/kind-config.yaml` with port mappings (80, 443)
- Helm deployment: `--set image.pullPolicy=Never` for Kind (uses loaded image)
- Test scripts: Run in Kind cluster using curl pods for in-cluster endpoint testing

## Project-Specific Conventions

### HTTP Response Structure (CRITICAL)
All JSON endpoints return `pkg/models/Response`:
```go
type Response struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    Timestamp string      `json:"timestamp"`  // RFC3339 format
}
```

### Handler Implementation Pattern
```go
// Example: internal/handlers/handlers.go
func (h *Handlers) Healthz(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")  // REQUIRED
    response := models.Response{
        Success:   true,
        Data:      healthData,
        Timestamp: time.Now().Format(time.RFC3339),     // REQUIRED
    }
    json.NewEncoder(w).Encode(response)                 // Use Encoder, not Marshal
}
```
**Exception**: `/ping` returns plain text "pong" (not JSON)

### Middleware Chain (Order is Critical)
```go
// internal/router/router.go
r.Use(middleware.LoggingMiddleware)       // 1st: Logs all requests (unless GO_ENV=test)
r.Use(middleware.CORSMiddleware)          // 2nd: Handles CORS preflight
r.Use(middleware.SecurityHeadersMiddleware) // 3rd: X-Frame-Options, CSP, etc.
```
**Why order matters**: Logging must capture CORS responses; security headers applied last

### Environment Variable Pattern
```go
// Always provide defaults in code
port := os.Getenv("PORT")
if port == "" {
    port = "8080"  // Default here, not in .env files
}
```

### Test File Pattern
```go
// REQUIRED in all test files
func TestXxx(t *testing.T) {
    os.Setenv("GO_ENV", "test")  // Disables middleware logging
    // ... test code
}
```

## Integration Points & External Dependencies

### Third-Party Libraries (go.mod)
- **gorilla/mux** v1.8.1: HTTP router - chosen for pattern matching & vars extraction
- **gopsutil/v3** v3.24.x: System metrics for `/healthz` (CPU%, memory, goroutines)

### Kubernetes/Helm Chart (k8s/chart/)
- **Default config**: 1 replica, ClusterIP service on port 3000
- **Probes**: `readinessProbe` & `livenessProbe` hit `/healthz` (initialDelaySeconds: 10)
- **HTTPRoute**: Optional Gateway API routing (disabled by default)
- **Deployment**: Uses `app.kubernetes.io/name=learn-go` label (NOT `app=learn-go`)

### Docker Build (Dockerfile)
```dockerfile
# Stage 1: golang:1.25-alpine → builds static binary (CGO_ENABLED=0)
# Stage 2: FROM scratch → minimal runtime (~8MB image)
EXPOSE 8080  # Container listens on 8080 (regardless of PORT env var)
```

### CI/CD Secrets Required
- **GitHub**: `FLUX_STAGING_RECEIVER_URL`, `FLUX_STAGING_WEBHOOK_SECRET`
- **GitHub**: `FLUX_PRODUCTION_RECEIVER_URL`, `FLUX_PRODUCTION_WEBHOOK_SECRET`
- **Auto-available**: `GITHUB_TOKEN` (for GHCR push)

## Common Pitfalls & Important Notes

1. **Test isolation**: ALWAYS set `os.Setenv("GO_ENV", "test")` in test setup - prevents middleware log spam
2. **Handler init**: `NewHandlers()` reads env vars at startup (cached) - changes require restart
3. **Middleware order**: Logging MUST be first to capture CORS preflight responses
4. **Docker FROM scratch**: Binary must be static (`CGO_ENABLED=0`) - no libc or shell available
5. **Kind testing**: Image loaded with `--set image.pullPolicy=Never` - doesn't pull from registry
6. **Pod labels**: Use `app.kubernetes.io/name=learn-go` NOT `app=learn-go` in kubectl selectors
7. **Workflow artifact**: Image built ONCE, shared via tar - don't rebuild in test-deployment job

## Quick Reference: Key Files & Commands

### Essential Files
- **Entry**: `cmd/api/main.go` (35 lines) - just calls `server.NewServer()` and `Start()`
- **Routes**: `internal/router/router.go` - all 6 endpoints + middleware chain
- **Handlers**: `internal/handlers/handlers.go` (250+ lines) - all business logic
- **CI/CD**: `.github/workflows/full-workflow.yml` (unified pipeline)
- **Tests**: `scripts/{smoke,e2e,integration}-test.sh` - bash testing scripts
- **Helm**: `k8s/chart/values.yaml` - deployment configuration

### Daily Commands
```bash
# Development
make dev                           # Start server (port 8080)
GO_ENV=test go test ./... -v      # Run all tests with output

# Docker local testing
docker build -t learn-go . && docker run -p 8080:8080 learn-go

# Integration testing (auto-detects env)
./scripts/integration-test.sh     # Local: tests Go + Docker | K8s: tests cluster

# Makefile targets
make help                          # Show all available targets
make test-coverage                 # Generate coverage.html
```

### Debugging Kind Cluster Issues
```bash
# Check if cluster exists
kind get clusters

# Load image manually
kind load docker-image learn-go:test --name test-cluster

# Check pods in Kind
kubectl get pods -l app.kubernetes.io/name=learn-go
kubectl logs -l app.kubernetes.io/name=learn-go

# Port-forward for local testing
kubectl port-forward service/learn-go 8080:3000
```
````
