# Seedable Config Architecture + Changelog Versioning (also known as CW Config)

> **Version:** 3.0.0  
> **Created:** 2026-02-01  
> **Updated:** 2026-04-03  
> **Status:** Active  
> **AI Confidence:** Production-Ready  
> **Ambiguity:** Low  
> **Purpose:** Reusable pattern for version-controlled configuration with automatic changelog updates and initial seeding

---

## Keywords

`configuration` · `seeding` · `changelog` · `versioning` · `sqlite` · `json-schema` · `semver` · `merge-strategy`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | Low |
| Health Score | 100/100 (A+) |

---

## Summary

The **Seedable Config Architecture + Changelog Versioning** (commonly referred to as **CW Config**) defines a pattern for managing application configuration where:

1. **First-run seeding** populates SQLite DB from `config.seed.json`
2. **Every config change updates the version**
3. **Every version change logs to changelog.md**
4. **Subsequent runs respect version** to avoid duplicate seeds

This ensures configuration is always traceable, auditable, and version-aware.

---

## Document Inventory

| # | File | Description |
|---|------|-------------|
| 00 | `00-overview.md` | This file — master index |
| 01 | `01-fundamentals.md` | Core concepts, configuration files, version flow, merge strategies |
| 02 | `02-features/00-overview.md` | Feature index |
| 02.01 | `02-features/01-rag-chunk-settings.md` | RAG chunk size and overlap configuration |
| 02.02 | `02-features/02-rag-validation-helpers.md` | Go validation patterns for RAG config |
| 02.03 | `02-features/03-rag-validation-tests.md` | Unit test specifications for validators |
| 02.04 | `02-features/04-rag-test-coverage-matrix.md` | Test coverage matrix for RAG validation |
| 02.05 | `02-features/05-validation-data-seeding.md` | CW Config → Root DB seeding pattern |
| 02.06 | `02-features/06-update-check-keys.md` | `Update.*` and `Storage.Backend` keys for update-check subsystem |
| 03 | `03-issues/00-overview.md` | Issues tracker |
| 97 | `97-acceptance-criteria.md` | Acceptance criteria |
| 97b | `97-changelog.md` | Changelog |
| 98 | `98-acceptance-criteria.md` | Extended acceptance criteria |
| 99 | `99-consistency-report.md` | Consistency report |

---

## Folder Structure

```
06-seedable-config-architecture/
├── 00-overview.md                    ← This file
├── 01-fundamentals.md                ← Core concepts & architecture
├── 02-features/
│   ├── 00-overview.md                ← Feature index
│   ├── 01-rag-chunk-settings.md
│   ├── 02-rag-validation-helpers.md
│   ├── 03-rag-validation-tests.md
│   ├── 04-rag-test-coverage-matrix.md
│   └── 05-validation-data-seeding.md
├── 03-issues/
│   └── 00-overview.md                ← Issues tracker
├── 97-acceptance-criteria.md
├── 97-changelog.md
├── 98-acceptance-criteria.md
└── 99-consistency-report.md
```

---

## Cross-References

| Reference | Description |
|-----------|-------------|
| [Split DB Architecture](../05-split-db-architecture/00-overview.md) | Database organization patterns |
| [App Project Template](../01-spec-authoring-guide/05-app-project-template.md) | Template this spec follows |

---

*Overview — updated: 2026-04-03*

---

## Verification

_Auto-generated section — see `spec/06-seedable-config-architecture/97-acceptance-criteria.md` for the full criteria index._

### AC-CFG-000: Seedable-config conformance: Overview

**Given** Diff the running config tree against `config.seed.json` after a SemVer-aware GORM merge.  
**When** Run the verification command shown below.  
**Then** Merged keys preserve user overrides; new seed keys are added; removed seed keys are pruned; merge is idempotent on a second pass.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
