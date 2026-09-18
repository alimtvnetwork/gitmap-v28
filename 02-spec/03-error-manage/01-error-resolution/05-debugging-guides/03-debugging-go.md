# Go Debugging Guide

> **Version:** 1.0.0
> **Created:** 2026-02-04
> **Applies To:** GSearch CLI, BRun CLI, AI Bridge CLI, Nexus Flow CLI, AI Transcribe CLI

---

## Overview

This guide covers debugging patterns for Go applications, particularly CLI tools with HTTP servers. It includes structured logging, error handling, health check verification, and common troubleshooting scenarios.

---

## Initialization Order (CRITICAL)

Go CLI services follow a strict initialization order to prevent runtime errors:

1. **STEP 1: Configuration** - Load environment variables and config files FIRST
2. **STEP 2: Directories** - Ensure all required directories exist
3. **STEP 3: Database** - Initialize database connections (Split DB pattern)
4. **STEP 4: Services** - Initialize business logic services
5. **STEP 5: HTTP Server** - Start the server ONLY AFTER all dependencies are ready

```go
func main() {
    // Step 1: Configuration
    cfg, err := config.Load()

    if err != nil {
        log.Fatal().Err(err).Msg("Failed to load configuration")
    }

    // Step 2: Directories
    if err := ensureDirectories(cfg); err != nil {
        log.Fatal().Err(err).Msg("Failed to ensure directories")
    }

    // Step 3: Database
    db, err := database.Connect(cfg.DatabasePath)

    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to database")
    }

    defer db.Close()

    // Step 4: Services
    svc := services.New(db, cfg)

    // Step 5: HTTP Server
    server := api.NewServer(svc)
    log.Info().Int("port", cfg.Port).Msg("Starting server")
    log.Fatal().Err(server.ListenAndServe()).Msg("Server stopped")
}
```

---

## Structured Logging with zerolog

### Logger Setup

```go
package logger

import (
    "os"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func Init(debug bool) {
    // Human-readable output for development
    if debug {
        log.Logger = log.Output(zerolog.ConsoleWriter{
            Out:        os.Stderr,
            TimeFormat: time.RFC3339,
        })
        zerolog.SetGlobalLevel(zerolog.DebugLevel)

        return
    }

    // JSON output for production
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
```

### Logging Levels

| Level | Usage | Example |
|-------|-------|---------|
| `Trace` | Extremely detailed, usually off | Loop iterations |
| `Debug` | Development information | Variable values, flow |
| `Info` | Normal operation | Server started, request completed |
| `Warn` | Recoverable issues | Retry attempts, deprecation |
| `Error` | Failures requiring attention | Database error, API failure |
| `Fatal` | Unrecoverable, exits program | Config missing, port in use |

### Logging Patterns

```go
// Simple message
log.Info().Msg("Server started")

// With fields
log.Info().
    Str("method", r.Method).
    Str("path", r.URL.Path).
    Int("status", status).
    Dur("duration", time.Since(start)).
    Msg("Request completed")

// With error
log.Error().
    Err(err).
    Str("operation", "database_query").
    Str("table", "settings").
    Msg("Failed to execute query")

// Contextual logger
reqLogger := log.With().
    Str("request_id", requestId).
    Str("user_id", userId).
    Logger()
reqLogger.Info().Msg("Processing request")
```

---

## Error Handling Pattern

### Standard Error Response

All Go backends MUST return the standard envelope:

```go
type Response[T any] struct {
    Success bool
    Data    T          `json:",omitempty"`
    Error   *ErrorInfo `json:",omitempty"`
}

type ErrorInfo struct {
    Code    int
    Message string
    Details string `json:",omitempty"`
}
```

### Error Response Helper

```go
func respondError(w http.ResponseWriter, code int, errCode int, message string, err error) {
    details := ""

    if err != nil {
        details = err.Error()
        log.Error().
            Err(err).
            Int("error_code", errCode).
            Str("message", message).
            Msg("API error response")
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(Response{
        Success: false,
        Error: &ErrorInfo{
            Code:    errCode,
            Message: message,
            Details: details,
        },
    })
}

func respondSuccess[T any](w http.ResponseWriter, data T) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(Response[T]{
        Success: true,
        Data:    data,
    })
}
```

---

## Health Check Implementation

### Correct Health Endpoint

```go
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
    // Check database connection
    if err := h.db.Ping(); err != nil {
        respondError(w, http.StatusServiceUnavailable, 5001, "Database unavailable", err)

        return
    }

    // HealthStatus is a typed response for the health endpoint.
    type HealthStatus struct {
        Status    string `json:"status"`
        Timestamp string `json:"timestamp"`
        Version   string `json:"version"`
    }

    respondSuccess(w, HealthStatus{
        Status:    "ok",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Version:   h.version,
    })
}
```

