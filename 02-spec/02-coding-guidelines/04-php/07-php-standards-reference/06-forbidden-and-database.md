# PHP Coding Standards — Forbidden patterns, database wrapper (AI Execution Prompt)

> **/goal** Catalog and enforce all prohibited PHP anti-patterns, and mandate typed database access via the `TypedQuery` wrapper and `DbResult` envelope hierarchy.
> **/learn** Master the comprehensive forbidden patterns catalog (catch Throwable, no magic strings, no nested if, 15-line limit, strict typing) and type-safe PDO query abstraction with row mappers.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Eliminate all forbidden patterns (magic strings, raw negations, untyped parameters/returns, nested `if`) across the codebase.
- [ ] `/learn` Never use raw PDO queries or `wpdb` calls directly in service layers; route queries through `TypedQuery`.
- [ ] `/goal` Enforce typed result envelopes (`DbResult<T>`, `DbResultSet<T>`, `DbExecResult`) for all database operations.
- [ ] `/learn` Use static domain model factories (e.g., `PluginInfo::fromRow($row)`) inside type-safe mapper closures.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [PHP Coding Standards](./readme.md)
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

```php
// ❌ FORBIDDEN: Catching Exception instead of Throwable, and using magic strings
try {
    do_action('custom_hook');
} catch (Exception $e) {
    error_log($e->getMessage());
}

// ✅ REQUIRED: Catching Throwable and using typed enum constants
try {
    do_action(HookType::CustomAction->value);
} catch (Throwable $throwableErr) {
    Logger::error('Failed custom action', ['error' => $throwableErr->getMessage()]);
}
```

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

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PHP-REF-006: PHP Forbidden Patterns Catalog and Typed Database Access

**Given** PHP standards reference files and companion plugin implementations.
**When** Audited against this reference specification.
**Then** Zero violations are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php/07-php-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.
