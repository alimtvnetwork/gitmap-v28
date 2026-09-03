# Consolidated Enum Standards — All Languages

**Version:** 3.2.0
**Updated:** 2026-04-16
**Source Specs:**

- [Go Enum Specification](../02-coding-guidelines/03-golang/01-enum-specification/01-index.md)
- [TypeScript Enums](../02-coding-guidelines/02-typescript/01-index.md)
- [PHP Enums](../02-coding-guidelines/04-php/02-enums.md)
- [Rust Naming Conventions](../02-coding-guidelines/05-rust/02-naming-conventions.md)
- [Enum Naming Quick Reference](../02-coding-guidelines/06-ai-optimization/07-enum-naming-quick-reference.md)

---

## Universal Rules (All Languages)

| Rule | Requirement |
|------|-------------|
| Type Suffix | Enum types must end with `Type` in camelCase languages (e.g. `UserType`), or `_type` in snake_case languages |
| No magic strings | Never compare against raw string literals — always use enum constants |
| PascalCase values | Enum members/cases use PascalCase (`Production`, not `PRODUCTION` or `production`) |
| One definition | String representations defined **once**, co-located with the enum type |
| Exhaustive switch | Every `switch`/`match` on an enum must have a `default` branch |
| No string unions | Use proper `enum` syntax — never `type Foo = 'a' | 'b'` |

---

## Go Enums

### Declaration

```go
package environmenttype          // package name = grouping (no type prefix on constants)

type Variant byte                // always byte, never string or int

const (
    Invalid     Variant = iota   // always first (zero value) — never "Unknown"
    Production
    Staging
    Development
)
```

### Naming Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Package name | Lowercase, `type` suffix | `environmenttype`, `providertype` |
| Type name | Always `Variant` | `type Variant byte` |
| Constants | PascalCase, **no type prefix** | `Production` (not `EnvironmentProduction`) |
| Zero value | Always `Invalid` | Never `Unknown` or `None` |
| Lookup table | `variantLabels` (unexported) | Single array, PascalCase values |

### Required Methods (every enum)

| Method | Signature | Purpose |
|--------|-----------|---------|
| `String()` | `(v Variant) String() string` | Serialization/logging |
| `Label()` | `(v Variant) Label() string` | Human-readable label |
| `Is{Value}()` | `(v Variant) IsSerpApi() bool` | Type-safe comparison |
| `IsValid()` | `(v Variant) IsValid() bool` | Bounds check |
| `MarshalJSON()` | `(v Variant) MarshalJSON() ([]byte, error)` | JSON output |
| `UnmarshalJSON()` | `(v *Variant) UnmarshalJSON(b []byte) error` | JSON input |
| `Parse()` | `func Parse(s string) (Variant, error)` | Case-insensitive string→enum |
| `All()` | `func All() []Variant` | Returns all valid variants |
| `ByIndex()` | `func ByIndex(i int) Variant` | Get variant by index |
| `IsOther()` | `(v Variant) IsOther(other Variant) bool` | True if NOT the given variant |
| `IsAnyOf()` | `(v Variant) IsAnyOf(others ...Variant) bool` | True if receiver matches any |

### Folder Structure

```
internal/enums/
├── environmenttype/
│   └── variant.go
├── providertype/
│   └── variant.go
└── platformtype/
    └── variant.go
```

### Full Example — `providertype`

```go
package providertype

import (
    "encoding/json"
    "fmt"
    "strings"
)

type Variant byte

const (
    Invalid     Variant = iota
    SerpApi
    MapsScraper
    Colly
)

var variantLabels = [...]string{
    Invalid:     "Invalid",
    SerpApi:     "SerpApi",
    MapsScraper: "MapsScraper",
    Colly:       "Colly",
}

func (v Variant) String() string {
    if int(v) < len(variantLabels) {
        return variantLabels[v]
    }

    return "Invalid"
}

func (v Variant) Label() string  { return v.String() }
func (v Variant) IsValid() bool  { return v > Invalid && int(v) < len(variantLabels) }
func (v Variant) IsSerpApi() bool     { return v == SerpApi }
func (v Variant) IsMapsScraper() bool { return v == MapsScraper }
func (v Variant) IsColly() bool       { return v == Colly }

func (v Variant) IsOther(other Variant) bool { return v != other }

func (v Variant) IsAnyOf(others ...Variant) bool {
    for _, o := range others {
        if v == o {
            return true
        }
    }

    return false
}

func All() []Variant {
    return []Variant{SerpApi, MapsScraper, Colly}
}

func ByIndex(i int) Variant {
    if i > 0 && i < len(variantLabels) {
        return Variant(i)
    }

    return Invalid
}

func Parse(s string) (Variant, error) {
    lower := strings.ToLower(s)
    for i, label := range variantLabels {
        if i == 0 {
            continue
        }
        if strings.ToLower(label) == lower {
            return Variant(i), nil
        }
    }

    return Invalid, fmt.Errorf("unknown provider: %s", s)
}

func (v Variant) MarshalJSON() ([]byte, error) {
    return json.Marshal(v.String())
}

func (v *Variant) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    parsed, err := Parse(s)
    if err != nil {
        return err
    }
    *v = parsed

    return nil
}
```

