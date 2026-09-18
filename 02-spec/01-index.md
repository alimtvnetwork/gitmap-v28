# Specification Root

. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose

Root index for the entire specification tree. Each top-level folder contains a domain-specific specification module with its own overview, acceptance criteria, and consistency report.

---

## Module Inventory

### Core Fundamentals (01–20)

| # | Module | Description |
|---|--------|-------------|
| 01 | [Spec Authoring Guide](./01-spec-authoring-guide/01-index.md) | Rules for writing and maintaining spec documents |
| 02 | [Coding Guidelines](./02-coding-guidelines/01-index.md) | Cross-language coding standards (Go, TS, PHP, Rust, C#) |
| 03 | [Error Management](./03-error-manage/01-index.md) | Error capture, modal UI, and resolution workflows |
| 04 | [Database Conventions](./04-database-conventions/01-index.md) | Naming, schema design, ORM, REST API format |
| 05 | [Split DB Architecture](./05-split-db-architecture/01-index.md) | SQLite partitioning and migration patterns |
| 06 | [Seedable Config (CW Config)](./06-seedable-config-architecture/01-index.md) | Configuration seeding and feature management |
| 07 | [Design System](./07-design-system/01-index.md) | Theme variables, typography, spacing, and component patterns |
| 08 | [Docs Viewer UI](./08-docs-viewer-ui/01-index.md) | Specification viewer React application |
| 09 | [Code Block System](./09-code-block-system/01-index.md) | Rich code block rendering — syntax highlighting, interactions, styling |
| 10 | [Research](./10-research/01-index.md) | Comparative studies, technology evaluations, exploratory notes |
| 11 | [PowerShell Integration](./11-powershell-integration/01-index.md) | PowerShell scripting conventions, cross-platform automation |
| 12 | [CI/CD Pipeline Workflows](./12-cicd-pipeline-workflows/01-index.md) | CI/CD pipeline specs, deployment workflows, automation |
| 13 | [Generic CLI Standards](./13-generic-cli/01-index.md) | Generic CLI terminal UX, colors, flags, and help contracts |
| 14 | [Self-Update & App Update](./14-update/01-index.md) | Rename-first deployment, release pipeline, cross-compilation |
| 15 | [Distribution and Runner](./15-distribution-and-runner/01-index.md) | Cross-platform installers, runners, and repository forwarding |
| 16 | [Generic Release Pipeline](./16-generic-release/01-index.md) | Universal release pipeline blueprint and asset matrices |
| 17 | [Consolidated Guidelines](./17-consolidated-guidelines/01-index.md) | AI-readable summaries of every major spec module |
| 18 | [WP Plugin Conventions](./18-wp-plugin-how-to/01-index.md) | WordPress plugin architecture, admin UI, and REST API |
| 19 | [Main Worker Service Architecture](./19-main-worker-service/01-index.md) | Split-tier architecture, credential-blind proxy, and backup nodes |

### App-Specific (21+)

| # | Module | Description |
|---|--------|-------------|
| 21 | [App](./21-app/01-index.md) | App-specific specs: features, workflows, architecture |
| 22 | [App Issues](./22-app-issues/01-index.md) | App bug analysis, root cause analysis, fix documentation |
| 23 | [App DB](./23-app-db/01-index.md) | App-specific data model, table designs, migration strategies |
| 24 | [App UI — Design System](./24-app-ui-design-system/01-index.md) | App-specific UI, design system, theming, component patterns |

---

## Supporting Files

| File | Purpose |
|------|---------|
| [folder-structure-root.md](./folder-structure-root.md) | Redirect to canonical folder structure spec |
| [spec-index.md](./spec-index.md) | Flat index of all spec files |
| [health-dashboard.md](./health-dashboard.md) | Spec tree health metrics and broken link report |
| [dashboard-data.json](./dashboard-data.json) | Machine-readable health data |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Folder Structure (canonical) | `./01-spec-authoring-guide/02-folder-structure.md` |
| Spec Authoring Guide | `./01-spec-authoring-guide/01-index.md` |
| Coding Guidelines | `./02-coding-guidelines/01-index.md` |
| Error Management | `./03-error-manage/01-index.md` |
| Design System | `./07-design-system/01-index.md` |
