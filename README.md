# Mercury

<p align="center">
  <img src="resources/logo.png" width="280" alt="Mercury Logo" />
</p>

<p align="center">
  A high-performance, self-hosted shortlink manager for developers.<br />
  Fast, secure, and customizable redirections with an Apple-inspired dark dashboard.
</p>

<p align="center">
  <video src="resources/demo_features.mp4" width="800" controls alt="Mercury Feature Demo"></video>
</p>

---

## Features

- ⚡️ **High-Performance:** Sub-millisecond redirection engine built in Go with Fiber.
- 🎨 **Modern Dashboard:** High-fidelity, minimalist dark UI with instant creation, filtering, and live state updates.
- 📦 **Docker Ready:** Tiny, secure distroless containers with persistent volume mounts.
- 💾 **SQLite Backed:** Local database persistence with automatic unique key constraint enforcement.

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

## Documentation & Reference

* **[Architecture](docs/architecture.md):** Detailed breakdown of system design, Go modules, and repository architecture.
* **[Directory Structure](docs/directory.md):** Full repository map.
* **[API Reference](docs/api.md):** Request and response schemas for all Fiber HTTP endpoints.
