# Info-Object Pattern for Go Enums

**Version:** 3.2.0  
**Status:** Complete  
**Updated:** 2026-04-16

---

## Purpose

This document defines the **info-object pattern** — the Go-idiomatic way to attach rich metadata (labels, icons, CSS classes, descriptions, sort order) to enum variants. This is the Go equivalent of PHP's `match`-based metadata methods.

> **PHP equivalent:** Each metadata field is a `match` method on the enum. See [WP Plugin Spec — Enum Metadata Pattern](../../../18-wp-plugin-how-to/02-enums-and-coding-style/02-enum-metadata-pattern.md).

---

## 1. Core Concept

Instead of scattering metadata across multiple methods, Go enums centralise all variant metadata in a **single info struct** and a **single lookup map**.

```
PHP approach (match per method):          Go approach (info-object):
──────────────────────────────            ──────────────────────────
label()   → match($this) {...}            variantInfoMap[variant].Label
icon()    → match($this) {...}            variantInfoMap[variant].Icon
cssClass() → match($this) {...}           variantInfoMap[variant].CssClass
```

**Why info-object in Go?**
- Go has no `match` expression — `switch` is verbose for per-variant metadata
- A single map lookup is more efficient than multiple switch statements
- The info struct is compile-time typed — missing fields are caught immediately
- Adding a new metadata field requires changing one struct + one map, not N methods

---

## 2. Info Struct Definition

Define a struct that holds all metadata for a single variant:

```go
package statustype

// VariantInfo holds metadata for a single status variant.
type VariantInfo struct {
    Label       string
    Description string
    Icon        string
    CssClass    string
    SortOrder   int
    IsTerminal  bool
}
```

**Rules:**
- Struct name: `VariantInfo` (always this name, in every enum package)
- All fields are **exported** (PascalCase)
- Fields are value types only — no pointers, no interfaces
- Add fields as needed per domain — not every enum needs all fields

---

## 3. Info Map — Single Source of Truth

Create a package-level map from `Variant` to `VariantInfo`:

```go
var variantInfoMap = map[Variant]VariantInfo{
    Invalid: {
        Label:       "Invalid",
        Description: "Unknown or unset status",
        Icon:        "❌",
        CssClass:    "status-invalid",
        SortOrder:   0,
        IsTerminal:  false,
    },
    Pending: {
        Label:       "Pending",
        Description: "Awaiting processing",
        Icon:        "⏳",
        CssClass:    "status-pending",
        SortOrder:   1,
        IsTerminal:  false,
    },
    Complete: {
        Label:       "Complete",
        Description: "Successfully finished",
        Icon:        "✅",
        CssClass:    "status-complete",
        SortOrder:   2,
        IsTerminal:  true,
    },
    Failed: {
        Label:       "Failed",
        Description: "Terminated with error",
        Icon:        "❌",
        CssClass:    "status-failed",
        SortOrder:   3,
        IsTerminal:  true,
    },
}
```

**Rules:**
- Map name: `variantInfoMap` (unexported — accessed via methods only)
- Every variant in `All()` MUST have an entry — enforced by unit test
- `Invalid` variant MUST have an entry (used as fallback)

---

## 4. Accessor Methods

Expose metadata through methods on `Variant`. Each method delegates to the info map:

```go
// Info returns the full metadata object for this variant.
func (v Variant) Info() VariantInfo {
    if info, ok := variantInfoMap[v]; ok {
        return info
    }

    return variantInfoMap[Invalid]
}

// Label returns the human-readable label.
func (v Variant) Label() string {
    return v.Info().Label
}

// Description returns the longer description.
func (v Variant) Description() string {
    return v.Info().Description
}

// Icon returns the display icon.
func (v Variant) Icon() string {
    return v.Info().Icon
}

// CssClass returns the CSS class for UI rendering.
func (v Variant) CssClass() string {
    return v.Info().CssClass
}

// IsTerminal returns whether this status is a final state.
func (v Variant) IsTerminal() bool {
    return v.Info().IsTerminal
}
```

**Rules:**
- `Info()` is the single access point — all other methods delegate to it
- Unknown/invalid variants fall back to the `Invalid` entry (never panic)
- Method names match the `VariantInfo` field names exactly

---

## 5. Relationship to variantLabels

The existing `variantLabels` array (used by `String()` and `Parse()`) remains unchanged. The info-object pattern **extends** it — it does not replace it.

```go
// variantLabels — used by String() and Parse() (required by enum spec)
var variantLabels = [...]string{
    Invalid:  "Invalid",
    Pending:  "Pending",
    Complete: "Complete",
    Failed:   "Failed",
}

// variantInfoMap — used by Info(), Label(), Icon(), etc. (info-object pattern)
var variantInfoMap = map[Variant]VariantInfo{
    // ... rich metadata
}
```

