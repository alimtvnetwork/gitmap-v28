# Golang Coding Standards — File organization, naming conventions, negations, guards

> **Parent:** [Golang Coding Standards](./00-overview.md)  
> **Version:** 3.7.0  
> **Updated:** 2026-03-31

---

## File Naming & Organization

### File Naming Rules

| Rule | Convention | Example |
|------|-----------|---------|
| File name | `snake_case.go` | `server_config.go`, `status_type.go` |
| Maps to primary type | File name derived from its exported type | `ServerConfig` → `server_config.go` |
| One exported type per file | Each struct/interface/enum gets its own file | Don't combine `Config` + `ServerConfig` |
| Related methods stay together | All methods on a type live in its file | `StatusType.IsValid()` stays in `status_type.go` |
| Suffix convention | Split large types using suffixes | `_crud.go`, `_helpers.go`, `_validation.go` |
| Package directory | `snake_case` for multi-word | `site_health/`, `search_mode/` |

### Splitting Convention (When Files Exceed 300 Lines)

| Suffix | Purpose | Example |
|--------|---------|---------|
| `{type}.go` | Struct + constructors | `config.go` |
| `{type}_crud.go` | Database CRUD operations | `plugin_crud.go` |
| `{type}_helpers.go` | Private utility functions | `config_helpers.go` |
| `{type}_validation.go` | Input/business rule validation | `upload_validation.go` |
| `{type}_json.go` | JSON marshal/unmarshal methods | `error_json.go` |

### Package Naming Rules

Go packages follow strict naming conventions:

| Rule | Convention | Example |
|------|-----------|---------|
| Lowercase only | No underscores, no mixed case | `sitehealth`, not `site_health` or `SiteHealth` |
| Short and clear | Prefer one word when possible | `auth`, `config`, `sync` |
| No stuttering | Package name must not repeat in exported identifiers | `config.Config` ❌ → `config.Settings` ✅ |
| No generic names | Avoid `util`, `common`, `misc`, `helpers` as package names | Use domain-specific names instead |
| Directory matches package | Directory name = package declaration | `internal/backup/` → `package backup` |

```go
// ❌ FORBIDDEN: Package name stutters in exported type
package config
type ConfigManager struct {} // config.ConfigManager stutters

// ✅ REQUIRED: No stutter
package config
type Manager struct {} // config.Manager reads cleanly

// ❌ FORBIDDEN: Uppercase or underscores in package name
package SiteHealth
package site_health

// ✅ REQUIRED: Lowercase, no separators
package sitehealth
```

### Package Directory Naming

```
// ✅ Correct — enum packages end with 'type' suffix, no underscores
internal/enums/logleveltype/
internal/enums/snapshotmodetype/
internal/services/sitehealth/

// ❌ Wrong
internal/enums/logLevel/
internal/enums/log_level/
internal/enums/SnapshotMode/
internal/services/SiteHealth/
```

### File-to-Type Mapping Examples

```
// ✅ One type per file, name matches
config.go           → type Config struct
server_config.go    → type ServerConfig struct
watcher_config.go   → type WatcherConfig struct
status_type.go      → type StatusType byte + methods
error_json.go       → MarshalJSON/UnmarshalJSON on AppError

// ❌ Wrong: multiple unrelated types in one file
config.go           → Config + ServerConfig + WatcherConfig + BackupConfig
```

---


---

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Package names | Lowercase, single word preferred | `wordpress`, `publish`, `apperror` |
| Package directories | Lowercase, no separators | `sitehealth`, `logleveltype` |
| File names | `snake_case.go`, maps to primary type | `server_config.go`, `status_type.go` |
| Exported functions | PascalCase, verb-led | `EnablePlugin`, `FetchStatus` |
| Unexported functions | camelCase, verb-led | `resolveNamespace`, `parseStackTrace` |
| Constructors | `New` + struct name | `NewConfig`, `NewPluginService` |
| Interfaces | PascalCase, `-er` suffix for single-method | `Publisher`, `PluginStore` |
| Constants | PascalCase | `MaxRetryAttempts`, `DefaultTimeout` |
| Error variables | `Err` prefix | `ErrPluginNotFound`, `ErrUploadFailed` |
| Boolean functions | Positive naming only | `IsValid()`, `HasPermission()` |

### Constructor Naming — `New[StructName]`

Every struct constructor **must** follow the `New` + struct name pattern. The constructor returns a pointer to the struct (or a `Result` for fallible construction).

