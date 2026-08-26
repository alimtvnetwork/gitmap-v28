# Enum Specification

**Version:** 3.2.0  
**Status:** Complete  
**Updated:** 2026-04-16    
**AI Confidence:** High  
**Ambiguity:** None
**Error Range:** N/A (Cross-cutting standard)

---


## Keywords

`coding`, `enum`, `golang`, `guidelines`, `specification`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |


## Purpose

This specification defines the **universal enum pattern** for all Go-based CLI applications in the ecosystem. All enums must follow this standard to ensure consistency, type safety, and maintainability.

---

## Index

| File | Purpose |
|------|---------|
| [01-enum-pattern.md](01-enum-pattern.md) | Core byte-based enum pattern |
| [02-required-methods.md](02-required-methods.md) | Mandatory methods for all enums |
| [03-folder-structure.md](03-folder-structure.md) | Directory layout standard |
| [04-validation-checklist.md](04-validation-checklist.md) | Compliance audit checklist |
| [05-info-object-pattern.md](05-info-object-pattern.md) | Rich metadata via info-object map pattern |

---

## Quick Reference

### Enum Declaration

```go
package providertype

type Variant byte

const (
    Invalid Variant = iota
    SerpApi
    MapsScraper
    Colly
)
```

> **Convention:** All enum packages end with `type` suffix (e.g., `providertype`, `httpmethodtype`). See [03-folder-structure.md](03-folder-structure.md) §Package Naming Convention.

### Required Methods

| Method | Signature | Purpose |
|--------|-----------|---------|
| `String` | `(v Variant) String() string` | String representation |
| `Label` | `(v Variant) Label() string` | Human-readable label |
| `Is{Value}` | `(v Variant) IsSerpApi() bool` | Type check for each variant |
| `IsOther` | `(v Variant) IsOther(other Variant) bool` | Inverse check — true if NOT the given variant |
| `IsAnyOf` | `(v Variant) IsAnyOf(others ...Variant) bool` | True if receiver matches any in the list |
| `All` | `All() []Variant` | Returns all valid variants |
| `ByIndex` | `ByIndex(i int) Variant` | Get variant by index |
| `Parse` | `Parse(s string) (Variant, error)` | Parse string to variant |
| `IsValid` | `(v Variant) IsValid() bool` | Check if variant is valid |
| `MarshalJSON` | `(v Variant) MarshalJSON() ([]byte, error)` | JSON serialization |
| `UnmarshalJSON` | `(v *Variant) UnmarshalJSON([]byte) error` | JSON deserialization |

### Key Rules

| Rule | Description |
|------|-------------|
| Zero value | Always `Invalid Variant = iota` (never `Unknown`) |
| variantLabels | Single lookup table for serialization, parsing, and display |

### Folder Structure

```
internal/enums/
├── providertype/
│   └── variant.go
├── platformtype/
│   └── variant.go
├── enginetype/
│   └── variant.go
└── searchmodetype/
    └── variant.go
```

---

## Applies To

| CLI / Project | Status | Score | Audit Report |
|---------------|--------|-------|--------------|
| GSearch CLI | ✅ Compliant | 50/50 | `.lovable/audits/gsearch-cli-enum-compliance-audit-2026-02-06.md` |
| BRun CLI | ✅ Compliant | 50/50 | `.lovable/audits/brun-cli-enum-compliance-audit-2026-02-06.md` |
| AI Bridge CLI | ✅ Compliant | 50/50 | `.lovable/audits/ai-bridge-cli-enum-compliance-audit-2026-02-06.md` |
| Nexus Flow CLI | ✅ Compliant | 50/50 | `.lovable/audits/nexus-flow-cli-enum-compliance-audit-2026-02-06.md` |
| Spec Reverse CLI | ✅ Compliant | 50/50 | `.lovable/audits/spec-reverse-cli-enum-compliance-audit-2026-02-06.md` |
| WP SEO Publish CLI | ✅ Compliant | 50/50 | `.lovable/audits/wp-seo-publish-cli-enum-compliance-audit-2026-02-06.md` |
| AI Transcribe CLI | ✅ Compliant | 50/50 | `.lovable/audits/ai-transcribe-cli-enum-compliance-audit-2026-02-06.md` |
| WP Plugin Builder | ✅ Compliant | 50/50 | `.lovable/audits/wp-plugin-builder-cli-enum-compliance-audit-2026-02-06.md` |
| Spec Management | ✅ Compliant | 50/50 | `.lovable/audits/spec-management-enum-compliance-audit-2026-02-06.md` |
| **WP Plugin Publish** | 🔄 Migration In Progress | — | 11/12 migrated, 1 int-based exempt |

> **Note:** All 9 CLIs have been migrated to `Invalid` as zero value per spec v2.0.0 (completed 2026-02-11).  
> **Note:** WP Plugin Publish backend — 11 of 12 string-based enums migrated to byte-based `internal/enums/` pattern (2026-02-21). `HttpStatusType` remains `int`-based (exempt).

### WP Plugin Publish — Enum Migration Tracker

| Old Type (`wordpress/`) | New Package (`enums/`) | Status |
|---|---|---|
| `StatusType` | `statustype.Variant` | ✅ Migrated |
| `PluginStatusType` | `pluginstatustype.Variant` | ✅ Migrated |
| `PostStatusType` | `poststatustype.Variant` | ✅ Migrated |
| `ActionType` | `actiontype.Variant` | ✅ Migrated |
| `ContentTypeValue` | `contenttypetype.Variant` | ✅ Migrated |
| `EndpointType` | `endpointtype.Variant` | ✅ Migrated |
| `HeaderType` | `headertype.Variant` | ✅ Migrated |
| `ResponseKeyType` | `responsekeytype.Variant` | ✅ Migrated |
| `ResponseMessageType` | `responsemessagetype.Variant` | ✅ Migrated |
| `SnapshotErrorType` | `snapshoterrrortype.Variant` | ✅ Migrated |
| `UploadSourceType` | `uploadsourcetype.Variant` | ✅ Migrated |
| `HttpStatusType` | — | ⏭️ Exempt (`int`) |

> **Pending:** Consumer import updates across handlers and services.

---

## Document Inventory

| File |
|------|
| 99-consistency-report.md |


## Cross-References

- Error Code Registry <!-- external: spec/03-error-manage/03-error-code-registry/01-registry.md -->
- Split DB Architecture <!-- external: spec/05-split-db-architecture/00-overview.md -->
- Coding Guidelines Memory <!-- external: .lovable/memories/constraints/coding-guidelines.md -->

---

*Universal enum standard for Go CLI ecosystem.*
