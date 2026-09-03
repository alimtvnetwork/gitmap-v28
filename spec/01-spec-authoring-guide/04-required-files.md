# Required Files

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

Every spec module (top-level folder under `spec/`) must contain a minimum set of files to pass the health dashboard audit. This document defines those requirements and explains the purpose of each mandatory file.

---

## Mandatory Files (All Modules)

| File | Purpose | Scored? |
|------|---------|---------|
| `01-index.md` | Module index: metadata, file inventory table, cross-references | ✅ 25% of health score |
| `99-consistency-report.md` | Structural health report: file inventory, naming compliance, link validation | ✅ 25% of health score |

---

## Strongly Recommended Files

| File | Purpose | When Required |
|------|---------|---------------|
| `97-acceptance-criteria.md` | Testable acceptance criteria (numbered, specific, measurable) | All implementation-ready modules |
| `98-changelog.md` | Chronological record of changes | Active modules with ongoing updates |
| `error-codes.json` | Error code definitions for the module | Modules with error code allocations |

---

## Health Score Formula

The health dashboard scores each module on four equally weighted criteria:

| Criterion | Weight | Check |
|-----------|--------|-------|
| `01-index.md` present | 25% | File exists in module root |
| `99-consistency-report.md` present | 25% | File exists in module root |
| Lowercase kebab-case naming | 25% | All files and folders comply |
| Unique numeric sequence prefixes | 25% | No duplicate prefixes within a folder |

**100/100 = A+** — All four criteria met.

---

## File Templates

### 01-index.md Template

```markdown

# Module Name

**Version:** 3.2.0
**Updated:** YYYY-MM-DD
**Status:** Draft | Planned | Active | Complete
**AI Confidence:** Low | Medium | High | Production-Ready
**Ambiguity:** None | Low | Medium | High | Critical

---

## Overview

[1-2 paragraph description of what this module covers, its purpose, and who benefits]

---

## Keywords

`keyword-1` · `keyword-2` · `keyword-3` · `category-tag`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Low / Medium / High / Production-Ready |
| Ambiguity | None / Low / Medium / High / Critical |
| Health Score | NN/100 |

```

**Files table template** *(copy into your `01-index.md` — replace placeholder names)*:

```markdown
| # | File | Description |
|---|------|-------------|
| 01 | [01-component.md](./01-component.md) | Description of component spec |
| 02 | [02-other.md](./02-other.md) | Description of other spec |
```

---

```markdown

## Cross-References

| Reference | Location |
|-----------|----------|
| Related Module | `../XX-module-name/01-index.md` |
| Coding Guidelines | `../02-coding-guidelines/01-index.md` |
```

### 97-acceptance-criteria.md Template

```markdown

# Acceptance Criteria — Module Name

**Version:** 3.2.0
**Updated:** YYYY-MM-DD

---

## Overview

[Total count] testable criteria across [N] areas.

---

## Area 1: [Name]

| # | Criterion | Source |
|---|-----------|--------|
| AC-001 | [Specific, testable requirement] | `01-component.md` |
| AC-002 | [Another requirement] | `01-component.md` |

---

## Area 2: [Name]

| # | Criterion | Source |
|---|-----------|--------|
| AC-003 | [Requirement] | `02-other.md` |
```

### 98-changelog.md Template

```markdown

# Changelog — Module Name

**Version:** 3.2.0
**Updated:** YYYY-MM-DD

---

## History

| Date | Version | Changes |
|------|---------|---------|
| YYYY-MM-DD | 1.0.0 | Initial specification created |
```

### 99-consistency-report.md Template

```markdown

# Consistency Report — Module Name

**Version:** 3.2.0
**Last Updated:** YYYY-MM-DD

---

## Module Health

| Criterion | Status |
|-----------|--------|
| `01-index.md` present | ✅ |
| `99-consistency-report.md` present | ✅ |
| Lowercase kebab-case naming | ✅ |
| Unique numeric sequence prefixes | ✅ |

**Health Score:** 100/100 (A+)

---

## File Inventory

| # | File | Status |
|---|------|--------|
| 00 | `01-index.md` | ✅ Present |
| 01 | `01-component.md` | ✅ Present |
| 97 | `97-acceptance-criteria.md` | ✅ Present |
| 99 | `99-consistency-report.md` | ✅ Present |

**Total:** N files

---

## Cross-Reference Validation

All internal links verified valid. ✅

---

## Validation History

| Date | Version | Action |
|------|---------|--------|
| YYYY-MM-DD | 1.0.0 | Initial consistency report created |
```

---

## Subfolder Required Files

Subfolders within a module (e.g., `01-backend/`, `02-frontend/`) also require:

| File | Required? |
|------|-----------|
| `01-index.md` | ✅ Always |
| `99-consistency-report.md` | ✅ For subfolders with 3+ files |
| `97-acceptance-criteria.md` | Optional (can consolidate at parent level) |

---

## Version Bumping Rules

| Change Type | Version Bump | Example |
|-------------|-------------|---------|
| Add/remove files | Major (X.0.0) | Adding `08-new-feature.md` |
| Content update to existing file | Minor (x.Y.0) | Updating acceptance criteria |
| Typo/formatting fix | Patch (x.y.Z) | Fixing a broken link |
| Consistency report refresh | Major (X.0.0) | Resynchronizing file inventory |
