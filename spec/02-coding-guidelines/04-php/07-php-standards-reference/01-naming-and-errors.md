# PHP Coding Standards — Naming conventions, error handling, structured responses

> **Parent:** [PHP Coding Standards](./00-overview.md)  
> **Version:** 5.1.0  
> **Updated:** 2026-03-31

---

## Naming Conventions

> **Baseline:** [PSR-12 PHP Naming Conventions](../03-naming-conventions.md)  
> The table below lists **project-specific overrides** that take precedence over PSR-12 defaults.

| Element | Convention | Example | Override reason |
|---------|-----------|---------|-----------------|
| Class names | PascalCase | `EnvelopeBuilder`, `SnapshotFactory` | _(matches PSR-12)_ |
| Method names | camelCase | `buildResponse()`, `getPluginInfo()` | Internal consistency (overrides WordPress snake_case) |
| Constants | UPPER_SNAKE_CASE (no `RISEUP_` prefix) | `REST_NAMESPACE`, `ACTION_UPLOAD` | _(matches PSR-12)_ |
| File names (classes) | `{PascalCase}.php` (PSR-4) | `EnvelopeBuilder.php`, `SnapshotFactory.php` | PSR-4 autoloading |
| File names (enums) | `{DefinitionName}Type.php` (PascalCase, PSR-4) | `UploadSourceType.php`, `CapabilityType.php` | PSR-4 in `includes/Enums/` |
| Variables | camelCase | `$pluginSlug`, `$stackTraceFrames` | _(matches PSR-12)_ |
| Enum types | PascalCase, **`Type` suffix required** | `UploadSourceType`, `CapabilityType`, `HttpMethodType` | PHP 8.1+ native backed enums |

> **PSR-4 file naming convention:**
> - **All classes** under `includes/` use PascalCase filenames matching the class name (e.g., `EnvelopeBuilder.php`, `SnapshotFactory.php`)
> - **Enums** in `includes/Enums/` use `{DefinitionName}Type.php` — PascalCase with `Type` suffix, no prefix, no hyphens, no underscores (e.g., `UploadSourceType.php` contains `enum UploadSourceType: string`)

---


---

## Error Handling — Safe Execution Strategy

### Rule: Catch `Throwable`, not just `Exception`

PHP 7+ introduces `Error` and `TypeError` that are **not** subclasses of `Exception`. All endpoint handlers must catch `Throwable`:

```php
// ❌ FORBIDDEN: Misses PHP 7+ Errors (e.g., missing class)
try {
    $result = $manager->process();
} catch (Exception $e) {
    wp_send_json_error($e->getMessage());
}

// ✅ REQUIRED: Catches all throwables
try {
    $result = $manager->process();

} catch (Throwable $e) {
    $this->logger->logException($e, 'process_failed');
    wp_send_json_error([
        'message'          => $e->getMessage(),
        'stackTrace'       => $e->getTraceAsString(),
        'stackTraceFrames' => $this->formatStackFrames($e),
    ], 500);
}
```

### Safe Execute Wrapper

All REST endpoint handlers must be wrapped in `safeExecute`:

```php
// ✅ Pattern: safeExecute wrapper
public function handleUpload(WP_REST_Request $request): WP_REST_Response {
    return $this->safeExecute(function() use ($request) {
        // Business logic here
        return $this->envelope->success($result);
    });
}

private function safeExecute(callable $callback): WP_REST_Response {
    try {
        return $callback();
    } catch (Throwable $e) {
        $this->logger->logException($e, 'endpoint_error');

        return $this->envelope->error($e->getMessage(), 500);
    }
}
```

### Global Shutdown Handler

Register a shutdown handler to catch fatal errors. **Delegate the type-check to `ErrorChecker`** which uses `ErrorType::FATAL_TYPES` (see [PHP Enum Spec](../01-enums.md)):

```php
// ❌ FORBIDDEN: Inline error-type checking
register_shutdown_function(function() {
    $error = error_get_last();
    if ($error && in_array($error['type'], [E_ERROR, E_CORE_ERROR, E_COMPILE_ERROR])) {
        // ...
    }
});

// ✅ REQUIRED: Use ErrorChecker for readable, centralized fatal-error detection
register_shutdown_function(function() {
    $error = error_get_last();
    if (ErrorChecker::isFatalError($error)) {
        // Log to fatal-errors.log via PathHelper::getFatalErrorLog()
        // Include memory_get_peak_usage() for diagnostics
        // Send JSON response before process dies (if REST_REQUEST)
    }
});
```

> **Implementation:** `ErrorChecker` delegates to `ErrorType::FATAL_TYPES` for the constant list. Use `ErrorChecker::getTypeLabel($error['type'])` to convert any `E_*` integer to a human-readable string (e.g., `'E_ERROR'`) — this replaces all inline type-mapping arrays. See [enums.md](../01-enums.md) for the full `ErrorChecker`, `ErrorType`, and `TYPE_LABELS` implementations.

---


---

## Structured Error Responses

### Required Fields

Every error response must include:

```json
{
  "message": "Human-readable error description",
  "StackTrace": "Full trace as string (debug_backtrace with unlimited depth)",
  "StackTraceFrames": [
    {
      "file": "/path/to/file.php",
      "line": 42,
      "function": "methodName",
      "class": "ClassName"
    }
  ]
}
```

### Stack Trace Logging

The logger captures two outputs for every error:

1. **Structured frames** — `stackTraceFrames` array in JSON responses
2. **Raw backtrace** — Written to `stacktrace.txt` with `debug_backtrace(0, 0)` (unlimited depth)

```php
// ✅ Dual logging: structured + raw
public function logException(Throwable $e, string $context = '') {
    // Structured frames for JSON responses
    $frames = $this->formatStackFrames($e);
    
    // Raw backtrace to file (unlimited depth)
    $backtrace = debug_backtrace(0, 0);
    file_put_contents($this->stacktraceFile, $this->formatBacktrace($backtrace), FILE_APPEND);
}
```

### Cross-References

- [Boolean Flag Method Splitting](../../01-cross-language/24-boolean-flag-methods.md) — Split bool-flag methods into two named methods (PHP examples included)

---
