# Consolidated Review Guide — Condensed

> One-liner rules. No fluff. Scan before every PR.

---

## Size Limits

- Function body: **≤ 15 lines** (target 8–10)
- File: **≤ 300 lines** (hard max 400)
- Type/struct: **≤ 120 lines**
- Parameters: **≤ 3** per function (use options object if more)
- Returns: **1 value only** (use Result/wrapper)

## Naming

- Classes/types: `PascalCase` — `OrderService`
- Functions: `camelCase` (TS/PHP), `PascalCase` (Go exported)
- Variables: `camelCase` — `orderTotal`
- Booleans (vars **and** functions): **must** start with `is` or `has` (99% of cases) — `should` only for recommendations/preferences
- **Never** use `not`/`no` in boolean names — use inverse: `isInactive`, `isMissing`, `isDisabled`
- Bool functions: `checkPermission()` → `hasPermission()`, `validateToken()` → `isTokenValid()`
- Abbreviations: `Id`, `Url`, `Json`, `Api`, `Http`, `Sql` (first letter only)
- Enums: PascalCase + `Type` suffix — `StatusType`

## Conditionals

- **Zero nested `if`** — invert condition, exit early, continue on happy path
- **No `!` on function calls** — use `isInvalid()`, `isMissing()`, `isInactive()`
- **Max 2 conditions** per `if`
- **Never mix `&&` and `||`** — extract into named booleans
- **Never mix positive + negative** (`isX && !isY` → extract)
- **No `else` after `return`/`throw`/`break`/`continue`**
- **Keep cyclomatic complexity low** — flat guard clauses, one exit path

## Values

- No magic strings → enum
- No magic numbers → named constant
- Exceptions: `0`, `1`, `-1`, `""`, `true`, `false`, `null`/`nil`

## Errors — CODE RED 🔴

> **ZERO TOLERANCE: Never swallow an error.** Empty `catch {}`, bare `return nil`, ignored `Result` = **automatic rejection.**

- 🔴 **Never swallow an error** — log or return, always. Silent failures are the #1 production incident cause
- Go: `apperror.Wrap()` at **first contact** with stdlib/third-party errors — **stack trace mandatory**
- Go: no `fmt.Errorf()`, no `errors.New()` — use `apperror.New("E1001", "msg")`
- Go: after wrapping, propagate via `.AppError()` — never re-wrap
- Go: `err.FullString()` for logs, `err.ToClipboard()` for AI paste
- PHP: catch `Throwable`, not `Exception`
- Check `HasError()` / `hasError()` before `.Value()`
- Structured responses: message + stack trace + frames
- **Never cache errors as success** — invalidate cache on failure, don't store empty/stale values

## Types — CODE RED

- No `any` / `interface{}` / `unknown` / `object` in business logic
- **Always prefer generics** over `any`/`interface{}`/`unknown` — create a concrete generic type first, then reuse it
- **Use concrete types** for everything else — no type assertions in Go business logic
- **No inline return types** — extract named types/interfaces so they can be reused
  - ❌ `func(): { text: string; store: Record<string, string> }`
  - ✅ `type ChecklistResult = { text: string; store: Record<string, string> }` → `func(): ChecklistResult`
- `unknown` only at parse/deserialization boundaries with type guard
- **No inline types in discriminated unions** — extract named interface per variant, use generic `TypedAction<T, P>` for shared shapes
- **Enum values: PascalCase only** — `AddToast`, never `ADD_TOAST`; dot notation only (`ActionType.AddToast`), never index access (`ActionType["AddToast"]`)
- `interface{}` / `any` only: SQL args, logger fields, third-party interfaces

### Go Generic Pattern — `Result[T]`

Build a concrete generic wrapper once, reuse everywhere:

```go
// ✅ Define the generic type once
type Result[T any] struct {
    value    T
    appError *AppError
}

func Ok[T any](value T) Result[T] {
    return Result[T]{value: value}
}

func Fail[T any](err *AppError) Result[T] {
    return Result[T]{appError: err}
}

func (r Result[T]) HasError() bool { return r.appError != nil }
func (r Result[T]) Value() T      { return r.value }
```

