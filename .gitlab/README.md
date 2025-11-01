# GitLab CI/CD Pipeline Documentation

## Overview

This directory contains configuration files for the GitLab CI/CD pipeline.

## Pipeline Configuration

The main pipeline configuration is located in `.gitlab-ci.yml` at the project root.

### Pipeline Variables

The following CI/CD variables should be configured in GitLab:

#### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CI_REGISTRY_USER` | Container registry username | `your-username` |
| `CI_REGISTRY_PASSWORD` | Container registry password/token | `your-token` |
| `CI_REGISTRY` | Container registry URL | `registry.gitlab.com` |
| `CI_REGISTRY_IMAGE` | Full image path | `registry.gitlab.com/dxas90/learn-go` |

#### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `KUBECTL_VERSION` | Kubectl version | `v1.28.4` |
| `KIND_VERSION` | Kind version | `v0.20.0` |
| `HELM_VERSION` | Helm version | `v3.13.0` |
| `TRIVY_VERSION` | Trivy scanner version | `0.47.0` |

## Pipeline Stages

### 1. Lint (lint)

**Image**: `golang:1.25-alpine`

**Purpose**: Ensures code quality and consistency

**Steps**:
- Install linting tools (golint, staticcheck)
- Run `go fmt` to check formatting
- Run `go vet` for static analysis
- Run `golint` for Go-specific linting
- Run `staticcheck` for advanced checks
- Verify module dependencies

**Triggers**: merge_requests, main, develop branches

**Duration**: ~1-2 minutes

### 2. Test (test)

**Image**: `golang:1.25-alpine`

**Purpose**: Run unit and integration tests with coverage

**Steps**:
- Download dependencies
- Run tests with coverage (`go test ./... -coverprofile=coverage.out`)
- Build the application binary
- Start server and perform health checks
- Generate coverage reports

**Coverage Target**: >79%

**Artifacts**:
- `coverage.out` - Coverage data
- `bin/` - Built binaries
- Coverage reports (Cobertura format)

**Triggers**: merge_requests, main, develop branches

**Duration**: ~2-3 minutes

### 3. Security Scan (security-scan)

**Image**: `aquasec/trivy:latest`

**Purpose**: Scan for security vulnerabilities

**Steps**:
- Build Docker image
- Scan filesystem for vulnerabilities
- Scan Docker image for vulnerabilities
- Fail on HIGH or CRITICAL vulnerabilities

**Artifacts**:
- `gl-dependency-scanning-report.json`
- `gl-container-scanning-report.json`

**Triggers**: merge_requests, main, develop branches

**Duration**: ~3-5 minutes

### 4. Build (build)

**Image**: `docker:28`

**Services**: `docker:28-dind`

**Purpose**: Build and push Docker images

**Steps**:
- Login to container registry
- Set version (Git tag or commit SHA)
- Build Docker image with build args
- Tag with version and latest
- Push to registry

**Build Arguments**:
- `BUILD_DATE` - Build timestamp
- `VCS_REF` - Git commit SHA
- `VERSION` - Application version

**Triggers**: main branch, tags only

**Duration**: ~3-5 minutes

### 5. Deploy Staging (deploy-staging)

**Image**: `alpine/k8s:1.34.1`

**Purpose**: Deploy to staging Kubernetes environment

**Environment**:
- Name: staging
- URL: https://staging.learn-go.example.com

**Steps**:
- Install Kind (Kubernetes in Docker)
- Create local Kubernetes cluster
- Deploy application with 2 replicas
- Configure health probes
- Expose service as LoadBalancer
- Wait for rollout completion

**Configuration**:
```yaml
Replicas: 2
Container Port: 8080
Environment:
  - PORT: 8080
  - HOST: 0.0.0.0
  - GO_ENV: staging
  - APP_VERSION: <commit-sha>

Probes:
  Liveness: GET /ping (every 10s)
  Readiness: GET /healthz (every 5s)
```

**Triggers**: main branch (manual)

**Duration**: ~5-8 minutes

### 6. Deploy Production (deploy-production)

