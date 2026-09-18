# 25 — Chrome Profile Management Specification

**Version:** 1.0.0
**Updated:** 2026-09-05
**Status:** Canonical
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose & Scope

The Chrome Profile Management module specifies the architectural standards, Local State schemas, process synchronization guards, and profile picker registration protocols for importing, exporting, copying, and managing Google Chrome / Chromium browser profiles in GitMap.

When a user imports or copies a Chrome profile snapshot via GitMap CLI commands (`gitmap chrome profile import`, `gitmap chrome-profile-copy`), the restored profile directory must be immediately and reliably visible inside Google Chrome's Profile Picker UI (`chrome://settings/manageProfile`), preserving extension bindings, bookmarks, and display metadata without risking in-memory state overwrites or profile corruption.

---

## Mandatory Module Scoring

- **AI Confidence Score:** `Production-Ready` (100% contracts defined, zero ambiguities in schema)
- **Ambiguity Score:** `None` (exact Chromium `Local State` and `Preferences` schemas documented)
- **Health Score (0-100):** `100` (25% `00-overview.md` + 25% `99-consistency-report.md` + 25% lowercase kebab-case naming + 25% numeric prefixes)

---

## Keywords

`chrome`, `chromium`, `profile-picker`, `local-state`, `info-cache`, `profiles-order`, `preferences`, `process-lock`, `reconciliation`, `smart-import`

---

## File Inventory

| File | Type | Description |
|------|------|-------------|
| [`00-overview.md`](./00-overview.md) | Entry Point | Module architecture, metadata, scoring, and inventory |
| [`01-profile-registration-and-picker.md`](./01-profile-registration-and-picker.md) | Specification | In-depth Chrome Local State schemas, picker attributes, and process guards |
| [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) | Verification | Given/When/Then acceptance criteria and CLI test vectors |
| [`99-consistency-report.md`](./99-consistency-report.md) | Audit | Cross-reference integrity and consistency evaluation |

---

## High-Level Architecture

```
                                 [ GitMap CLI ]
                                       │
                ┌──────────────────────┴──────────────────────┐
                ▼                                             ▼
     [ Process Guard Check ]                       [ Snapshot Ingestion ]
   (Detect active chrome.exe)                      (Parse Bookmarks/Prefs)
                │                                             │
                ▼                                             ▼
     [ Target Allocation ]                        [ Target Dir Population ]
 (Email match or next Profile N)               (Write Bookmarks, Prefs, Exts)
                │                                             │
                └──────────────────────┬──────────────────────┘
                                       ▼
                       [ Profile Sanitize & Align ]
                     (Preferences: profile.name, GAIA scrub)
                                       ▼
                       [ Local State Registration ]
                     (info_cache full schema + profiles_order)
                                       ▼
                       [ Profile Picker Recognition ]
                     (Visible tile in Chrome UI upon launch)
```

---

## Cross-References

- [`Coding Guidelines Golang`](../02-coding-guidelines/03-golang/00-overview.md)
- [`Error Management Architecture`](../03-error-manage/02-error-architecture/00-overview.md)
- [`Generic CLI Design`](../13-generic-cli/00-overview.md)
- [`Generic Update & Maintenance`](../14-update/00-overview.md)