### Health Check Verification

```bash

# Test health endpoint

curl -s http://localhost:8080/api/v1/health | jq .

# Expected response:

{
  "success": true,
  "data": {
    "status": "ok",
    "timestamp": "2026-02-04T12:00:00Z",
    "version": "1.0.0"
  }
}
```

---

## HTTP Request Logging Middleware

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Generate request ID
        requestId := r.Header.Get("X-Request-ID")

        if requestId == "" {
            requestId = uuid.New().String()
        }

        // Wrap response writer to capture status
        wrapped := &statusResponseWriter{ResponseWriter: w}

        // Add request ID to context
        ctx := context.WithValue(r.Context(), "request_id", requestId)

        // Set response header
        w.Header().Set("X-Request-ID", requestId)

        // Log request start
        log.Debug().
            Str("request_id", requestId).
            Str("method", r.Method).
            Str("path", r.URL.Path).
            Str("remote_addr", r.RemoteAddr).
            Msg("Request started")

        // Call next handler
        next.ServeHTTP(wrapped, r.WithContext(ctx))

        // Log request completion
        log.Info().
            Str("request_id", requestId).
            Str("method", r.Method).
            Str("path", r.URL.Path).
            Int("status", wrapped.status).
            Dur("duration", time.Since(start)).
            Msg("Request completed")
    })
}

type statusResponseWriter struct {
    http.ResponseWriter
    status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}
```

---

## Common Issues and Solutions

### Issue: "Connection refused" on startup

**Symptoms:**

- Server appears to start but clients can't connect
- Health check returns "connection refused"

**Check:**

1. Is the server binding to the correct address?
   ```go
   // ❌ Wrong - only localhost
   server.Addr = "localhost:8080"

   // ✅ Correct - all interfaces
   server.Addr = ":8080"
   ```

2. Is the port already in use?
   ```bash
   lsof -i :8080
   ```

3. Is there a firewall blocking the port?

### Issue: 404 on API base URL

**Symptoms:**

- `GET /api/v1` returns 404
- Frontend shows "Backend disconnected"

**Check:**

1. Is there a handler for the base URL?
   ```go
   // Add index route
   r.Get("/api/v1", h.Index)
   r.Get("/api/v1/health", h.Health)
   ```

2. Is the router prefix correct?
   ```go
   r.Route("/api/v1", func(r chi.Router) {
       r.Get("/", h.Index)     // Handles /api/v1
       r.Get("/health", h.Health)
   })
   ```

### Issue: Response format mismatch

**Symptoms:**

- Backend returns 200 OK
- Frontend still shows "disconnected"

**Check:**

1. Is the response using the standard envelope?
   ```go
   // ❌ Wrong
   json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

   // ✅ Correct
   respondSuccess(w, map[string]string{"status": "ok"})
   ```

2. Does the frontend expect the correct structure?
   - Check frontend detection logic
   - Verify it uses HTTP status codes (2xx) as primary indicator

### Issue: CORS errors

**Symptoms:**

- Browser console shows CORS errors
- Requests work from curl but not browser

**Solution:**
```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")

        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)

            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Issue: Database "locked" errors

**Symptoms:**

- SQLite returns "database is locked"
- Concurrent requests fail

**Solution:**
```go
// Configure connection pool for SQLite
db.SetMaxOpenConns(1)  // SQLite only supports one writer
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(time.Hour)

// Use WAL mode for better concurrency
_, err := db.Exec("PRAGMA journal_mode=WAL")
```

---

## Debugging Commands

### Check if server is running

```bash

# Check process

ps aux | grep "cli-name"

# Check port

lsof -i :8080
netstat -tlnp | grep 8080
```

### Test endpoints

```bash

# Health check

curl -s http://localhost:8080/api/v1/health | jq .

# With verbose output

curl -v http://localhost:8080/api/v1/health

# Check response headers

curl -I http://localhost:8080/api/v1/health
```

### Enable debug logging

```bash

# Set environment variable

export DEBUG=true
./cli-name serve

# Or pass flag

./cli-name serve --debug
```

### Check logs

```bash

# If using file logging

tail -f logs/app.log

# If using journald (systemd)

journalctl -u cli-name -f

# Filter by level

journalctl -u cli-name -f | grep -E '"level":"error"'
```

---

## Performance Profiling

### Enable pprof

```go
import _ "net/http/pprof"

func main() {
    // Expose pprof on separate port
    go func() {
        log.Info().Msg("pprof available at http://localhost:6060/debug/pprof/")
        http.ListenAndServe("localhost:6060", nil)
    }()

    // ... rest of server setup
}
```

