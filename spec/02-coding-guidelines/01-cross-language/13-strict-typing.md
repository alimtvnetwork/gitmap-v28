# Strict Typing — Cross-Language Type Declaration Rules

> **Version:** 1.0.0  
> **Updated:** 2026-03-09  
> **Applies to:** PHP, TypeScript, Go

---

## Rule

Every function parameter, return value, and class property **must** have an explicit type declaration wherever the language supports it. Untyped signatures are forbidden in new code.

---

## PHP (7.4+ / 8.0+)

### Required Type Locations

| Location | Required | Since |
|----------|----------|-------|
| Function parameters | ✅ | PHP 7.0 |
| Return types | ✅ | PHP 7.0 |
| Class properties | ✅ | PHP 7.4 |
| Constructor promotion | ✅ (preferred) | PHP 8.0 |
| Nullable types (`?Type`) | ✅ | PHP 7.1 |
| Union types (`int\|string`) | ✅ | PHP 8.0 |

### Examples

```php
// ❌ FORBIDDEN: Untyped parameters and return
public function logException($e, $context = '') { ... }

// ✅ REQUIRED: Fully typed
public function logException(Throwable $e, string $context = ''): bool { ... }
```

```php
// ❌ FORBIDDEN: Untyped class property
class Manager {
    private $logger;
    private $initialized;
}

// ✅ REQUIRED: Typed properties
class Manager {
    private RiseupFileLogger $logger;
    private bool $isInitialized;
}
```

```php
// ✅ Constructor promotion (PHP 8.0+)
class User {
    public function __construct(
        public string $name,
        public int $age
    ) {}
}
```

### Limitation

PHP does not support type declarations on local variables. Types apply to parameters, returns, properties, and constants only.

---

## TypeScript

Already enforced by the generics-first rule and `strict: true` in tsconfig. Key reinforcements:

- `any` is **prohibited** everywhere (see [TypeScript Standards](../02-typescript/08-typescript-standards-reference.md))
- `unknown` only at parse boundaries with immediate narrowing
- All function signatures must have explicit parameter and return types

---

## Go

Already statically typed. Key reinforcements:

