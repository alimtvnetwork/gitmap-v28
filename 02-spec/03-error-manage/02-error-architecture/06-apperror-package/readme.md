# AppError Package

> [!IMPORTANT]
> **Modernized Standard: `pkg/appfault` (`*appfault.AppError`) & Concrete Types in `types.go`:**
> The structured error package has been modernized from legacy `apperror` to `pkg/appfault` (`*appfault.AppError` / `appfault.Fault`).
> All functions returning results MUST use concrete named types defined in `types.go` (e.g. `type ConfigResult = result.Wrap[*Config]`) rather than leaking raw generic wrappers.

. **CRITICAL AI INSTRUCTION:** This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.3.0
**Status:** Active
**Updated:** 2026-09-22
**AI Confidence:** High
**Ambiguity:** None

---

## Keywords

`error`, `resolution`, `apperror`, `package`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `readme.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

## Purpose

Application error package specification.

---

## Document Inventory

| File | Purpose |
|------|---------|
| 01-apperror-reference.md | AppError struct, Result types, usage patterns |
| 01-apperror-reference/ | Subfolder with split reference docs (incl. 05-apperrtype-enums.md) |
| 99-consistency-report.md | Structural health |

---

## Cross-References

_See parent folder's `readme.md` for broader context._
