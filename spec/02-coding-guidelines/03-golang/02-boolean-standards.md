# Go Boolean Standards — Positive Logic & Naming

> **Version**: 1.4.0
> **Last updated**: 2026-02-28

## 1. Positive Boolean Naming (Rule P1)

All boolean-returning functions and variables **must** use positive semantic names with `Is` or `Has` prefixes. Boolean fields, local variables, and struct properties follow the same rule — bare adjectives without `is`/`has` are **prohibited**.

```go
// ✅ Positive naming
func IsValid() bool
func HasPermission() bool
func IsActive() bool

// ✅ Variables and fields
isCacheEnabled := true
isForceRefresh := false

type Config struct {
    IsCacheEnabled bool
    IsForceRefresh bool
}

// ❌ Negative naming — PROHIBITED
func IsNotValid() bool
func HasNoPermission() bool
func IsDisabled() bool

// ❌ Missing is/has prefix — PROHIBITED
cacheEnabled := true     // → isCacheEnabled
forceRefresh := false    // → isForceRefresh

type Config struct {
    CacheEnabled bool     // → IsCacheEnabled
    ForceRefresh bool     // → IsForceRefresh
}
```

**Exception**: Enum variant checkers where the variant itself has a negative-sounding name are permitted (e.g., `IsNotFound()` for the `NotFound` variant, `IsUnknown()` for the `Unknown` variant).

## 2. Negation Elimination (Rule P2)

### 2.1 — Named Boolean Variables

Replace inline `!` negation with named positive-logic variables:

```go
// ❌ Inline negation
if !user.IsAdmin() && !request.IsInternal() {
    return ErrForbidden
}

// ✅ Named positive logic
isExternalNonAdmin := user.IsRegular() && request.IsExternal()
if isExternalNonAdmin {
    return ErrForbidden
}
```

### 2.2 — `IsDefined()` and `IsDefinedAndValid()` Guards

Every struct that can be nil or absent **must** implement `IsDefined()` for positive existence checks. If validation logic exists, also implement `IsDefinedAndValid()`:

```go
// ❌ FORBIDDEN — raw nil check
if config != nil && config.IsValid() {
    applyConfig(config)
}

// ✅ REQUIRED — positive combined guard
if config.IsDefinedAndValid() {
    applyConfig(config)
}
```

Implementation pattern:

```go
func (c *Config) IsDefined() bool {
    return c != nil
}

func (c *Config) IsDefinedAndValid() bool {
    return c != nil && c.validate() == nil
}
```

> **Note:** On `apperror.Result[T]`, `IsDefined()` is already built-in. `IsSafe()` serves the same purpose as `IsDefinedAndValid()` (value exists AND no error).

### 2.3 — Positive Counterpart Variables (Rule P3)

When a negated boolean (`!isX`) must be used in any expression, **first** create a positive counterpart variable on a separate line. The negated form must **never** appear directly in a compound condition.

**The priority is always: determine what the positive meaning of the negation is, name it, then use it.**

```go
// ❌ FORBIDDEN — negated boolean used directly in compound condition
isLiveRunWithDeletions := !isDryRun && totalDeleted > 0

// ❌ FORBIDDEN — explicit == false is low priority, still obscure intent
isLiveRunWithDeletions := isDryRun == false && totalDeleted > 0

// ✅ REQUIRED — create the positive counterpart FIRST, then compose
isLiveRun := !isDryRun
hasDeletions := totalDeleted > 0
isLiveRunWithDeletions := isLiveRun && hasDeletions

if isLiveRunWithDeletions {
    commitDeletions()
}
```

**Key principle:** Ask yourself: *"What does `!isDryRun` actually mean?"* — it means the run is live. Name it `isLiveRun`. This makes the compound condition read as plain English: `isLiveRun && hasDeletions`.

More examples:

