# Error Handling — Cross-Stack Specification

> **Version:** 2.0.0  
> **Updated:** 2026-03-09  
> **Status:** Active  
> **Applies to:** Go backend, React/TypeScript frontend, PHP WordPress plugin, any delegated 3rd-party server

---

## Overview

The project implements a **three-tier error handling architecture** spanning the full React → Go → Delegated Server (PHP/3rd-party) request chain. Every error is captured with structured diagnostics, stack traces, and contextual metadata to enable deep debugging from the frontend Global Error Modal.

**v2.0.0 change:** Added `DelegatedRequestServer` structured error block — the Go backend now captures and propagates full request/response/stack details from any downstream server it proxies to.

---

## Error Flow Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Frontend (React/TypeScript)                   │
│  ┌─────────────────┐    ┌──────────────────┐    ┌───────────────┐   │
│  │ API Client       │───▸│ Error Store       │───▸│ Global Error  │   │
│  │ (parseEnvelope)  │    │ (captureError)    │    │ Modal (tabs)  │   │
│  └─────────────────┘    └──────────────────┘    └───────────────┘   │
│         │                       │                                    │
│         │ Envelope.Errors       │ executionChain                    │
│         │ Envelope.MethodsStack │ clickPath                         │
│         │ Envelope.SessionId    │ componentContext                  │
│         │ Envelope.Errors       │                                    │
│         │  .DelegatedRequest    │                                    │
│         │   Server ◀────────── NEW (v2.0.0)                         │
└─────────────────────────────────────────────────────────────────────┘
                              ▲
                              │ Universal Response Envelope
┌─────────────────────────────────────────────────────────────────────┐
│                        Backend (Go)                                  │
│  ┌─────────────────┐    ┌──────────────────┐    ┌───────────────┐   │
│  │ apperror.Wrap() │───▸│ Session Logger    │───▸│ error.log.txt │   │
│  │ + .WithContext() │    │ (per-request ID)  │    │ (deduped)     │   │
│  └─────────────────┘    └──────────────────┘    └───────────────┘   │
│         │                       │                                    │
│         │ stack trace           │ buildDelegatedRequestServer()      │
│         │ error code            │ fetchAndAttachRemotePHPErrors      │
│         │                       │                                    │
│         │ ┌─────────────────────────────────────────────────────┐    │
│         │ │ DelegatedRequestServer Builder (NEW v2.0.0)         │    │
│         │ │  • Captures: endpoint, method, statusCode           │    │
│         │ │  • Captures: requestBody, response, stackTrace      │    │
│         │ │  • Injects into Envelope.Errors block               │    │
│         │ │  • Writes to error.log.txt "Delegated Server Info"  │    │
│         │ └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────┘
                              ▲
                              │ REST API (JSON) — any downstream server
┌─────────────────────────────────────────────────────────────────────┐
│              Delegated Server (PHP / Chrome Extension / Other)       │
│  ┌─────────────────┐    ┌──────────────────┐    ┌───────────────┐   │
│  │ safe_execute()  │───▸│ FileLogger        │───▸│ stacktrace.txt│   │
│  │ catch Throwable │    │ (6-frame backtrace)│   │ fatal-errors  │   │
│  │                 │    │                    │    │ error.txt     │   │
│  │ OR any 3rd-party│    │ OR structured     │    │               │   │
│  │ error format    │    │ error response     │    │               │   │
│  └─────────────────┘    └──────────────────┘    └───────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### Request Chain (3-Hop)

```
  React Frontend          Go Backend           Delegated Server
       │                      │                      │
       │  GET /api/v1/sites   │                      │
       │  /1/snapshots/       │                      │
       │  settings            │                      │
       │─────────────────────▸│                      │
       │                      │  GET /wp-json/       │
       │                      │  riseup.../v1/       │
       │                      │  snapshots/settings  │
       │                      │─────────────────────▸│
       │                      │                      │
       │                      │  HTTP 403            │
       │                      │  { code, message,    │
       │                      │    stackTrace,       │
       │                      │    plugin_version }  │
       │                      │◀─────────────────────│
       │                      │                      │
       │                      │ Build envelope:      │
       │                      │ • Errors.Backend[]   │
       │                      │ • Errors.Delegated   │
       │                      │   RequestServer{}    │
       │                      │ • MethodsStack       │
       │                      │                      │
       │  HTTP 500            │                      │
       │  Universal Envelope  │                      │
       │◀─────────────────────│                      │
       │                      │                      │
       │ parseEnvelope()      │                      │
       │ captureError()       │                      │
       │ openErrorModal()     │                      │
```

