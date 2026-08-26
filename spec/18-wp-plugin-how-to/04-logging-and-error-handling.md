# Phase 4 — Logging and Error Handling

> **Purpose:** Define the complete logging architecture, error handling strategy, debug-mode gating, and stack trace formatting in full detail so any AI can implement it correctly from scratch.

---

## 4.1 Two-Tier Logging Architecture

The plugin uses two independent logging tiers:

| Tier | Class | When available | Purpose |
|------|-------|---------------|---------|
| Tier 1 — PHP native | `error_log()` / `ErrorLogHelper` | Always | Fallback when FileLogger is not ready; also emits to WP_DEBUG log |
| Tier 2 — FileLogger | `FileLogger` (singleton) | After autoloader loads | Primary structured log with rotation, dedup, and stack trace separation |

### Why two tiers

The autoloader and bootstrap run before any plugin classes are available. If they fail, Tier 1 (native `error_log()`) captures the failure. Once the plugin initialises, Tier 2 (FileLogger) handles all logging with structured output and file management.

---

## 4.2 Debug Mode — The Master Gate

Every plugin defines a debug mode constant in its main plugin file:

```
/** Enable debug mode — exposes stack traces in API responses. */
define('MY_PLUGIN_DEBUG', false);
```

### What debug mode controls

| Feature | Debug ON | Debug OFF |
|---------|----------|-----------|
| Stack traces in API error responses | ✅ Full frames included in `Errors.Backend` | ❌ `Errors.Backend` omitted entirely |
| Verbose log entries | ✅ `debug()` writes to `info.log` | ❌ `debug()` calls are silently skipped |
| Error response detail | ✅ Full exception message in `Errors.BackendMessage` | ⚠️ Generic message: `"An internal error occurred"` |
| Performance timing in logs | ✅ Included in context | ❌ Omitted |

### Checking debug mode

Use the `PluginConfigType` enum, not the raw constant:

```
enum PluginConfigType: string
{
    case DebugConstant = 'MY_PLUGIN_DEBUG';

    /** Check if the plugin is running in debug mode. */
    public static function isDebugMode(): bool
    {
        $constantName = self::DebugConstant->value;
        $isDefined = defined($constantName);

        return $isDefined && constant($constantName) === true;
    }
}
```

**Usage:**
```
$isDebug = PluginConfigType::isDebugMode();
```

---

## 4.3 FileLogger — Complete Specification

### Singleton access

```
$logger = FileLogger::getInstance();
```

### Log files

The logger writes to three separate files, all under `wp-content/uploads/{plugin-slug}/logs/`:

| File | Contains | Written by |
|------|----------|-----------|
| `info.log` | All log entries (debug, info, warn, error) | Every log call |
| `error.log` | Only warn and error entries | `warn()` and `error()` calls |
| `stacktrace.log` | Full stack traces for errors | `logException()` and error-level calls |

### Public API — Full Signatures

```
class FileLogger
{
    /**
     * Log a debug-level message. Only writes if debug mode is enabled.
     *
     * @param string               $message  Human-readable description
     * @param array<string, mixed> $context  Structured key-value context data
     */
    public function debug(string $message, array $context = []): void;

    /**
     * Log an informational message.
     *
     * @param string               $message  Human-readable description
     * @param array<string, mixed> $context  Key-value pairs (e.g., ['version' => '2.31.0', 'timeMs' => 1.23])
     */
    public function info(string $message, array $context = []): void;

    /**
     * Log a warning. Writes to both info.log and error.log.
     * Also writes a stack trace entry for diagnostic context.
     *
     * @param string               $message  Warning description
     * @param array<string, mixed> $context  Key-value pairs with diagnostic data
     */
    public function warn(string $message, array $context = []): void;

    /**
     * Log an error. Writes to info.log, error.log, and stacktrace.log.
     *
     * @param string               $message  Error description
     * @param array<string, mixed> $context  Key-value pairs (e.g., ['endpoint' => '/status', 'userId' => 1])
     */
    public function error(string $message, array $context = []): void;

    /**
     * Log an exception with full stack trace extraction.
     *
     * @param Throwable            $exception  The caught exception
     * @param string               $context    Human-readable context string (e.g., 'Route registration')
     */
    public function logException(Throwable $exception, string $context = ''): void;

    /**
     * Log a critical exception and re-throw it. Return type is `never`.
     * Use in infrastructure code where silent failure causes cascading breakage.
     *
     * @param Throwable $exception  The caught exception
     * @param string    $context    Human-readable context string
     *
     * @throws Throwable Always re-throws the original exception
     */
    public function logCriticalException(Throwable $exception, string $context = ''): never;
}
```

