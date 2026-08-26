# Error Code Registry

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None  
**Scope:** Cross-project utility

---

## Keywords

`error-codes` · `registry` · `cross-project` · `debugging` · `integer-codes` · `prefixed-codes` · `collision-prevention`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Centralized error code registry ensuring no collisions between projects, consistent structure for debugging, machine-parseable error codes, and human-readable messages.

---

## Error Code Formats

| Format | Used By | Example | Pattern |
|--------|---------|---------|---------|
| `XX-NNN-NN` | General specs, PHP plugins | `SM-400-01` | `^[A-Z]{2,4}-[0-9]{3}-[0-9]{2}$` |
| Integer | Go CLI tools | `7001`, `9301` | `^[0-9]{4,5}$` |

---

## Registered Ranges (Quick Reference)

| Prefix | Project | Range |
|--------|---------|-------|
| `GEN` | General/Shared | 0-999 |
| `SM` | Spec Management | 2000-2999 |
| `GS` | GSearch (all modules) | 7000-7919 |
| `BR` | BRun | 7100-7599 |
| `NF` | Nexus Flow | 8000-8099 |
| `AB` | AI Bridge (all modules) | 9000-9999, 19000-19049 |
| `PS` | PowerShell | 9500-9540 |
| `WPB` | WP Plugin Builder | 10000-10499 |
| `SRC` | Spec Reverse | 11000-11999 |
| `WSP` | WP SEO Publish | 12000-12599 |
| `WPP` | WP Plugin Publish | 13000-13499 |
| `AIT` | AI Transcribe | 14000-14499 |
| `EQM` | Exam Manager | 14500-14999 |
| `LM` | Link Manager | 15000-15999 |
| `SM-CG` | SM Code Generation | 16000-16799 |
| `SM-PE` | SM Project Editor | 17000-17999 |
| `SM-GS` | SM GSearch (Ecosystem Remap) | 18000-18249 |
| `AB-LR` | AB Lovable Reasoning | 19000-19049 |
| `AB-TR` | AB Non-Vector RAG | 20000-20999 |

> See `01-registry.md` for the complete master list with sub-ranges, collision resolution log, and range allocation map.

---

## Document Inventory

| # | File | Purpose |
|---|------|---------|
| 01 | [01-registry.md](./01-registry.md) | Master list of all registered codes |
| 02 | [02-integration-guide.md](./02-integration-guide.md) | How to add codes to your project |
| 03 | [03-collision-resolution-summary.md](./03-collision-resolution-summary.md) | Consolidated before/after table of all 13 resolutions |
| 04 | [04-error-code-utilization-report.md](./04-error-code-utilization-report.md) | Auto-generated range utilization report |
| 05 | [05-overlap-validator.md](./05-overlap-validator.md) | Automated overlap detection tool spec |
| 06 | [error-codes-master.json](./error-codes-master.json) | Machine-readable master index |
| — | 99-consistency-report.md | — |

### Subfolders

| # | Folder | Description | Files |
|---|--------|-------------|-------|
| 07 | [07-schemas/](./07-schemas/00-overview.md) | JSON schemas for validation | 2 |
| 08 | [08-linter-scripts/](./08-linter-scripts/00-overview.md) | Automation scripts | 4 |
| 09 | [09-templates/](./09-templates/00-overview.md) | Templates for project error docs | 1 |

| — | 99-consistency-report.md | — |
---

## Quick Reference

**To register new codes:**
1. Check the Range Allocation Map in `01-registry.md`
2. Claim a project prefix
3. Add your codes following category offsets
4. Document in your project's spec folder
5. Update `01-registry.md` in the same commit

---

## Cross-References

- [Parent Overview](../00-overview.md) — Error Management root
- [Error Resolution](../01-error-resolution/00-overview.md) — Debugging and verification
- [Error Architecture](../02-error-architecture/00-overview.md) — Cross-stack error handling
