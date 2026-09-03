# Debugging Cheat Sheet

> **Version:** 1.0.0
> **Created:** 2026-02-04
> **Quick Reference for:** PHP, Go, TypeScript

---

## Universal Principles

| Principle | Description |
|-----------|-------------|
| **Never Assume** | Always verify endpoints exist before implementing |
| **HTTP Status First** | Use HTTP status codes (2xx) as primary indicator, not response body |
| **Standard Envelope** | All backends return `{success, data, error}` format |
| **Diagnostics** | Always show raw env vars vs resolved values |

---

## Initialization Order (ALL Languages)

```
1. Configuration    → Load env vars and config files FIRST
2. Directories      → Ensure all required directories exist
3. Database         → Initialize connections (only after dirs exist)
4. Services         → Initialize business logic components
5. Server/App       → Start ONLY after all dependencies ready
```

---

## PHP Quick Reference

### Log Locations

```
wp-content/uploads/{plugin-slug}/logs/debug.log
wp-content/uploads/{plugin-slug}/logs/error.log
```

### Enable/Disable Logging

```php
define('PLUGIN_DEBUG_LOGGING', true);   // Trace every step
define('PLUGIN_ERROR_LOGGING', true);   // Errors + stack traces
```

### Common Commands

```php
// Log debug message
Logger::debug('Message', ['context' => $value]);

// Log error with trace
Logger::error('Error message', ['exception' => $e]);

// Check PDO extension
extension_loaded('pdo_sqlite'); // Returns true/false
```

### Debugging Steps

1. Clear logs → Delete `debug.log` and `error.log`
2. Trigger action → Try to activate plugin
3. Check debug.log → Find last successful step
4. Check error.log → Find exception with file/line

### Common Issues

| Issue | Check |
|-------|-------|
| Plugin won't activate | PDO SQLite extension, directory permissions |
| Database connection fails | Directory path, permissions, PDO exceptions |
| Components not initializing | Class loading, missing dependencies |

---

## Go Quick Reference

### Logging Levels

```
Trace → Debug → Info → Warn → Error → Fatal
```

### Common Commands

```bash

# Check if server running

ps aux | grep "cli-name"
lsof -i :8080

# Test endpoints

curl -s http://localhost:8080/api/v1/health | jq .
curl -v http://localhost:8080/api/v1/health

# Enable debug mode

export DEBUG=true && ./cli-name serve
./cli-name serve --debug

# View logs

tail -f logs/app.log
journalctl -u cli-name -f | grep '"level":"error"'
```

### Structured Logging (zerolog)

```go
// Simple message
log.Info().Msg("Server started")

// With fields
log.Info().
    Str("method", r.Method).
    Int("status", status).
    Dur("duration", time.Since(start)).
    Msg("Request completed")

// With error
log.Error().Err(err).Str("operation", "db_query").Msg("Failed")

// Contextual logger
logger := log.With().Str("request_id", id).Logger()
```

### Response Helpers

```go
// Error response
respondError(w, http.StatusBadRequest, 7001, "Invalid input", err)

// Success response — use typed struct, not map[string]any
type StatusResponse struct {
    Status string
}

respondSuccess(w, StatusResponse{Status: "ok"})
```

### Common Issues

| Issue | Solution |
|-------|----------|
| Connection refused | Check `server.Addr = ":8080"` (not `localhost:8080`) |
| 404 on base URL | Add handler for `/api/v1` route |
| Response format mismatch | Use `respondSuccess()` helper |
| CORS errors | Add CORS middleware with proper headers |
| Database locked | `db.SetMaxOpenConns(1)` + WAL mode |

### Profiling (pprof)

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof http://localhost:6060/debug/pprof/heap
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

---

## TypeScript/React Quick Reference

### Console Logging

```typescript
// Structured logger
logger.debug('Message', data);        // Dev only
logger.info('Message', data);         // Normal operation
logger.warn('Message', data);         // Recoverable issues
logger.error('Message', error);       // Failures
logger.api('GET', '/api/v1/health', 200, 45); // API calls
```

