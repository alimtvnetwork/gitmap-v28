# Consolidated Code Review Guide

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Scope:** All Languages (Go, TypeScript, PHP, C#, Rust)

---

## Table of Contents

1. [Workflow & Process](#1-workflow--process)
2. [Function & File Size](#2-function--file-size)
3. [Parameters & Returns](#3-parameters--returns)
4. [Naming Conventions](#4-naming-conventions)
5. [Boolean & Conditionals](#5-boolean--conditionals)
6. [Enums & Constants](#6-enums--constants)
7. [Error Handling](#7-error-handling)
8. [Type Safety](#8-type-safety)
9. [Parallel Execution](#9-parallel-execution)
10. [Database Conventions](#10-database-conventions)
11. [SQL Safety](#11-sql-safety)
12. [Logging](#12-logging)
13. [Security (OWASP)](#13-security-owasp)

---

## 1. Workflow & Process

- **Spec first, code second.** Before any task, write/update the spec document.
- **Bug? Issue file first.** Create `spec/.../03-issues/{NN}-{name}.md` with root cause analysis, then fix.
- **Folder structure matters.** Follow `00-overview.md` + `99-consistency-report.md` convention.

---

## 2. Function & File Size

| Metric | Limit | Hard Max |
|--------|-------|----------|
| Function body | 8–15 lines | 15 lines |
| File length | 200–300 lines | 400 lines (with `// NOTE: Needs refactor`) |
| Type/struct definition | — | 120 lines |

- If a function exceeds 15 lines → extract helpers.
- If a file exceeds 300 lines → split by concern (`*_crud`, `*_helpers`, `*_validation`).

```go
// ❌ 30-line function
func ProcessOrder(order Order) Result {
    // ... 30 lines of mixed logic
}

// ✅ Top-level orchestrator + small helpers
func ProcessOrder(order Order) Result {
    validated := validateOrder(order)
    enriched := enrichOrder(validated)

    return saveOrder(enriched)
}
```

---

## 3. Parameters & Returns

- **Max 3 parameters** per function. More → use an options struct/object.
- **Single return value.** Use a Result/wrapper type, never multiple loose returns.

```go
// ❌ Too many params + multi-return
func CreateUser(name string, email string, age int, role string, org string) (User, error) { ... }

// ✅ Options struct + Result
func CreateUser(opts CreateUserInput) apperror.Result[User] { ... }
```

```typescript
// ❌ Multiple params
function sendEmail(to: string, subject: string, body: string, cc: string, bcc: string): void { ... }

// ✅ Options object
function sendEmail(opts: SendEmailInput): Result<void> { ... }
```

---

## 4. Naming Conventions

### Universal Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Classes / Types | PascalCase | `OrderService`, `UserRole` |
| Functions / Methods | camelCase (TS/PHP), PascalCase (Go exported) | `getUser()`, `GetUser()` |
| Variables | camelCase | `orderTotal`, `userName` |
| Constants | UPPER_SNAKE_CASE or PascalCase | `MAX_RETRIES`, `DefaultTimeout` |
| Files (Go) | PascalCase | `OrderService.go` |
| Files (TS) | kebab-case or PascalCase | `order-service.ts` |
| Enum types | PascalCase + `Type` suffix | `StatusType`, `RoleType` |

### Abbreviation Overrides (PascalCase Always)

| Abbreviation | Correct | Wrong |
|-------------|---------|-------|
| ID | `UserId`, `OrderId` | `userID`, `orderID` |
| URL | `BaseUrl`, `ApiUrl` | `baseURL`, `apiURL` |
| JSON | `JsonPayload` | `JSONPayload` |
| API | `ApiClient` | `APIClient` |
| HTTP | `HttpMethod` | `HTTPMethod` |
| SQL | `SqlQuery` | `SQLQuery` |

> Only first letter capitalized for abbreviations — consistent across all languages.

### Boolean Naming (All Languages — TS, Go, PHP, C#, Rust)

Every boolean variable **and** every function/method that returns a boolean **must** start with `is` or `has` (99% of cases). Use `should` only when expressing a recommendation or preference (e.g., `shouldRetry`). Never use `can`, `was`, or `will`.

**Never use `not`, `no`, or any negative word** in a boolean name. Use the semantic inverse instead.

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `active` | `isActive` | Missing prefix |
| `loaded` | `isLoaded` | Missing prefix |
| `visible` | `isVisible` | Missing prefix |
| `error` | `hasError` | Missing prefix |
| `checkPermission()` | `hasPermission()` | Function returns bool — needs prefix |
| `validateToken()` | `isTokenValid()` | Function returns bool — needs prefix |
| `existsInDb()` | `isPersistedInDb()` | Function returns bool — needs prefix |

### Semantic Inverse Pairs (Never Use `not` or `no`)

When you need the opposite of a boolean, **do not negate** — use the inverse name:

| ❌ Negative (banned) | ✅ Positive Inverse |
|----------------------|-------------------|
| `isNotActive` | `isInactive` |
| `isNotReady` | `isPending` |
| `isNotValid` | `isInvalid` |
| `isNotFound` | `isMissing` |
| `isNotAllowed` | `isForbidden` |
| `isNotConnected` | `isDisconnected` |
| `hasNoPermission` | `isUnauthorized` |
| `isNotEmpty` | `hasContent` |
| `isNotComplete` | `isIncomplete` |
| `isNotEnabled` | `isDisabled` |
| `notBlocked` | `isBlocked` / `isAllowed` |
| `noResults` | `isResultEmpty` |

> **Rule of thumb:** If you're typing `not`, `no`, or `!` in a name — stop and find the inverse word.

---

## 5. Boolean & Conditionals

### Rules

1. **No nested `if`** — flatten with early returns (guard clauses).
2. **No negation (`!`)** on function calls — use semantic inverse methods.
3. **Max 2 conditions** per `if` expression.
4. **Never mix `&&` and `||`** in one expression.
5. **Never mix positive and negative** (`isX && !isY` is forbidden).
6. **No `else` after `return`/`throw`/`continue`/`break`**.
7. **Keep cyclomatic complexity low** — one path through the function, guard and exit early.

### Zero Nested `if` — The Inverse Guard Pattern

The idea: **invert the condition, exit the function, continue on the happy path.** This keeps every function a single straight line — no indentation pyramids, no cognitive load.

```go
// ❌ CODE RED — Nested if = pyramid of doom
func ProcessPlugin(ctx context.Context, id int64) apperror.Result[Plugin] {
    plugin := s.repo.GetById(ctx, id)
    if plugin.IsSafe() {
        if plugin.Value().IsActive() {
            if plugin.Value().HasValidLicense() {
                return buildOutput(plugin.Value())
            }
        }
    }

    return apperror.FailNew[Plugin]("E9001", "processing failed")
}

// ✅ REQUIRED — Flat guards, inverse conditions, early exit
func ProcessPlugin(ctx context.Context, id int64) apperror.Result[Plugin] {
    plugin := s.repo.GetById(ctx, id)

    if plugin.HasError() {
        return apperror.Fail[Plugin](plugin.AppError())
    }

    if plugin.Value().IsInactive() {
        return apperror.FailNew[Plugin]("E9002", "plugin is inactive")
    }

    if plugin.Value().HasInvalidLicense() {
        return apperror.FailNew[Plugin]("E9003", "license is invalid")
    }

    return buildOutput(plugin.Value())
}
```

```typescript
// ❌ CODE RED — Nested conditions
function processOrder(order: Order | null): Result<Receipt> {
    if (order) {
        if (order.isValid()) {
            if (order.items.length > 0) {
                return createReceipt(order);
            }
        }
    }

    return fail("invalid order");
}

// ✅ REQUIRED — Guard, exit, continue
function processOrder(order: Order | null): Result<Receipt> {
    if (!order) {
        return fail("order is missing");
    }

    if (order.isInvalid()) {
        return fail("order is invalid");
    }

    if (order.items.length === 0) {
        return fail("order has no items");
    }

    return createReceipt(order);
}
```

### Boolean Complexity

```go
// ❌ Mixed operators — impossible to read
if (isReady && hasPermission) || (isAdmin && !isBlocked) { ... }

// ✅ Named booleans — one concern each
isAuthorizedUser := isReady && hasPermission
isPrivilegedAdmin := isAdmin && isUnblocked

if isAuthorizedUser {
    process()
}
if isPrivilegedAdmin {
    process()
}
```

```typescript
// ❌ Complex condition
if (isReady && !isBlocked && (hasPermission || isAdmin)) { ... }

// ✅ Extract to named boolean
const canProceed = isReady && hasPermission;
const isPrivileged = isReady && isAdmin;

if (!canProceed && !isPrivileged) {
    return;
}
```

---

## 6. Enums & Constants

- **Groups of related strings → Enum** (language-appropriate pattern).
- **Standalone values → Named constant.**
- **Never use magic strings or magic numbers.**
- **Exceptions:** `0`, `1`, `-1`, `""`, `true`, `false`, `null`/`nil`.

```go
// ❌ Magic string
if status == "active" { ... }

// ✅ Enum constant
if status.Is(StatusActive) { ... }
```

```typescript
// ❌ Magic number
if (retries > 3) { ... }

// ✅ Named constant
const MAX_RETRIES = 3;
if (retries > MAX_RETRIES) { ... }
```

---

## 7. Error Handling

### Code-Red Rules

> 🔴 **CODE RED — ZERO TOLERANCE: Never swallow an error.**
> An empty `catch {}`, a bare `return nil`, or an ignored result is an **automatic rejection**.
> Every error — without exception — must be **logged, returned, or rethrown**.
> Silent failures are the #1 cause of production incidents that take hours to debug.

1. **Never swallow an error.** Every error must be logged or returned to the caller. Empty `catch` blocks, bare `return nil` after an error check, and ignored `Result` values are all **Code Red violations**.
2. **Always include full context:** message, stack trace, source location.
3. **Go: Always capture stack trace** via `apperror.New()` / `apperror.Wrap()`. Without a stack trace, debugging is impossible.
4. **No `fmt.Errorf()`** in Go — use `apperror.Wrap()`.
5. **PHP: Catch `Throwable`**, not just `Exception`.
6. **Structured error responses:** message + stack trace + frames.
7. **Cache errors appropriately.** Failed lookups that are cached without error state cause **cascading silent failures** — always cache the error alongside the result or invalidate the cache entry on error.

```go
// ❌ CODE RED — Silent swallow (error disappears)
result, err := service.Process()
if err != nil {
    return nil
}

// ❌ CODE RED — No stack trace, useless in production
return fmt.Errorf("process failed: %w", err)

// ❌ CODE RED — errors.New has no stack trace
return errors.New("something failed")

// ✅ REQUIRED — Wrap at first contact with full stack trace
return apperror.Wrap(err, "E5001", "service.Process failed")
```

```typescript
// ❌ Silent catch — CODE RED
try { await fetchData(); } catch (e) { /* nothing */ }

// ✅ Always log or rethrow
try {
    await fetchData();
} catch (error) {
    logger.error("fetchData failed", { error, context: "UserService" });
    throw error;
}
```

### Go Stack Trace — The Wrap-Immediately Pattern

Every error from stdlib or third-party **must be wrapped immediately** into `apperror`. After wrapping, propagate as `AppError` — never re-wrap.

```go
// ✅ REQUIRED — Wrap stdlib error at first contact
func (s *PluginService) Upload(ctx context.Context, req UploadRequest) apperror.Result[UploadResult] {
    // Framework boundary — raw error, wrap immediately
    file, err := os.Open(req.Path)

    if err != nil {
        return apperror.FailWrap[UploadResult](err, "E4001", "failed to open plugin file")
    }

    defer file.Close()

    // Application boundary — already AppError, propagate via .AppError()
    meta := s.metaService.GetById(ctx, req.PluginId)

    if meta.HasError() {
        return apperror.Fail[UploadResult](meta.AppError())
    }

    return apperror.Ok(UploadResult{Plugin: meta.Value(), File: file})
}
```

**Stack trace gives you:**
- `err.FullString()` → code + message + full stack + cause chain
- `err.ToClipboard()` → markdown-formatted for AI paste
- `err.CallerLine()` → `"Upload.go:42"` compact reference

**Result pattern:** Always `apperror.Result[T]` — check `HasError()` before `.Value()`.

---

## 8. Type Safety

### Absolute Ban

| Forbidden | Allowed Alternative |
|-----------|-------------------|
| `any` (TS) | Explicit type or generic `<T>` |
| `interface{}` / `any` (Go) | Concrete struct or `[T any]` generic |
| `unknown` (TS) | Only at parse boundaries with type guard |
| `object` (TS) | Typed interface |
| Type assertions (Go) | `typecast.CastOrFail[T]()` or concrete types |

### Use Generics Instead

```go
// ❌ FORBIDDEN — any in exported API
func ProcessData(data interface{}) interface{} { ... }
func FetchResults() (any, error) { ... }

// ✅ REQUIRED — Concrete types or generics
func ProcessData(data PluginDetails) apperror.Result[PluginSummary] { ... }
func FetchResults[T any]() apperror.Result[T] { ... }
```

```typescript
// ❌ FORBIDDEN — any loses all type safety
function process(data: any): any { ... }
function fetchData(): Promise<any> { ... }

// ✅ REQUIRED — Generic with constraint
function process<T extends Processable>(data: T): Result<T> { ... }
function fetchData<T>(endpoint: string): Promise<T> { ... }
```

### Concrete Type Aliases for Generics

When a generic instantiation appears more than once, **create a named type alias**. This improves readability and provides a single place to update if the underlying generic changes.

```go
// ❌ Repeated generic syntax
func GetUser(ctx context.Context, id int64) apperror.Result[User] { ... }
func GetOrder(ctx context.Context, id int64) apperror.Result[Order] { ... }

// ✅ Named aliases
type UserResult = apperror.Result[User]
type OrderResult = apperror.Result[Order]

func GetUser(ctx context.Context, id int64) UserResult { ... }
func GetOrder(ctx context.Context, id int64) OrderResult { ... }
```

```typescript
// ❌ Verbose generics repeated
function fetchUser(id: string): Promise<ApiResponse<User>> { ... }

// ✅ Named alias
type UserResponse = ApiResponse<User>;
function fetchUser(id: string): Promise<UserResponse> { ... }
```

> **Rule:** If a generic instantiation appears more than once → create a named alias.

### Where `any` / `interface{}` Is Acceptable

1. **SQL query args** — `args ...any` in `dbutil` (framework boundary)
2. **Logger variadic params** — `map[string]any` for structured fields (internal only)
3. **Third-party interfaces** — when a library requires `interface{}`

### Discriminated Unions — Named Interfaces Required

Every variant in a discriminated union **must** be a named interface — inline `{ type: ...; payload: ... }` blocks are prohibited. Enum values must use **PascalCase** and be accessed via **dot notation** only.

| Rule | ❌ Prohibited | ✅ Required |
|------|-------------|------------|
| Union variants | Inline `{ type: ...; }` | Named interface per variant |
| Enum values | `ADD_TOAST`, `REMOVE_TOAST` | `AddToast`, `RemoveToast` |
| Enum access | `ActionType["AddToast"]` | `ActionType.AddToast` |

```typescript
// ❌ PROHIBITED — inline types in union
type ToastAction =
  | { type: "ADD_TOAST"; toast: ToasterToast }
  | { type: "REMOVE_TOAST"; toastId?: string };

// ✅ REQUIRED — named interfaces with PascalCase enum
enum ActionType {
  AddToast = "AddToast",
  UpdateToast = "UpdateToast",
  DismissToast = "DismissToast",
  RemoveToast = "RemoveToast",
}

interface AddToastAction {
  type: ActionType.AddToast;
  toast: ToasterToast;
}

interface UpdateToastAction {
  type: ActionType.UpdateToast;
  toast: Partial<ToasterToast>;
}

interface DismissToastAction {
  type: ActionType.DismissToast;
  toastId?: string;
}

interface RemoveToastAction {
  type: ActionType.RemoveToast;
  toastId?: string;
}

type ToastAction =
  | AddToastAction
  | UpdateToastAction
  | DismissToastAction
  | RemoveToastAction;
```

> **Rule:** If you're writing `| { type: ... }` inline — stop and extract a named interface.

---

## 9. Parallel Execution

- If tasks have **no dependency**, run them in **parallel**.
- Wait for all to complete before proceeding.

```typescript
// ❌ Sequential — wastes time
const users = await fetchUsers();
const orders = await fetchOrders();
const settings = await fetchSettings();

// ✅ Parallel — no dependency between calls
const [users, orders, settings] = await Promise.all([
    fetchUsers(),
    fetchOrders(),
    fetchSettings(),
]);
```

```go
// ✅ Go parallel with goroutines + errgroup
g, ctx := errgroup.WithContext(ctx)

var users []User
var orders []Order

g.Go(func() error { users, err = fetchUsers(ctx); return err })
g.Go(func() error { orders, err = fetchOrders(ctx); return err })

if err := g.Wait(); err != nil {
    return apperror.Wrap(err, "parallel fetch failed")
}
```

---

## 10. Database Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Table names | PascalCase | `UserAccounts`, `OrderItems` |
| Column names | PascalCase | `FirstName`, `CreatedAt` |
| Primary key | `{TableName}Id` | `UserId`, `OrderId` |
| Primary key type | `INTEGER PRIMARY KEY AUTOINCREMENT` | Never UUID unless explicitly requested |
| Foreign key | Same name as source PK | `UserId` in `Orders` table |
| JSON keys | PascalCase | `{ "UserId": 1, "FirstName": "John" }` |
| Index names | PascalCase | `IdxUsersEmail` |
| Enum values | PascalCase | `StatusActive`, `RoleAdmin` |

```sql
-- ✅ Correct table definition
-- linter-waive: MISSING-DESC-001 reason="Cross-language review example; not a real schema"
CREATE TABLE User (
    UserId     INTEGER PRIMARY KEY AUTOINCREMENT,
    FirstName   TEXT NOT NULL,
    Email       TEXT NOT NULL UNIQUE,
    IsActive    INTEGER NOT NULL DEFAULT 1,
    CreatedAt   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Foreign key references the exact PK name
-- linter-waive: MISSING-DESC-001 reason="Cross-language review example; not a real schema"
CREATE TABLE Order (
    OrderId    INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId     INTEGER NOT NULL REFERENCES User(UserId),
    TotalAmount INTEGER NOT NULL
);
```

---

## 11. SQL Safety

- **All queries with user input → parameterized** (no string concatenation).
- **Joins → use database views.** Name views with `Vw` prefix: `VwUserOrders`.
- **All SQL must be tested** — unit tests for queries, integration tests for views.
- **SQL injection validation** on every input boundary.

```go
// ❌ SQL injection risk
query := "SELECT * FROM User WHERE Email = '" + email + "'"

// ✅ Parameterized
query := "SELECT * FROM User WHERE Email = ?"
db.Query(query, email)
```

---

## 12. Logging

- **Add logger calls generously** — at function entry, before/after external calls, on errors.
- Include context: function name, relevant IDs, operation being performed.
- Use structured logging (key-value pairs), not string concatenation.

```go
// ✅ Structured logging with context
logger.Info("processing order", "OrderId", order.OrderId, "UserId", order.UserId)
```

---

## 13. Security (OWASP)

Quick checklist based on OWASP Top 10:

- [ ] **Injection:** All inputs parameterized / sanitized
- [ ] **Broken Auth:** Tokens validated server-side, no secrets in client code
- [ ] **Sensitive Data:** No PII in logs, encryption at rest and in transit
- [ ] **XXE:** XML parsing disabled or restricted
- [ ] **Broken Access Control:** RLS / authorization checked on every endpoint
- [ ] **Misconfiguration:** No default credentials, debug mode off in production
- [ ] **XSS:** All user content sanitized before rendering
- [ ] **Deserialization:** No untrusted deserialization without validation
- [ ] **Vulnerable Dependencies:** Pinned versions, regular audit
- [ ] **Logging:** Security events logged, no sensitive data in logs

---

## 14. Caching

### When to Cache

- **Repeated reads** of the same data within a request or short time window.
- **Expensive computations** (aggregations, transformations) that produce the same result for the same input.
- **External API responses** where the source is rate-limited or slow.

### Code-Red Rules

1. **Never cache errors as success.** If a fetch fails, do not store a stale/empty value — either skip caching or cache a typed error state.
2. **Always set a TTL.** Unbounded caches grow indefinitely and serve stale data. Define explicit expiration.
3. **Invalidate on mutation.** After any write (create/update/delete), invalidate or update the relevant cache entries immediately.
4. **Cache key must be deterministic.** Build keys from stable inputs only — never from timestamps, random values, or user-session state (unless session-scoped cache is intended).
5. **Use typed cache entries.** Never `cache.set(key, any)` — use a typed wrapper so consumers know the shape.

### Patterns

```typescript
// ❌ CODE RED — Caching an error as empty data (silent failure)
try {
    const data = await fetchUsers();
    cache.set("users", data);
} catch {
    cache.set("users", []); // Consumers think there are zero users
}

// ✅ REQUIRED — Skip cache on error, let consumers see the failure
try {
    const data = await fetchUsers();
    cache.set("users", data, { ttl: 300_000 }); // 5 min TTL
} catch (error) {
    cache.delete("users"); // Invalidate stale entry
    logger.error("fetchUsers failed — cache invalidated", { error });

    throw error;
}
```

```go
// ❌ CODE RED — Ignoring error, returning stale cache
result := s.repo.GetById(ctx, id)
if result.HasError() {
    return cachedValue // Silently serving stale data without logging
}

// ✅ REQUIRED — Log staleness, return error to caller
result := s.repo.GetById(ctx, id)
if result.HasError() {
    s.cache.Invalidate(id)
    logger.Warn("cache invalidated due to fetch error", "id", id)

    return apperror.Fail[Plugin](result.AppError())
}

s.cache.Set(id, result.Value(), 5*time.Minute)
```

### React Query as Cache

When using React Query (TanStack Query), the query cache **is** the cache layer:

- Set `staleTime` explicitly — never rely on the default (`0`).
- Set `gcTime` (garbage collection) to control memory.
- Use `queryClient.invalidateQueries()` after mutations.
- Use `placeholderData` for optimistic UI — never cache `undefined` as a real value.

```typescript
// ✅ Explicit cache configuration
useQuery({
    queryKey: ["users", filters],
    queryFn: () => fetchUsers(filters),
    staleTime: 5 * 60 * 1000,  // 5 minutes
    gcTime: 10 * 60 * 1000,    // 10 minutes
});
```

---

## Quick Checklist (Pre-Merge)

- [ ] Functions ≤ 15 lines
- [ ] Files ≤ 300 lines (hard max 400)
- [ ] Max 3 parameters per function
- [ ] Single return value (Result/wrapper)
- [ ] No nested `if` — flat guards only
- [ ] No `!` on function calls — semantic inverse
- [ ] No magic strings/numbers — enums or constants
- [ ] Booleans start with `is`/`has` (99%) or `should` (rare)
- [ ] No mixed `&&`/`||` in one expression
- [ ] 🔴 All errors logged or returned — **never swallowed** (Code Red)
- [ ] Go errors use `apperror` with stack trace
- [ ] Independent async calls use `Promise.all` / goroutines
- [ ] No `any`/`interface{}`/`unknown` in business logic
- [ ] Discriminated unions use named interfaces — no inline variants
- [ ] DB tables/columns/keys in PascalCase
- [ ] Primary keys: `{TableName}Id` + `INTEGER AUTOINCREMENT`
- [ ] SQL parameterized — no string concatenation
- [ ] Joins use views
- [ ] Cache entries have TTL — no unbounded caches
- [ ] Cache invalidated on mutation — no stale reads after writes
- [ ] Never cache errors as success — invalidate or skip
- [ ] OWASP checklist reviewed
- [ ] Spec/issue file written before code
- [ ] Logger calls at key points