---

## Tier 1: Delegated Server Error Handling

> The delegated server can be **any downstream service** the Go backend proxies to — a WordPress/PHP plugin, a Node.js microservice, a Chrome extension, a Python API, etc. For the sake of example, this spec uses **PHP (WordPress plugin)** as the delegated language. When implementing in another language, follow that language's best practices for structured error responses, stack traces, and naming conventions.

### Safe Execution Pattern (PHP Example)

Every REST endpoint handler is wrapped in `safeExecute`:

```php
public function handleRequest(WP_REST_Request $request): WP_REST_Response {
    return $this->safeExecute(function() use ($request) {
        // Business logic

        return $this->envelope->success($result);
    });
}
```

The wrapper catches `Throwable` (not just `Exception`) to capture PHP 7+ Errors like missing classes.

### Structured Error Response

```json
{
  "message": "Class 'PDO' not found",
  "StackTrace": "#0 /path/file.php(42): PluginManager->connect()\n#1 {main}",
  "StackTraceFrames": [
    { "file": "/path/file.php", "line": 42, "function": "connect", "class": "PluginManager" }
  ]
}
```

### REST API Error Enrichment (PHP Example)

The plugin uses a `rest_post_dispatch` filter to inject metadata into all error responses. Other delegated languages should implement equivalent response enrichment using their framework's middleware or interceptor pattern.

```php
add_filter('rest_post_dispatch', function($response, $server, $request) {
    if ($response->is_error()) {
        $data = $response->get_data();
        $data['plugin_version'] = PluginConfigType::Version->value;
        $data['timestamp'] = gmdate('c');
        $data['log_hint'] = $this->getLogHint($response->get_status());
        $response->set_data($data);
    }

    return $response;
}, 10, 3);
```

This ensures the Go backend always receives structured metadata when a delegated request fails. **Any delegated language** should return equivalent structured JSON on failure.

### Logging Outputs

| File | Content | Depth |
|------|---------|-------|
| `error.txt` | Structured error entries with context metadata | Last N entries |
| `log.txt` | General diagnostic log | All operations |
| `stacktrace.txt` | Raw PHP backtraces (`debug_backtrace(0, 0)`) | Unlimited |
| `fatal-errors.log` | Fatal errors caught by shutdown handler | With memory usage |

### Global Shutdown Handler (PHP Example)

Use `ErrorChecker::isFatalError()` to centralize fatal error detection. `ErrorChecker` delegates to `ErrorTypeEnum::FATAL_TYPES` (see [PHP Enum Spec](../../02-coding-guidelines/04-php/01-enums.md) for full implementation). Other delegated languages should implement equivalent uncaught-exception handlers (e.g., Node.js `process.on('uncaughtException')`, Python `sys.excepthook`).

```php
register_shutdown_function(function() {
    $error = error_get_last();

    if (ErrorChecker::isFatalError($error)) {
        // Log to fatal-errors.log via PathHelper::getFatalErrorLog()
        // Include memory_get_peak_usage() for diagnostics
        // Send JSON response before process terminates (if REST_REQUEST)
    }

});
```

### Context Enrichment

Every `error()` and `logException()` call automatically captures:
- 6-frame backtrace
- HTTP method and endpoint
- User-agent and IP
- Memory usage
- Request body (truncated)

---

## Tier 2: Go Backend Error Handling

### `apperror` Package

All errors crossing service boundaries must use `apperror`:

```go
// Wrap existing errors with code and context
return apperror.Wrap(err, apperror.ErrSyncCheck, "failed to upload plugin").
    WithPluginContext(pluginId, pluginSlug).
    WithEndpoint(requestUrl)

// Create new errors
return apperror.New(
    apperror.ErrFileRead, "invalid plugin slug",
)
```

