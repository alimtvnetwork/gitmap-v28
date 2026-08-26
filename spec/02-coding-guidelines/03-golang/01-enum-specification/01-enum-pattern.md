# Enum Pattern

**Version:** 3.2.0  
**Status:** Complete  
**Updated:** 2026-04-16

---

## Core Pattern

All enums MUST use `byte` as the underlying type for memory efficiency and performance.

### Declaration

```go
package provider

// Variant represents the provider type
type Variant byte

const (
    // Invalid is the zero value (always first)
    Invalid Variant = iota
    
    // SerpApi is the SerpAPI provider
    SerpApi
    
    // MapsScraper is the gosom maps scraper
    MapsScraper
    
    // Colly is the Colly web scraper
    Colly
)
```

---

## Rules

### 1. Always Start with Invalid

The first constant MUST be `Invalid` as the zero-value representing an invalid state:

```go
const (
    Invalid Variant = iota  // Always first
    // ... other values
)
```

### 2. Use `byte` Type

```go
// ✅ Correct
type Variant byte

// ❌ Wrong
type Variant string
type Variant int
```

### 3. Use `iota` for Sequential Values

```go
// ✅ Correct
const (
    Invalid Variant = iota
    TypeA
    TypeB
)

// ❌ Wrong - explicit values
const (
    Invalid Variant = 0
    TypeA   Variant = 1
    TypeB   Variant = 2
)
```

### 4. PascalCase for Variant Names and Values

> **Rule:** For non-string enums (byte/int-based), both the constant name AND its serialized label MUST use **PascalCase**. This applies to GoLang, TypeScript, and C#. `camelCase` and `_camelCase` are strictly forbidden. For Rust or other languages, follow the ecosystem's dominant convention.

```go
// ✅ Correct — PascalCase names
const (
    SerpApi     Variant = iota
    MapsScraper
)

// ✅ Correct — PascalCase labels in lookup table
var variantLabels = [...]string{
    SerpApi:     "SerpApi",      // ← PascalCase value
    MapsScraper: "MapsScraper",  // ← PascalCase value
}

// ❌ Wrong — SCREAMING_CASE
const (
    SERP_API     Variant = iota
    maps_scraper
)

// ❌ Wrong — lowercase/kebab-case labels
var variantLabels = [...]string{
    SerpApi:     "serp-api",       // ← wrong: kebab-case
    MapsScraper: "maps_scraper",   // ← wrong: snake_case
}
```

#### Cross-Language Enum Value Rules

| Enum Type | Value Format | Example |
|-----------|-------------|---------|
| Non-string (byte/int/iota) | **PascalCase** | `SerpApi`, `MapsScraper`, `HttpGet` |
| String enum (rare, discouraged) | **PascalCase** | `"SerpApi"`, `"MapsScraper"` |
| Slug/URL representation | **kebab-case** | `serp-api`, `maps-scraper` |

> **Important:** When an enum value needs to appear in a URL or slug context, convert it to kebab-case at the serialization boundary — never store it as kebab-case in the enum itself.

### 5. Document Each Variant

```go
const (
    // Invalid represents an unspecified provider
    Invalid Variant = iota
    
    // SerpApi uses the commercial SerpAPI service
    SerpApi
    
    // MapsScraper uses gosom/google-maps-scraper
    MapsScraper
)
```

---

## Internal Lookup Table

Use a single unexported array for all lookups, serialization, and display. All values MUST be **PascalCase** matching the constant name:

```go
var variantLabels = [...]string{
    Invalid:     "Invalid",
    SerpApi:     "SerpApi",
    MapsScraper: "MapsScraper",
    Colly:       "Colly",
}
```

> **Note:** A single `variantLabels` table replaces the previous dual-table pattern (`variantStrings` + `variantLabels`). `Label()` delegates to `String()`. `Parse()` uses `strings.EqualFold()` for case-insensitive matching.

### PascalCase Rule

| ❌ Forbidden | ✅ Required |
|-------------|-----------|
| `"invalid"` | `"Invalid"` |
| `"serpapi"` | `"SerpApi"` |
| `"SerpAPI"` | `"SerpApi"` |
| `"maps_scraper"` | `"MapsScraper"` |
| `"per_table"` | `"PerTable"` |
| `"baseURL"` | `"BaseUrl"` |

**Exception:** Protocol-driven enums (`content_type`, `endpoint`, `header`, `response_key`, `response_message`) preserve their functional values (e.g., `"application/json"`, `"X-Riseup-Auth"`).

---

## Why `byte`?

| Aspect | `byte` | `string` | `int` |
|--------|--------|----------|-------|
| Memory | 1 byte | 16+ bytes | 8 bytes |
| Comparison | O(1) | O(n) | O(1) |
| Switch | Jump table | String compare | Jump table |
| JSON size | Uses String() | Direct | Uses String() |
| Type safety | ✅ Strong | ⚠️ Weak | ✅ Strong |

---

## Anti-Patterns

### ❌ String-Based Enums

```go
// DON'T DO THIS
type Provider string

const (
    SerpApi Provider = "serpapi"
)
```

### ❌ Hardcoded Strings in Logic

```go
// DON'T DO THIS
if provider == "serpapi" { ... }

// DO THIS
if provider.IsSerpApi() { ... }
```

### ❌ Switch Without Exhaustive Check

```go
// DON'T DO THIS
switch p {
case SerpApi:
    // ...
}

// DO THIS
switch p {
case SerpApi:
    // ...
case MapsScraper:
    // ...
case Colly:
    // ...
default:
    return fmt.Errorf("invalid provider: %s", p)
}
```

### ❌ Using "Unknown" as Zero Value

```go
// DON'T DO THIS
const (
    Unknown Variant = iota
)

// DO THIS
const (
    Invalid Variant = iota
)
```

---

*Core enum pattern specification.*