| Concern | Source |
|---------|--------|
| `String()` / `Parse()` | `variantLabels` array (fast, index-based) |
| Rich metadata (icon, CSS, description) | `variantInfoMap` map |
| `Label()` | Delegates to `Info().Label` (same value as `variantLabels` entry) |

> **Rule:** `variantLabels[v]` and `variantInfoMap[v].Label` MUST return the same string for every variant. Enforce via unit test.

---

## 6. When to Use Info-Object vs Plain Enum

| Scenario | Pattern |
|----------|---------|
| Enum only needs `String()` / `Parse()` / `Is*()` | **Plain enum** — `variantLabels` only |
| Enum has 2+ metadata fields (icon, CSS, description, sort) | **Info-object** — add `VariantInfo` + `variantInfoMap` |
| Metadata is used in UI rendering or API responses | **Info-object** — always |
| Enum is a simple internal discriminator | **Plain enum** — keep it simple |

---

## 7. Testing

### 7.1 Every Variant Has an Info Entry

```go
func TestAllVariantsHaveInfo(t *testing.T) {
    for _, v := range All() {
        info := v.Info()
        assert.NotEmpty(t, info.Label, "variant %s must have a label", v.String())
    }
}
```

### 7.2 Label Consistency

```go
func TestInfoLabelMatchesStringLabel(t *testing.T) {
    for _, v := range All() {
        assert.Equal(
            t,
            v.String(),
            v.Info().Label,
            "Info().Label must match String() for variant %s",
            v.String(),
        )
    }
}
```

---

## 8. Complete Reference Implementation

```go
package statustype

// Variant represents a processing status.
type Variant byte

const (
    Invalid  Variant = iota
    Pending
    Complete
    Failed
)

// --- String/Parse (standard enum spec) ---

var variantLabels = [...]string{
    Invalid:  "Invalid",
    Pending:  "Pending",
    Complete: "Complete",
    Failed:   "Failed",
}

func (v Variant) String() string {
    if !v.IsValid() {
        return variantLabels[Invalid]
    }

    return variantLabels[v]
}

func (v Variant) IsValid() bool {
    return v >= Pending && v <= Failed
}

// --- Info-object pattern ---

type VariantInfo struct {
    Label       string
    Description string
    Icon        string
    CssClass    string
    SortOrder   int
    IsTerminal  bool
}

var variantInfoMap = map[Variant]VariantInfo{
    Invalid: {
        Label:       "Invalid",
        Description: "Unknown or unset status",
        Icon:        "❌",
        CssClass:    "status-invalid",
        SortOrder:   0,
        IsTerminal:  false,
    },
    Pending: {
        Label:       "Pending",
        Description: "Awaiting processing",
        Icon:        "⏳",
        CssClass:    "status-pending",
        SortOrder:   1,
        IsTerminal:  false,
    },
    Complete: {
        Label:       "Complete",
        Description: "Successfully finished",
        Icon:        "✅",
        CssClass:    "status-complete",
        SortOrder:   2,
        IsTerminal:  true,
    },
    Failed: {
        Label:       "Failed",
        Description: "Terminated with error",
        Icon:        "❌",
        CssClass:    "status-failed",
        SortOrder:   3,
        IsTerminal:  true,
    },
}

func (v Variant) Info() VariantInfo {
    if info, ok := variantInfoMap[v]; ok {
        return info
    }

    return variantInfoMap[Invalid]
}

func (v Variant) Label() string       { return v.Info().Label }
func (v Variant) Description() string { return v.Info().Description }
func (v Variant) Icon() string        { return v.Info().Icon }
func (v Variant) CssClass() string    { return v.Info().CssClass }
func (v Variant) IsTerminal() bool    { return v.Info().IsTerminal }
```

---

## Cross-Language Comparison

| Language | Metadata Pattern | Lookup Structure | Access |
|----------|-----------------|------------------|--------|
| **Go** | Info-object (this document) | `map[Variant]VariantInfo` | `v.Info().Field` |
| **PHP** | `match` expression per method | Built into language | `$v->label()` |
| **TypeScript** | Record lookup map | `Record<EnumValue, EnumInfo>` | `enumInfo[v].field` |

---

## Cross-References

- [01-enum-pattern.md](01-enum-pattern.md) — Core byte-based enum pattern and `variantLabels`
- [02-required-methods.md](02-required-methods.md) — Mandatory methods (String, Parse, Is*, All)
- [PHP Enum Metadata Pattern](../../../18-wp-plugin-how-to/02-enums-and-coding-style/02-enum-metadata-pattern.md) — PHP `match`-based equivalent

---

*Info-object pattern for Go enums — v3.2.0 — 2026-04-16*