**Forbidden:** `fmt.Errorf` for errors leaving a service (no stack trace).

### DelegatedRequestServer Injection (NEW v2.0.0)

When the Go backend proxies a request to any downstream server and the request fails (status ≥ 400), it builds a `DelegatedRequestServer` object and injects it into the envelope's `Errors` block.

#### Construction Flow

```go
// In the HTTP client wrapper (e.g., wordpress.(*Client).doRequest)
func (c *Client) doRequest(
	context stdctx.Context,
	method string,
	url string,
	body any,
) apperror.Result[*http.Response] {
    // ... execute request ...
    
    if resp.StatusCode >= 400 {
        // Build DelegatedRequestServer from the failed response
        delegated := &DelegatedRequestServer{
            DelegatedEndpoint:  url,
            Method:             method,
            StatusCode:         resp.StatusCode,
            RequestBody:        body,        // what we sent
            Response:           respBody,    // what they returned (parsed JSON)
            StackTrace:         extractStackTrace(respBody),  // from response if available
            AdditionalMessages: extractMessage(respBody),
        }
        
        // Attach to the request context for envelope builder to pick up
        context = stdctx.WithValue(context, delegatedServerKey, delegated)
        
        return resp, apperror.Wrap(err, apperror.ErrWpConnect, "delegated request failed").
            WithEndpoint(url).
            WithStatusCode(resp.StatusCode)
    }
}
```

#### Envelope Builder Integration

```go
func (b *EnvelopeBuilder) BuildErrorResponse(context stdctx.Context, err error) *Response {
    resp := &Response{
        Status: Status{IsFailed: true, Code: getStatusCode(err), ...},
        Errors: &Errors{
            BackendMessage: err.Error(),
            Backend:        getGoStackTrace(err),
        },
    }

    b.injectDelegatedServer(context, resp)

    return resp
}

func (b *EnvelopeBuilder) injectDelegatedServer(context stdctx.Context, resp *Response) {
    delegated, ok := context.Value(delegatedServerKey).(*DelegatedRequestServer)

    if !ok {
        return
    }

    resp.Errors.DelegatedRequestServer = delegated

    // Populate legacy field for backward compatibility
    if len(delegated.StackTrace) > 0 {
        resp.Errors.DelegatedServiceErrorStack = delegated.StackTrace
    }
}
```

#### Go Struct Definition

```go
type DelegatedRequestServer struct {
    DelegatedEndpoint  string
    Method             string
    StatusCode         int
    RequestBody        json.RawMessage `json:",omitempty"`
    Response           json.RawMessage `json:",omitempty"`
    StackTrace         []string        `json:",omitempty"`
    AdditionalMessages string          `json:",omitempty"`
}
```

#### error.log.txt Format

When `DelegatedRequestServer` is present, the error log includes a `Delegated Server Info:` section:

```
[2026-02-12 00:53:34] HTTP 500 GET FAILED
  Requested To: GET http://localhost:8080/api/v1/sites/1/snapshots/settings
  Duration: 4.1904552s
  Error Code: 500
  Error Message: [E3001] failed to fetch snapshot settings...
  Backend Error: [E3025] [E3001] ...
  Go Backend Stack:
    handler_factory.go:107 handlers.init.handleSiteActionById.func63
    ...
  Delegated Server Info:
    Endpoint: "https://example.com/wp-json/riseup.../v1/snapshots/settings"
    Method: "GET"
    Status: 403
    Stacktrace:
        #0 riseup-asia-uploader.php(1098): FileLogger->error()
        #1 class-wp-hook.php(341): Plugin->enrichErrorResponse()
        ...
    RequestBody:
        (none — GET request)
    Additional Message:
        Endpoint 'snapshots' is not enabled in plugin settings.
  Response Body:
    { "Status": { ... }, "Errors": { ..., "DelegatedRequestServer": { ... } } }
```

### Remote PHP Error Injection (Legacy)

When a remote WordPress operation fails, the Go backend automatically:

