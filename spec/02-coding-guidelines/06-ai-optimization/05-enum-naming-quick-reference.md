# Enum Naming Quick Reference — All Languages

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Purpose:** Single-page AI reference for enum declaration, naming, and usage rules across Go, TypeScript, PHP, and Rust

---

## Universal Rules (All Languages)

| Rule | Requirement |
|------|-------------|
| No magic strings | Never compare against raw string literals — always use enum constants |
| PascalCase values | Enum members/cases use PascalCase (`Production`, not `PRODUCTION` or `production`) |
| One definition | String representations defined **once**, co-located with the enum type |
| Exhaustive switch | Every `switch`/`match` on an enum must have a `default` branch |
| No string unions | Use proper `enum` syntax — never `type Foo = 'a' | 'b'` |

---

## Go

### Declaration

```go
package environmenttype          // ← package name = grouping (no type prefix on constants)

type Variant byte                // ← always byte, never string or int

const (
    Invalid     Variant = iota   // ← always first (zero value)
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
| String constants | `{Value}Str` (package-scoped) | `ProductionStr` (not `EnvironmentProductionStr`) |
| Zero value | Always `Invalid` | Never `Unknown` or `None` |
| Lookup table | `variantLabels` (unexported) | Single array, PascalCase values |

### Required Methods (every enum)

| Method | Signature | Purpose |
|--------|-----------|---------|
| `String()` | `func (v Variant) String() string` | Serialization/logging |
| `Label()` | `func (v Variant) Label() string` | Delegates to `String()` |
| `Is{Value}()` | `func (v Variant) IsSerpApi() bool` | Type-safe comparison |
| `IsValid()` | `func (v Variant) IsValid() bool` | Bounds check |
| `MarshalJSON()` | `func (v Variant) MarshalJSON() ([]byte, error)` | JSON output |
| `UnmarshalJSON()` | `func (v *Variant) UnmarshalJSON(b []byte) error` | JSON input |
| `Parse()` | `func Parse(s string) (Variant, error)` | Case-insensitive string→enum |

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

### ❌ Forbidden Patterns

```go
type Provider string                    // ❌ string-based enum
if provider == "serpapi" { ... }        // ❌ magic string comparison
EnvironmentProductionStr = "production" // ❌ type-prefixed constant name
const ( Unknown Variant = iota )       // ❌ "Unknown" as zero value
```

---

## TypeScript

### Declaration

```typescript
// src/lib/enums/connection-status.ts

export enum ConnectionStatus {       // ← PascalCase name, string enum
  Connected = "CONNECTED",
  Disconnected = "DISCONNECTED",
  Connecting = "CONNECTING",
  Error = "ERROR",
}
```

### Naming Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Enum name | PascalCase | `ConnectionStatus`, `LogLevel` |
| File name | kebab-case | `connection-status.ts`, `log-level.ts` |
| Members | PascalCase | `Connected`, `Disconnected` |
| Values | UPPER_SNAKE string | `"CONNECTED"`, `"DISCONNECTED"` |
| Folder | `src/lib/enums/` | One file per enum |

### Usage Patterns

```typescript
// ✅ CORRECT
if (ws.status === ConnectionStatus.Connected) { ... }

// ❌ FORBIDDEN
if (ws.status === 'connected') { ... }

// ✅ CORRECT type definition
interface WsState { status: ConnectionStatus; }

// ❌ FORBIDDEN string union
interface WsState { status: 'connected' | 'disconnected'; }
```

### Defined Enums

| Enum | File | Values |
|------|------|--------|
| `ConnectionStatus` | `01-connection-status-enum.md` | Connected, Disconnected, Connecting, Reconnecting, Error |
| `EntityStatus` | `02-entity-status-enum.md` | Active, Inactive, Pending, Archived |
| `ExecutionStatus` | `03-execution-status-enum.md` | Pending, Running, Completed, Failed, Cancelled |
| `ExportStatus` | `04-export-status-enum.md` | Pending, Processing, Completed, Failed |
| `HttpMethod` | `05-http-method-enum.md` | Get, Post, Put, Patch, Delete |
| `MessageStatus` | `06-message-status-enum.md` | Pending, Streaming, Completed, Error |
| `LogLevel` | `10-log-level-enum.md` | Debug, Info, Warn, Error, Fatal |

---

## PHP

### Declaration

```php
// includes/Enums/HttpMethodType.php

namespace RiseupAsia\Enums;

