# Mercury Shortlink Service

Mercury is a high-performance, self-hosted shortlink manager built in Go using the Fiber web framework and SQLite. It provides a web dashboard and API endpoints for creating, deleting, and listing shortlinks, alongside a high-speed redirection engine that intercepts subdomain-specific traffic.

---

## Architecture

For a detailed breakdown of the system architecture, core components, and standardized constructor patterns, please refer to the [Architecture Documentation](docs/architecture.md).

---

## Directory Structure

```
.
├── README.md                      # Documentation
├── go.mod                         # Go module definition
├── go.sum                         # Go dependency checksums
├── main.go                        # Application entry point
├── mercury.db                     # SQLite database file
└── internal
    ├── controllers
    │   ├── dashboard.go           # Dashboard view controller
    │   └── shortlink.go           # Shortlink API controller
    ├── database
    │   ├── queries.go             # Embedded query registry
    │   ├── redirect.go            # SQLite client & schema migration
    │   ├── redirect_test.go       # Database unit tests
    │   └── sql
    │       ├── delete_link.sql    # Delete query
    │       ├── get_link.sql       # Get shortlink query
    │       ├── insert_link.sql    # Insert shortlink query
    │       ├── list_links.sql     # List all shortlinks query
    │       └── migrate.sql        # Database initialization schema
    ├── middleware
    │   └── redirect.go            # Hostname-based routing middleware
    ├── server
    │   ├── routes.go              # Route registrations
    │   ├── server.go              # Server configuration and setup
    │   └── server_test.go         # Integration test suite for HTTP server
    ├── service
    │   ├── redirect.go            # Redirection service
    │   ├── redirect_test.go       # Redirect service test cases
    │   ├── shortlink.go           # Management service
    │   └── shortlink_test.go      # Shortlink service test cases
    ├── structs
    │   └── shortlink.go           # Core Domain struct
    └── templates
        ├── 404.html               # Not Found embedded static template
        ├── index.html             # Dashboard frontend embedded template
        └── templates.go           # Embedded templates registry
```

---

## API Reference

### Shortlinks

#### List Shortlinks
* **Endpoint:** `GET /api/links`
* **Response Status:** `200 OK`
* **Response Body:**
```json
[
  {
    "key": "github",
    "url": "https://github.com",
    "domain": "localhost:45800",
    "created_at": "2026-05-23T19:51:13Z"
  }
]
```

#### Create Shortlink
* **Endpoint:** `POST /api/shorten`
* **Request Body:**
```json
{
  "key": "google",
  "url": "https://google.com",
  "domain": "localhost:45800"
}
```
* **Response Status:** `201 Created`
* **Response Body:**
```json
{
  "message": "shortlink created successfully"
}
```

#### Delete Shortlink
* **Endpoint:** `DELETE /api/links/:key`
* **Response Status:** `200 OK`
* **Response Body:**
```json
{
  "message": "shortlink deleted successfully"
}
```

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