### Forbidden Patterns (Go)

```go
type Provider string                    // ❌ string-based enum
if provider == "serpapi" { ... }        // ❌ magic string comparison
EnvironmentProductionStr = "production" // ❌ type-prefixed constant name
const ( Unknown Variant = iota )       // ❌ "Unknown" as zero value
```

---

## TypeScript Enums

### Declaration

```typescript
// src/lib/enums/connection-status.ts

export enum ConnectionStatus {
  Connected = "CONNECTED",
  Disconnected = "DISCONNECTED",
  Connecting = "CONNECTING",
  Reconnecting = "RECONNECTING",
  Error = "ERROR",
}
```

### Naming Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Enum name | PascalCase | `ConnectionStatus`, `LogLevel` |
| File name | kebab-case | `connection-status.ts` |
| Members | PascalCase | `Connected`, `Disconnected` |
| Values | UPPER_SNAKE string | `"CONNECTED"`, `"DISCONNECTED"` |
| Folder | `src/lib/enums/` | One file per enum |

### Usage Patterns

```typescript
// ✅ CORRECT — enum constant
if (ws.status === ConnectionStatus.Connected) { ... }

// ❌ FORBIDDEN — magic string
if (ws.status === 'connected') { ... }

// ✅ CORRECT — typed interface
interface WsState { status: ConnectionStatus; }

// ❌ FORBIDDEN — string union
interface WsState { status: 'connected' | 'disconnected'; }
```

### Conditional Rendering (React)

```typescript
// ✅ CORRECT
{wsStatus === ConnectionStatus.Disconnected && <ReconnectBanner />}

// ❌ FORBIDDEN
{wsStatus === 'disconnected' && <ReconnectBanner />}
```

### Defined Enums

| Enum | Values |
|------|--------|
| `ConnectionStatus` | Connected, Disconnected, Connecting, Reconnecting, Error |
| `EntityStatus` | Active, Inactive, Pending, Archived |
| `ExecutionStatus` | Pending, Running, Completed, Failed, Cancelled |
| `ExportStatus` | Pending, Processing, Completed, Failed |
| `HttpMethod` | Get, Post, Put, Patch, Delete |
| `MessageStatus` | Pending, Streaming, Completed, Error |
| `LogLevel` | Debug, Info, Warn, Error, Fatal |

---

## PHP Enums

### Declaration

```php
// includes/Enums/HttpMethodType.php

namespace RiseupAsia\Enums;

enum HttpMethodType: string
{
    case Get    = 'GET';
    case Post   = 'POST';
    case Put    = 'PUT';
    case Patch  = 'PATCH';
    case Delete = 'DELETE';

    public function isEqual(self $other): bool
    {
        return $this === $other;
    }
}
```

### Naming Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Enum name | PascalCase + `Type` suffix | `HttpMethodType`, `HookType` |
| File name | Matches enum name | `HttpMethodType.php` |
| Cases | PascalCase | `case RestApi`, not `case REST_API` |
| Namespace | `RiseupAsia\Enums` | All enums in same namespace |
| Folder | `includes/Enums/` | One file per enum |

### Required Methods

| Method | Purpose |
|--------|---------|
| `isEqual(self $other): bool` | Type-safe comparison (mandatory on every backed enum) |
| `validValues(): array` | Static — returns all valid string values |
| `isValid(string $input): bool` | Static — check if a raw string maps to a case |

### Comparison — `isEqual()` Only

```php
// ❌ FORBIDDEN — raw === comparison
if ($status === StatusType::Success) { ... }

// ✅ REQUIRED — isEqual() method
if ($status->isEqual(StatusType::Success)) { ... }
```

### Internal Helpers Delegate to `isEqual()`

```php
enum StatusType: string
{
    case Success = 'success';
    case Failed  = 'failed';

    public function isEqual(self $other): bool
    {
        return $this === $other;
    }

    // ❌ FORBIDDEN — direct === in helpers
    public function isSuccess(): bool
    {
        return $this === self::Success;
    }

    // ✅ REQUIRED — delegate to isEqual()
    public function isSuccess(): bool
    {
        return $this->isEqual(self::Success);
    }
}
```

### Compound Checks

