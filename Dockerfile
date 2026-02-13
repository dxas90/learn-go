# syntax=docker/dockerfile:1

# Multi-platform build optimized with cross-compilation
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine AS builder

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
COPY . /build/
ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod tidy && go mod vendor

# Copy OpenAPI spec for embedding
RUN mkdir -p internal/apispec && cp api/openapi.yaml internal/apispec/openapi.yaml

# Cross-compile for target platform (fast on any builder platform)
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o main ./cmd/api

FROM alpine:3.23 AS production
ARG CREATED="0000-00-00T00:00:00Z"

# Install ca-certificates for HTTPS and curl for health checks
RUN apk --no-cache add ca-certificates curl

LABEL org.opencontainers.image.authors="Daniel Ramirez <dxas90@gmail.com>" \
    org.opencontainers.image.created=${CREATED} \
    org.opencontainers.image.description="A container image to learn." \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.source="https://github.com/dxas90/learn-go" \
    org.opencontainers.image.title="learn Image" \
    org.opencontainers.image.version="1.0.0"

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001 -G appuser

WORKDIR /app
COPY --from=builder /build/main /app/

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Health check using curl
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl http://127.0.0.1:8080/healthz || exit 1

EXPOSE 8080
ENTRYPOINT [ "/app/main" ]
