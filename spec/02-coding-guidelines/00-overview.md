# Coding Guidelines

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

Consolidated coding standards and conventions organized by category. This folder is the **single canonical location** for all language-specific and cross-language coding guidelines, including file naming, security policies, database design conventions, PowerShell integration, and research.

---

> 🔴 **MANDATORY — AI Agents Must Commit These Rules to Memory**
>
> After reading this coding guideline, you **MUST** internalize the following rules and apply them to **every single code change** without exception. Do not proceed with implementation until you have understood and committed these to your working memory:
>
> 1. **No Generated Code or Artifacts (Never Commit)** - Never commit generated code (e.g., ORM models, gRPC clients), test results, test reports, or compiled binaries. They belong in build artifacts or CI, never in source control. See [No Generated Code](./01-cross-language/29-no-generated-artifacts.md).
> 2. **Error Management is the #1 priority** - Error handling from [03-error-manage/](../03-error-manage/00-overview.md) must be implemented from the **very first line of code**. Never write business logic without proper error handling wrapping it. This is non-negotiable.
> 3. **Boolean naming** - All booleans use `is`/`has`/`should` prefixes and are **positively named only** (`IsActive`, never `IsDisabled`). Extract multi-part conditions into named variables.
> 4. **if/else and nesting** - Zero nesting. Use early returns and guard clauses. No nested `if` blocks.
> 5. **Database conventions** - Singular table names (`User` not `Users`), PascalCase everywhere, `{TableName}Id` as `INTEGER PRIMARY KEY AUTOINCREMENT`, FK uses the exact PK name. See [Database Conventions](../04-database-conventions/00-overview.md).
> 6. **Never hallucinate** - If a requirement is unclear or missing, **ask a clarifying question** instead of guessing. Wrong assumptions cause rewrites.
> 7. **Function metrics** - Functions: 8-15 lines. Files: < 300 lines. React components: < 100 lines.
>
> These rules are **CODE RED** — violations are treated as bugs and must be fixed before merge.

---

## ⚠️ Naming Convention Policy — AI Critical Instruction

```
STOP — EVERY AI AGENT MUST READ THIS SECTION BEFORE GENERATING CODE.

This project uses a HYBRID naming convention strategy. Most languages follow the
project-wide PascalCase mandate. Rust is an INTENTIONAL EXCEPTION.

╔══════════════════════════════════════════════════════════════════════╗
║  LANGUAGE         IDENTIFIER CONVENTION         DATABASE    ENUM   ║
║                                                 COLUMNS    VALUES  ║
╠══════════════════════════════════════════════════════════════════════╣
║  Go               PascalCase (exported)         PascalCase PascalCase ║
║  TypeScript        PascalCase (keys/values)      PascalCase PascalCase ║
║  PHP               PascalCase (keys/values)      PascalCase PascalCase ║
║  C#                PascalCase (methods/props)     PascalCase PascalCase ║
║  Rust              snake_case (community std)     PascalCase PascalCase ║
╚══════════════════════════════════════════════════════════════════════╝

KEY INSIGHT:
- Go, TypeScript, PHP, C# → PascalCase is the DEFAULT for identifiers, keys, JSON.
- Rust → snake_case is the DEFAULT per community conventions (RFC 430).
- ALL LANGUAGES (including Rust) → PascalCase for DATABASE and ENUM STRING VALUES.

WHY RUST IS DIFFERENT:
Rust's compiler enforces snake_case for functions/variables via lint warnings.
Fighting the compiler and ecosystem to force PascalCase is impractical.
Instead, Rust follows its community standard EXCEPT at cross-system boundaries
(database and enum serialization) where PascalCase is mandatory for interop.

WHEN GENERATING RUST CODE:
- Functions → snake_case:        fn get_active_window()
- Variables → snake_case:        let session_id = 42;
- Constants → SCREAMING_SNAKE:   const MAX_RETRIES: u32 = 3;
- Types/Enums → PascalCase:      struct BrowserActivity (Rust standard)
- DB columns → PascalCase:       "SELECT SessionId FROM Sessions"
- Enum strings → PascalCase:     "TabChange" (serde default for PascalCase variants)
- JSON keys → PascalCase:        #[serde(rename_all = "PascalCase")]

WHEN GENERATING Go/TS/PHP/C# CODE:
- Follow the PascalCase mandate from 01-cross-language/11-key-naming-pascalcase.md
- Database columns → PascalCase (same as Rust)
- Enum values → PascalCase (same as Rust)

See 05-rust/01-naming-conventions.md for the complete Rust naming reference.
See 01-cross-language/11-key-naming-pascalcase.md for the general PascalCase mandate.
```

---

## ⚠️ Numbering Convention — AI Instruction

```
IMPORTANT — AI INSTRUCTION:

1. Folders 01–20 are RESERVED for core fundamentals only.
   - Language standards, cross-cutting principles, naming, security, database, integrations, research.
   - No app-specific content may appear in this range.

2. Folders 21+ are for APP-SPECIFIC content.
   - 21-app: Application feature specs, workflows, architecture decisions.
   - 22-app-issues: App bug analysis, root cause analysis, fix documentation.

3. Decision guide for placement:
   - Reusable, foundational, or principle-driven → 01–20 (core fundamentals)
   - Exploratory, comparative, or evaluative → 10-research
   - App feature or workflow definition → 21-app
   - App bug/failure/root cause analysis → 22-app-issues

4. New core fundamental folders use the next available number within 01–20.
5. New app folders use the next available number after 22.
```