```go
// ❌ FORBIDDEN: Non-standard constructor names
func CreateConfig(path string) *Config { ... }
func MakeService(db *sql.DB) *PluginService { ... }
func InitManager() *Manager { ... }

// ✅ REQUIRED: New + struct name
func NewConfig(path string) *Config { ... }
func NewPluginService(db *sql.DB) *PluginService { ... }
func NewManager() *Manager { ... }

// ✅ For fallible constructors — return Result
func NewConnection(dsn string) apperror.Result[*Connection] { ... }
```

---


---

## No Raw Negations — Use Positive Guard Functions

> **Canonical source:** [No Raw Negations](../../01-cross-language/12-no-negatives.md)

```go
// ❌ FORBIDDEN
if !fileExists(path) { ... }
if !strings.Contains(s, substr) { ... }

// ✅ REQUIRED
if IsFileMissing(path) { ... }
if IsMissingSubstring(s, substr) { ... }
```

---


---

## `IsDefined` and `IsDefinedAndValid` — Positive Nil/Existence Guards

### Rule: Never use negated nil checks — use `IsDefined()` instead

Raw `!= nil` combined with negation or nested validity checks creates cognitive overhead. Use positive guard methods that express intent clearly.

### `IsDefined()` — Value Existence Check

Returns `true` when the value has been set (is not nil/zero). This replaces `!= nil` checks and negated null patterns.

```go
// ❌ FORBIDDEN: Negated nil check
if config != nil {
    applyConfig(config)
}

// ❌ FORBIDDEN: Double negation
if !isNil(config) {
    applyConfig(config)
}

// ✅ REQUIRED: Positive existence check
if config.IsDefined() {
    applyConfig(config)
}
```

### `IsDefinedAndValid()` — Existence + Validity Combined

Returns `true` when the value exists AND passes its own validation rules. This replaces the common pattern of checking nil then checking validity in a nested or compound condition.

```go
// ❌ FORBIDDEN: Nested nil + validity check
if config != nil {
    if config.IsValid() {
        applyConfig(config)
    }
}

// ❌ FORBIDDEN: Compound with nil check (nested if ban applies too)
if config != nil && config.IsValid() {
    applyConfig(config)
}

// ✅ REQUIRED: Single positive guard
if config.IsDefinedAndValid() {
    applyConfig(config)
}
```

### Implementation Pattern

Every struct that can be nil or absent should implement both methods:

```go
// On pointer receiver types
func (c *Config) IsDefined() bool {
    return c != nil
}

func (c *Config) IsDefinedAndValid() bool {
    return c != nil && c.validate() == nil
}

// On Result[T] wrapper (already built-in)
result.IsDefined()  // true when value was set
result.IsSafe()     // true when value exists AND no error (equivalent to IsDefinedAndValid for results)
```

### Guard Function Table

| Guard Method | Replaces | Description |
|-------------|----------|-------------|
| `IsDefined()` | `!= nil`, `x != nil` | Value exists (not nil/zero) |
| `IsDefinedAndValid()` | `!= nil && IsValid()` | Value exists AND passes validation |
| `IsEmpty()` | `== nil`, `x == nil` | No value set (absent) |
| `IsInvalid()` | `!IsValid()` | Value fails validation |

### On `apperror.Result[T]` (Already Built-in)

| Method | Meaning |
|--------|---------|
| `IsDefined()` | Value was set (regardless of error state) |
| `IsSafe()` | Value exists AND no error — equivalent to `IsDefinedAndValid()` |
| `IsEmpty()` | No value was set |
| `HasError()` | Operation failed |

### Real-World Examples

```go
// Service layer — checking optional input
func (s *SiteService) Update(ctx context.Context, input UpdateSiteInput) apperror.Result[Site] {
    if input.Config.IsDefined() {
        if input.Config.IsDefinedAndValid() {
            applyConfig(input.Config)
        } else {
            return apperror.FailNew[Site](
                "E3010",
                "invalid site config",
            )
        }
    }
    // ...
}

// Handler layer — checking optional query parameter
func (h *Handler) ListPlugins(w http.ResponseWriter, r *http.Request) {
    filter := parseFilter(r)

    if filter.IsDefinedAndValid() {
        plugins = h.plugins.ListFiltered(r.Context(), filter)
    } else {
        plugins = h.plugins.ListAll(r.Context())
    }
}
```

---