```go
// ❌ FORBIDDEN — what does !isPending mean? Name it.
if !isPending && hasResults {
    processResults()
}

// ✅ REQUIRED — isPending negated = isCompleted or isProcessed
isProcessed := !isPending
isReadyToProcess := isProcessed && hasResults

if isReadyToProcess {
    processResults()
}
```

### 2.3.1 — Dual Boolean Field Rule (Rule P3b)

When a single boolean source is consumed in both its positive and negative forms within a scope, **both named variables must be declared together** as a pair on consecutive lines before any conditional usage. This is the "dual boolean field" pattern.

**Rationale:** Scattered negations (`!isX`) throughout a function make it impossible to audit the polarity of conditions at a glance. Declaring both forms upfront creates a self-documenting "boolean vocabulary" for the scope.

```go
// ❌ FORBIDDEN — negation scattered, dual form not declared upfront
isProjectExists := pathutil.IsDir(projectDir)

if isProjectExists {
    loadProject(projectDir)
}

// ... 20 lines later ...
if !isProjectExists {   // reader must mentally negate
    createProject(projectDir)
}

// ✅ REQUIRED — dual boolean fields declared together
isProjectExists := pathutil.IsDir(projectDir)
isProjectMissing := !isProjectExists

if isProjectExists {
    loadProject(projectDir)
}

// ... 20 lines later ...
if isProjectMissing {    // instantly clear
    createProject(projectDir)
}
```

**Struct fields follow the same rule.** If a struct exposes a boolean, the opposite meaning must be available via a method — never force callers to negate:

```go
// ❌ FORBIDDEN — callers forced to negate
type Session struct {
    IsAuthenticated bool
}

// caller writes: if !session.IsAuthenticated { ... }  ← negation!

// ✅ REQUIRED — dual accessors
type Session struct {
    IsAuthenticated bool
}

func (s *Session) IsAnonymous() bool {
    return !s.IsAuthenticated
}

// caller writes: if session.IsAnonymous() { ... }  ← positive!
```

**When to create dual fields:**

| Scenario | Required? |
|----------|-----------|
| Both `isX` and `!isX` used in same function | ✅ Yes — declare both upfront |
| Only positive form used | ❌ No — single variable sufficient |
| Only negative form used | ✅ Yes — declare positive first, then negate |
| Struct boolean accessed by multiple callers | ✅ Yes — provide dual accessor methods |

### 2.4 — Positive Counterpart Methods

When a type has an `IsX()` method and code frequently uses `!IsX()`, add a positive counterpart method:

```go
// pathutil package
func IsDirMissing(path string) bool { return !IsDir(path) }

// dbutil.Result[T]
func (r Result[T]) IsEmpty() bool { return !r.defined }  // already exists ✅
```

### 2.5 — Enum Comparisons

Use `IsOther(val)` or `IsInvalid()` instead of `!=` or `!IsValid()`:

```go
// ❌ Negated comparison
if !v.IsValid() {
    return variantLabels[Invalid]
}

// ✅ Positive counterpart
if v.IsInvalid() {
    return variantLabels[Invalid]
}
```

### 2.6 — Named Numeric Comparisons (Rule P5)

Raw numeric comparisons in `if` conditions are **prohibited** when they represent a domain concept. Extract to a named boolean that describes the intent.

```go
// ❌ FORBIDDEN — raw numeric comparison hides intent
if statusCode < 400 {
    return response
}

// ✅ REQUIRED — named boolean explains the domain meaning
isSuccessResponse := statusCode < 400

if isSuccessResponse {
    return response
}
```

```go
// ❌ FORBIDDEN — what does > 0 mean in this context?
if retryCount > 0 && isTransientError {
    retry()
}

// ✅ REQUIRED — name the intent
hasRetriesRemaining := retryCount > 0
isRetryable := hasRetriesRemaining && isTransientError

if isRetryable {
    retry()
}
```

**Exemptions:** Simple loop bounds (`i < len(items)`) and trivial guards (`if count == 0 { return }`) are exempt.

