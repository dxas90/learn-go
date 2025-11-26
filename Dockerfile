# syntax=docker/dockerfile:1

# Multi-platform build optimized with cross-compilation
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
COPY . /build/
ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod tidy && go mod vendor

# Cross-compile for target platform (fast on any builder platform)
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o main ./cmd/api

FROM scratch AS production
ARG CREATED="0000-00-00T00:00:00Z"
LABEL org.opencontainers.image.authors="Daniel Ramirez <dxas90@gmail.com>" \
    org.opencontainers.image.created=${CREATED} \
    org.opencontainers.image.description="A container image to learn." \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.source="https://github.com/dxas90/learn-go" \
    org.opencontainers.image.title="learn Image" \
    org.opencontainers.image.version="1.0.0"
COPY --from=builder /build/main /app/
EXPOSE 8080
WORKDIR /app
ENTRYPOINT [ "/app/main" ]
