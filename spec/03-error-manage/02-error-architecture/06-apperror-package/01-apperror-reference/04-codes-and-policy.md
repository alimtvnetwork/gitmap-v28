# AppError Package Reference — Error code convention, stack trace skip rules, file size

> **Parent:** [AppError Package Reference](./00-overview.md)  
> **Version:** 1.3.0  
> **Updated:** 2026-03-31

---

## 6. Error Code Convention

Error codes are defined as string constants in `codes.go`. **No magic strings.**

| Range | Category |
|-------|----------|
| E1xxx | Configuration errors |
| E2xxx | Database errors |
| E3xxx | WordPress API errors |
| E4xxx | File system errors |
| E5xxx | Sync errors |
| E6xxx | Backup errors |
| E7xxx | Git errors |
| E8xxx | Build errors |
| E9xxx | General errors |
| E10xxx | E2E test errors |
| E11xxx | Publish errors |
| E12xxx | Version errors |
| E13xxx | Session errors |
| E14xxx | Crypto errors |

### 6.1 Error Type Enum Pattern (Best Practice)

Instead of passing raw string codes and duplicate messages, define error types as **byte enums** in the `apperrtype` package. Each enum maps to both a code and a default message via a detail map.

**Three maturity levels:**

```go
// Level 1 — raw strings (acceptable for prototyping, not for production)
apperror.New("E2010", "site not found")

// Level 2 — enum code, manual message (better — no magic code strings)
apperror.New(apperrtype.SiteNotFound, "site not found")

// Level 3 — enum with built-in message (best — zero duplication)
apperror.NewType(apperrtype.SiteNotFound)
```

**Structure (v2.0 — single `Variation uint16` enum):**

```go
// apperrtype/variant_structure.go
package apperrtype

// VariantStructure holds the name, code, message, and variant for each error type.
type VariantStructure struct {
    Name    string     // "SiteNotFound"
    Code    string     // "E2010"
    Message string     // "site not found"
    Variant Variation  // the enum value itself
}

// ErrorType is the interface all error type enums must implement.
type ErrorType interface {
    Code() string
    Message() string
    Name() string
}
```

```go
// apperrtype/variation.go — single enum for ALL error types
package apperrtype

type Variation uint16

const (
    NoError Variation = iota // 0
    // ... E1xxx–E18xxx domains ...
    SiteNotFound        // E2010
    SiteBlocked         // E2011
    PluginSlugMissing   // E2012
    PluginNotFound      // E2013
    PluginAlreadyActive // E2014
    // ...
    MaxError // sentinel
)
```

```go
// apperrtype/variant_registry.go — global registry
var variantRegistry = map[Variation]VariantStructure{
    SiteNotFound:      {Name: "SiteNotFound",      Code: "E2010", Message: "site not found",      Variant: SiteNotFound},
    SiteBlocked:       {Name: "SiteBlocked",        Code: "E2011", Message: "site is blocked",     Variant: SiteBlocked},
    PluginSlugMissing: {Name: "PluginSlugMissing",  Code: "E2012", Message: "plugin slug required", Variant: PluginSlugMissing},
    // ... all domains in one map
}
```

```go
// apperror/constructors.go — NewType creates AppError from any ErrorType enum
func NewType(errType apperrtype.ErrorType) *AppError {
    return New(errType.Code(), errType.Message())
}
```

**Rules:**
- Single `Variation uint16` enum in `variation.go` — all domains in one file
- Global `variantRegistry` maps each `Variation` → `VariantStructure{Name, Code, Message, Variant}`
- `Variation` implements `ErrorType` interface (`Code()` + `Message()` + `Name()`)
- `StringToVariantMap` provides reverse lookup from name strings to `Variation` values
- `apperror.NewType()` / `apperror.WrapType()` create `AppError` from any `Variation`
- All types live in `types/apperrtype/`
- Never pass raw string codes (`"E2012"`) when an `apperrtype` variant exists

> 📖 Full enum specification: [05-apperrtype-enums.md](./05-apperrtype-enums.md)

---


---

## 7. Stack Trace Skip Rules

Understanding skip values is critical for accurate error attribution.

The table below shows what each constructor passes to its underlying `CaptureStack` call. `WrapWithSkip` has a base of `3` and `NewWithSkip` has a base of `2` because `Wrap` delegates through one extra internal frame.

| Constructor | Delegates To | `skip` Passed | Effective `CaptureStack` | Reason |
|-------------|-------------|---------------|--------------------------|--------|
| `New()` | `CaptureStack(2)` | — | 2 | Skips `CaptureStackN` + `CaptureStack` + `New` |
| `Wrap()` | `WrapWithSkip(…, 0)` | 0 | 3 | Skips through `Wrap` → `WrapWithSkip` → `CaptureStack` chain |
| `NewWithSkip()` | `CaptureStack(2+skip)` | caller-provided | 2 + skip | Additional skip on top of `New` base |
| `WrapWithSkip()` | `CaptureStack(3+skip)` | caller-provided | 3 + skip | Additional skip on top of `Wrap` base |
| `FailWrap()` | `WrapWithSkip(…, 0)` | 0 | 3 | Same depth as `Wrap` — replaces it, doesn't nest |
| `FailSliceWrap()` | `WrapWithSkip(…, 0)` | 0 | 3 | Same depth as `Wrap` — replaces it, doesn't nest |
| `FailMapWrap()` | `WrapWithSkip(…, 0)` | 0 | 3 | Same depth as `Wrap` — replaces it, doesn't nest |
| `FailNew()` | `NewWithSkip(…, 1)` | 1 | 3 | One frame deeper than `New` (FailNew → NewWithSkip) |
| `FailSliceNew()` | `NewWithSkip(…, 1)` | 1 | 3 | One frame deeper than `New` (FailSliceNew → NewWithSkip) |
| `FailMapNew()` | `NewWithSkip(…, 1)` | 1 | 3 | One frame deeper than `New` (FailMapNew → NewWithSkip) |

> **Key insight:** `FailWrap` calls `WrapWithSkip` directly (same as `Wrap` does), so it sits at the **same depth** and passes `skip=0`. `FailNew` calls `NewWithSkip` directly (one frame deeper than `New`), so it passes `skip=1`.

---


---

## 8. File Size Policy

| Target | Soft Limit | Description |
|--------|-----------|-------------|
| 300 lines | 400 lines | All files target 300 lines. Up to 400 is acceptable but marked `// NOTE: Needs refactor — exceeds 300-line target` at the top. |

---