---

## Keywords

`coding-standards` · `cross-language` · `typescript` · `golang` · `php` · `rust` · `csharp` · `naming-conventions` · `boolean-patterns` · `dry` · `strict-typing` · `file-naming` · `folder-naming` · `security` · `dependency-pinning` · `database` · `orm` · `schema-design` · `slug-conventions` · `powershell` · `research`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

---

## Categories

### Core Fundamentals (01–20)

#### Language & Cross-Language Standards

| # | Category | Description | Files |
|---|----------|-------------|-------|
| 01 | [Cross-Language](./01-cross-language/00-overview.md) | Language-agnostic rules: DRY, naming, booleans, typing, complexity, lazy eval, regex, mutation, null safety, nesting, slugs | 29 |
| 02 | [TypeScript](./02-typescript/00-overview.md) | TypeScript enum patterns, type safety, promise/await patterns | 13 |
| 03 | [Golang](./03-golang/00-overview.md) | Go coding standards, enum specification, boolean rules, defer, internals, severity | 16 |
| 04 | [PHP](./04-php/00-overview.md) | PHP coding standards, enums, forbidden patterns, naming, spacing/imports, ResponseKeyType | 12 |
| 05 | [Rust](./05-rust/00-overview.md) | Rust standards: naming, error handling, async, memory safety, FFI | 10 |
| 06 | [AI Optimization](./06-ai-optimization/00-overview.md) | Anti-hallucination rules, AI quick-reference checklist, common AI mistakes, enum naming reference | 8 |
| 07 | [C#](./07-csharp/00-overview.md) | C# standards: naming, method design, error handling, type safety | 5 |

#### Infrastructure & Convention Standards

| # | Category | Description | Files |
|---|----------|-------------|-------|
| 08 | [File & Folder Naming](./08-file-folder-naming/00-overview.md) | Per-language file and folder naming conventions (PHP/WordPress, Go, TS/JS, Rust, C#) | 7 |
| 09 | [PowerShell Integration](./09-powershell-integration/00-overview.md) | PowerShell scripting conventions and cross-platform automation | 0 |
| 10 | [Research](./10-research/00-overview.md) | Comparative studies, technology evaluations, exploratory technical notes | 0 |
| 11 | [Security](./11-security/00-overview.md) | Security policies, dependency pinning (Axios), vulnerability tracking | 6 |
| 12–20 | _Reserved_ | Available for future core fundamental topics | — |

### App-Specific (21+)

| # | Category | Description | Files |
|---|----------|-------------|-------|
| 21 | [App](./21-app/00-overview.md) | App-specific specs: features, workflows, architecture decisions | 0 |
| 22 | [App Issues](./22-app-issues/00-overview.md) | App bug analysis, root cause analysis, fix documentation | 0 |
| 23 | [App DB](./23-app-db/00-overview.md) | App-specific data model, table designs, migration strategies | 0 |
| 24 | [App UI — Design System](./24-app-ui-design-system/00-overview.md) | App-specific design system, theming, component patterns, layout | 0 |

---

## Consolidation Status

✅ **Complete.** All unique content from 5 legacy sources has been merged into this canonical location.

---

## Migration History

| Date | Change |
|------|--------|
| 2026-04-16 | **Flattened structure** — removed nested `03-coding-guidelines-spec/` folder, moved all subfolders to root level |
| 2026-04-09 | Restructured: 09→PowerShell, 10→Research, 09-security→11, 10-database→12, added 21-app, 22-app-issues |
| 2026-04-02 | Added `10-database-conventions/` (8 files: schema design, ORM, views, testing, REST API format) |
| 2026-04-02 | Added `09-security/` and moved Axios version control from `spec/01-app/` |
| 2026-04-02 | Added `08-file-folder-naming/` (per-language conventions) |
| 2026-04-02 | Added `28-slug-conventions.md` to cross-language |
| 2026-03-31 | Consolidated 5 guideline sources into this canonical location |

---

## Document Inventory

| File | Description |
|------|-------------|
| [consolidated-review-guide.md](./consolidated-review-guide.md) | Full code review guide with examples (all languages) |
| [consolidated-review-guide-condensed.md](./consolidated-review-guide-condensed.md) | One-liner bullet-point checklist for quick scanning |
| 97-acceptance-criteria.md | Testable criteria across guideline categories |
| 99-consistency-report.md | Module health and file inventory |

---

## Cross-References

- [Spec Authoring Guide](../01-spec-authoring-guide/00-overview.md)
- [Error Management](../03-error-manage/00-overview.md) — **Highest priority spec. Read first.**
- [Database Conventions](../04-database-conventions/00-overview.md) — Naming, schema, key design

---

*Coding guidelines v1.1.0 — 2026-04-16*

---

## Verification

_Auto-generated section — see `spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-000: Coding guideline conformance: Overview

**Given** Run the cross-language coding-guidelines validator against `src/` and language-specific source roots.  
**When** Run the verification command shown below.  
**Then** Zero CODE-RED violations are reported (functions ≤ 15 lines, files ≤ 300 lines, no nested ifs, max 2 boolean operands).

**Verification command:**

```bash
go run linter-scripts/validate-guidelines.go --path spec --max-lines 15 && python3 linter-scripts/validate-guidelines.py spec
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
