# syntax=docker/dockerfile:1

# Build stage using a Debian-based Go image
FROM golang:1.25-bookworm AS builder

# Set the working directory
WORKDIR /workspace

# Copy dependency files first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies using a persistent BuildKit cache mount.
# The module cache is preserved across builds, so deps are only downloaded once.
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the rest of the application source code
COPY . .

# Build the application as a statically linked binary.
# Two cache mounts are used:
#   - /go/pkg/mod: reuses the already-downloaded module cache
#   - /root/.cache/go-build: reuses the Go build cache, skipping unchanged packages
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mercury .

# Final stage using a minimal, secure distroless image
FROM gcr.io/distroless/static-debian12

# Copy the statically compiled binary from the builder stage
COPY --from=builder /workspace/mercury /mercury

# Expose the port that the Fiber server listens on
EXPOSE 45800

# Set the entrypoint to run the service
ENTRYPOINT ["/mercury"]