### Collect profiles

```bash

# CPU profile (30 seconds)

go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profile

go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile

go tool pprof http://localhost:6060/debug/pprof/goroutine
```

---

## Request Tracing Pattern

```go
// Add trace ID to all logs for a request
func (h *Handler) ProcessRequest(w http.ResponseWriter, r *http.Request) {
    traceId := r.Header.Get("X-Trace-ID")

    if traceId == "" {
        traceId = uuid.New().String()
    }

    logger := log.With().
        Str("trace_id", traceId).
        Str("operation", "process_request").
        Logger()

    logger.Info().Msg("Starting request processing")

    // Step 1
    logger.Debug().Msg("Step 1: Validating input")

    if err := h.validateInput(r); err != nil {
        logger.Error().Err(err).Msg("Validation failed")
        respondError(w, http.StatusBadRequest, 4001, "Invalid input", err)

        return
    }

    // Step 2
    logger.Debug().Msg("Step 2: Processing data")
    result, err := h.processData(r.Context())

    if err != nil {
        logger.Error().Err(err).Msg("Processing failed")
        respondError(w, http.StatusInternalServerError, 5001, "Processing error", err)

        return
    }

    logger.Info().Msg("Request completed successfully")
    respondSuccess(w, result)
}
```

---

## Database Error Stack Traces

### Mandatory Requirements

All database operations MUST use the centralized `DBOperation` wrapper from `pkg/database`. This ensures:

1. **Automatic Stack Trace Capture** - Every error includes the full caller chain
2. **Affected Rows Validation** - Compare expected vs actual rows affected
3. **Table Name Logging** - Every log entry includes the table being operated on
4. **Duration Tracking** - Time taken for each operation

### Stack Trace Format

```text
[ERROR] Database operation completed
  Table: User
  Operation: Create
  ExpectedRows: 1
  AffectedRows: 0
  Duration: 5.678ms
  Error: UNIQUE constraint failed: User.Email
  Stack:
    -> user_repository.go:45 (CreateUser)
    -> auth_service.go:112 (RegisterUser)
    -> auth_handler.go:78 (HandleRegister)
```

### Required Log Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Table` | string | Always | Table name being operated on |
| `Operation` | string | Always | Create/Read/Update/Delete |
| `AffectedRows` | int64 | For writes | Actual rows changed |
| `ExpectedRows` | int | For writes | Expected rows for validation |
| `Duration` | duration | Always | Operation execution time |
| `Stack` | []string | On error | Full caller chain with file:line |
| `Error` | string | On error | Error message |

### Usage Pattern

```go
// All database operations MUST use this pattern
op := database.NewDbOperation("User", database.OpCreate).
    ExpectRows(1)

// EXEMPTED: op.Execute callback uses (int64, error) as internal framework boundary
result := op.Execute(func() (int64, error) {
    tx := r.db.Create(user)

    return tx.RowsAffected, tx.Error
})

if result.Error != nil {
    // Stack trace is automatically logged
    return result.Error
}
```

### Anti-Patterns

```go
// ❌ WRONG: Direct GORM without wrapper
result := r.db.Create(user)

// ❌ WRONG: Raw SQL (forbidden except FTS5/vectors)
r.db.Exec("INSERT INTO User (Id, Email) VALUES (?, ?)", id, email)

// ❌ WRONG: Missing expected rows for write operations
op := database.NewDbOperation("User", database.OpUpdate)
// Missing: .ExpectRows(1)
```

---

## ORM-Only Policy

### Rule

Use GORM for 99% of database operations. Raw SQL is forbidden except for:

- FTS5 virtual tables (SQLite limitation)
- Vector storage operations (if no ORM support)
- Complex CTEs/recursive queries

### Relationship-First Pattern

Always find the parent model first, then manipulate relationships:

```go
// ❌ WRONG: Raw SQL INSERT
db.Exec("INSERT INTO File (ProjectId, Name) VALUES (?, ?)", projectId, name)

// ✅ CORRECT: Relationship-first
var project Project
db.First(&project, "Id = ?", projectId)
project.Files = append(project.Files, File{Name: name})
db.Save(&project)
```

See full specification: `02-spec/02-spec-management-software/13-shared-packages/08-pkg-database-operations.md`

---

## Cross-Reference

- [Error Resolution Overview](../../01-index.md)
- [Frontend-Backend Sync Verification](../04-verification-patterns/02-frontend-backend-sync.md)
- [Error Code Registry](../../03-error-code-registry/01-index.md)
- Split DB Architecture *(external spec)*
- DBOperation Wrapper Spec *(external spec)*
- ORM-Only Policy — enforced project-wide (no raw SQL outside migrations)
