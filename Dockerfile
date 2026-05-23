# Build stage using a Debian-based Go image
FROM golang:1.25-bookworm AS builder

# Set the working directory
WORKDIR /workspace

# Copy dependency files first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application as a statically linked binary
# - CGO_ENABLED=0 ensures no dynamic linking to glibc
# - ldflags "-s -w" shrinks the binary size by removing symbols and debug info
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mercury .

# Final stage using a minimal, secure distroless image
FROM gcr.io/distroless/static-debian12

# Copy the statically compiled binary from the builder stage
COPY --from=builder /workspace/mercury /mercury

# Expose the port that the Fiber server listens on
EXPOSE 45800

# Set the entrypoint to run the service
ENTRYPOINT ["/mercury"]