### Method behaviour matrix

| Method | Level | Writes to info.log | Writes to error.log | Writes stacktrace | Dedup enabled | Skipped in non-debug |
|--------|-------|--------------------|---------------------|--------------------|---------------|---------------------|
| `debug()` | Debug | Yes (if debug) | No | No | Yes (persistent) | ✅ Yes |
| `info()` | Info | Yes | No | No | Yes (persistent) | No |
| `warn()` | Warn | Yes | Yes | Yes | No | No |
| `error()` | Error | Yes | Yes | Yes | No | No |
| `logException()` | Error | Yes | Yes | Yes (from exception) | No | No |
| `logCriticalException()` | Error | Yes | Yes | Yes | No — re-throws | No |

---

## 4.4 Log Entry Format

Every log line follows this exact format:

```
[{timestamp} v{version}] [{Level}] {message} ({file}:{line}) {json_context}
```

| Component | Source | Example |
|-----------|--------|---------|
| Timestamp | `DateHelper::nowLogDisplay()` | `07-Apr-26 2:30 PM` |
| Version | `PluginConfigType::Version->value` | `2.31.0` |
| Level | `LogLevelType` case value | `Info`, `Error` |
| Message | Passed by caller | `Plugin initialized` |
| File:Line | Extracted from `debug_backtrace()` | `Plugin.php:107` |
| Context | JSON-encoded associative array | `{"version":"2.31.0","timeMs":1.23}` |

### Examples of actual log lines

```
[07-Apr-26 2:30 PM v2.31.0] [Info] Plugin initialized successfully (Plugin.php:107) {"version":"2.31.0","timeMs":12.5}
[07-Apr-26 2:30 PM v2.31.0] [Warn] Request exceeded timeout threshold (UploadHandlerTrait.php:45) {"endpoint":"/upload","durationMs":5200,"thresholdMs":5000}
[07-Apr-26 2:31 PM v2.31.0] [Error] Failed to activate plugin on remote site (ActivateHandlerTrait.php:78) {"site":"site-a","httpStatus":502,"responseBody":"Bad Gateway"}
```

### Caller resolution

The logger uses `debug_backtrace(DEBUG_BACKTRACE_IGNORE_ARGS, 3)` to capture the actual caller's file and line, not the logger's own location. The skip depth is:

```
[0] = logAtLevel() (internal)
[1] = actual caller ← this is captured
[2] = caller's caller
```

---

## 4.5 Stack Trace File Format

Stack traces are written to a dedicated file with visual separators:

```
================================================================================
[07-Apr-26 2:30 PM v2.31.0] Error message here (SomeFile.php:42)
Exception: RuntimeException
Message: Cannot connect to remote endpoint
--------------------------------------------------------------------------------
#0 ActivateHandlerTrait.php(78): PluginName\Traits\Activate\ActivateHandlerTrait->executeActivation()
#1 ResponseTrait.php(35): PluginName\Traits\Core\ResponseTrait->safeExecute()
#2 WP_REST_Server.php(1181): WP_REST_Server->dispatch()
#3 rest-api.php(407): rest_do_request()
================================================================================

```