## 2.7 — No Mixed-Polarity Conditions (Rule P6)

Never combine a positive boolean with a negated boolean in the same `if` condition. A negated boolean (`!isX`) may only appear **alone** in a condition — never combined with other terms via `&&` or `||`.

```go
// ❌ FORBIDDEN — mixed polarity in condition
if s.isCacheEnabled && !isForceRefresh {
    return cachedResult
}

// ❌ FORBIDDEN — negation directly in compound expression
if !isDryRun && totalDeleted > 0 {
    commitDeletions()
}

// ✅ REQUIRED — extract negation to positive counterpart, then compose
isNormalRefresh := !isForceRefresh
isCacheHit := s.isCacheEnabled && isNormalRefresh

if isCacheHit {
    return cachedResult
}

// ✅ REQUIRED — negation used ALONE is permitted
if !isValid {
    return ErrInvalid
}
```

**Rule summary:**
- `!isX` alone in a condition → ✅ Permitted
- `!isX && isY` or `isY && !isX` → ❌ Prohibited — extract `!isX` to a named positive counterpart first
- `!isX && !isY` → ❌ Prohibited — two negations is never acceptable

See [Boolean Principles P6](../01-cross-language/02-boolean-principles/03-parameters-and-conditions.md#principle-6-never-mix-positive-and-negative-booleans-in-a-single-condition) for the cross-language rule.

## 2.8 — No Inline Statements in `if` Conditions (Rule P7)

Go allows semicolon-separated inline statements in `if` conditions (e.g., `if x := compute(); x > 0 {`). This pattern is **prohibited** in application code because it:
- Hides variable assignment inside control flow
- Makes the condition harder to read and debug
- Encourages coupling unrelated operations (filesystem check + boolean logic)

The **only exemption** is the idiomatic comma-ok pattern (`if v, ok := m[k]; ok {`) and type assertions (`if v, ok := x.(T); ok {`).

```go
// ❌ FORBIDDEN — inline os.Stat inside if condition
if _, err := os.Stat(projectDir); err == nil && !isOverwrite {
    return fmt.Errorf("project exists, use isOverwrite=true to replace")
}

// ❌ FORBIDDEN — inline statement with pre-computed boolean
isProjectConflict := err == nil && !isOverwrite
if _, err := os.Stat(projectDir); isProjectConflict {
    return fmt.Errorf("project exists, use isOverwrite=true to replace")
}

// ✅ REQUIRED — separate computation from condition
isProjectExists := pathutil.IsDir(projectDir)
isReadOnly := !isOverwrite
isProjectConflict := isProjectExists && isReadOnly

if isProjectConflict {
    return apperror.FailNew[ProjectResult](
        errors.ErrFsConflict,
        "project already exists; use isOverwrite=true to replace",
    )
}
```

### Exemptions

```go
// ✅ Exempt — idiomatic error propagation
if err := validateUpload(req); err != nil {
    return err
}

// ✅ Exempt — idiomatic comma-ok
if v, ok := myMap[key]; ok {
    process(v)
}

// ✅ Exempt — idiomatic type assertion
if concrete, ok := iface.(MyType); ok {
    concrete.DoWork()
}

// ✅ Exempt — recover in deferred panic handler
defer func() {
    if r := recover(); r != nil {
        handlePanic(r)
    }

}()
```

> **Note:** The `if err := fn(); err != nil` pattern is exempt because it is idiomatic Go error propagation where the variable is only used within the error check. The prohibition targets non-error value extraction patterns like `if cached := getFromCache(url); cached != nil`.

## 2.9 — No Raw Filesystem Calls in Application Code (Rule P8)

Application code **must not** call raw `os.Stat`, `os.MkdirAll`, `os.Remove`, `os.ReadFile`, or any `os` package filesystem function directly. Instead, use `pathutil` wrapper functions that:

1. Return `*apperror.AppError` with proper error codes (not raw `error`)
2. Provide positive-named boolean helpers (`IsDir`, `IsDirMissing`, `IsFile`, `IsFileMissing`)
3. Handle edge cases (permissions, symlinks) consistently

```go
// ❌ FORBIDDEN — raw os.Stat in application code
if _, err := os.Stat(projectDir); err == nil {
    // exists...
}

// ❌ FORBIDDEN — raw os.MkdirAll
if err := os.MkdirAll(outputDir, 0755); err != nil {
    return fmt.Errorf("failed to create dir: %w", err)
}

// ✅ REQUIRED — use pathutil wrappers
isProjectExists := pathutil.IsDir(projectDir)

// ✅ REQUIRED — use pathutil with apperror
if err := pathutil.EnsureDir(outputDir); err != nil {
    return apperror.Fail[OutputResult](err)
}
```

### Required `pathutil` Inventory

| Function | Returns | Description |
|----------|---------|-------------|
| `IsDir(path)` | `bool` | True if path exists and is a directory |
| `IsDirMissing(path)` | `bool` | True if path does not exist as directory |
| `IsFile(path)` | `bool` | True if path exists and is a regular file |
| `IsFileMissing(path)` | `bool` | True if path does not exist as file |
| `EnsureDir(path)` | `*apperror.AppError` | Creates directory if missing; returns structured error |
| `Remove(path)` | `*apperror.AppError` | Removes file/dir; returns structured error |
| `Stat(path)` | `(os.FileInfo, *apperror.AppError)` | Wraps `os.Stat` with `apperror` |

### Comprehensive Example — All Rules Combined

This example demonstrates violations of **P6** (mixed polarity), **P7** (inline statement), and **P8** (raw `os.Stat`) — and the correct fix:

```go
// ❌ FORBIDDEN — 4 violations in 3 lines:
//   1. P8: raw os.Stat (must use pathutil)
//   2. P7: inline statement in if (semicolon)
//   3. P6: mixed polarity (err == nil && !isOverwrite)
//   4. Raw fmt.Errorf instead of apperror
isProjectConflict := err == nil && !isOverwrite
if _, err := os.Stat(projectDir); isProjectConflict {
    return fmt.Errorf("project exists, use isOverwrite=true to replace")
}

// ✅ REQUIRED — clean, readable, all rules applied:
//   1. P8: pathutil.IsDir() wraps os.Stat
//   2. P7: no inline statement; all variables computed before if
//   3. P6: mixed polarity extracted to single-intent boolean
//   4. apperror.FailNew returns structured *apperror.AppError
isProjectExists := pathutil.IsDir(projectDir)
isReadOnly := !isOverwrite
isProjectConflict := isProjectExists && isReadOnly

if isProjectConflict {
    return apperror.FailNew[ProjectResult](
        errors.ErrFsConflict,
        "project already exists; use isOverwrite=true to replace",
    )
}
```

## 2.10 — No Compound Error Conditions (Rule P9)

**Combining `err != nil` with any other condition using `&&` or `||` is PROHIBITED.** Error checking must be a single, isolated guard. If you need to check an error alongside another condition, use `*apperror.AppError` methods which handle nil-safety through pointer receivers.

### Prohibited: Double or Compound `err != nil`

```go
// ❌ FORBIDDEN — two err != nil checks combined
if err != nil && otherErr != nil {
    return fmt.Errorf("both failed")
}

// ❌ FORBIDDEN — err != nil combined with domain condition
if err != nil && !os.IsNotExist(err) {
    return err
}

// ❌ FORBIDDEN — err != nil mixed with positive condition
if err != nil && mood.IsDefined() {
    return handleMoodError(err, mood)
}

// ❌ FORBIDDEN — error == nil combined with other condition (mixed polarity)
if err == nil && isReady {
    proceed()
}
```

### Required: Use `*apperror.AppError` Methods

`*apperror.AppError` pointer receiver methods are **nil-safe** — calling `HasError()`, `IsDefined()`, or any method on a nil `*AppError` returns the appropriate zero value without panicking. This eliminates the need for raw `err != nil` checks entirely in application code.

```go
// ✅ REQUIRED — single error guard, no compound condition
if appErr.HasError() {
    return appErr
}

// ✅ REQUIRED — check errors separately, never combine
if firstErr.HasError() {
    return firstErr
}

if secondErr.HasError() {
    return secondErr
}

// ✅ REQUIRED — combine errors using AppError methods, not && operators
combinedErr := apperror.Combine(firstErr, secondErr)

if combinedErr.HasError() {
    return combinedErr
}
```

### Error + Domain Condition — Separate Guards

When you need to check both an error and a domain condition, they **must** be separate `if` blocks:

```go
// ❌ FORBIDDEN — compound error + domain check
if appErr.HasError() || mood.IsUndefined() {
    return apperror.FailNew[MoodResult](
        errors.ErrInvalidMood,
        "mood check failed",
    )
}

// ✅ REQUIRED — separate guards, each with clear intent
if appErr.HasError() {
    return apperror.Fail[MoodResult](appErr)
}

if mood.IsUndefined() {
    return apperror.FailNew[MoodResult](
        errors.ErrInvalidMood,
        "mood is undefined",
    )
}
```

### Single `err != nil` Exemption

A **single** `err != nil` check (not combined with anything) remains exempt for raw error propagation at library boundaries:

```go
// ✅ Exempt — single err != nil for raw error wrapping at boundary
if err != nil {
    return apperror.WrapNew[T](err, errors.ErrCodeHere, "context message")
}
```

> **Note:** In application code, prefer `appError.HasError()` even for single checks. The `err != nil` exemption exists only for library boundary wrapping where you receive a raw `error` from external packages.

## 3. Idiomatic Go Exemptions

The following patterns are **exempt** from negation elimination:

### 3.1 — Comma-ok Pattern

The comma-ok return value **must** be renamed to a semantically meaningful positive boolean. The bare `ok` variable name is **prohibited** — always name it to describe what "ok" means in context (e.g., `isExists`, `isFound`, `isLoaded`).

If the negative case is needed, create a positive counterpart on the next line:

```go
// ❌ PROHIBITED — bare `ok` hides meaning
value, ok := someMap[key]
if !ok {
    return ErrNotFound
}

// ✅ REQUIRED — semantic name describes the positive case
value, isExists := someMap[key]
isMissing := !isExists

if isMissing {
    return ErrNotFound
}

// ✅ Also acceptable — positive guard when you only need the positive path
value, isExists := someMap[key]

if isExists {
    process(value)
}
```

More examples:

```go
// ❌ PROHIBITED
conn, ok := connections[id]
if !ok { ... }

// ✅ REQUIRED
conn, isFound := connections[id]
isNotFound := !isFound

if isNotFound {
    return apperror.FailNew[ConnResult](
        errors.ErrConnNotFound,
        "connection missing",
    )
}

// ❌ PROHIBITED
cached, ok := cache.Load(key)
if !ok { ... }

// ✅ REQUIRED
cached, isLoaded := cache.Load(key)
isCacheMiss := !isLoaded

if isCacheMiss {
    cached = fetchFromSource(key)
}
```

> **Note:** The inline comma-ok in `if` conditions (`if v, isExists := m[k]; isExists {`) remains exempt from Rule P7 but **must** still use a semantic name instead of `ok`.

### 3.2 — Handler Guard Returns

Early-return guards in HTTP handlers that return false on failure:

```go
// ✅ Exempt — handler guard pattern
if !requireService(w, Services.SyncService, "Sync service") {
    return
}

if !decodeJSON(w, r, &input) {
    return
}
```

### 3.3 — Error-nil Check

A **single** `err != nil` is exempt for idiomatic Go error propagation. However, compound error conditions are **prohibited** (see Rule P9 in §2.10).

```go
// ✅ Exempt — single err != nil for propagation
if err != nil {
    return err
}

// ✅ PREFERRED in application code — use appError methods
if appErr.HasError() {
    return appErr
}

// ❌ PROHIBITED — compound error condition (Rule P9)
if err != nil && !os.IsNotExist(err) {
    return err
}

// ❌ PROHIBITED — double error check (Rule P9)
if err1 != nil && err2 != nil {
    return errors.Join(err1, err2)
}
```

### 3.4 — Standard Library Returns

Direct `!` on stdlib function returns where no wrapper exists:

```go
// ✅ Exempt — stdlib call
if !strings.HasPrefix(path, "/api/") {
    return
}
```

However, if the same stdlib negation appears 3+ times, extract a named boolean or helper:

```go
// When repeated, extract:
isNonApiRoute := !strings.HasPrefix(r.URL.Path, "/api/")

if isNonApiRoute {
    next.ServeHTTP(w, r)

    return
}
```

## 4. Variable Naming Rules

| Pattern | Example | Status |
|---------|---------|--------|
| `is` + PositiveAdjective | `isValid`, `isActive`, `isReady` | ✅ Required |
| `has` + PositiveNoun | `hasPermission`, `hasRows`, `hasError` | ✅ Required |
| `is` + NegativeResult | `isDirMissing`, `isMkdirFailed` | ✅ Permitted |
| `isDefined` | Positive nil/existence check | ✅ Required on nullable structs |
| `isDefinedAndValid` | Existence + validation combined | ✅ Required when validation exists |
| `not` prefix | `notFound`, `notReady` | ❌ Prohibited |
| `no` prefix | `noResults`, `noPermission` | ❌ Prohibited |
| Bare `ok` | `value, ok := map[key]` | ❌ Prohibited — use semantic name |

## 5. Enforcement

- **Automated**: `linter-scripts/lint-negative.sh` flags `IsNot*`, `HasNo*` function declarations
- **Manual review**: Inline `!` negation in compound boolean expressions
- **Enum exemption**: Variant checkers matching their constant name (e.g., `IsNotFound` for `NotFound` variant) are auto-excluded
- **Compound error check**: Any `err != nil &&` or `err != nil ||` pattern is flagged

## 6. Rule Summary

| Rule | ID | Summary |
|------|----|---------|
| Positive Naming | P1 | All booleans use `is`/`has` positive prefixes |
| Negation Elimination | P2 | Replace `!` with named positive variables |
| Positive Counterpart Variables | P3 | Negated booleans must be assigned to a positive-named variable before use in compounds |
| Dual Boolean Fields | P3b | Both positive and negative named forms declared together upfront; structs provide dual accessor methods |
| Named Numeric Comparisons | P5 | Raw numeric comparisons → named booleans |
| No Mixed Polarity | P6 | `!isX` only alone; never combined with `&&`/`||` |
| No Inline Statements | P7 | No semicolon assignments in `if`; exemptions for comma-ok, type assertions, error propagation |
| No Raw Filesystem | P8 | Use `pathutil` wrappers, not raw `os` calls |
| No Compound Errors | P9 | `err != nil` never combined with other conditions; use `appError.HasError()` |
| Semantic Comma-ok | — | Rename `ok` to meaningful name (`isExists`, `isFound`, etc.) |

## 7. Cross-Language Alignment

This standard mirrors the cross-language [Boolean Principles](../01-cross-language/02-boolean-principles/00-overview.md) (P1–P6) and [No-Negatives](../01-cross-language/12-no-negatives.md) with Go-specific exemptions for idiomatic patterns (comma-ok, handler guards, error-nil checks) and Go-specific additions (P3b, P5, P7–P9). See [PHP Standards](../04-php/03-naming-conventions.md) for the PHP counterpart.
