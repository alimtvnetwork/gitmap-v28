# How to Write a WordPress Plugin — Gold Standard Specification

> **Updated:** 2026-04-14  
> **Purpose:** Comprehensive, self-contained guide for building WordPress plugins following the Gold Standard architecture.  
> **Audience:** AI code generators and human developers.

---

## Phase Index

| # | File | Title | Description |
|---|------|-------|-------------|
| 00 | [00-quick-start.md](00-quick-start.md) | **Quick Start Guide** | Condensed onboarding — read this first. 5 starter files, critical patterns, decision matrix, common mistakes |
| 01 | [01-foundation-and-architecture.md](01-foundation-and-architecture.md) | Foundation & Architecture | PHP 8.1+ requirements, canonical folder structure, bootstrap file, autoloader, Core Plugin singleton, ABSPATH guard, global type imports, version tracking |
| 02 | [02-enums-and-coding-style/](02-enums-and-coding-style/00-overview.md) | Enums & Coding Style | Backed enum architecture (4 categories), `match`-based metadata pattern, `PluginConfigType`, `PhpNativeType`, coding style rules (conditionals, formatting, naming, PHPDoc) |
| 03 | [03-traits-and-composition.md](03-traits-and-composition.md) | Traits & Composition | Trait folder structure, anatomy, `ResponseTrait` (safeExecute), `RouteRegistrationTrait`, `AuthTrait`, `TypeCheckerTrait`, composition in Plugin.php |
| 04 | [04-logging-and-error-handling.md](04-logging-and-error-handling.md) | Logging & Error Handling | Two-tier logging (FileLogger + ErrorLogHelper), debug mode gate, log format, stack traces, rotation, deduplication, `safeExecute()`, shutdown handler, `DateHelper` |
| 05 | [05-helpers-responses-and-integration.md](05-helpers-responses-and-integration.md) | Helpers, Response Envelope & Integration | Helper classes (PathHelper, BooleanHelpers, InitHelpers, HttpConfigType), `EnvelopeBuilder` fluent API, standard response format, integration checklist, Split DB concept, Database trait decomposition |
| 06 | [06-input-validation-patterns.md](06-input-validation-patterns.md) | Input Validation Patterns | Guard clauses, body validation, enum-based validation, sanitization rules |
| 07 | [07-reference-implementations.md](07-reference-implementations.md) | Reference Implementations | Complete working examples: bootstrap, autoloader, Plugin.php, handler trait, enum, helper, migration, settings page, REST endpoint |
| 08 | [08-wordpress-integration-patterns.md](08-wordpress-integration-patterns.md) | WordPress Integration Patterns | Admin pages & settings API, AJAX handlers, WP-Cron scheduled tasks, file upload handling, SQLite migrations, database seeding from JSON |
| 09 | [09-testing-patterns.md](09-testing-patterns.md) | Testing Patterns | Internal test framework, assertion helpers, database tests, validation tests, seed data verification |
| 10 | [10-deployment-patterns.md](10-deployment-patterns.md) | Deployment Patterns | Build pipeline, `.distignore`, self-update with rollback, `SelfUpdateStatusType` enum, `UpdateResolver` trait decomposition, uninstall cleanup |
| 11 | [11-frontend-and-template-patterns.md](11-frontend-and-template-patterns.md) | Frontend & Template Patterns | 200-line file limit, orchestrator pattern, partial templates, JS/CSS enqueuing, optional React integration, source maps, decision matrix |
| 12 | [12-design-system.md](12-design-system.md) | Design System | CSS variables, typography, color tokens, shadow hierarchy, animation library, badge patterns, card patterns, button variants, form input styling |
| 13 | [13-admin-ui-patterns.md](13-admin-ui-patterns.md) | Admin UI Patterns | Page layout, filter bars, table patterns with date grouping, badge system, modal anatomy, notices, empty/loading states, stats bars, toggle switches |
| 14 | [14-rest-api-conventions.md](14-rest-api-conventions.md) | REST API Conventions | Route naming, namespace patterns, HTTP method selection, pagination, category grouping, `endpoints.json` format, controller organisation |
| 15 | [15-settings-architecture.md](15-settings-architecture.md) | Settings Architecture | `OptionNameType`/`SettingsKeyType` enums, settings groups, defaults, validation, toggle switches, conditional display, save feedback, action buttons |
| 16 | [16-error-handling-extraction.md](16-error-handling-extraction.md) | Error Handling & Diagnostics | `ErrorType` classification, bootstrap vs FileLogger capture, `safeExecute` wrapper, admin error sessions, flash banners |
| 17 | [17-data-file-patterns.md](17-data-file-patterns.md) | Data File Patterns | JSON schemas for `colors.json`/`endpoints.json`/`openapi.json`, enum-driven access via `ColorGroupType`, static caching, validation |
| 18 | [18-frontend-javascript-patterns.md](18-frontend-javascript-patterns.md) | Frontend JavaScript Patterns | Localized object pattern (`wp_localize_script`), button state management, modal interactions, table rendering, i18n conventions |
| 19 | [19-micro-orm-and-root-db.md](19-micro-orm-and-root-db.md) | Micro-ORM & Root Database | Fluent query builder, TypedQuery with Go-style DbResult/DbResultSet/DbExecResult wrappers, FileCache with scan/store trait decomposition, cross-plugin RootDb manifest |
| 20 | [20-end-to-end-walkthrough.md](20-end-to-end-walkthrough.md) | **End-to-End Walkthrough** | Complete "Task Tracker" plugin built from scratch — 14 steps covering all 19 phases, with phase coverage matrix and final checklist |
| 21 | [21-ping-endpoint.md](21-ping-endpoint.md) | **Ping Endpoint** | Standard `/ping` health-check endpoint — returns author, company, version. Supports authorized and non-authorized modes via `PluginConfigType` |