Each trace entry is a self-contained block with `=` separators for easy visual scanning. The file is separate from `info.log` and `error.log` to avoid cluttering operational logs.

---

## 4.6 Log Rotation

The FileLogger automatically rotates log files when they exceed a size threshold.

| Setting | Default | Range |
|---------|---------|-------|
| Max file size | 512 KB | 64 KB – 10 MB |
| Max rotations | 10 | 1 – 100 |

### Rotation process

1. Before each write, check if the target file exceeds `maxLogSizeBytes`
2. If yes, move the file to `logs/archive/{NNN}/{filename}` where NNN is a zero-padded sequential index
3. If archive folder count reaches `maxRotations`, delete the oldest folders first
4. The current file is now empty and ready for new writes

### Archive structure

```
logs/
├── info.log           ← current
├── error.log          ← current
├── stacktrace.log     ← current
└── archive/
    ├── 001/
    │   ├── info.log
    │   └── error.log
    ├── 002/
    │   └── info.log
    └── 003/
        └── stacktrace.log
```

---

## 4.7 Deduplication

The logger has two dedup layers to prevent repetitive log entries:

### In-memory dedup (per-request)

- Hashes `level + message + file + line`
- If the same hash appears again in the same PHP request, the entry is silently skipped
- Prevents loops from flooding logs

### Persistent dedup (cross-request)

- Stores hashes in a JSON file (`dedup-registry.json`) in the logs directory
- Used only for `debug()` and `info()` level entries
- Maximum 500 entries; oldest entries are pruned when limit is reached
- Prevents boot/init messages from repeating on every request

---

## 4.8 Error Handling — Mandatory Rules

### Rule 1: Always catch Throwable

Every try-catch block catches `Throwable`, never `Exception`. This captures both standard exceptions and PHP fatal errors (TypeError, Error, etc.).

```
try {
    // code
} catch (Throwable $e) {
    // handle
}
```

### Rule 2: Every catch block must log with stack trace

Every `error_log()` call inside a catch block that has access to `$e` **must** append the trace:

```
error_log($context . ' ' . $e->getMessage() . "\n" . $e->getTraceAsString());
```

Logging only `$e->getMessage()` without the trace is a **critical defect**.

### Rule 3: Use ErrorLogHelper for native logging

When FileLogger is not available (autoloader, bootstrap), use the `ErrorLogHelper` static class:

| Method | Behaviour |
|--------|-----------|
| `ErrorLogHelper::log($e, 'Context:')` | Logs message + trace to `error_log()` |
| `ErrorLogHelper::logAndThrow($e, 'Context:')` | Logs and re-throws (return type `never`) |

### Rule 4: safeExecute wraps all endpoints

Every public REST handler method must be wrapped in `$this->safeExecute()`. Direct try-catch in endpoint handlers is not allowed — delegate to the ResponseTrait infrastructure.

### Rule 5: Stack trace frames are debug-mode gated

Error responses include structured stack trace frames **only when debug mode is enabled**. In production, the `Errors.Backend` field is omitted entirely to prevent information leakage.

### Rule 6: No error swallowed — Forbidden Patterns

The following patterns are **critical defects**. They MUST never appear in any plugin codebase.

```php
// ❌ NEVER: Empty catch — error is silently lost
catch (Throwable $e) {
}

// ❌ NEVER: Catch without logging — error is silently lost
catch (Throwable $e) {
    return false;
}

// ❌ NEVER: Log message without stack trace
catch (Throwable $e) {
    error_log($e->getMessage());  // Missing: "\n" . $e->getTraceAsString()
}

// ❌ NEVER: Catch Exception instead of Throwable
catch (Exception $e) {   // Misses TypeError, Error, ParseError
    // ...
}

// ❌ NEVER: Swallow error in boolean/null return without logging
catch (Throwable $e) {
    return null;
}

// ❌ NEVER: Generic error_log without context prefix
catch (Throwable $e) {
    error_log($e->getMessage() . "\n" . $e->getTraceAsString());
    // Missing: '[PluginName] Context:' prefix
}
```

