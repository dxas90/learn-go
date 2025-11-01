# learn-go - Complete Documentation

## 📚 Overview

This document provides comprehensive documentation for the learn-go microservice, including architecture, API reference, testing, and development guidelines.

## 🏗️ Architecture

The project follows Go best practices with a clean architecture:

```
learn-go/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware (CORS, logging, security)
│   ├── router/           # Route definitions
│   └── server/           # HTTP server setup
├── pkg/models/           # Shared data models
├── configs/              # Configuration files
├── api/                  # API specifications (OpenAPI)
└── bin/                  # Compiled binaries
```

### Design Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Dependency Injection**: Handlers receive dependencies via constructors
3. **Interface-Based Design**: Facilitates testing and mocking
4. **Error Handling**: Consistent error responses across all endpoints
5. **Middleware Pattern**: Cross-cutting concerns handled uniformly

## 🔌 API Reference

### Base URL
```
http://localhost:8080
```

### Endpoints

#### 1. Index - `GET /`
**Description**: Welcome page with API documentation

**Response**:
```json
{
  "success": true,
  "data": {
    "message": "Welcome to learn-go API",
    "description": "A simple Go microservice for learning Kubernetes and Docker",
    "documentation": {
      "swagger": null,
      "postman": null
    },
    "links": {
      "repository": "https://github.com/dxas90/learn-go",
      "issues": "https://github.com/dxas90/learn-go/issues"
    },
    "endpoints": [
      {"path": "/", "method": "GET", "description": "Welcome page"},
      {"path": "/ping", "method": "GET", "description": "Simple health check"},
      {"path": "/healthz", "method": "GET", "description": "Detailed health check"},
      {"path": "/info", "method": "GET", "description": "Application information"},
      {"path": "/version", "method": "GET", "description": "Application version"},
      {"path": "/echo", "method": "POST", "description": "Echo back the request body"}
    ]
  },
  "timestamp": "2025-11-01T09:00:00Z"
}
```

#### 2. Ping - `GET /ping`
**Description**: Simple ping-pong health check

**Response**:
```
pong
```

**Headers**:
- `Content-Type: text/plain`

#### 3. Health Check - `GET /healthz`
**Description**: Detailed health check with system metrics

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "uptime": 123.456,
    "timestamp": "2025-11-01T09:00:00Z",
    "memory": {
      "rss_mb": 15.234,
      "vms_mb": 1234.567,
      "used_percent": 45.67
    },
    "version": "1.0.0",
    "environment": "development"
  },
  "timestamp": "2025-11-01T09:00:00Z"
}
```

**Metrics Included**:
- `uptime`: Server uptime in seconds
- `memory.rss_mb`: Resident Set Size in megabytes
- `memory.vms_mb`: Virtual Memory Size in megabytes
- `memory.used_percent`: System memory usage percentage

#### 4. Application Info - `GET /info`
**Description**: Comprehensive system and runtime information

**Response**:
```json
{
  "success": true,
  "data": {
    "app": {
      "name": "learn-go",
      "version": "1.0.0",
      "environment": "development"
    },
    "system": {
      "os": "linux",
      "arch": "amd64",
      "cpu_count": 8,
      "cpu_percent": [12.5],
      "memory_total_mb": 16384.0,
      "memory_used_mb": 8192.0,
      "memory_percent": 50.0
    },
    "runtime": {
      "go_version": "go1.2",
      "goroutines": 10,
      "uptime": 123.456
    },
    "process": {
      "pid": 12345,
      "memory_rss_mb": 15.234,
      "memory_vms_mb": 1234.567
    }
  },
  "timestamp": "2025-11-01T09:00:00Z"
}
```

#### 5. Version - `GET /version`
**Description**: Application version information

**Response**:
```json
{
  "success": true,
  "data": {
    "version": "1.0.0",
    "name": "learn-go",
    "environment": "development"
  },
  "timestamp": "2025-11-01T09:00:00Z"
}
```

#### 6. Echo - `POST /echo`
**Description**: Echo back the request body with metadata

**Request**:
```json
{
  "message": "hello",
  "value": 123
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "echo": {
      "message": "hello",
      "value": 123
    },
    "headers": {
      "Content-Type": "application/json",
      "User-Agent": "curl/7.68.0"
    },
    "method": "POST"
  },
  "timestamp": "2025-11-01T09:00:00Z"
}
```

**Error Response** (Invalid JSON):
```json
{
  "error": true,
  "message": "Invalid JSON",
  "statusCode": 400,
  "timestamp": "2025-11-01T09:00:00Z"
}
```

## 🧪 Testing

### Test Coverage

Current test coverage: **79.5%**

| Package | Coverage |
|---------|----------|
| `internal/handlers` | 97.8% |
| `internal/middleware` | 92.3% |
| `internal/router` | 86.7% |
| `internal/server` | 70.0% |
| `cmd/api` | 0.0% (main function, not typically tested) |

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests with detailed coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run tests verbosely
go test ./... -v

# Run specific package tests
go test ./internal/handlers -v
```

### Test Structure

