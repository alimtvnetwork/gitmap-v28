# Function Naming — No Boolean Flag Parameters

> **Version:** 1.0.0  
> **Updated:** 2026-03-09  
> **Applies to:** PHP, TypeScript, Go

---

## Rule

When a boolean parameter changes the **meaning** of an operation (not a minor detail), split the behavior into separate explicitly named methods. The call site must communicate intent without inspecting parameter values.

### Exception

Boolean parameters are acceptable when they adjust a **minor detail** within the same operation (e.g., a formatting option that does not alter the core behavior). When in doubt, split.

---

## Anti-Pattern

```php
// ❌ FORBIDDEN: Boolean flag hides intent
$logger->log("Payment failed", true);   // What does `true` mean?
$logger->log("User saved", false);      // What does `false` mean?
```

The call site is unreadable. The reader must inspect the method signature to understand the behavior.

---

## Preferred Pattern: Explicitly Named Methods

Split the behavior so the method name describes what happens.

### PHP

```php
final class Logger
{
    public function info(string $message): void
    {
        error_log($message);
    }

    public function infoWithTrace(string $message): void
    {
        $trace = (new \Exception())->getTraceAsString();
        error_log($message . "\n" . $trace);
    }
}

// ✅ Call sites are self-documenting
$logger->info("User saved");
$logger->infoWithTrace("Payment failed");
```

### TypeScript

```typescript
// ❌ FORBIDDEN
function logMessage(message: string, includeStack: boolean): void { ... }
logMessage("Failed", true);

// ✅ REQUIRED
function logMessage(message: string): void { ... }
function logMessageWithStack(message: string): void { ... }

logMessage("User saved");
logMessageWithStack("Payment failed");
```

### Go

```go
// ❌ FORBIDDEN
func LogMessage(message string, includeStack bool) { ... }
LogMessage("Failed", true)

// ✅ REQUIRED
func LogMessage(message string) { ... }
func LogMessageWithStack(message string) { ... }

LogMessage("User saved")
LogMessageWithStack("Payment failed")
```

---

## Real-World Example: ErrorResponse

```php
// ✅ Two clearly named static methods
final class ErrorResponse
{
    // Standard: uses exception's native file/line
    public static function logAndReturn(
        FileLogger $logger,
        Throwable $e,
        string $context = ''
    ): array { ... }

    // With trace: uses debug_backtrace() with frame skipping
    public static function logAndReturnWithTrace(
        FileLogger $logger,
        Throwable $e,
        string $context = '',
        int $skipFrames = 1
    ): array { ... }
}

// ✅ Call sites
return ErrorResponse::logAndReturn($this->fileLogger, $e, 'List posts exception');

return ErrorResponse::logAndReturnWithTrace($this->fileLogger, $e, 'Middleware error', 2);
```

---

## Guidelines

| Situation | Action |
|-----------|--------|
| Boolean changes operation meaning | Split into named methods |
| Boolean adjusts a minor formatting detail | Acceptable as parameter (confirm first) |
| More than one boolean flag | Always split — combinatorial APIs are unreadable |
| Flag count growing over time | Refactor to named methods immediately |

---

## Cross-References

- [PHP Standards](../04-php/07-php-standards-reference/00-overview.md)
- [TypeScript Standards](../02-typescript/08-typescript-standards-reference.md)
- [Go Standards](../03-golang/04-golang-standards-reference/00-overview.md)
- [Cross-Language Code Style](./04-code-style/00-overview.md)

---

*Function naming specification v1.0.0 — 2026-02-14*