---

## Phase 2 Subfiles

| File | Description |
|------|-------------|
| [01-enum-architecture.md](02-enums-and-coding-style/01-enum-architecture.md) | Core enum pattern, 4 standard categories, comparison methods, coding style rules |
| [02-enum-metadata-pattern.md](02-enums-and-coding-style/02-enum-metadata-pattern.md) | `match`-based metadata methods (`label`, `icon`, `cssClass`), `is*()` helpers |
| [03-self-update-status-enum.md](02-enums-and-coding-style/03-self-update-status-enum.md) | `SelfUpdateStatusType` — reference implementation (17 cases, deployment domain) |
| [04-action-type-enum.md](02-enums-and-coding-style/04-action-type-enum.md) | `ActionType` — reference implementation (40+ cases, transaction logging) |

---

## Reading Order

1. **Phase 1** — structural rules and folder layout
2. **Phase 2** — enum patterns (used everywhere)
3. **Phases 3–5** — traits, logging, helpers (core runtime)
4. **Phase 6** — input validation
5. **Phase 7** — reference implementations (ties 1–6 together)
6. **Phases 8–10** — WordPress integration, testing, deployment
7. **Phases 11–13** — frontend templates, design system, UI patterns
8. **Phases 14–18** — REST API, settings, error diagnostics, data files, JS patterns
9. **Phase 19** — micro-ORM query builder and cross-plugin Root Database
10. **Phase 20** — end-to-end walkthrough (ties everything together with a working plugin)

---

## Key Cross-References

| From | To | Topic |
|------|----|-------|
| Phase 1, §1.8 | Phase 2 | `PluginConfigType` version tracking |
| Phase 2, §2.3 | Phase 3, §3.8 | `TypeCheckerTrait` for syntax-validator-safe type checks |
| Phase 3, §3.4 | Phase 4, §4.9 | `safeExecute()` error boundary |
| Phase 5 | Phase 4, §4.11/§4.13 | `ErrorLogHelper` and `DateHelper` specs |
| Phase 7 | Phase 8, §8.5.1 | Database seeding patterns |
| Phase 7 | Phase 11, §11.1/§11.3/§11.8/§11.10 | File limits, orchestrator pattern, source maps, React decision |
| Phase 10, §10.5 | Phase 2, 03-self-update-status-enum | `SelfUpdateStatusType` enum |
| Phase 13 | Phase 12 | Design system tokens |
| Phase 14 | Phases 3, 4, 5, 6 | Trait anatomy, error handling, envelope, validation |
| Phase 15 | Phases 12, 13 | Design system + UI patterns for settings pages |
| Phase 19 | Phase 5, §5.5/§5.5.1 | ORM integrates with Database trait decomposition and Split DB |

---

## Changelog

See [changelog.md](changelog.md) for version history and all P1/P2 fixes applied.

---

*Maintain this index when adding or renaming spec phases.*