### Correct patterns for every scenario

| Scenario | Pattern |
|----------|---------|
| REST endpoint handler | `$this->safeExecute(fn() => ..., 'endpoint-name')` |
| Non-endpoint method with FileLogger | `ErrorResponse::logAndReturn($this->fileLogger, $e, 'Context')` |
| Non-endpoint returning false | `ErrorResponse::logAndReturnFalse($this->fileLogger, $e, 'Context')` |
| Non-endpoint returning WP_Error | `ErrorResponse::logAndReturnWpError($this->fileLogger, $e, 'Context')` |
| Bootstrap / no FileLogger | `ErrorLogHelper::log($e, '[PluginName] Context:')` |
| Infrastructure (must re-throw) | `ErrorLogHelper::logAndThrow($e, '[PluginName] Context:')` |

---

## 4.9 safeExecute() — Complete Specification

This is the most critical error-handling function in the plugin. It is defined in `ResponseTrait` and serves as the universal error boundary for all REST endpoints.

### Full implementation pattern

```
/**
 * Execute a callback with comprehensive error handling.
 *
 * @param callable $callback     The business logic to execute (must return WP_REST_Response)
 * @param string   $endpointName A human-readable name for logging (e.g., 'activate-plugin')
 *
 * @return WP_REST_Response Always returns a valid response, even on failure
 */
protected function safeExecute(
    callable $callback,
    string $endpointName,
): WP_REST_Response {
    try {
        return $callback();
    } catch (Throwable $e) {
        // Tier 1: Always log to PHP error_log (available even if FileLogger fails)
        error_log(
            "[{$this->getPluginSlug()}] safeExecute error in '{$endpointName}': "
            . $e->getMessage() . "\n"
            . $e->getTraceAsString()
        );

        // Tier 2: Log via FileLogger if available
        $hasLogger = ($this->fileLogger !== null);

        if ($hasLogger) {
            $this->fileLogger->logException($e, "safeExecute:{$endpointName}");
        }

        // Build error response with debug-mode gating
        return $this->buildErrorResponse($e, $endpointName);
    }
}
```

### buildErrorResponse — Debug-Mode Gating

```
/**
 * Build an error response with stack trace conditionally included.
 *
 * @param Throwable $e            The caught exception
 * @param string    $endpointName The endpoint name for context
 *
 * @return WP_REST_Response Formatted error envelope
 */
private function buildErrorResponse(
    Throwable $e,
    string $endpointName,
): WP_REST_Response {
    $isDebug = PluginConfigType::isDebugMode();

    // In debug mode: real message + stack trace
    // In production: generic message, no trace
    $errorMessage = $isDebug
        ? $e->getMessage()
        : 'An internal error occurred';

    $builder = EnvelopeBuilder::error($errorMessage, 500);

    if ($isDebug) {
        $builder->setStackTrace($this->formatStackFrames($e));
    }

    return $builder
        ->setRequestedAt($endpointName)
        ->toResponse();
}
```

### formatStackFrames — Structured Trace Extraction

```
/**
 * Extract structured stack trace frames from an exception.
 *
 * @param Throwable $e The exception to extract frames from
 *
 * @return array<int, string> Formatted trace lines
 */
private function formatStackFrames(Throwable $e): array
{
    $rawTrace = $e->getTraceAsString();
    $lines = explode("\n", $rawTrace);
    $frames = [];

    foreach ($lines as $line) {
        $trimmedLine = trim($line);
        $hasContent = ($trimmedLine !== '');

        if ($hasContent) {
            $frames[] = $trimmedLine;
        }
    }

    return $frames;
}
```

---

## 4.10 Structured Stack Trace Transport Format

When the Go backend or any REST consumer parses error responses, stack trace frames MUST use a structured array format — not just the raw string from `getTraceAsString()`. This enables structured display in error modals and log viewers.

