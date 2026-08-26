# Spec Authoring Guide — Acceptance Criteria

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

---

## Overview

18 testable criteria across 4 areas covering spec structure, naming, content, and tooling.

---

## AC-01: Folder Structure & Required Files

| # | Criterion | Source |
|---|-----------|--------|
| AC-001 | Every spec module has `00-overview.md` at root | `03-required-files.md` |
| AC-002 | Every spec module has `99-consistency-report.md` at root | `03-required-files.md` |
| AC-003 | CLI modules follow 3-folder pattern (`01-backend/`, `02-frontend/`, `03-deploy/`) | `04-cli-module-template.md` |
| AC-004 | Subfolders with 3+ files include their own `00-overview.md` | `03-required-files.md` |

---

## AC-02: Naming Conventions

| # | Criterion | Source |
|---|-----------|--------|
| AC-005 | All files use lowercase kebab-case naming | `02-naming-conventions.md` |
| AC-006 | All folders use lowercase kebab-case naming | `02-naming-conventions.md` |
| AC-007 | All spec files have unique numeric sequence prefixes within their folder | `02-naming-conventions.md` |
| AC-008 | Reserved prefixes (00, 97, 98, 99) used only for their designated purposes | `02-naming-conventions.md` |

---

## AC-03: Overview Content Standards

| # | Criterion | Source |
|---|-----------|--------|
| AC-009 | Every `00-overview.md` includes Version and Updated metadata | `00-overview.md` |
| AC-010 | Every `00-overview.md` includes AI Confidence score | `00-overview.md` |
| AC-011 | Every `00-overview.md` includes Ambiguity score | `00-overview.md` |
| AC-012 | Every `00-overview.md` includes Keywords section | `00-overview.md` |
| AC-013 | Every `00-overview.md` includes Scoring table | `00-overview.md` |
| AC-014 | Every `00-overview.md` includes numbered file inventory table | `00-overview.md` |
| AC-015 | Every `00-overview.md` includes Cross-References table | `00-overview.md` |

---

## AC-04: Cross-References & Validation

| # | Criterion | Source |
|---|-----------|--------|
| AC-016 | All cross-references use relative paths (never root-relative or absolute) | `08-cross-references.md` |
| AC-017 | All linked files include `.md` extension | `08-cross-references.md` |
| AC-018 | Zero broken links reported by dashboard scanner | `08-cross-references.md` |

---

## Cross-References

- [Overview](./00-overview.md)
- [Required Files](./03-required-files.md)
- [Naming Conventions](./02-naming-conventions.md)
- [Cross-References Guide](./08-cross-references.md)