1. Calls `fetchAndAttachRemotePHPErrors` on the target site
2. Retrieves the 10 most recent PHP errors from remote SQLite database
3. Retrieves `stacktrace.txt` content
4. Injects this data into Go session logs and the envelope's `Errors` block

> **Note:** `DelegatedRequestServer` (v2.0.0) supersedes the legacy `fetchAndAttachRemotePHPErrors` approach for inline error data. The legacy system remains for retrieving historical PHP errors not present in the immediate response.

### Error Log Deduplication

The backend uses MD5 hashing to suppress identical error log entries:

```
Hash = MD5(action + siteId + plugin + endpoint + statusCode + responseBody)
```

A "Clear Dedup Hashes" button in Settings resets the in-memory hash map.

### Redefined Log Format

Every failure log entry follows this structure:

1. **Site Request URL** — Full compiled endpoint on the target WordPress site
2. **Site Identification** — Site name and URL
3. **Backend Endpoint** — The Go endpoint hit by the frontend
4. **Delegated Request** — Method, delegated server endpoint, full JSON request body
5. **Delegated Response** — Status code and body
6. **Delegated Server Info** — Endpoint, method, status, stack trace, request body, additional messages (NEW v2.0.0)
7. **Error Summary** — Concise error description
8. **Guard Rail** — Blocks unauthorized direct mutations to `/wp/v2/plugins/*`

### Session-Based Logging

Every HTTP request gets a unique session ID. Full request/response data is captured:
- Headers (with Authorization redacted)
- Bodies (truncated at 50KB)
- Timing
- Error extraction for status ≥ 400
- DelegatedRequestServer data (if delegated request failed)

Storage: `backend/data/request-sessions/{date}/{hour}/{uuid}.json`

---

## Tier 3: Frontend Error Handling

### Error Store (`errorStore.ts`)

Centralized Zustand store that captures:

```typescript
interface CapturedError {
  // Identity
  id: string;
  code: string;
  level: 'error' | 'warn' | 'info';
  message: string;
  details?: string;
  createdAt: string;
  
  // API request context
  endpoint?: string;
  method?: string;
  requestBody?: unknown;
  responseStatus?: number;
  
  // Trigger context
  triggerComponent?: string;
  triggerAction?: string;
  invocationChain?: string[];
  
  // Session-based logging
  sessionId?: string;
  sessionType?: string;
  
  // Universal Envelope diagnostic fields
  requestedAt?: string;
  requestDelegatedAt?: string;
  envelopeErrors?: EnvelopeErrors;       // Contains DelegatedRequestServer (v2.0.0)
  envelopeMethodsStack?: EnvelopeMethodsStack;
  
  // Frontend diagnostics
  parsedFrames?: StackFrame[];
  uiClickPath?: ClickEvent[];
  executionLogs?: ExecutionLogEntry[];
  executionLogsFormatted?: string;
}

// EnvelopeErrors contains DelegatedRequestServer (NEW v2.0.0)
interface EnvelopeErrors {
  BackendMessage: string;
  DelegatedServiceErrorStack?: string[];  // Legacy delegated server stack lines
  Backend?: string[];
  Frontend?: string[];
  DelegatedRequestServer?: {              // NEW v2.0.0 — structured delegated error
    DelegatedEndpoint: string;
    Method: string;
    StatusCode: number;
    RequestBody?: unknown;
    Response?: unknown;
    StackTrace?: string[];
    AdditionalMessages?: string;
  };
}
```

> **Note:** `DelegatedRequestServer` is accessed via `error.envelopeErrors?.DelegatedRequestServer`, not as a top-level field on `CapturedError`. See `./04-error-modal/03-error-modal-reference.md` for the full interface.

### Envelope Parsing

The API client's `parseEnvelope` detects failed responses and extracts:
- `Errors.BackendMessage` — Primary error text
- `Errors.DelegatedServiceErrorStack` — Delegated server stack trace lines (legacy)
- `Errors.DelegatedRequestServer` — Full delegated server error details (NEW v2.0.0)
- `Errors.Backend` — Go stack trace lines
- `MethodsStack.Backend` — Go call chain with file:line
- `Attributes.SessionId` — Links to session-level diagnostics