### Frame Structure (PascalCase keys)

Each frame in the `StackTraceFrames` array contains:

```json
{
  "File": "/var/www/html/wp-content/plugins/my-plugin/includes/Traits/ActivateHandlerTrait.php",
  "FileBase": "ActivateHandlerTrait.php",
  "Line": 78,
  "Function": "executeActivation",
  "Class": "PluginName\\Traits\\Activate\\ActivateHandlerTrait"
}
```

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `File` | string | `$frame['file']` | Full absolute path |
| `FileBase` | string | `basename($frame['file'])` | Filename only — for compact display |
| `Line` | int | `$frame['line']` | Line number in source file |
| `Function` | string | `$frame['function']` | Method or function name |
| `Class` | string\|null | `$frame['class']` | Fully-qualified class name (null for global functions) |

### Debug mode ON — response with structured frames

```json
{
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 500,
    "Message": "Cannot connect to remote endpoint: Connection refused",
    "Timestamp": "2026-04-07T14:30:00Z"
  },
  "Attributes": {
    "RequestedAt": "/my-plugin-api/v1/activate",
    "TotalRecords": 0
  },
  "Results": [],
  "Errors": {
    "BackendMessage": "Cannot connect to remote endpoint: Connection refused",
    "ExceptionType": "RuntimeException",
    "Backend": [
      "#0 ActivateHandlerTrait.php(78): PluginName\\Traits\\Activate\\ActivateHandlerTrait->executeActivation()",
      "#1 ResponseTrait.php(35): PluginName\\Traits\\Core\\ResponseTrait->safeExecute()"
    ],
    "StackTraceFrames": [
      {
        "File": "/var/www/html/wp-content/plugins/my-plugin/includes/Traits/Activate/ActivateHandlerTrait.php",
        "FileBase": "ActivateHandlerTrait.php",
        "Line": 78,
        "Function": "executeActivation",
        "Class": "PluginName\\Traits\\Activate\\ActivateHandlerTrait"
      },
      {
        "File": "/var/www/html/wp-content/plugins/my-plugin/includes/Traits/Core/ResponseTrait.php",
        "FileBase": "ResponseTrait.php",
        "Line": 35,
        "Function": "safeExecute",
        "Class": "PluginName\\Traits\\Core\\ResponseTrait"
      }
    ]
  }
}
```

### Debug mode OFF — no Errors key at all

```json
{
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 500,
    "Message": "An internal error occurred",
    "Timestamp": "2026-04-07T14:30:00Z"
  },
  "Attributes": {
    "RequestedAt": "/my-plugin-api/v1/activate",
    "TotalRecords": 0
  },
  "Results": []
}
```

### Frame extraction implementation

```php
/**
 * Extract structured stack trace frames from an exception.
 *
 * @return array<int, array{File: string, FileBase: string, Line: int, Function: string, Class: string|null}>
 */
private function extractStructuredFrames(Throwable $e): array
{
    $trace = $e->getTrace();
    $frames = [];

    foreach ($trace as $frame) {
        $hasFile = isset($frame['file']);

        $frames[] = [
            'File'     => $hasFile ? $frame['file'] : '[internal]',
            'FileBase' => $hasFile ? basename($frame['file']) : '[internal]',
            'Line'     => $frame['line'] ?? 0,
            'Function' => $frame['function'] ?? '',
            'Class'    => $frame['class'] ?? null,
        ];
    }

    return $frames;
}
```

### Rules

1. `Errors.Backend` (string array) is ALWAYS included for backward compatibility when debug mode is ON
2. `Errors.StackTraceFrames` (object array) is the **preferred** format for structured consumers
3. Both fields are omitted entirely when debug mode is OFF
4. Frame extraction uses `$e->getTrace()` (structured), not `$e->getTraceAsString()` (string)
5. `FileBase` is always computed — never trust the consumer to parse paths

---

## 4.11 API Error Response Format — Additional Examples