#### ⚠️ Acceptable (not best practice) — inline generic usage

```go
// Works but repeats Result[User] everywhere — no concrete alias
func GetUser(id int64) Result[User]     { ... }
func GetPlugin(id int64) Result[Plugin] { ... }
func ListOrders() Result[[]Order]       { ... }
```

#### ✅ Best Practice — concrete type aliases for reusability

Create a named type alias from `Result[T]` for each domain type. This gives you a single concrete type to reference, import, and document — no generic syntax at call sites:

```go
// ✅ Define concrete result types once per domain model
type UserResult = Result[User]
type PluginResult = Result[Plugin]
type OrderListResult = Result[[]Order]

// ✅ Functions return the concrete alias — clean, reusable, searchable
func GetUser(id int64) UserResult     { ... }
func GetPlugin(id int64) PluginResult { ... }
func ListOrders() OrderListResult     { ... }

// ✅ Callers use the concrete type — no generic noise
func handleUser() {
    result := GetUser(42)

    if result.HasError() { ... }
    user := result.Value()
}
```

**Rule:** When tempted to use `any` or `interface{}`, ask: _can I make this generic?_ Almost always yes. Then create a concrete type alias so every call site is clean.

## Async / Parallel

- Independent calls → `Promise.all()` (TS) / `errgroup` (Go)
- Never sequential when parallel is possible

## Database

- Tables/columns: **PascalCase** — `Users`, `FirstName`
- Primary key: `{TableName}Id` + `INTEGER PRIMARY KEY AUTOINCREMENT`
- Foreign key: same name as source PK
- JSON keys: **PascalCase**
- No UUID unless explicitly requested
- Joins → database views (`VwUserOrders`)

## SQL Safety

- All user input → parameterized queries
- All SQL under test
- SQL injection validation on every input

## Logging

- Log at: function entry, external calls, errors
- Structured key-value pairs, not string concat

## Security (OWASP Top 10)

- Inputs sanitized, queries parameterized, tokens server-validated
- No PII in logs, no secrets in client code
- Dependencies pinned + audited
- XSS: sanitize before render

## Newline Conventions

- **Blank line before `return`/`throw`** — only when preceded by a statement
- **No blank line** when `return` is the only line inside `{ }` (single-line body)
- **No blank line at function/class start** — first line of code goes immediately after `{`
- **Blank line after `}`** when followed by more code
- **No double blank lines** — ever
- **No trailing blank line** at end of function/class

```ts
// ✅ Single statement — no blank line needed
function getId(): number {
  return this.id;
}

// ✅ Multi-line — blank line before return
function process(): number {
  doSomething();

  return result;
}

// ❌ WRONG — blank line at function start
function bad(): void {

  doSomething();
}

// ❌ WRONG — no blank line before return in multi-line
function bad(): number {
  doSomething();
  return result;
}
```

## File Naming

- **Markdown files**: lowercase with numeric prefix — `00-overview.md`, `01-setup.md`, `readme.md`
- **All other files/values**: PascalCase (types, enums, JSON keys, DB columns)

## Caching — CODE RED 🔴

- **Never cache errors as success** — empty `catch { cache.set(key, []) }` is a silent failure (Code Red)
- **Always set TTL** — unbounded caches serve stale data and grow without limit
- **Invalidate on mutation** — after create/update/delete, invalidate or update cache entries immediately
- **Cache keys must be deterministic** — stable inputs only, no timestamps or random values
- **Use typed cache entries** — never `cache.set(key, any)`
- React Query: set `staleTime` explicitly, use `invalidateQueries()` after mutations

## Workflow

- **Spec first → code second**
- **Bug → issue file first** (`spec/.../03-issues/{NN}-{name}.md`) → root cause → fix
- Follow folder naming: `00-overview.md` + `99-consistency-report.md`