### API Response Validation

```typescript
// Standard envelope
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: { code: number; message: string; details?: string };
}

// Detection logic (USE HTTP STATUS FIRST!)
if (!response.ok) throw new Error(`HTTP ${response.status}`);
const body = await response.json();

if (!body.success) throw new Error(body.error?.message);

return body.data;
```

### Health Check Pattern

```typescript
const response = await fetch(`${baseUrl}/api/v1/health`, {
  signal: AbortSignal.timeout(5000),
});
const connected = response.ok; // PRIMARY indicator
```

### React Query Debugging

```typescript
// Add DevTools
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
<ReactQueryDevtools initialIsOpen={false} />

// Check query state
const { status, fetchStatus, isLoading, isFetching } = useQuery({...});
console.log({ status, fetchStatus, isLoading, isFetching });
```

### WebSocket Debugging

```typescript
ws.onclose = (event) => console.log({
  code: event.code,
  reason: event.reason,
  wasClean: event.wasClean
});
```

### Browser DevTools Filters

```

# Network tab

/api/v1/                        # Filter API calls
status-code:4xx                 # Filter 4xx errors
larger-than:100k                # Large responses

# Console filtering

-:extension:                    # Hide extension logs
url:/api/                       # API-related only
```

### Common Issues

| Issue | Solution |
|-------|----------|
| Backend disconnected | Check API URL, CORS, response format |
| Works in dev, fails in prod | Check env vars, relative URLs, HTTPS |
| State not updating | Use `queryClient.invalidateQueries()` |
| TypeScript errors | Define explicit response types, use Zod |

---

## Quick Diagnosis Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    ISSUE REPORTED                               │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  1. VERIFY ENDPOINT EXISTS IN BACKEND                           │
│     • Check router registration                                 │
│     • Check handler implementation                              │
│     • Test with curl                                            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. CHECK RESPONSE FORMAT                                       │
│     • Returns {success, data, error}?                           │
│     • HTTP status correct?                                      │
│     • Content-Type: application/json?                           │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. VERIFY FRONTEND DETECTION LOGIC                             │
│     • Uses response.ok as PRIMARY indicator?                    │
│     • Handles error envelope correctly?                         │
│     • Shows diagnostics on failure?                             │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  4. CHECK ENVIRONMENT                                           │
│     • Raw vs resolved env vars match?                           │
│     • CORS configured correctly?                                │
│     • Ports/paths correct?                                      │
└─────────────────────────────────────────────────────────────────┘
```

---

## Error Code Ranges

| System | Range | Example |
|--------|-------|---------|
| General | 1000-1999 | 1001: Generic error |
| Spec Management | 2000-2999 | 2001: Config error |
| GSearch CLI | 7000-7099 | 7001: Search failed |
| BRun CLI | 7100-7599 | 7101: Build failed |
| Nexus Flow CLI | 8000-8399 | 8001: Flow error |
| AI Bridge CLI | 9000-9499 | 9001: Bridge error |
| Spec Publish | 9500-9999 | 9501: Publish failed |
| WP Plugin Builder | 10000-10999 | 10001: Plugin error |
| Spec Reverse | 11000-11999 | 11001: Parse error |
| WP SEO Publish | 12000-12599 | 12001: SEO error |
| AI Transcribe | 14000-14499 | 14001: Transcribe error |

---

## Cross-References

| Guide | Target |
|-------|--------|
| [PHP Debugging](./05-debugging-guides/02-debugging-php.md) | WordPress, PHP backends |
| [Go Debugging](./05-debugging-guides/03-debugging-go.md) | CLI tools, Go backends |
| [TypeScript Debugging](./05-debugging-guides/04-debugging-typescript.md) | React frontends |
| [Frontend-Backend Sync](./04-verification-patterns/02-frontend-backend-sync.md) | Integration verification |
| [Cross-Reference Diagram](./03-cross-reference-diagram.md) | Architecture overview |

---

*Quick reference for debugging across the ecosystem. For detailed guides, see the linked documents.*
