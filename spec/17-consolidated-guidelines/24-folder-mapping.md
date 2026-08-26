# Folder-by-Folder Mapping — Source Specs ↔ Consolidated Guidelines

**Version:** 1.0.0  
**Updated:** 2026-04-22  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

This file is a **bidirectional index** between every numbered source-spec subfolder under `spec/` and the consolidated guideline files under `spec/17-consolidated-guidelines/`. It exists so that any AI handed the consolidated folder can answer two questions instantly:

1. **"Which source folder owns this rule?"** — given a topic, find the canonical spec.
2. **"Where in the consolidated docs is this folder summarized?"** — given a source folder, find every consolidated reference and any blind-spots.

Use this mapping as the **first stop** before implementing anything from the consolidated section. It surfaces gaps, cross-folder dependencies, and ownership boundaries.

---

## Keywords

`folder-mapping` · `source-of-truth` · `cross-reference` · `coverage-matrix` · `ai-onboarding` · `gap-analysis`

---

## How to Read This Document

Each row in the master matrix below contains:

| Column | Meaning |
|--------|---------|
| **Source Folder** | The canonical `spec/NN-name/` directory — single source of truth |
| **Contributes** | The rules, contracts, schemas, or patterns this folder owns |
| **Primary Consolidated File** | The consolidated `.md` that summarizes this folder end-to-end |
| **Secondary References** | Other consolidated files that reference rules from this folder |
| **Coverage** | ✅ Full · 🟡 Partial · 🔴 Reference-only — how completely the consolidated docs reflect the source |

---

## Master Matrix (All 22 Numbered Source Folders)