enum HttpMethodType: string {        // ← string-backed, Type suffix
    case Get    = 'GET';
    case Post   = 'POST';
    case Put    = 'PUT';
    case Delete = 'DELETE';
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
| No prefix | No `RISEUP_` prefix | Namespace provides scoping |

### Required Methods

| Method | Purpose |
|--------|---------|
| `isEqual(self $other): bool` | Type-safe comparison (mandatory on every backed enum) |
| `validValues(): array` | Static — returns all valid values |

### Parsing — No Manual Switch

```php
// ✅ REQUIRED: Use built-in backed enum parsing
$method = HttpMethodType::from($input);       // throws ValueError
$method = HttpMethodType::tryFrom($input);    // returns null

// ❌ FORBIDDEN: Manual match with raw strings
match ($input) {
    'GET' => HttpMethodType::Get,
    'POST' => HttpMethodType::Post,
}
```

---

## Rust

### Declaration

```rust
// src/enums/environment_type.rs

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]       // ← MANDATORY for JSON/API serialization
pub enum EnvironmentType {
    Production,
    Staging,
    Development,
}
```

### Naming Rules

| Element | Convention | Example |
|---------|-----------|---------|
| Enum name | PascalCase | `EnvironmentType`, `LogLevel` |
| File name | snake_case | `environment_type.rs` |
| Variants | PascalCase | `Production`, `Staging` |
| **String values** | **PascalCase (mandatory)** | `"Production"` via `#[serde(rename_all = "PascalCase")]` |
| Module name | snake_case | `mod environment_type;` |
| Folder | `src/enums/` | One file per enum |

> **🔴 AI CRITICAL:** Rust identifiers (functions, variables, modules) use snake_case per RFC 430, BUT enum string values serialized to JSON/DB/API MUST be PascalCase. Use `#[serde(rename_all = "PascalCase")]` for serialization and `#[sqlx(rename_all = "PascalCase")]` for DB mappings.

### Required Derives & Traits

| Derive/Trait | Purpose |
|-------------|---------|
| `Debug` | Logging and debugging output |
| `Clone` | Value duplication |
| `PartialEq` | Comparison via `==` |
| `Serialize` / `Deserialize` | JSON/API serialization (serde) |
| `Display` (optional) | Custom `to_string()` output |

### DB Mapping Example

```rust
#[derive(Debug, Clone, PartialEq, sqlx::Type)]
#[sqlx(rename_all = "PascalCase")]          // ← DB column values = PascalCase
pub enum StatusType {
    Active,
    Inactive,
    Pending,
}
```

### ❌ Forbidden Patterns

```rust
#[serde(rename_all = "snake_case")]         // ❌ snake_case string values
#[serde(rename_all = "SCREAMING_SNAKE_CASE")] // ❌ UPPER_SNAKE string values
if status == "active" { ... }               // ❌ magic string comparison
type Status = &str;                         // ❌ string alias instead of enum
```

---

## Cross-Language Comparison

| Feature | Go | TypeScript | PHP | Rust |
|---------|-----|-----------|-----|------|
| Underlying type | `byte` (iota) | String enum | String-backed enum | Sum type (derive) |
| Type name | `Variant` (in package) | PascalCase | PascalCase + `Type` suffix | PascalCase |
| Zero value | `Invalid` | N/A | N/A | N/A (use `Default` trait) |
| String parsing | `Parse()` function | `Object.values().find()` | `::from()` / `::tryFrom()` | `serde` / `FromStr` |
| Comparison | `Is{Value}()` method | `=== Enum.Member` | `isEqual()` method | `== Enum::Variant` |
| Location | `internal/enums/{name}type/` | `src/lib/enums/` | `includes/Enums/` | `src/enums/` |
| Grouping | Package name | Enum name | Namespace + enum name | Module name |
| **String value casing** | **PascalCase** | **UPPER_SNAKE** | **PascalCase** | **PascalCase** |
| **Identifier casing** | **PascalCase** | **PascalCase** | **PascalCase** | **snake_case** |

---

## AI Validation Checklist

Before generating any enum-related code:

- [ ] Used `enum` syntax (not string union or `const` object)
- [ ] PascalCase for all enum members/cases/variants
- [ ] No raw string literals in `switch`/`case`/`match`/`if` comparisons
- [ ] Go: `byte` type, `Invalid` zero value, `iota`, package-scoped constants
- [ ] Go: No type-prefixed constant names — package name provides grouping
- [ ] TypeScript: String enum with UPPER_SNAKE values
- [ ] PHP: `Type` suffix, string-backed, `isEqual()` method present
- [ ] **Rust: `#[serde(rename_all = "PascalCase")]` on every serialized enum**
- [ ] **Rust: `#[sqlx(rename_all = "PascalCase")]` on every DB-mapped enum**
- [ ] **Rust: `Debug`, `Clone`, `PartialEq` derived on every enum**
- [ ] `default` branch in every switch/match on enum values

---

## Cross-References

- [Go Enum Specification](../03-golang/01-enum-specification/00-overview.md) — Full Go enum pattern, methods, folder structure
- [TypeScript Enums](../02-typescript/00-overview.md) — All TypeScript enum definitions
- [PHP Enums](../04-php/01-enums.md) — PHP backed enum patterns and rules
- [Rust Naming Conventions](../05-rust/01-naming-conventions.md) — Rust hybrid naming strategy and enum string rules
- [Code Style — Switch Exemption](../01-cross-language/04-code-style/00-overview.md) — Switch-based enum parser exemption from nesting ban
- [AI Quick Reference Checklist](./02-ai-quick-reference-checklist.md) — Broader code validation checklist

---

*Enum naming quick reference v1.1.0 — 2026-04-11*