```php
// ✅ Multiple isEqual() calls
public function isLifecycle(): bool
{
    return $this->isEqual(self::Enable)
        || $this->isEqual(self::Disable)
        || $this->isEqual(self::Delete);
}
```

### Parsing — No Manual Switch

```php
// ✅ REQUIRED — use built-in backed enum parsing
$method = HttpMethodType::from($input);       // throws ValueError
$method = HttpMethodType::tryFrom($input);    // returns null

// ❌ FORBIDDEN — manual match with raw strings
match ($input) {
    'GET' => HttpMethodType::Get,
    'POST' => HttpMethodType::Post,
}
```

### Full Example — `UploadSourceType`

```php
namespace RiseupAsia\Enums;

enum UploadSourceType: string
{
    case Script  = 'upload_script';
    case RestApi = 'rest_api';
    case AdminUi = 'admin_ui';
    case WpCli   = 'wp_cli';

    public function isEqual(self $other): bool
    {
        return $this === $other;
    }

    public static function validValues(): array
    {
        return array_column(self::cases(), 'value');
    }

    public static function isValid(string $source): bool
    {
        return self::tryFrom($source) !== null;
    }
}
```

---

## Rust Enums

### Declaration

```rust
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ScreenshotTrigger {
    Periodic,
    TabChange,
    AppSwitch,
    Idle,
    Manual,
}
// Serializes to: "TabChange" — correct PascalCase
```

### Naming Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Enum name | PascalCase | `ScreenshotTrigger`, `AppError` |
| Variants | PascalCase | `TabChange`, `AppSwitch` |
| File name | snake_case | `screenshot_trigger.rs` |
| Module | Re-export from `mod.rs` | `pub use screenshot_trigger::ScreenshotTrigger;` |
| Derives | Always include `Debug, Clone, PartialEq` | Plus `Serialize, Deserialize` for API types |

### Error Enums with `thiserror`

```rust
use thiserror::Error;

#[derive(Debug, Error)]
pub enum ConfigError {
    #[error("config file not found: {0}")]
    FileNotFound(String),

    #[error("parse error: {0}")]
    ParseError(#[from] serde_json::Error),

    #[error("invalid value for {field}: {reason}")]
    InvalidValue { field: String, reason: String },
}
```

### Pattern Matching — Must Be Exhaustive

```rust
// ✅ CORRECT — exhaustive match
match trigger {
    ScreenshotTrigger::Periodic  => schedule_next(),
    ScreenshotTrigger::TabChange => capture_tab(),
    ScreenshotTrigger::AppSwitch => capture_app(),
    ScreenshotTrigger::Idle      => capture_idle(),
    ScreenshotTrigger::Manual    => capture_now(),
}

// ❌ FORBIDDEN — magic string
if trigger_str == "TabChange" { ... }
```

### Data-Carrying Variants

```rust
#[derive(Debug)]
pub enum Command {
    Start { port: u16 },
    Stop,
    Status,
    Capture { output: PathBuf, format: ImageFormat },
}

// Pattern matching extracts data
match cmd {
    Command::Start { port } => start_server(port),
    Command::Stop => stop_server(),
    Command::Status => print_status(),
    Command::Capture { output, format } => take_capture(output, format),
}
```

### Folder Structure

```
src/
├── models/
│   ├── mod.rs
│   ├── activity.rs          # ActivityEvent enum
│   ├── session.rs
│   └── screenshot.rs        # ScreenshotTrigger enum
├── event.rs                 # EventSender type alias
└── config.rs
```

---

## Cross-Language Comparison

| Feature | Go | TypeScript | PHP | Rust |
|---------|-----|-----------|-----|------|
| Underlying type | `byte` (iota) | String enum | String-backed enum | Algebraic (zero-cost) |
| Type name | `Variant` (in package) | PascalCase | PascalCase + `Type` suffix | PascalCase |
| Zero value | `Invalid` | N/A | N/A | N/A (no implicit default) |
| String parsing | `Parse()` function | `Object.values().find()` | `::from()` / `::tryFrom()` | `serde` / `FromStr` |
| Comparison | `Is{Value}()` method | `=== Enum.Member` | `isEqual()` method | `==` (derives `PartialEq`) |
| Location | `internal/enums/{name}type/` | `src/lib/enums/` | `includes/Enums/` | `src/models/` or `src/enums/` |
| Data variants | Not supported (struct fields) | Not supported | Not supported | Supported natively |
| JSON support | `MarshalJSON/UnmarshalJSON` | Built-in (string) | `->value` | `Serialize/Deserialize` derive |

---

## AI Validation Checklist

Before generating any enum-related code:

