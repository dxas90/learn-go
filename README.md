# learn-go

[![Build and Test](https://github.com/dxas90/learn-go/workflows/Docker%20Build%20and%20Security%20Scan/badge.svg)](https://github.com/dxas90/learn-go/actions)
[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A simple Go microservice for learning Kubernetes, Docker, and modern Go development practices.

## 🚀 Features

- **RESTful API** with multiple endpoints
- **Health checks** and monitoring endpoints
- **CORS support** for cross-origin requests
- **Security headers** (X-Frame-Options, CSP, etc.)
- **Docker support** with multi-stage builds
- **Kubernetes ready** with deployment configurations
- **CI/CD pipelines** (GitLab CI, GitHub Actions)
- **Comprehensive testing** with Go testing
- **Production-ready** with proper logging and error handling

## 📋 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Welcome page with API documentation |
| `/ping` | GET | Simple ping-pong health check |
| `/healthz` | GET | Detailed health check with system metrics |
| `/info` | GET | Application and system information |
| `/version` | GET | Application version information |
| `/echo` | POST | Echo back the request body |

## 🛠️ Quick Start

### Prerequisites

- Go 1.25 or higher
- Docker (optional)
- make (optional, for using Makefile commands)

### Local Development

1. **Clone the repository**
   ```sh
   git clone https://github.com/dxas90/learn-go.git
   cd learn-go
   ```

2. **Install dependencies**
   ```sh
   go mod tidy
   ```

3. **Run the application**
   ```sh
   go run ./cmd/api
   ```

4. **Test the API**
   ```sh
   curl http://localhost:8080/
   curl http://localhost:8080/ping
   curl http://localhost:8080/healthz
   ```

### Running Tests

```sh
go test ./...
```

### Docker Deployment

```sh
# Build the image
docker build -t learn-go .

# Run the container
docker run -p 8080:8080 learn-go
```

## 📊 Environment Variables

- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: 127.0.0.1)
- `GO_ENV`: Environment (development/production/test)
- `APP_VERSION`: Application version (default: 0.0.1)
- `CORS_ORIGIN`: CORS allowed origin (default: *)

## 🏗️ Project Structure

```
learn-go/
├── cmd/
│   └── api/
│       └── main.go                  # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── handlers.go              # HTTP request handler implementations
│   │   └── handlers_test.go         # Handler unit tests
│   ├── middleware/
│   │   └── middleware.go            # CORS, logging, security middleware
│   ├── router/
│   │   └── router.go                # Route definitions and setup
│   └── server/
│       └── server.go                # HTTP server configuration
├── pkg/
│   └── models/
│       └── responses.go             # API response data models
├── configs/
│   └── config.yaml                  # Application configuration
├── api/
│   └── openapi.yaml                 # OpenAPI/Swagger specification
├── scripts/
│   └── run-local.sh                 # Local development startup script
├── Dockerfile                       # Multi-stage Docker build
├── go.mod                           # Go module dependencies
├── go.sum                           # Go dependencies checksums
├── Makefile                         # Build automation targets
├── .gitlab-ci.yml                   # GitLab CI/CD pipeline
├── .github/workflows/               # GitHub Actions workflows
│   ├── dockerimage.yml              # Docker build and security scan
│   └── k8s-deployment.yml           # Kubernetes deployment pipeline
├── .gitignore                       # Git ignore rules
└── README.md                        # Project documentation
```

## 🔄 CI/CD Pipelines

The project includes comprehensive CI/CD pipelines for both GitLab CI and GitHub Actions.

### GitHub Actions Workflows

#### 1. Docker Build and Security Scan (`.github/workflows/dockerimage.yml`)

**Triggers**: Push to main/develop, Pull requests, Tags

**Jobs**:
- **Lint**: Go code quality checks (fmt, vet, golint, staticcheck)
- **Test**: Run tests with coverage on Go 1.20, 1.21, 1.22
- **Security Scan**: govulncheck and gosec security scanning
- **Build**: Multi-arch Docker image build (amd64/arm64) and push to GHCR
- **Deploy Staging**: Deploy to Kind cluster on main branch
- **Deploy Production**: Deploy on version tags

#### 2. Kubernetes Deployment (`.github/workflows/k8s-deployment.yml`)

**Triggers**: Push to main, Pull requests, Tags

**Jobs**:
- **Validate K8s**: Validate Kubernetes manifests with kubeval
- **Test Deployment**: Deploy to Kind cluster and test endpoints

### GitLab CI Pipeline (`.gitlab-ci.yml`)

**Pipeline Stages**:

1. **Lint** - Code quality checks
   - `go fmt` - Format checking
   - `go vet` - Static analysis
   - `golint` - Linting
   - `staticcheck` - Advanced static analysis
   - `go mod verify` - Dependency verification

2. **Test** - Unit and integration tests
   - Run all tests with coverage
   - Build verification
   - Health check validation
   - Coverage reporting (target: >79%)

3. **Security** - Security scanning
   - Filesystem scanning with Trivy
   - Docker image scanning
   - Vulnerability detection (HIGH/CRITICAL fails the build)
   - Dependency scanning reports

4. **Build** - Docker image building
   - Multi-stage Docker builds
   - Version tagging (commit SHA or Git tag)
   - Push to container registry

5. **Deploy Staging** - Staging deployment
   - Kind (Kubernetes in Docker) cluster
   - Automated deployment with 2 replicas
   - Health probes (liveness & readiness)
   - Manual trigger required

6. **Deploy Production** - Production deployment
   - Manual trigger required
   - Tag-based deployments only
   - Production environment configuration

7. **Cleanup** - Resource cleanup
   - Automatic cleanup on pipeline completion
   - Docker system pruning
   - Kind cluster removal

### Running CI Checks Locally

To simulate the CI pipeline locally:

```bash
# Lint
go fmt ./...
go vet ./...
go mod verify

# Test
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out

# Build
docker build -t learn-go:local .

# Run security scan (requires Trivy)
trivy image learn-go:local
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Commit your changes
7. Push to the branch
8. Submit a pull request

## 📄 License

This project is open source and available under the MIT License.