**Image**: `alpine/k8s:1.34.1`

**Purpose**: Deploy to production Kubernetes cluster

**Environment**:
- Name: production
- URL: https://learn-go.example.com

**Configuration**: Currently a placeholder for actual production deployment

**Triggers**: tags only (manual)

**Duration**: Depends on implementation

### 7. Cleanup (cleanup)

**Image**: `alpine/k8s:1.34.1`

**Purpose**: Clean up resources after pipeline completion

**Steps**:
- Prune Docker system
- Delete Kind staging cluster

**Triggers**: Always runs (even on failure)

**Allow Failure**: Yes

**Duration**: ~1 minute

## Pipeline Flow

### For Merge Requests

```
lint → test → security-scan
```

### For Main Branch

```
lint → test → security-scan → build → deploy-staging (manual)
```

### For Tags

```
lint → test → security-scan → build → deploy-production (manual)
```

## Setting Up GitLab CI

### 1. Configure Variables

Go to **Settings > CI/CD > Variables** and add:

```bash
CI_REGISTRY_USER=<your-username>
CI_REGISTRY_PASSWORD=<your-token>
CI_REGISTRY=registry.gitlab.com
CI_REGISTRY_IMAGE=registry.gitlab.com/<your-namespace>/learn-go
```

### 2. Enable Container Registry

Ensure the GitLab Container Registry is enabled for your project.

### 3. Enable Shared Runners

Go to **Settings > CI/CD > Runners** and enable shared runners.

### 4. Configure Protected Branches

- **main**: Protected, allow maintainers to push
- **develop**: Protected, allow developers to push

### 5. Configure Protected Tags

Set tag protection pattern: `v*` (for version tags like v1.0.0)

## Local Testing

### Test Lint Stage

```bash
# Install tools
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Run checks
go fmt ./...
go vet ./...
golint ./...
staticcheck ./...
go mod verify
```

### Test Build Stage

```bash
# Build image
docker build -t learn-go:local .

# Run container
docker run -p 8080:8080 learn-go:local
```

### Test Security Scan

```bash
# Install Trivy
wget https://github.com/aquasecurity/trivy/releases/download/v0.47.0/trivy_0.47.0_Linux-64bit.tar.gz
tar zxvf trivy_0.47.0_Linux-64bit.tar.gz
sudo mv trivy /usr/local/bin/

# Scan
trivy image learn-go:local
```

### Test Kubernetes Deployment

```bash
# Install Kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Create cluster
kind create cluster --name test

# Deploy (requires manifests)
kubectl apply -f k8s/

# Cleanup
kind delete cluster --name test
```

## Troubleshooting

### Pipeline Fails on Lint

**Issue**: Formatting or linting errors

**Solution**:
```bash
# Fix formatting
go fmt ./...

# Check for issues
go vet ./...
golint ./...
```

### Pipeline Fails on Tests

**Issue**: Tests failing or coverage too low

**Solution**:
```bash
# Run tests locally
go test ./... -v

# Check coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Pipeline Fails on Security Scan

**Issue**: HIGH or CRITICAL vulnerabilities found

**Solution**:
- Update dependencies: `go get -u ./...`
- Check for security advisories
- Update base Docker image

### Build Stage Fails

**Issue**: Docker login or build errors

**Solution**:
- Verify CI variables are set correctly
- Check Dockerfile syntax
- Ensure base images are accessible

### Deploy Fails

**Issue**: Kubernetes deployment errors

**Solution**:
- Check Kind cluster creation logs
- Verify image is pushed to registry
- Check deployment manifests
- Review kubectl logs

## Best Practices

1. **Always run tests locally** before pushing
2. **Keep dependencies updated** regularly
3. **Monitor security scan results** and fix vulnerabilities promptly
4. **Use semantic versioning** for tags (v1.0.0, v1.1.0, etc.)
5. **Test staging deployment** before production
6. **Review pipeline logs** for warnings
7. **Keep Docker images small** using multi-stage builds
8. **Document any manual steps** required for deployment

## References

- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Trivy Documentation](https://aquasecurity.github.io/trivy/)