- [ ] Used `enum` syntax (not string union or `const` object)
- [ ] PascalCase for all enum members/cases
- [ ] No raw string literals in comparisons
- [ ] Go: `byte` type, `Invalid` zero value, `iota`, package-scoped constants, no type prefix
- [ ] TypeScript: String enum with UPPER_SNAKE values
- [ ] PHP: `Type` suffix, string-backed, `isEqual()` method present
- [ ] Rust: Derive `Debug, Clone, PartialEq` minimum; `Serialize/Deserialize` for API types
- [ ] `default` branch in every switch/match on enum values

---

*Consolidated enum standards v1.0.0 — Go, TypeScript, PHP, Rust — 2026-04-10*

---

## §11 Cross-Language Enum Generator — Source-of-Truth Workflow

This section closes the gap on **how** enums stay synchronized across Go, TypeScript, PHP, and Rust. Without this, a blind AI will add an enum value to one language and ship a partial implementation that compiles in CI but breaks downstream language ports.

### 11.1 The Source of Truth

Every cross-language enum is defined **once** in a YAML manifest:

```
spec/<module>/enums/<EnumName>.yaml
```

Example: `02-spec/03-error-manage/enums/ErrorSeverity.yaml`

```yaml
name: ErrorSeverity
description: Severity classification for thrown errors
package: apperror
values:
  - name: Info
    ordinal: 0
    description: Informational only; no user impact
  - name: Warning
    ordinal: 1
    description: Degraded behavior; recoverable
  - name: Error
    ordinal: 2
    description: Failure; user-visible
  - name: Critical
    ordinal: 3
    description: System-level failure; alert ops
```

**Required fields:** `name`, `description`, `package`, `values[].name`, `values[].ordinal`, `values[].description`.

### 11.2 The Generator Scripts

| Generator | Reads | Emits | Output Path Pattern |
|-----------|-------|-------|---------------------|
| `gen-go-enums.mjs` | `spec/**/enums/*.yaml` | Go file with `iota` block, `String()`, `ParseEnum()`, `MarshalJSON`, `UnmarshalJSON` | `internal/<package>/<enum_name>_generated.go` |
| `gen-ts-enums.mjs` | same | TS union type, `parse()`, `is<Name>()` guard, JSON I/O | `src/lib/enums/<enumName>.generated.ts` |
| `gen-php-enums.mjs` | same | PHP 8.1 backed enum with `from()`, `tryFrom()`, `cases()` | `src/Enums/<EnumName>.php` |
| `gen-rust-enums.mjs` | same | Rust `enum` with `FromStr`, `Display`, `Serialize`, `Deserialize` derives | `src/enums/<enum_name>_generated.rs` |

All four generators are invoked from `scripts/codegen/`.

### 11.3 Invocation

```bash

# Single enum, all languages:

node scripts/codegen/gen-all-enums.mjs --enum ErrorSeverity

# All enums in a module, all languages:

node scripts/codegen/gen-all-enums.mjs --module 03-error-manage

# Full repo regeneration (used in CI):

node scripts/codegen/gen-all-enums.mjs --all
```

### 11.4 Drift Detection in CI

CI runs `gen-all-enums.mjs --all` and then `git diff --exit-code` on all generated paths. **Any uncommitted generator output fails CI** with:

```
Error: enum codegen drift detected in <file>
Hint: run `node scripts/codegen/gen-all-enums.mjs --all` and commit the result
```

### 11.5 The 7-Step Workflow for Adding/Modifying an Enum

```
1. Edit spec/<module>/enums/<EnumName>.yaml — add/modify a value
2. node scripts/codegen/gen-all-enums.mjs --enum <EnumName>
3. Verify the 4 emitted files compile in their respective languages
4. Update tests that reference the enum (one per language)
5. Run language-specific test suites (go test / npm test / phpunit / cargo test)
6. node scripts/sync-version.mjs && node scripts/sync-spec-tree.mjs
7. Commit YAML + 4 generated files + tests + version.json in ONE atomic commit
```

**Skipping any step** results in either CI drift failure (skip step 2 or 6) or runtime breakage in one language (skip step 3 or 4).

### 11.6 Forbidden Patterns

A blind AI must **never**:

- Hand-edit a `*_generated.{go,ts,php,rs}` file — changes will be overwritten on next codegen.
- Define an enum directly in a language file — it will collide with future codegen.
- Use a numeric value directly in business code — use the enum constant via `ParseEnum`/`from`/`parse`.
- Add a new value with a non-sequential `ordinal` — codegen requires monotonic ordinals.

### 11.7 Cross-Reference

- §27 (Error Code Registry) uses the same generator pattern for error codes.
- §10 (this file) documents the parsing methods that the generated code provides.

---

*Cross-Language Enum Generator added — v1.1.0 — 2026-04-22*
