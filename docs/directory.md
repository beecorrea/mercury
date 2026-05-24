# Directory Structure

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
