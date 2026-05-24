# Mercury Shortlink Service

Mercury is a high-performance, self-hosted shortlink manager built in Go using the Fiber web framework and SQLite. It provides a web dashboard and API endpoints for creating, deleting, and listing shortlinks, alongside a high-speed redirection engine that intercepts subdomain-specific traffic.

---

## Local Development & Installation

### Prerequisites
* Go 1.25 or higher

### Running the Application
1. Download Go dependencies:
   ```bash
   go mod download
   ```
2. Run the application:
   ```bash
   go run main.go
   ```
3. Open the browser:
   * Dashboard: `http://localhost:45800`
   * Target shortlink redirects occur when hostnames start with `mercury.`, such as `http://mercury.localhost:45800/key`.

### Running with Docker

You can build and run the application in a secure, minimal container. To persist the SQLite database across container restarts and recreations, mount a host volume directory and specify the database path via the `MERCURY_DB_PATH` environment variable:

1. Build the Docker image:
   ```bash
   docker build -t mercury .
   ```

2. Run the container with a volume mount:
   ```bash
   docker run -d \
     -p 45800:45800 \
     -v /path/to/local/data:/data \
     -e MERCURY_DB_PATH=/data/mercury.db \
     --name mercury \
     mercury
   ```

### Running with Docker Compose

A [docker-compose.yml](file:///Users/bianca/projects/mercury/docker-compose.yml) is provided in the repository root. This orchestrates the build process, exposes the required port, sets the `MERCURY_DB_PATH` environment variable, and configures a persistent named volume `mercury_data` to ensure the SQLite database persists across container restarts and recreations.

To start the service using Docker Compose:
```bash
docker compose up -d --build
```

To stop the service:
```bash
docker compose down
```

### Running Tests
Unit and integration test coverage can be verified by running:
```bash
go test -v ./...
```

---

## Feature Demo

Watch a short demonstration video showcasing the core features (adding a shortlink, redirecting via a short subdomain URL, and deleting active redirection keys):

![Mercury Features Demo](resources/demo_features.webm)

---

## Documentation & Reference

### Architecture
For a detailed breakdown of the system architecture, core components, and standardized constructor patterns, please refer to the [Architecture Documentation](docs/architecture.md).

### Directory Structure
For the full project directory layout, see [docs/directory.md](docs/directory.md).

### API Reference
For all endpoint specs and request/response payloads, see [docs/api.md](docs/api.md).