### Debug mode ON — full details

```json
{
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 500,
    "Message": "Cannot connect to remote endpoint: Connection refused",
    "Timestamp": "2026-04-07T14:30:00Z"
  },
  "Attributes": {
    "RequestedAt": "/my-plugin-api/v1/activate",
    "TotalRecords": 0
  },
  "Results": [],
  "Errors": {
    "BackendMessage": "Cannot connect to remote endpoint: Connection refused",
    "ExceptionType": "RuntimeException",
    "Backend": [
      "#0 ActivateHandlerTrait.php(78): PluginName\\Traits\\Activate\\ActivateHandlerTrait->executeActivation()",
      "#1 ResponseTrait.php(35): PluginName\\Traits\\Core\\ResponseTrait->PluginName\\Traits\\Core\\{closure}()",
      "#2 ResponseTrait.php(42): PluginName\\Traits\\Core\\ResponseTrait->safeExecute()",
      "#3 WP_REST_Server.php(1181): WP_REST_Server->dispatch()",
      "#4 rest-api.php(407): rest_do_request()"
    ]
  }
}
```

### Debug mode OFF — safe for production

```json
{
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 500,
    "Message": "An internal error occurred",
    "Timestamp": "2026-04-07T14:30:00Z"
  },
  "Attributes": {
    "RequestedAt": "/my-plugin-api/v1/activate",
    "TotalRecords": 0
  },
  "Results": []
}
```

Note: No `Errors` key at all in production. The error is fully logged server-side (both Tier 1 and Tier 2), but the API consumer only sees a generic message. This prevents:
- Exposing internal file paths
- Leaking class/method names
- Revealing PHP version or WordPress internals

### 400-level errors — always include message (not sensitive)

Validation errors (400, 401, 403, 404) always include a descriptive message regardless of debug mode, because they contain no internal implementation details:

```json
{
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 400,
    "Message": "Missing required field: plugin_slug",
    "Timestamp": "2026-04-07T14:30:00Z"
  },
  "Attributes": {
    "RequestedAt": "/my-plugin-api/v1/activate",
    "TotalRecords": 0
  },
  "Results": []
}
```

---

## 4.12 ErrorLogHelper — Complete Specification

A minimal static class for Tier 1 logging when FileLogger is unavailable.

### Full implementation pattern

```
class ErrorLogHelper
{
    /**
     * Log an exception with full context and stack trace to PHP's error_log.
     *
     * @param Throwable $exception The caught exception
     * @param string    $context   Human-readable context (e.g., 'Autoloader:')
     */
    public static function log(Throwable $exception, string $context): void
    {
        $message = sprintf(
            '%s %s in %s:%d\n%s',
            $context,
            $exception->getMessage(),
            $exception->getFile(),
            $exception->getLine(),
            $exception->getTraceAsString(),
        );

        error_log($message);
    }

    /**
     * Log an exception and re-throw it. Use in infrastructure code
     * where silent failure causes cascading breakage.
     *
     * @param Throwable $exception The caught exception
     * @param string    $context   Human-readable context
     *
     * @throws Throwable Always re-throws the original exception
     */
    public static function logAndThrow(Throwable $exception, string $context): never
    {
        self::log($exception, $context);

        throw $exception;
    }
}
```

### When to use ErrorLogHelper vs FileLogger

| Scenario | Use |
|----------|-----|
| Inside autoloader (`vendor/autoload.php`) | `ErrorLogHelper::log()` |
| Inside `Plugin::boot()` before FileLogger init | `ErrorLogHelper::logAndThrow()` |
| Inside route registration (FileLogger may fail) | `ErrorLogHelper::log()` as fallback |
| Inside any trait handler method | `$this->fileLogger->logException()` |
| Inside a static helper class | `ErrorLogHelper::log()` |

---

## 4.13 Shutdown Handler (Fatal Errors)

Register a global shutdown handler to catch fatal errors that bypass try-catch:

