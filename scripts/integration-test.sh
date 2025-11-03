#!/bin/bash
set -euo pipefail

# Integration test script for the learn-go application
# This script can be run locally or in CI to validate application integration

# Configuration
APP_NAME="learn-go"
SERVICE_NAME="learn-go"
NAMESPACE="default"
LOCAL_PORT=8080
TIMEOUT_SECONDS=60

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running in Kubernetes or locally
check_environment() {
    if kubectl cluster-info >/dev/null 2>&1; then
        echo "k8s"
    else
        echo "local"
    fi
}

# Test local Go application
test_local_application() {
    log_info "Testing local Go application..."

    # Start the application in the background
    log_info "Starting application..."
    PORT=${LOCAL_PORT} GO_ENV=test go run ./cmd/api &
    APP_PID=$!

    # Wait for application to start
    log_info "Waiting for application to start..."
    sleep 5

    # Test endpoints
    test_endpoint "http://localhost:${LOCAL_PORT}/healthz" "Local Health endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/ping" "Local Ping endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/" "Local Root endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/info" "Local Info endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/version" "Local Version endpoint"

    # Test POST endpoint
    log_info "Testing Local Echo endpoint (POST)..."
    if curl -s -f -X POST -H "Content-Type: application/json" -d '{"message":"test"}' "http://localhost:${LOCAL_PORT}/echo" >/dev/null; then
        log_success "Local Echo endpoint test passed"
    else
        log_error "Local Echo endpoint test failed"
        kill $APP_PID || true
        exit 1
    fi

    # Kill the application
    kill $APP_PID || true
    wait $APP_PID 2>/dev/null || true

    log_success "Local application testing completed"
}

# Test Docker container
test_docker_container() {
    log_info "Testing Docker container..."

    # Build Docker image
    log_info "Building Docker image..."
    docker build -t learn-go-test .

    # Run container
    log_info "Starting Docker container..."
    docker run -d --name learn-go-test-container -p ${LOCAL_PORT}:8080 -e GO_ENV=test learn-go-test

    # Wait for container to start
    log_info "Waiting for container to start..."
    sleep 5

    # Test endpoints
    test_endpoint "http://localhost:${LOCAL_PORT}/healthz" "Docker Health endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/ping" "Docker Ping endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/" "Docker Root endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/info" "Docker Info endpoint"
    test_endpoint "http://localhost:${LOCAL_PORT}/version" "Docker Version endpoint"

    # Test POST endpoint
    log_info "Testing Docker Echo endpoint (POST)..."
    if curl -s -f -X POST -H "Content-Type: application/json" -d '{"message":"test"}' "http://localhost:${LOCAL_PORT}/echo" >/dev/null; then
        log_success "Docker Echo endpoint test passed"
    else
        log_error "Docker Echo endpoint test failed"
        docker stop learn-go-test-container || true
        docker rm learn-go-test-container || true
        docker rmi learn-go-test || true
        exit 1
    fi

    # Cleanup
    docker stop learn-go-test-container || true
    docker rm learn-go-test-container || true
    docker rmi learn-go-test || true

    log_success "Docker container testing completed"
}

# Test Kubernetes deployment
test_k8s_application() {
    log_info "Testing Kubernetes deployment..."

    # Check if service exists
    if ! kubectl get service ${SERVICE_NAME} -n ${NAMESPACE} >/dev/null 2>&1; then
        log_error "Service ${SERVICE_NAME} not found in namespace ${NAMESPACE}"
        exit 1
    fi

    # Port forward to access the service
    log_info "Setting up port forwarding..."
    kubectl port-forward -n ${NAMESPACE} service/${SERVICE_NAME} ${LOCAL_PORT}:8080 &
    PF_PID=$!
    sleep 5

    # Test endpoints
    test_endpoint "http://localhost:${LOCAL_PORT}/healthz" "Health endpoint (K8s)"
    test_endpoint "http://localhost:${LOCAL_PORT}/ping" "Ping endpoint (K8s)"
    test_endpoint "http://localhost:${LOCAL_PORT}/" "Root endpoint (K8s)"
    test_endpoint "http://localhost:${LOCAL_PORT}/info" "Info endpoint (K8s)"
    test_endpoint "http://localhost:${LOCAL_PORT}/version" "Version endpoint (K8s)"

    # Test POST endpoint
    log_info "Testing Echo endpoint (POST) (K8s)..."
    if curl -s -f -X POST -H "Content-Type: application/json" -d '{"message":"test"}' "http://localhost:${LOCAL_PORT}/echo" >/dev/null; then
        log_success "Echo endpoint (K8s) test passed"
    else
        log_error "Echo endpoint (K8s) test failed"
        kill $PF_PID || true
        exit 1
    fi

    # Kill port forwarding
    kill $PF_PID || true
    wait $PF_PID 2>/dev/null || true

    log_success "Kubernetes application testing completed"
}

# Generic endpoint testing function
test_endpoint() {
    local url=$1
    local description=$2
    local max_retries=5
    local retry_count=0

    log_info "Testing ${description}: ${url}"

    while [ $retry_count -lt $max_retries ]; do
        if curl -s -f "${url}" >/dev/null; then
            log_success "${description} test passed"
            return 0
        else
            retry_count=$((retry_count + 1))
            log_info "Retry ${retry_count}/${max_retries} for ${description}"
            sleep 2
        fi
    done

    log_error "${description} test failed after ${max_retries} retries"
    return 1
}

# Run unit tests
run_unit_tests() {
    log_info "Running unit tests..."

    if GO_ENV=test go test -v ./...; then
        log_success "Unit tests passed"
    else
        log_error "Unit tests failed"
        return 1
    fi
}

# Run integration tests based on environment
run_integration_tests() {
    local environment=$(check_environment)

    log_info "Running integration tests in ${environment} environment"

    case $environment in
        "k8s")
            test_k8s_application
            ;;
        "local")
            test_local_application
            test_docker_container
            ;;
        *)
            log_error "Unknown environment: ${environment}"
            return 1
            ;;
    esac
}

# Generate integration report
generate_integration_report() {
    log_info "Generating integration test report..."

    cat <<EOF
===========================================
     INTEGRATION TEST REPORT
===========================================

Environment: $(check_environment)

EOF

    if kubectl cluster-info >/dev/null 2>&1; then
        cat <<EOF
Kubernetes Context: $(kubectl config current-context)
Cluster: $(kubectl cluster-info | head -1)
Namespace: ${NAMESPACE}
Service: ${SERVICE_NAME}
EOF
    else
        cat <<EOF
Local environment
Go version: $(go version)
Docker version: $(docker --version)
EOF
    fi

    cat <<EOF

===========================================
EOF

    log_success "Integration testing completed successfully!"
}

# Main execution
main() {
    log_info "Starting integration testing for ${APP_NAME}..."

    # Check dependencies
    if ! command -v curl >/dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi

    if ! command -v go >/dev/null; then
        log_error "Go is required but not installed"
        exit 1
    fi

    # Run tests
    run_unit_tests
    run_integration_tests
    generate_integration_report

    log_success "All integration tests passed! 🎉"
}

# Cleanup function
cleanup() {
    # Kill any background processes
    jobs -p | xargs -r kill 2>/dev/null || true

    # Cleanup Docker containers
    docker stop learn-go-test-container 2>/dev/null || true
    docker rm learn-go-test-container 2>/dev/null || true
    docker rmi learn-go-test 2>/dev/null || true
}

# Set up trap for cleanup
trap cleanup EXIT

# Run main function
main "$@"
