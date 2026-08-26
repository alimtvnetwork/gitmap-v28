# WordPress Plugin How-To — Overview

**Version:** 1.0.0  
**Updated:** 2026-04-16  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

Comprehensive blueprint for building **WordPress plugins to the Gold Standard** used across the Riseup Asia stack. Covers PHP foundation and architecture, enum-driven domain modeling, trait composition, structured logging, REST API conventions, settings architecture, frontend/admin UI patterns, deployment, testing, and end-to-end walkthroughs. Any new WordPress plugin MUST follow the contracts and patterns in this folder.

---

## Keywords

`wordpress` · `wp-plugin` · `php` · `enums` · `traits` · `rest-api` · `gold-standard` · `micro-orm` · `settings-api` · `admin-ui` · `gutenberg`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

---

## Architecture (One-Liner)

```
PSR-4 autoload → Enum-driven domain → Trait-composed services → Micro-ORM (root DB)
  → REST API (typed responses) → Admin UI (settings) → Frontend templates → Test/Deploy
```

The defining property: **enums are the source of truth.** Every state, action type, status code, and configuration option is modeled as a backed PHP enum with metadata methods, never as a magic string.

---

## File Inventory

| # | File | Description | Status |
|---|------|-------------|--------|
| 00 | [00-quick-start.md](./00-quick-start.md) | 5-minute quick start: scaffold, install, activate | ✅ Active |
| 01 | [01-foundation-and-architecture.md](./01-foundation-and-architecture.md) | Plugin bootstrap, PSR-4 layout, namespaces, lifecycle hooks | ✅ Active |
| 02 | [02-enums-and-coding-style/](./02-enums-and-coding-style/) | Enum architecture, metadata pattern, action/status enums (subfolder) | ✅ Active |
| 03 | [03-traits-and-composition.md](./03-traits-and-composition.md) | Trait composition over inheritance, shared service mixins | ✅ Active |
| 04 | [04-logging-and-error-handling.md](./04-logging-and-error-handling.md) | Structured logging, error envelopes, apperror parity | ✅ Active |
| 05 | [05-helpers-responses-and-integration.md](./05-helpers-responses-and-integration.md) | Response helpers, third-party integration patterns | ✅ Active |
| 06 | [06-input-validation-patterns.md](./06-input-validation-patterns.md) | Sanitization, validation rules, schema-driven input checks | ✅ Active |
| 07 | [07-reference-implementations.md](./07-reference-implementations.md) | Annotated reference plugin code | ✅ Active |
| 08 | [08-wordpress-integration-patterns.md](./08-wordpress-integration-patterns.md) | WP hooks, filters, action priorities, plugin interop | ✅ Active |
| 09 | [09-testing-patterns.md](./09-testing-patterns.md) | PHPUnit, WP test suite, fixtures, integration tests | ✅ Active |
| 10 | [10-deployment-patterns.md](./10-deployment-patterns.md) | Release packaging, .org repo, private distribution | ✅ Active |
| 11 | [11-frontend-and-template-patterns.md](./11-frontend-and-template-patterns.md) | Template overrides, theme compatibility, asset enqueuing | ✅ Active |
| 12 | [12-design-system.md](./12-design-system.md) | CSS tokens, admin/frontend theming parity | ✅ Active |
| 13 | [13-admin-ui-patterns.md](./13-admin-ui-patterns.md) | Settings pages, list tables, meta boxes, screen options | ✅ Active |
| 14 | [14-rest-api-conventions.md](./14-rest-api-conventions.md) | REST namespaces, route registration, typed responses, auth | ✅ Active |
| 15 | [15-settings-architecture.md](./15-settings-architecture.md) | Options API, settings groups, sections, defaults, migration | ✅ Active |
| 16 | [16-error-handling-extraction.md](./16-error-handling-extraction.md) | Extracting error patterns into reusable trait | ✅ Active |
| 17 | [17-data-file-patterns.md](./17-data-file-patterns.md) | JSON/YAML data files, seed-and-merge, version pinning | ✅ Active |
| 18 | [18-frontend-javascript-patterns.md](./18-frontend-javascript-patterns.md) | Frontend JS, wp-scripts, Gutenberg blocks, modules | ✅ Active |
| 19 | [19-micro-orm-and-root-db.md](./19-micro-orm-and-root-db.md) | Micro-ORM layer, root DB integration, query patterns | ✅ Active |
| 20 | [20-end-to-end-walkthrough.md](./20-end-to-end-walkthrough.md) | Full plugin build walkthrough from scratch | ✅ Active |
| 21 | [21-ping-endpoint.md](./21-ping-endpoint.md) | Health-check REST endpoint pattern | ✅ Active |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding guidelines (cross-language) | `../02-coding-guidelines/00-overview.md` |
| Enum standards (consolidated) | `../17-consolidated-guidelines/04-enum-standards.md` |
| Error management | `../03-error-manage/00-overview.md` |
| Database conventions | `../04-database-conventions/00-overview.md` |
| Design system (core) | `../07-design-system/00-overview.md` |
| Consolidated WP plugin summary | `../17-consolidated-guidelines/20-wp-plugin-conventions.md` |

---

## Placement Rules

```
AI INSTRUCTION:

1. ALL WordPress plugin authoring guidance belongs in this folder.
2. App-specific WP plugin code (e.g., a single product's plugin) goes in 21-app/, not here.
3. Cross-language enum/coding rules live in 02-coding-guidelines/; this folder applies them to PHP/WP.
4. Each file follows the standard {NN}-{kebab-case-name}.md naming convention.
5. Subfolders are allowed when a topic has 3+ files (see 02-enums-and-coding-style/).
6. Add new files to the Feature Inventory above and update 99-consistency-report.md.
```

---

*Overview — updated: 2026-04-16*

---

## Verification

_Auto-generated section — see `spec/18-wp-plugin-how-to/97-acceptance-criteria.md` for the full criteria index._

### AC-WP-000: WordPress plugin conformance: Overview

**Given** Static-analyze the plugin source against the documented enum, trait, and REST conventions.  
**When** Run the verification command shown below.  
**Then** Enums are `enum X: string` with metadata methods; REST routes use the `/wp-json/<plugin>/v1/` namespace; nonces are verified on every mutating request.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
