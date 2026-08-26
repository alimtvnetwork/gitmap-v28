# PHP Coding Standards — Forbidden patterns, database wrapper

> **Parent:** [PHP Coding Standards](./00-overview.md)  
> **Version:** 5.1.0  
> **Updated:** 2026-03-31

---

## Forbidden Patterns

| Pattern | Why | Alternative |
|---------|-----|-------------|
| `catch (Exception $e)` | Misses PHP 7+ `Error` types | `catch (Throwable $e)` |
| Magic strings in hooks | Unmaintainable, typo-prone | `HookType::*->value` enum cases |
| Inline concatenation at call site | Hard to read, duplicated | Compose a named constant first |
| Magic strings in handlers | Unmaintainable | `constants.php` |
| `wp_die()` in REST handlers | Breaks JSON responses | `wp_send_json_error()` |
| Manual path concatenation | Fragile paths | `PathHelper` fully-typed accessors |
| `getDataDir() . '/file.db'` | Partial accessor, still magic | Add a typed accessor to `PathHelper` |
| Constructor WordPress calls | Load order issues | Lazy initialization |
| `error_log()` for diagnostics | No structure | Use `FileLogger` / `Logger` |
| Inline `!class_exists('PDO')` checks | Duplicated logic | `ErrorChecker::isInvalidPdoExtension()` |
| Nested `if` | **Zero tolerance** — absolute ban | Flatten with early returns or combined conditions |
| Functions > 15 lines | Hard to read, test, review | Extract helpers |
| `return` without blank line after statements | Poor readability | Blank line before `return` |
| Single-line `if (...) return;` | Easy to miss, inconsistent | Always use braces `{ }` |
| Inline multi-part `if` condition (2+ operators) | Hard to read, not reusable | Extract to named `$is_*` variable or method |
| `BooleanHelpers::isFalsy/isTruthy/...` | Trivial wrappers (deprecated) | Native PHP operators |
| `!$obj->isActive()` | Easy to miss negation | `$obj->isDisabled()` |
| `!file_exists()` / `!is_dir()` | Raw negation | `isFileMissing()` / `isDirMissing()` |
| `current_user_can('manage_options')` | Magic string | `CapabilityType::ManageOptions->value` |
| `'POST'` in routes | Inconsistent | `HttpMethodType::Post->value` |
| Untyped function parameters | No runtime safety | Add type declarations (see [Strict Typing](../../01-cross-language/13-strict-typing.md)) |
| Untyped return values | No contract enforcement | Add return type declarations |
| Redundant `@param` on typed signatures | Noisy duplication | Remove; keep summary only (see [Strict Typing](../../01-cross-language/13-strict-typing.md)) |
| Boolean flag changing operation meaning | Unreadable call sites | Split into named methods (see [Function Naming](../../01-cross-language/10-function-naming.md)) |

---


---

## Database Wrapper — `TypedQuery`

All database queries SHOULD use the generic `TypedQuery` class. It wraps `PDO` and returns typed result envelopes with automatic stack traces.

### Result Types

| Class | Purpose | Key Methods |
|-------|---------|-------------|
| `DbResult<T>` | Single-row query | `isDefined()`, `isEmpty()`, `hasError()`, `isSafe()`, `value()`, `error()`, `stackTrace()` |
| `DbResultSet<T>` | Multi-row query | `hasAny()`, `isEmpty()`, `count()`, `hasError()`, `isSafe()`, `items()`, `first()`, `error()`, `stackTrace()` |
| `DbExecResult` | INSERT/UPDATE/DELETE | `isEmpty()`, `hasError()`, `isSafe()`, `affectedRows()`, `lastInsertId()`, `error()`, `stackTrace()` |

### Usage

```php
$tq = new TypedQuery($pdo);

// Single row — returns DbResult<PluginInfo>
$result = $tq->queryOne(
    'SELECT * FROM plugins WHERE id = :id',
    [':id' => $id],
    fn(array $row): PluginInfo => PluginInfo::fromRow($row),
);

if ($result->hasError()) { /* handle */ }

if ($result->isEmpty()) { /* not found */ }
$plugin = $result->value();

// Multiple rows — returns DbResultSet<SiteInfo>
$set = $tq->queryMany(
    'SELECT * FROM sites ORDER BY name',
    [],
    fn(array $row): SiteInfo => SiteInfo::fromRow($row),
);

foreach ($set->items() as $site) { /* ... */ }

// Exec — returns DbExecResult
$res = $tq->exec('DELETE FROM plugins WHERE id = :id', [':id' => $id]);

if ($res->hasError()) { /* handle */ }
echo $res->affectedRows();
```

### Mapper Closures

Callers provide a `Closure(array): T` mapper for type-safe row mapping (equivalent to Go's scanner functions). Use static `fromRow()` factory methods on domain models for consistency.

---
