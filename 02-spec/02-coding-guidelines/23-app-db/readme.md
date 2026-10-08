# App DB (AI Execution Prompt)

> **/goal** Define and govern application-specific database schemas, table conventions, migration patterns, and query architectures.
> **/learn** Master Split-DB boundaries, PascalCase table names, positive boolean column naming, and ORM entity mapping standards.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce PascalCase table naming and affirmative boolean column definitions across all app databases.
- [ ] `/learn` Maintain strict schema isolation between operational app tables and system/pipeline databases.
- [ ] `/goal` Standardize migration scripts and query patterns to eliminate raw unparameterized SQL operations.
- [ ] `/learn` Verify zero cross-database joins and enforce strict relative paths in all database design specifications.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Draft
**Ambiguity:** None

---

## Purpose

App-specific database (App DB) specifications under coding guidelines. Covers data model decisions, table designs, and query patterns unique to this application — whether the project is a web app, Chrome extension, CLI, plugin, or otherwise.

---

## Contents

_No content yet. Add database design documents as numbered files within this folder._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding Guidelines Overview | [../readme.md](../readme.md) |
| Consolidated Database Conventions | [../../17-consolidated-guidelines/21-database-conventions.md](../../17-consolidated-guidelines/21-database-conventions.md) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-APP-003: Application Database Standards & Split SQLite Architecture

**Given** Application database schemas, migration patterns, and Split SQLite specifications in `02-spec/02-coding-guidelines/23-app-db/`.
**When** Audited against this reference specification and coding guidelines.
**Then** PascalCase tables, positive boolean fields, and Split-DB boundaries are enforced with zero violations and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/23-app-db --check-only
```
**Expected:** exit 0. Zero violations.