### DelegatedRequestServer Frontend Extraction (NEW v2.0.0)

```typescript
// In parseEnvelope() or buildCapturedError()
if (envelope.Errors) {
  captured.envelopeErrors = {
    BackendMessage: envelope.Errors.BackendMessage,
    DelegatedServiceErrorStack: envelope.Errors.DelegatedServiceErrorStack,
    Backend: envelope.Errors.Backend,
    Frontend: envelope.Errors.Frontend,
    DelegatedRequestServer: envelope.Errors.DelegatedRequestServer,
  };
}
```

### Global Error Modal Tabs

| Tab | Content |
|-----|---------|
| **Overview** | Error message, component context, suggested fixes |
| **Log** | error.log.txt content |
| **Execution** | Go call chain table + session logs |
| **Stack** | Backend (Go) + Delegated Server stack traces + delegated language frames |
| **Session** | Session diagnostics: logs, request, response, stack trace |
| **Request** | HTTP request/response chain (3-hop: React → Go → Delegated) |
| **Traversal** | Endpoint flow + methods stack + delegated error stack |

### DelegatedRequestServer in Modal Tabs (NEW v2.0.0)

| Tab | What's Shown |
|-----|-------------|
| **Stack** | New "Delegated Server Stack" section (purple-themed) with StackTrace lines + Response JSON |
| **Request** | 3rd node in chain: Go → Delegated (with endpoint, method, status, request body, response) |
| **Traversal** | Endpoint flow extended to 3 hops; DelegatedRequestServer details below methods stack |
| **Overview** | AdditionalMessages shown as info banner below the error banner |

### Session Diagnostics Auto-Fetch

When `sessionId` is present, the modal automatically fetches session-level diagnostics from `GET /api/v1/sessions/{id}/diagnostics`, merging deep Go and delegated server stack traces into the Stack and Execution tabs.

### Error Reporting Bundle

The "Download Bundle" button exports:
- All diagnostic data as JSON
- Syntax-highlighted error report
- Execution chain and click path
- Full request/response data (including DelegatedRequestServer)

---

## Error Code Ranges

| Range | Category | Example |
|-------|----------|---------|
| E1000–E1999 | Connection/network errors | E1001: Backend unreachable |
| E2000–E2999 | Remote site errors | E2001: Invalid credentials |
| E3000–E3999 | Resource/data errors | E3001: Failed to fetch resource |
| E4000–E4999 | Client/validation errors | E4001: Invalid plugin slug |
| E5000–E5999 | Server/infrastructure errors | E5001: Upload failed |
| E6000–E6999 | Remote site (WordPress) errors | E6001: Plugin not found on site |
| E7000–E7999 | Scheduler/background job errors | E7001: Scheduled publish timeout |
| E8000–E8999 | Delegated server errors | E8001: Delegated server returned 5xx |
| E9000–E9999 | Frontend/UI errors | E9003: Unhandled API error, E9005: HTML instead of JSON |

---

## Fallback Visibility

The Errors page implements a 3-tier fallback:

1. **Live Backend API** — Primary source
2. **Global Error Store** — Session-captured errors
3. **Error-level Notifications** — Local errors with Eye icon for modal access

---

## Cross-References

- [Error Resolution Retrospectives](../01-error-resolution/03-retrospectives/)
- [Session-Based Logging](./07-logging-and-diagnostics/02-session-based-logging.md)
- [React Execution Logger](./07-logging-and-diagnostics/01-react-execution-logger.md)
- [Error Modal Spec](./04-error-modal/03-error-modal-reference.md)
- [Copy Format Samples](./04-error-modal/01-copy-formats.md)
- [Response Envelope Schema](./05-response-envelope/envelope.schema.json)
- [Envelope Configurability](./05-response-envelope/01-adr.md)
- [PHP Standards](../../02-coding-guidelines/04-php/07-php-standards-reference/00-overview.md)
- [Golang Standards](../../02-coding-guidelines/03-golang/04-golang-standards-reference/00-overview.md)

---

*Error handling specification v2.0.0 — updated: 2026-02-11*