1. Check `error_get_last()` for fatal error types
2. Log to a dedicated `fatal-errors.log` file (not through FileLogger, which may be compromised)
3. Include memory usage statistics (helps diagnose OOM kills)
4. Attempt JSON output if the response has not been sent

### Implementation pattern

```
register_shutdown_function(function (): void {
    $lastError = error_get_last();
    $hasError = ($lastError !== null);

    if (!$hasError) {
        return;
    }

    $fatalTypes = [E_ERROR, E_CORE_ERROR, E_COMPILE_ERROR, E_PARSE];
    $isFatal = in_array($lastError['type'], $fatalTypes, true);

    if (!$isFatal) {
        return;
    }

    $memoryUsage = memory_get_peak_usage(true);
    $memoryMb = round($memoryUsage / 1048576, 2);

    $logEntry = sprintf(
        "[FATAL] %s in %s:%d | Memory: %sMB | %s\n",
        $lastError['message'],
        $lastError['file'],
        $lastError['line'],
        $memoryMb,
        DateHelper::nowLogDisplay(),
    );

    // Write directly to file — FileLogger may be compromised
    $logPath = PathHelper::getLogsDir() . '/fatal-errors.log';
    file_put_contents($logPath, $logEntry, FILE_APPEND | LOCK_EX);
});
```

---

## 4.14 DateHelper — Timestamp Specification

All timestamps flow through a centralised `DateHelper` class:

| Method | Returns | Used for |
|--------|---------|----------|
| `nowUtc()` | `2026-04-07T14:30:00Z` | API responses, database storage |
| `nowIso()` | ISO 8601 with timezone | API metadata |
| `nowLogDisplay()` | `07-Apr-26 2:30 PM` | Log file entries |
| `formatInWpTimezone($format, $timestamp)` | Formatted string in WP timezone | All display timestamps |

### Timezone handling

- All storage is UTC
- All display converts to the WordPress-configured timezone (`Settings > General > Timezone`)
- The timezone is resolved once and cached for the request lifetime
- Supports both named timezones (`Asia/Kuala_Lumpur`) and GMT offset fallback

---

## 4.15 Complete Error Handling Flow — End to End

This shows exactly what happens when an exception occurs during an API request:

```
1. Client sends: POST /my-plugin-api/v1/activate { "plugin_slug": "some-plugin" }

2. WordPress routes to: ActivateHandlerTrait::handleActivate($request)

3. handleActivate() calls: $this->safeExecute(fn() => $this->executeActivation($request), 'activate')

4. executeActivation() throws: RuntimeException("Connection refused")

5. safeExecute() catches Throwable:
   a. Tier 1: error_log("[MyPlugin] safeExecute error in 'activate': Connection refused\n#0 ...")
   b. Tier 2: $this->fileLogger->logException($e, "safeExecute:activate")
      → Writes to info.log:    [07-Apr-26 2:31 PM v2.31.0] [Error] safeExecute:activate: Connection refused (ResponseTrait.php:35) {}
      → Writes to error.log:   [07-Apr-26 2:31 PM v2.31.0] [Error] safeExecute:activate: Connection refused (ResponseTrait.php:35) {}
      → Writes to stacktrace.log:
         ================================================================================
         [07-Apr-26 2:31 PM v2.31.0] safeExecute:activate: Connection refused (ResponseTrait.php:35)
         Exception: RuntimeException
         Message: Connection refused
         --------------------------------------------------------------------------------
         #0 ActivateHandlerTrait.php(78): ...->executeActivation()
         #1 ResponseTrait.php(35): ...->safeExecute()
         ================================================================================
   c. Calls buildErrorResponse($e, 'activate')
      → Checks PluginConfigType::isDebugMode()
      → If debug: includes Errors.Backend with trace frames + real message
      → If production: generic "An internal error occurred", no Errors key

6. Client receives: JSON envelope with Code 500
```
