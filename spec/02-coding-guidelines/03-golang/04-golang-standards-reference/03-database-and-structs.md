# Golang Coding Standards — Database naming, dbutil wrapper, struct design

> **Parent:** [Golang Coding Standards](./00-overview.md)  
> **Version:** 3.7.0  
> **Updated:** 2026-03-31

---

## Database Naming Convention — PascalCase

> **Canonical source:** [Database Naming Convention](../../01-cross-language/07-database-naming.md)

All custom SQLite table names, column names, and index names MUST use **PascalCase**. Go struct `db` tags must match column names. JSON tags follow the redundancy rule below — omit when the field name already matches.

```go
// ✅ PascalCase table and column names — no redundant json tags
const queryList = `SELECT Id, ProjectId, DisplayName, CreatedAt FROM Project`

type Project struct {
    Id          int64  `db:"Id"`
    ProjectId   string `db:"ProjectId"`
    DisplayName string `db:"DisplayName"`
    CreatedAt   string `db:"CreatedAt"`
}
```

---


---

## Database Wrapper — `pkg/dbutil`

All database queries MUST use the generic `dbutil` package. Returns typed result envelopes with automatic `apperror` stack traces.

### Result Types

| Type | Purpose | Key Methods |
|------|---------|-------------|
| `Result[T]` | Single-row query | `IsDefined()`, `IsEmpty()`, `HasError()`, `IsSafe()`, `Value()`, `AppError()`, `StackTrace()` |
| `ResultSet[T]` | Multi-row query | `HasAny()`, `IsEmpty()`, `Count()`, `HasError()`, `IsSafe()`, `Items()`, `First()`, `AppError()`, `StackTrace()` |
| `ExecResult` | INSERT/UPDATE/DELETE | `IsEmpty()`, `HasError()`, `IsSafe()`, `AffectedRows`, `LastInsertId`, `AppError()`, `StackTrace()` |

> **Naming: `.AppError()` vs `.Error()`**
> Both `dbutil` result types and `apperror` result types use `.AppError()` (not `.Error()`) to return the underlying `*apperror.AppError`. This avoids collision with Go's built-in `error` interface method `.Error() string`. The `apperror.AppError` struct itself still implements the standard `error` interface via `.Error() string` (returns `"[code] message"`), but all result wrappers — whether in `dbutil` (`Result[T]`, `ResultSet[T]`, `ExecResult`) or `apperror` (`Result[T]`, `ResultSlice[T]`, `ResultMap[K,V]`) — expose `.AppError()` for the structured error accessor.

### Generic Query Functions

```go
// Single row — returns Result[T]
result := dbutil.QueryOne[Plugin](ctx, db, query, scanPlugin, pluginId)

// Multiple rows — returns ResultSet[T]
set := dbutil.QueryMany[Site](ctx, db, query, scanSite)

// Exec — returns ExecResult
res := dbutil.Exec(ctx, db, query, args...)
```

---


---

## Struct Design

### JSON Tags — Omit Redundant Tags

Go's `encoding/json` marshaler uses the **field name** by default. Since all our fields are PascalCase and JSON output is PascalCase, explicit `json:"FieldName"` tags are **redundant** and MUST be omitted. Only add a `json` tag when using `omitempty`:

```go
// ❌ WRONG — redundant json tags that repeat the field name
type PluginDetails struct {
    Id        int    `json:"Id"`
    Name      string `json:"Name"`
    Slug      string `json:"Slug"`
    Version   string `json:"Version"`
    IsActive  bool   `json:"IsActive"`
    UpdatedAt string `json:"UpdatedAt,omitempty"`
}

// ✅ CORRECT — no tags unless omitempty is needed
type PluginDetails struct {
    Id        int
    Name      string
    Slug      string
    Version   string
    IsActive  bool
    UpdatedAt string `json:",omitempty"`
}
```

**Rules:**
- **No tag:** Field marshals as its Go name (PascalCase) — this is the default
- **`json:",omitempty"`:** Only when zero-value fields should be excluded from output
- **`json:"-"`:** Only when a field must be excluded from JSON entirely
- **`db:"ColumnName"`:** Always required for database-mapped structs (db tags don't auto-derive)

### Function Parameters — Max 2-3

Functions should have **2-3 parameters maximum**. Use config/options structs for more:

```go
// ❌ Bad: Too many parameters
func StartSession(sessionType SessionType, pluginId, siteId int64, pluginName, siteName string) (string, error)

// ✅ Good: Use a struct
type StartSessionInput struct {
    Type       SessionType
    PluginId   int64
    SiteId     int64
    PluginName string
    SiteName   string
}
func StartSession(input StartSessionInput) apperror.Result[string]

// ✅ Acceptable: 2-3 essential parameters (context doesn't count)
func GetById(ctx context.Context, id int64) apperror.Result[Model]
```

---