- `interface{}` / `any` is **prohibited** in exported APIs (see [Go Standards](../03-golang/04-golang-standards-reference/00-overview.md))
- Use concrete types or constrained generics (`[T any]` in generic signatures is acceptable)
- All struct fields must use concrete types, not `map[string]interface{}`
- **Single return value:** Functions return `apperror.Result[T]` — never `(T, error)` or multi-value tuples (see [§7.1](./15-master-coding-guidelines/04-type-safety.md#71-single-return-value-rule-go))
- **No type assertions:** `.(Type)` casting is prohibited in business logic — use concrete structs (see [§7.2](./15-master-coding-guidelines/04-type-safety.md#72-no-type-assertions--casting-go))

---

## Docblock Rules

### Rule: Remove redundant `@param`/`@return` when types are declared

When the function signature already declares types, docblock `@param` and `@return` annotations that merely repeat the type are **redundant** and must be removed.

### When to Keep Docblocks

| Condition | Keep docblock? |
|-----------|---------------|
| Function body > 10 lines | ✅ Keep a summary comment |
| Complex behavior or side effects | ✅ Describe semantics |
| Non-obvious constraints (e.g., "must be positive") | ✅ Document constraint |
| Parameter type is already in signature, no extra semantics | ❌ Remove |
| Return type is already in signature, no extra semantics | ❌ Remove |

### Examples

```php
// ❌ FORBIDDEN: Redundant docblock duplicating types
/**
 * Log exception and return error array.
 *
 * @param RiseupFileLogger $logger  File logger instance.
 * @param Throwable        $e       The caught exception.
 * @param string           $context Context message.
 * @return array Error response array.
 */
public static function logAndReturn(
    RiseupFileLogger $logger,
    Throwable $e,
    string $context = '',
): array {

// ✅ REQUIRED: Brief summary only (types are in signature)
/** Log exception and return standardized error array. */
public static function logAndReturn(
    RiseupFileLogger $logger,
    Throwable $e,
    string $context = '',
): array {
```

```typescript
// ❌ FORBIDDEN: JSDoc duplicating TypeScript types
/**
 * @param message - The error message
 * @returns The formatted string
 */
function formatError(message: string): string { ... }

// ✅ REQUIRED: No redundant JSDoc
function formatError(message: string): string { ... }
```

---

## Parameter Count Rule

### Rule: Maximum 3 parameters per function

Functions must accept **3 or fewer** parameters. When more are needed, group them into a typed object or class.

### Exception

Utility, framework, or infrastructure functions (e.g., static helpers, middleware wrappers) may exceed 3 parameters when each parameter serves a distinct, well-understood role. Always confirm before adding a 4th parameter.

### Examples

```php
// ❌ FORBIDDEN: Too many parameters
public function createPost(string $title, string $content, string $status, int $authorId, array $meta): int { ... }

// ✅ REQUIRED: Grouped into a typed object
public function createPost(CreatePostParams $params): int { ... }
```

```typescript
// ❌ FORBIDDEN
function createUser(name: string, email: string, role: string, department: string): User { ... }

// ✅ REQUIRED
interface CreateUserParams {
  name: string;
  email: string;
  role: string;
  department: string;
}
function createUser(params: CreateUserParams): User { ... }
```

```go
// ❌ FORBIDDEN
func CreateUser(name string, email string, role string, dept string) apperror.Result[User] { ... }

// ✅ REQUIRED
type CreateUserParams struct {
    Name  string
    Email string
    Role  string
    Dept  string
}
func CreateUser(params CreateUserParams) apperror.Result[User] { ... }
```

---

## Result Guard Rule (Zero Silent Failures)

Every typed Result wrapper — `apperror.Result[T]` / `ResultSlice[T]` / `ResultMap[K, V]` (Go) or `DbResult` / `DbResultSet` / `DbExecResult` (PHP) — **MUST** have its error state checked before accessing the contained value. Accessing `.value()` / `.Value()` without a prior `hasError()` / `HasError()` or `isSafe()` / `IsSafe()` guard is a **spec violation**.

**Principle:** No error may ever be swallowed. If a result carries an error, it must be explicitly handled — logged, returned, or propagated.

### PHP — DbResult / DbResultSet / DbExecResult

```php
// ❌ WRONG: No guard — error silently swallowed
$result = $query->queryOne(...);
$result->value();

// ✅ CORRECT: Guard before access
$result = $query->queryOne(...);

if ($result->hasError()) {
    $this->logger->logException($result->error(), 'context');

    return null;
}

return $result->value();
```

```php
// ❌ WRONG: No guard on collection
$results = $query->queryAll(...);
foreach ($results->items() as $row) { ... }

// ✅ CORRECT: Guard before iteration
$results = $query->queryAll(...);

if ($results->hasError()) {
    $this->logger->logException($results->error(), 'query failed');

    return [];
}

return $results->items();
```

```php
// ❌ WRONG: No guard on write result
$execResult = $query->execute(...);
$execResult->affectedRows();

// ✅ CORRECT: Guard before access
$execResult = $query->execute(...);

if ($execResult->hasError()) {
    $this->logger->logException($execResult->error(), 'execute failed');

    return false;
}

return $execResult->affectedRows() > 0;
```

### Go — Propagation Rules

> `.AppError()` returns `*AppError` — stack trace and diagnostic context are always preserved. Named `AppError()` (not `Error()`) to avoid confusion with Go's native `error` interface.

```go
// ✅ Same-type → direct return (Result, ResultSlice, ResultMap)
result := svc.GetById(ctx, id)            // Result[Plugin]

if result.HasError() { return result }     // no re-wrapping
plugin := result.Value()

// ✅ Cross-type → Fail/FailSlice/FailMap IS needed
plugins := s.pluginService.List(ctx)       // ResultSlice[Plugin]

if plugins.HasError() {
    return apperror.FailSlice[SyncResult](plugins.AppError())
}

// ❌ WRONG — redundant (same type re-wrapped)
if result.HasError() { return apperror.Fail[Plugin](result.AppError()) }

// ✅ Collection access — early-return guard, then iterate
if result.HasError() { return result }

for _, item := range result.Items() { process(item) }

// ✅ Adapter unwrap — Result[T] → (*T, error) // EXEMPTED: framework boundary adapter
func (a *Adapter) GetById(ctx context.Context, id int64) (*models.Plugin, error) {
    result := a.Service.GetById(ctx, id)

    if result.HasError() { return nil, result.AppError() }
    v := result.Value()

    return &v, nil
}
```

> **Full examples with PHP/Go/TypeScript:** see [apperror § Result Guard Rule](../../03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/06-serialization-and-guards.md#12-result-guard-rule--mandatory-error-check-before-value-access)

### Enforcement Checklist

- [ ] Every `result.Value()` / `$result->value()` call is preceded by `HasError()` / `hasError()` or `IsSafe()` / `isSafe()`
- [ ] Every `result.Items()` / `$results->items()` call is preceded by a guard
- [ ] Every `result.Get(key)` on `ResultMap` is preceded by a guard
- [ ] Every `$execResult->affectedRows()` on `DbExecResult` is preceded by a guard
- [ ] No error is silently discarded — all errors are logged, returned, or propagated
- [ ] Cross-service callers guard results the same way

---

## Cross-References

- [PHP Standards](../04-php/07-php-standards-reference/00-overview.md)
- [TypeScript Standards](../02-typescript/08-typescript-standards-reference.md)
- [Go Standards](../03-golang/04-golang-standards-reference/00-overview.md)
- [Function Naming](./10-function-naming.md)
- Generic Enforce <!-- external: spec/31-generic-enforce/00-overview.md -->
- apperror Package — Result Guard Rule <!-- external: spec/03-error-manage/01-error-resolution/10-apperror-package/01-apperror-reference.md -->
- [Master Guidelines — Section 6.1](./15-master-coding-guidelines/03-code-style-and-errors.md#61--result-guard-rule-zero-silent-failures)

---

*Strict typing specification v1.1.0 — 2026-02-23*
