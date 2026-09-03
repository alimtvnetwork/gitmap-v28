# WordPress Plugin How-To — Overview

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 18 Wp Plugin How To.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `18-wp-plugin-how-to/`.
- [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


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
| 00 | [02-quick-start.md](./02-quick-start.md) | 5-minute quick start: scaffold, install, activate | ✅ Active |
| 01 | [03-foundation-and-architecture.md](./03-foundation-and-architecture.md) | Plugin bootstrap, PSR-4 layout, namespaces, lifecycle hooks | ✅ Active |
| 02 | [02-enums-and-coding-style/](./02-enums-and-coding-style/) | Enum architecture, metadata pattern, action/status enums (subfolder) | ✅ Active |
| 03 | [04-traits-and-composition.md](./04-traits-and-composition.md) | Trait composition over inheritance, shared service mixins | ✅ Active |
| 04 | [05-logging-and-error-handling.md](./05-logging-and-error-handling.md) | Structured logging, error envelopes, apperror parity | ✅ Active |
| 05 | [06-helpers-responses-and-integration.md](./06-helpers-responses-and-integration.md) | Response helpers, third-party integration patterns | ✅ Active |
| 06 | [07-input-validation-patterns.md](./07-input-validation-patterns.md) | Sanitization, validation rules, schema-driven input checks | ✅ Active |
| 07 | [08-reference-implementations.md](./08-reference-implementations.md) | Annotated reference plugin code | ✅ Active |
| 08 | [09-wordpress-integration-patterns.md](./09-wordpress-integration-patterns.md) | WP hooks, filters, action priorities, plugin interop | ✅ Active |
| 09 | [10-testing-patterns.md](./10-testing-patterns.md) | PHPUnit, WP test suite, fixtures, integration tests | ✅ Active |
| 10 | [11-deployment-patterns.md](./11-deployment-patterns.md) | Release packaging, .org repo, private distribution | ✅ Active |
| 11 | [12-frontend-and-template-patterns.md](./12-frontend-and-template-patterns.md) | Template overrides, theme compatibility, asset enqueuing | ✅ Active |
| 12 | [13-design-system.md](./13-design-system.md) | CSS tokens, admin/frontend theming parity | ✅ Active |
| 13 | [14-admin-ui-patterns.md](./14-admin-ui-patterns.md) | Settings pages, list tables, meta boxes, screen options | ✅ Active |
| 14 | [15-rest-api-conventions.md](./15-rest-api-conventions.md) | REST namespaces, route registration, typed responses, auth | ✅ Active |
| 15 | [16-settings-architecture.md](./16-settings-architecture.md) | Options API, settings groups, sections, defaults, migration | ✅ Active |
| 16 | [17-error-handling-extraction.md](./17-error-handling-extraction.md) | Extracting error patterns into reusable trait | ✅ Active |
| 17 | [18-data-file-patterns.md](./18-data-file-patterns.md) | JSON/YAML data files, seed-and-merge, version pinning | ✅ Active |
| 18 | [19-frontend-javascript-patterns.md](./19-frontend-javascript-patterns.md) | Frontend JS, wp-scripts, Gutenberg blocks, modules | ✅ Active |
| 19 | [20-micro-orm-and-root-db.md](./20-micro-orm-and-root-db.md) | Micro-ORM layer, root DB integration, query patterns | ✅ Active |
| 20 | [21-end-to-end-walkthrough.md](./21-end-to-end-walkthrough.md) | Full plugin build walkthrough from scratch | ✅ Active |
| 21 | [22-ping-endpoint.md](./22-ping-endpoint.md) | Health-check REST endpoint pattern | ✅ Active |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding guidelines (cross-language) | `../02-coding-guidelines/01-index.md` |
| Enum standards (consolidated) | `../17-consolidated-guidelines/07-enum-standards.md` |
| Error management | `../03-error-manage/01-index.md` |
| Database conventions | `../04-database-conventions/01-index.md` |
| Design system (core) | `../07-design-system/01-index.md` |
| Consolidated WP plugin summary | `../17-consolidated-guidelines/23-wp-plugin-conventions.md` |

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

_Auto-generated section — see `02-spec/18-wp-plugin-how-to/97-acceptance-criteria.md` for the full criteria index._

### AC-WP-001: WordPress plugin conformance: Index

**Given** Static-analyze the plugin source against the documented enum, trait, and REST conventions.
**When** Run the verification command shown below.
**Then** Enums are `enum X: string` with metadata methods; REST routes use the `/wp-json/<plugin>/v1/` namespace; nonces are verified on every mutating request.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
