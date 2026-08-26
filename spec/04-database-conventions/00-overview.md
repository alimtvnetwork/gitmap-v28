# Database Conventions

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`database` · `sqlite` · `split-db` · `orm` · `pascalcase` · `primary-key` · `foreign-key` · `views` · `testing` · `naming` · `schema-design`

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

Comprehensive database design and implementation conventions covering naming, schema design, key sizing, ORM usage, view patterns, relationship modeling, and testing strategies. This is the **single source of truth** for how databases are designed and used across all languages.

> 🔴 **MANDATORY — AI Agents Must Commit Database Rules to Memory**
>
> After reading this document, you **MUST** retain and enforce these database conventions in every schema, migration, model, and query you generate:
>
> 1. **Singular table names** — `User`, `Project`, `Transaction` — never plural (`Users`, `Projects`)
> 2. **PascalCase everything** — tables, columns, indexes, views, JSON response fields
> 3. **PK = `{TableName}Id`** — e.g., `UserId`, `ProjectId` — always `INTEGER PRIMARY KEY AUTOINCREMENT`, never UUID
> 4. **FK = exact PK name** — if `User` has PK `UserId`, any child table references it as `UserId` (not `user_id`, not `fk_user`)
> 5. **Booleans** — `Is`/`Has` prefix, positive-only names (`IsActive`, never `IsDisabled`)
>
> If any database requirement is ambiguous or conflicts with these rules, **ask a clarifying question** instead of guessing. Wrong schema decisions are expensive to fix.

---

## Golden Rules

> 1. **Singular table names** — `User`, `Project`, `Transaction` (not `Users`, `Projects`)
> 2. **PascalCase everything** — tables, columns, indexes, views
> 3. **PK = `{TableName}Id`** — `INTEGER PRIMARY KEY AUTOINCREMENT`, never UUID
> 4. **FK = exact PK name** — `UserId` in both `User` and `UserProfile` tables
> 5. **SQLite first** (Split DB pattern) — MySQL as fallback
> 6. **Always use ORMs** — never write raw SQL in business logic
> 7. **Smallest possible key type** — `INTEGER` over `BIGINT`, never UUID unless required
> 8. **Repeated values → separate table** — normalize with foreign key relationships
> 9. **Views for joins** — define DB views instead of on-the-fly joins in code
> 10. **Test with in-memory DB** — unit test schemas, integration test with real queries

---

## Document Index

| # | File | Description |
|---|------|-------------|
| 01 | [01-naming-conventions.md](./01-naming-conventions.md) | PascalCase rules for tables, columns, indexes — singular table names |
| 02 | [02-schema-design.md](./02-schema-design.md) | Key sizing, primary keys, foreign keys, normalization rules |
| 03 | [03-orm-and-views.md](./03-orm-and-views.md) | ORM-first approach, view patterns, no raw SQL in business logic |
| 04 | [04-testing-strategy.md](./04-testing-strategy.md) | Unit tests for schemas, integration tests with in-memory DB |
| 05 | [05-relationship-diagrams.md](./05-relationship-diagrams.md) | Visual relationship patterns and AI-readable schema diagrams |
| 06 | [06-rest-api-format.md](./06-rest-api-format.md) | PascalCase REST API response format, full CRUD sample, response envelope |
| 07 | [07-split-db-pattern.md](./07-split-db-pattern.md) | Split DB pattern — one SQLite file per bounded context |
| 99 | [99-consistency-report.md](./99-consistency-report.md) | Module health and validation |

---

## Quick Reference

| Topic | Rule |
|-------|------|
| Table names | **Singular** PascalCase: `User`, `Transaction`, `Project` |
| Column names | PascalCase: `PluginSlug`, `CreatedAt` |
| Primary key format | `{TableName}Id` (e.g., `UserId`, `ProjectId`) |
| Primary key type | `INTEGER PRIMARY KEY AUTOINCREMENT` |
| Foreign key format | Exact PK name from referenced table (e.g., `UserId` in child table) |
| UUID/GUID | ❌ Avoid unless explicitly required |
| Booleans | `Is`/`Has` prefix, positive-only (`IsActive`, not `IsDisabled`) |
| Repeated values | Normalize into separate table with FK |
| Joins in code | ❌ Use DB views instead |
| Raw SQL in business logic | ❌ Use ORM |
| Default database | SQLite (Split DB pattern) |
| Fallback database | MySQL |
| Schema testing | Unit test + integration test with in-memory DB |

---

## Database Engine Priority

| Priority | Engine | When to Use |
|----------|--------|-------------|
| 1st | **SQLite** (Split DB) | Default for all projects — embedded, zero-config, portable |
| 2nd | **MySQL** | When concurrent write-heavy loads or multi-server access is needed |

> The **Split DB** pattern uses multiple small SQLite databases per domain concern rather than one monolithic database. See [07-split-db-pattern.md](./07-split-db-pattern.md) for the full specification.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Spec Root | [../00-overview.md](../00-overview.md) |
| Coding Guidelines | [../02-coding-guidelines/00-overview.md](../02-coding-guidelines/00-overview.md) |
| Cross-Language DB Naming | [../02-coding-guidelines/01-cross-language/07-database-naming.md](../02-coding-guidelines/01-cross-language/07-database-naming.md) |
| Split DB Architecture | [../05-split-db-architecture/00-overview.md](../05-split-db-architecture/00-overview.md) |
| Consolidated DB Conventions | [../17-consolidated-guidelines/18-database-conventions.md](../17-consolidated-guidelines/18-database-conventions.md) |

---

*Single source of truth for database design and conventions across all languages.*

---

## Verification

_Auto-generated section — see `spec/04-database-conventions/97-acceptance-criteria.md` for the full criteria index._

### AC-DB-000: Database convention conformance: Overview

**Given** Run the SQL schema linter against your DDL files.  
**When** Run the verification command shown below.  
**Then** Every table is PascalCase singular; PK is `<TableName>Id INTEGER PRIMARY KEY AUTOINCREMENT`; columns are `NOT NULL` unless waived; no `createdAt`, `created_at`, `UUID` tokens.

**Verification command:**

```bash
python3 linter-scripts/check-forbidden-strings.py
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