Each package has corresponding test files:
- `handlers_test.go` - Tests for all HTTP handlers
- `middleware_test.go` - Tests for middleware functions
- `router_test.go` - Tests for route setup
- `server_test.go` - Tests for server initialization

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PORT` | Server port | `8080` | `3000` |
| `HOST` | Server host | `127.0.0.1` | `0.0.0.0` |
| `GO_ENV` | Environment | `development` | `production` |
| `APP_VERSION` | Application version | `0.0.1` | `1.0.0` |
| `CORS_ORIGIN` | CORS allowed origin | `*` | `https://example.com` |

### Setting Environment Variables

**Linux/macOS**:
```bash
export PORT=8080
export APP_VERSION=1.0.0
export GO_ENV=production
```

**Windows (PowerShell)**:
```powershell
$env:PORT="8080"
$env:APP_VERSION="1.0.0"
$env:GO_ENV="production"
```

**Docker**:
```bash
docker run -e PORT=8080 -e APP_VERSION=1.0.0 learn-go
```

## 🚀 Development

### Prerequisites

- Go 1.2 or higher
- Make (optional)
- Docker (optional)

### Build

```bash
# Build binary
go build -o bin/learn-go ./cmd/api

# Build with Makefile
make build

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/learn-go-linux ./cmd/api
GOOS=windows GOARCH=amd64 go build -o bin/learn-go.exe ./cmd/api
GOOS=darwin GOARCH=arm64 go build -o bin/learn-go-mac ./cmd/api
```

### Run

```bash
# Run directly
go run ./cmd/api

# Run built binary
./bin/learn-go

# Run with custom configuration
PORT=3000 APP_VERSION=2.0.0 ./bin/learn-go

# Run with script
./scripts/run-local.sh
```

### Docker

```bash
# Build image
docker build -t learn-go:latest .

# Run container
docker run -p 8080:8080 learn-go:latest

# Run with environment variables
docker run -p 8080:8080 \
  -e APP_VERSION=1.0.0 \
  -e GO_ENV=production \
  learn-go:latest

# Run detached
docker run -d -p 8080:8080 --name learn-go-server learn-go:latest

# View logs
docker logs learn-go-server

# Stop container
docker stop learn-go-server
```

## 🔒 Security Features

### Middleware

1. **CORS Middleware**
   - Configurable allowed origins
   - Supports preflight requests
   - Headers: `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`

2. **Security Headers Middleware**
   - `X-Frame-Options: DENY` - Prevents clickjacking
   - `X-Content-Type-Options: nosniff` - Prevents MIME sniffing
   - `Content-Security-Policy: default-src 'self'` - Restricts resource loading
   - `X-XSS-Protection: 1; mode=block` - XSS protection (legacy browsers)

3. **Logging Middleware**
   - Logs all incoming requests
   - Format: `[timestamp] METHOD path - User-Agent: agent`
   - Helps with debugging and auditing

## 📊 Monitoring

### Health Check Endpoint

The `/healthz` endpoint provides:
- Service status
- Uptime
- Memory usage (RSS, VMS, percentage)
- Timestamp

**Example Integration**:
```bash
# Kubernetes liveness probe
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 10

# Docker healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/ping || exit 1
```

### Metrics

The `/info` endpoint provides comprehensive metrics:
- System info (OS, architecture)
- CPU usage and count
- Memory usage (total, used, percentage)
- Go runtime info (version, goroutines)
- Process info (PID, memory)

## 🐛 Debugging

### Common Issues

1. **Port already in use**
   ```bash
   # Find process using port
   lsof -i :8080
   # Or change port
   PORT=3000 ./bin/learn-go
   ```

2. **Module errors**
   ```bash
   # Update dependencies
   go mod tidy
   go mod download
   ```

3. **Build errors**
   ```bash
   # Clean build cache
   go clean -cache
   go build -a ./cmd/api
   ```

### Logging

The application uses Go's standard `log` package. To see detailed logs:

```bash
# Run with verbose output
./bin/learn-go

# Redirect logs to file
./bin/learn-go > app.log 2>&1

# View logs in real-time
tail -f app.log
```

## 📝 Code Documentation

All exported functions and types are documented with GoDoc comments. To view:

```bash
# View package documentation
go doc github.com/dxas90/learn-go/internal/handlers

# View specific function
go doc github.com/dxas90/learn-go/internal/handlers.NewHandlers

# Generate HTML documentation
godoc -http=:6060
# Then visit http://localhost:6060/pkg/github.com/dxas90/learn-go/
```

## 🤝 Contributing

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Run `go vet` for static analysis
- Add tests for new functionality

### Commit Messages

Follow conventional commits:
```
feat: add new endpoint for user management
fix: correct memory leak in handler
docs: update API documentation
test: add tests for middleware
refactor: simplify router logic
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Commit your changes
6. Push to the branch
7. Open a Pull Request

## 📄 License

MIT License - See LICENSE file for details

## 🙏 Acknowledgments

- Inspired by the [learn-python](https://github.com/dxas90/learn-python) project
- Uses [gorilla/mux](https://github.com/gorilla/mux) for routing
- Uses [gopsutil](https://github.com/shirou/gopsutil) for system metrics