| # | Source Folder | Contributes | Primary Consolidated | Secondary Refs | Coverage |
|---|---------------|-------------|----------------------|----------------|----------|
| 01 | `01-spec-authoring-guide/` | Folder structure (01–20 core / 21+ app), naming conventions, required files, AI onboarding prompt, CLI/app/non-CLI templates, memory-folder guide, cross-reference rules, exceptions, mandatory linter infrastructure | [`01-spec-authoring.md`](./01-spec-authoring.md) | [`21-lovable-folder-structure.md`](./21-lovable-folder-structure.md), [`19-gap-analysis.md`](./19-gap-analysis.md) | ✅ Full |
| 02 | `02-coding-guidelines/` | Cross-language code style, language-specific rules (TS, Go, PHP, Rust, C#, PowerShell), AI optimization, CI/CD integration, file/folder naming, security, app-coding subsections | [`02-coding-guidelines.md`](./02-coding-guidelines.md) | [`04-enum-standards.md`](./04-enum-standards.md), [`11-research.md`](./11-research.md), [`13-app.md`](./13-app.md), [`14-app-issues.md`](./14-app-issues.md), [`16-app-design-system-and-ui.md`](./16-app-design-system-and-ui.md), [`22-app-database.md`](./22-app-database.md) | ✅ Full |
| 03 | `03-error-manage/` | 3-tier error architecture, `apperror` package, response envelope, Global Error Modal, error code registry, TypedQuery, retrospectives | [`03-error-management.md`](./03-error-management.md) | [`02-coding-guidelines.md`](./02-coding-guidelines.md) §error-handling, [`23-generic-cli.md`](./23-generic-cli.md) §07 | ✅ Full |
| 04 | `04-database-conventions/` | Naming conventions, schema design, ORM and views, testing strategy, relationship diagrams, REST API format, split-db pattern reference | [`18-database-conventions.md`](./18-database-conventions.md) | [`05-split-db-architecture.md`](./05-split-db-architecture.md), [`22-app-database.md`](./22-app-database.md) | ✅ Full |
| 05 | `05-split-db-architecture/` | Root/App/Session SQLite hierarchy, WAL pragma settings, RBAC via Casbin, schemas, reset API, user isolation, features and known issues | [`05-split-db-architecture.md`](./05-split-db-architecture.md) | [`18-database-conventions.md`](./18-database-conventions.md) §07, [`22-app-database.md`](./22-app-database.md) | ✅ Full |
| 06 | `06-seedable-config-architecture/` | SemVer GORM merge of `config.seed.json`, fundamentals, features, schema evolution, known issues | [`06-seedable-config.md`](./06-seedable-config.md) | [`05-split-db-architecture.md`](./05-split-db-architecture.md), [`23-generic-cli.md`](./23-generic-cli.md) §05 | ✅ Full |
| 07 | `07-design-system/` | CSS variable architecture, typography, spacing/layout, borders/shapes, motion/transitions, code blocks, header/nav, button system, sidebar, section patterns, page creation, WP migration | [`07-design-system.md`](./07-design-system.md) | [`16-app-design-system-and-ui.md`](./16-app-design-system-and-ui.md), [`08-docs-viewer-ui.md`](./08-docs-viewer-ui.md), [`09-code-block-system.md`](./09-code-block-system.md) | ✅ Full |
| 08 | `08-docs-viewer-ui/` | Documentation viewer: typography rules, keyboard navigation, fullscreen mode, search, rendering pipeline, features and issues | [`08-docs-viewer-ui.md`](./08-docs-viewer-ui.md) | [`07-design-system.md`](./07-design-system.md), [`09-code-block-system.md`](./09-code-block-system.md) | ✅ Full |
| 09 | `09-code-block-system/` | Code block architecture, HTML structure, syntax highlighting, interactions, styling, constants/maps, clipboard, checklist blocks, tree-structure rendering, selection bar | [`09-code-block-system.md`](./09-code-block-system.md) | [`07-design-system.md`](./07-design-system.md) §07, [`08-docs-viewer-ui.md`](./08-docs-viewer-ui.md) | ✅ Full |
| 10 | `10-research/` | Root-level research placement rules and references | [`12-root-research.md`](./12-root-research.md) | [`11-research.md`](./11-research.md) | 🟡 Partial — placeholder folder, expansion pending |
| 11 | `11-powershell-integration/` | PowerShell runner, pnpm PnP, pipeline steps, configuration schema, error codes, firewall rules, PHP known issues, multi-site deployment, examples and templates | [`10-powershell-integration.md`](./10-powershell-integration.md) | [`02-coding-guidelines.md`](./02-coding-guidelines.md) §powershell | ✅ Full |
| 12 | `12-cicd-pipeline-workflows/` | CI pipeline, release pipeline, vulnerability scanning, install-script generation, code signing, self-update mechanism, release body/changelog, installation flow, version & help, env-var setup, terminal output standards, binary icon branding, browser-extension and Go-binary deploy archetypes, reusable CI guards | [`15-cicd-pipeline-workflows.md`](./15-cicd-pipeline-workflows.md) | [`17-self-update-app-update.md`](./17-self-update-app-update.md), [`16-generic-release.md` source] | ✅ Full |
| 13 | `13-generic-cli/` | Generic CLI: project structure, subcommand architecture, flag parsing, configuration, output formatting, error handling, code style, help system, database, build/deploy, testing, checklist, date formatting, constants, verbose logging, progress tracking, batch execution, shell completion, terminal output design | [`23-generic-cli.md`](./23-generic-cli.md) | [`02-coding-guidelines.md`](./02-coding-guidelines.md), [`15-cicd-pipeline-workflows.md`](./15-cicd-pipeline-workflows.md) | ✅ Full |
| 14 | `14-update/` | Self-update overview, deploy-path resolution, rename-first deploy, build scripts, handoff mechanism, cleanup, console-safe handoff, repo-path sync, version verification, last-release detection, Windows icon embedding, code signing, release assets, checksums verification, release versioning, cross-compilation, release pipeline, install scripts, updater binary, **install-script version probe (20-repo ladder)** | [`17-self-update-app-update.md`](./17-self-update-app-update.md) | [`15-cicd-pipeline-workflows.md`](./15-cicd-pipeline-workflows.md) | ✅ Full |
| 15 | `15-distribution-and-runner/` | Install contract, runner contract, release pipeline, install config | [`17-self-update-app-update.md`](./17-self-update-app-update.md) §distribution | [`15-cicd-pipeline-workflows.md`](./15-cicd-pipeline-workflows.md) | 🟡 Partial — referenced inline, no dedicated consolidated file |
| 16 | `16-generic-release/` | Cross-compilation, release pipeline, install scripts, checksums verification, release assets, release metadata, known issues and fixes, version-pinned release installers | [`15-cicd-pipeline-workflows.md`](./15-cicd-pipeline-workflows.md) §generic-release | [`17-self-update-app-update.md`](./17-self-update-app-update.md) | 🟡 Partial — folded into CI/CD consolidated; no standalone file |
| 17 | `17-consolidated-guidelines/` | This folder itself — consolidated AI-readable digests and the gap analysis | (self) | All consolidated files | ✅ Full |
| 18 | `18-wp-plugin-how-to/` | WP plugin foundation/architecture, enums and coding style, traits/composition, logging/error handling, helpers/responses/integration, input validation, reference implementations, WP integration patterns, testing, deployment, frontend/template patterns, design system, admin UI, REST API conventions, settings architecture, error handling extraction, data file patterns, frontend JS patterns | [`20-wp-plugin-conventions.md`](./20-wp-plugin-conventions.md) | [`02-coding-guidelines.md`](./02-coding-guidelines.md) §php, [`03-error-management.md`](./03-error-management.md) | ✅ Full |
| 21 | `21-app/` | App-specific spec placement and decision guide | [`13-app.md`](./13-app.md) | [`14-app-issues.md`](./14-app-issues.md), [`22-app-database.md`](./22-app-database.md), [`16-app-design-system-and-ui.md`](./16-app-design-system-and-ui.md) | 🟡 Partial — placeholder, expanded inside app-* consolidated files |
| 22 | `22-app-issues/` | App bug analysis, issue file template, placement rules | [`14-app-issues.md`](./14-app-issues.md) | [`13-app.md`](./13-app.md) | 🟡 Partial — placeholder folder |
| 23 | `23-app-db/` | App-database conventions overlay on top of `04-database-conventions/` | [`22-app-database.md`](./22-app-database.md) | [`18-database-conventions.md`](./18-database-conventions.md), [`05-split-db-architecture.md`](./05-split-db-architecture.md) | ✅ Full |
| 24 | `24-app-ui-design-system/` | App-design overlay on top of `07-design-system/` | [`16-app-design-system-and-ui.md`](./16-app-design-system-and-ui.md) | [`07-design-system.md`](./07-design-system.md), [`08-docs-viewer-ui.md`](./08-docs-viewer-ui.md), [`09-code-block-system.md`](./09-code-block-system.md) | ✅ Full |

---

## Reverse Index — Consolidated File → Source Folders

| Consolidated File | Owns (Primary Source) | Pulls From (Secondary Sources) |
|-------------------|----------------------|--------------------------------|
| `01-spec-authoring.md` | `01-spec-authoring-guide/` | `10-research/`, `21-app/` |
| `02-coding-guidelines.md` | `02-coding-guidelines/` (all 16 subfolders) | `11-powershell-integration/`, `13-generic-cli/`, `18-wp-plugin-how-to/`, `03-error-manage/` |
| `03-error-management.md` | `03-error-manage/` | `02-coding-guidelines/11-security/`, `13-generic-cli/07-error-handling/`, `18-wp-plugin-how-to/04-logging-and-error-handling/` |
| `04-enum-standards.md` | `02-coding-guidelines/` (cross-language enum subsections) | `13-generic-cli/15-constants-reference/` |
| `05-split-db-architecture.md` | `05-split-db-architecture/` | `04-database-conventions/07-split-db-pattern/`, `23-app-db/` |
| `06-seedable-config.md` | `06-seedable-config-architecture/` | `05-split-db-architecture/`, `13-generic-cli/05-configuration/` |
| `07-design-system.md` | `07-design-system/` | `08-docs-viewer-ui/`, `09-code-block-system/`, `24-app-ui-design-system/` |
| `08-docs-viewer-ui.md` | `08-docs-viewer-ui/` | `07-design-system/`, `09-code-block-system/` |
| `09-code-block-system.md` | `09-code-block-system/` | `07-design-system/07-code-blocks/`, `08-docs-viewer-ui/` |
| `10-powershell-integration.md` | `11-powershell-integration/` | `02-coding-guidelines/09-powershell-integration/` |
| `11-research.md` | `02-coding-guidelines/10-research/` | `10-research/` |
| `12-root-research.md` | `10-research/` | `02-coding-guidelines/10-research/` |
| `13-app.md` | `21-app/` | `02-coding-guidelines/21-app/` |
| `14-app-issues.md` | `22-app-issues/` | `02-coding-guidelines/22-app-issues/`, `21-app/` |
| `15-cicd-pipeline-workflows.md` | `12-cicd-pipeline-workflows/` | `14-update/`, `15-distribution-and-runner/`, `16-generic-release/` |
| `16-app-design-system-and-ui.md` | `24-app-ui-design-system/` | `07-design-system/`, `02-coding-guidelines/24-app-ui-design-system/` |
| `17-self-update-app-update.md` | `14-update/` | `12-cicd-pipeline-workflows/06-self-update-mechanism/`, `15-distribution-and-runner/`, `16-generic-release/08-version-pinned-release-installers/` |
| `18-database-conventions.md` | `04-database-conventions/` | `05-split-db-architecture/`, `23-app-db/`, `02-coding-guidelines/23-app-db/` |
| `19-gap-analysis.md` | (meta — analyzes all consolidated files) | All folders |
| `20-wp-plugin-conventions.md` | `18-wp-plugin-how-to/` | `02-coding-guidelines/04-php/`, `07-design-system/`, `03-error-manage/` |
| `21-lovable-folder-structure.md` | `01-spec-authoring-guide/01-folder-structure.md` | `01-spec-authoring-guide/03-required-files/` |
| `22-app-database.md` | `23-app-db/` | `04-database-conventions/`, `05-split-db-architecture/` |
| `23-generic-cli.md` | `13-generic-cli/` | `02-coding-guidelines/`, `03-error-manage/`, `06-seedable-config-architecture/` |
| `24-folder-mapping.md` | (this file — meta-index) | All folders |

---

## Coverage Heatmap by Topic

| Topic | Owner Folder | Consolidated Reference | AI-Readiness |
|-------|--------------|------------------------|--------------|
| Folder structure & naming (01–20 core / 21+ app) | `01-spec-authoring-guide/01-folder-structure.md` | `01-spec-authoring.md` §2, `21-lovable-folder-structure.md` | ✅ |
| Cross-language coding standards | `02-coding-guidelines/01-cross-language/` | `02-coding-guidelines.md` §§1–10 | ✅ |
| TypeScript patterns & enums | `02-coding-guidelines/02-typescript/` | `02-coding-guidelines.md` §22, `04-enum-standards.md` | ✅ |
| Golang patterns & enums | `02-coding-guidelines/03-golang/` | `02-coding-guidelines.md` §23, `04-enum-standards.md` | ✅ |
| PHP patterns & enums | `02-coding-guidelines/04-php/` | `02-coding-guidelines.md` §24, `20-wp-plugin-conventions.md` | ✅ |
| Rust patterns & enums | `02-coding-guidelines/05-rust/` | `02-coding-guidelines.md` §25, `04-enum-standards.md` | ✅ |
| C# patterns & enums | `02-coding-guidelines/07-csharp/` | `02-coding-guidelines.md` §26 | ✅ |
| AI optimization | `02-coding-guidelines/06-ai-optimization/` | `02-coding-guidelines.md` §27 | ✅ |
| Security baseline | `02-coding-guidelines/11-security/` | `02-coding-guidelines.md` §31, `15-cicd-pipeline-workflows.md` §vulnerability-scanning | ✅ |
| File/folder naming (Zero-Underscore) | `02-coding-guidelines/08-file-folder-naming/` | `02-coding-guidelines.md` §29 | ✅ |
| Error architecture (`apperror`) | `03-error-manage/02-error-architecture/` | `03-error-management.md` §§1–17 | ✅ |
| Global Error Modal | `03-error-manage/02-error-architecture/` | `03-error-management.md` §20 | ✅ |
| Response envelope | `03-error-manage/02-error-architecture/` | `03-error-management.md` §21 | ✅ |
| Error code registry | `03-error-manage/03-error-code-registry/` | `03-error-management.md` §26 | ✅ |
| Database naming (PascalCase, `{TableName}Id`) | `04-database-conventions/01-naming-conventions.md` | `18-database-conventions.md` §§1–6 | ✅ |
| Free-text columns (Rules 10/11/12) | `04-database-conventions/02-schema-design.md` | `18-database-conventions.md` §18, linter rule `MISSING-DESC-001` | ✅ |
| Split DB hierarchy (Root/App/Session) | `05-split-db-architecture/` | `05-split-db-architecture.md`, `18-database-conventions.md` §07 | ✅ |
| Seedable config | `06-seedable-config-architecture/` | `06-seedable-config.md` | ✅ |
| Design tokens & themes (HSL CSS vars) | `07-design-system/02-theme-variable-architecture.md` | `07-design-system.md` §§1–5, `16-app-design-system-and-ui.md` | ✅ |
| Code block rendering | `09-code-block-system/` | `09-code-block-system.md`, `07-design-system.md` §07 | ✅ |
| Docs viewer UI | `08-docs-viewer-ui/` | `08-docs-viewer-ui.md` | ✅ |
| PowerShell runner | `11-powershell-integration/` | `10-powershell-integration.md` | ✅ |
| CI pipeline (lint, vuln, parallel tests) | `12-cicd-pipeline-workflows/01-ci-pipeline.md` | `15-cicd-pipeline-workflows.md` §1 | ✅ |
| Release pipeline | `12-cicd-pipeline-workflows/02-release-pipeline.md`, `14-update/17-release-pipeline.md` | `15-cicd-pipeline-workflows.md` §2, `17-self-update-app-update.md` | ✅ |
| Cross-compilation matrix | `14-update/16-cross-compilation.md`, `16-generic-release/01-cross-compilation.md` | `15-cicd-pipeline-workflows.md` §cross-compilation | ✅ |
| Code signing | `12-cicd-pipeline-workflows/05-code-signing.md`, `14-update/12-code-signing.md` | `15-cicd-pipeline-workflows.md` §code-signing, `17-self-update-app-update.md` | ✅ |
| Self-update (rename-first, atomicity) | `14-update/03-rename-first-deploy.md`, `14-update/05-handoff-mechanism.md` | `17-self-update-app-update.md` §§rename-first, handoff | ✅ |
| Console-safe handoff | `14-update/07-console-safe-handoff.md` | `17-self-update-app-update.md` §console-safe | ✅ |
| Install-script 20-repo version probe | `14-update/23-install-script-version-probe.md` | `17-self-update-app-update.md` §install-script-version-probe | ✅ |
| Checksums verification | `14-update/14-checksums-verification.md`, `16-generic-release/04-checksums-verification.md` | `15-cicd-pipeline-workflows.md` §checksums, `17-self-update-app-update.md` | ✅ |
| Windows icon embedding (`go-winres`) | `14-update/11-windows-icon-embedding.md`, `12-cicd-pipeline-workflows/09-binary-icon-branding.md` | `15-cicd-pipeline-workflows.md` §binary-branding | ✅ |
| Generic CLI architecture | `13-generic-cli/` | `23-generic-cli.md` | ✅ |
| WP plugin conventions | `18-wp-plugin-how-to/` | `20-wp-plugin-conventions.md` | ✅ |
| Distribution & runner contracts | `15-distribution-and-runner/` | `17-self-update-app-update.md` §distribution, `15-cicd-pipeline-workflows.md` | 🟡 No standalone consolidated file |
| Generic release archetype | `16-generic-release/` | `15-cicd-pipeline-workflows.md` §generic-release | 🟡 Folded into CI/CD consolidated |

---

## Known Blind-Spots (Tracked in `19-gap-analysis.md`)

| Gap | Source Folder(s) | Status |
|-----|------------------|--------|
| `15-distribution-and-runner/` has no standalone consolidated file — merged into `17-self-update-app-update.md` | `15-distribution-and-runner/` | 🟡 Acceptable — referenced inline; promote to standalone if it grows beyond 5 files |
| `16-generic-release/` folded into CI/CD consolidated rather than its own file | `16-generic-release/` | 🟡 Acceptable — overlaps heavily with `12-cicd-pipeline-workflows/` |
| `10-research/` and `21-app/`, `22-app-issues/` are placeholder folders with only `00-overview.md` | `10-research/`, `21-app/`, `22-app-issues/` | 🟡 Expected — content lives in app-* consolidated files |
| `02-coding-guidelines/06-cicd-integration/` cross-references CI/CD specs but is not summarized in `02-coding-guidelines.md` | `02-coding-guidelines/06-cicd-integration/` | ✅ Covered via `15-cicd-pipeline-workflows.md` cross-ref |

---

## Usage Instructions for AI Agents

When implementing **any** rule from the consolidated section:

1. **Locate** the rule in the relevant consolidated file (use the **Reverse Index** above).
2. **Identify** the owning source folder (use the **Master Matrix** above).
3. **Verify** the rule against the source spec **only if** the consolidated text lacks a code example, schema, or formula you need.
4. **If** the topic is marked 🟡 in the **Coverage** column, read the source folder directly — the consolidated file is summary-only.
5. **Never** modify a rule in a consolidated file without updating the source spec first. The source folder is the single source of truth.

---

## Cross-References

- [`00-overview.md`](./00-overview.md) — file inventory and scoring
- [`19-gap-analysis.md`](./19-gap-analysis.md) — formal gap analysis and implementability scores
- [`99-consistency-report.md`](./99-consistency-report.md) — validation history
- `../01-spec-authoring-guide/01-folder-structure.md` — canonical folder structure rules

---

## Validation History

| Date | Version | Action |
|------|---------|--------|
| 2026-04-22 | 1.0.0 | Initial folder-by-folder mapping; 22 source folders × 24 consolidated files indexed bidirectionally |

---

*Folder Mapping — v1.0.0 — 2026-04-22*
