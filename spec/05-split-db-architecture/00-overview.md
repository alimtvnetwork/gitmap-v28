# Split Database Architecture

> **Version:** 3.0.0  
> **Created:** 2026-02-01  
> **Updated:** 2026-04-03  
> **Status:** Active  
> **AI Confidence:** Production-Ready  
> **Ambiguity:** Low  
> **Purpose:** Reusable pattern for hierarchical SQLite database organization across all projects

---

## Keywords

`sqlite` · `split-database` · `hierarchical-storage` · `connection-pooling` · `wal-mode` · `backup` · `multi-project`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | Low |
| Health Score | 100/100 (A+) |

---

## CRITICAL: Naming Convention

**All field names use PascalCase. No underscores allowed.**

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `session_id` | `SessionId` |
| `created_at` | `CreatedAt` |
| `message_count` | `MessageCount` |

---

## Summary

The **Split DB Architecture** defines a pattern for organizing SQLite databases into a **multi-layer hierarchical structure** where a **Root DB** manages metadata about child databases, and item-specific databases are created dynamically as needed. This pattern enables efficient data isolation, improved performance, logical organization, and easy import/export via zip files.

---

## Document Inventory

| # | File | Description |
|---|------|-------------|
| 00 | `00-overview.md` | This file — master index |
| 01 | `01-fundamentals.md` | Core concepts, terminology, hierarchical structure, implementation patterns |
| 02 | `02-features/00-overview.md` | Feature index |
| 02.01 | `02-features/01-cli-examples.md` | Concrete examples for AI Bridge, GSearch, BRun, Nexus Flow |
| 02.02 | `02-features/02-reset-api-standard.md` | 2-step reset API standard (5-min TTL) |
| 02.03 | `02-features/03-database-flow-diagrams.md` | Visual architecture diagrams |
| 02.04 | `02-features/04-rbac-casbin.md` | Role-Based Access Control with Casbin |
| 02.05 | `02-features/05-user-scoped-isolation.md` | User-scoped database isolation patterns |
| 03 | `03-issues/00-overview.md` | Issues tracker |
| 97 | `97-acceptance-criteria.md` | Acceptance criteria |
| 97b | `97-changelog.md` | Changelog |
| 98 | `98-acceptance-criteria.md` | Extended acceptance criteria |
| 99 | `99-consistency-report.md` | Consistency report |

---

## Folder Structure

```
05-split-db-architecture/
├── 00-overview.md                    ← This file
├── 01-fundamentals.md                ← Core concepts & architecture
├── 02-features/
│   ├── 00-overview.md                ← Feature index
│   ├── 01-cli-examples.md
│   ├── 02-reset-api-standard.md
│   ├── 03-database-flow-diagrams.md
│   ├── 04-rbac-casbin.md
│   └── 05-user-scoped-isolation.md
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
| [Seedable Config](../06-seedable-config-architecture/00-overview.md) | Configuration seeding patterns |
| [App Project Template](../01-spec-authoring-guide/05-app-project-template.md) | Template this spec follows |

---

*Overview — updated: 2026-04-03*

---

## Verification

_Auto-generated section — see `spec/05-split-db-architecture/97-acceptance-criteria.md` for the full criteria index._

### AC-SDB-000: Split-DB architecture conformance: Overview

**Given** Inspect Root/App/Session DB lifecycle wiring and Casbin RBAC enforcement points.  
**When** Run the verification command shown below.  
**Then** Each tier opens its own SQLite handle (WAL mode), policy reload happens on Casbin policy change, and user-scope isolation is enforced by row filters.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
