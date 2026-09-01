# GitMap Application Specifications

**Version:** 3.2.0  
**Updated:** 2026-09-01  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Overview

The `spec/21-app/` directory is the canonical repository home for all **GitMap Application Specifications**. Following the repository-wide Coding Guidelines, all core foundational standards reside in `01–20`, while all application-level features, workflows, CLI surfaces, database models, and execution engines reside under `21-app/`.

GitMap is a CLI tool that scans directory trees for Git repositories, extracts clone URLs and branch information, and outputs structured data (terminal, CSV, JSON, folder-structure Markdown). It manages clones, clusters, SSH synchronization, macros, and releases.

---

## High-Level Components

1. **Core CLI & Scanner Specs (`01-overview.md` – `118-*.md`)**:
   - `01-overview.md`: Architecture overview and purpose
   - `02-cli-interface.md`: Complete command-line grammar, arguments, and aliases
   - `03-scanner.md`: Directory tree traversal and `.git` repository identification
   - `04-formatter.md`: Terminal glyph rendering, CSV, and JSON formatters
   - `05-cloner.md`: Re-cloning engine, hierarchy preservation, and speed optimization
   - `06-config.md`: Configuration loader, environment defaults, and user flags
   - `07-data-model.md`: Domain structures, database entities, and SQLite schemas
   - `12-release-command.md`: Release pipeline, semver bumping, and release ceremony
   - `50-ssh-keys.md` & `110-clone-ssh-flag.md`: SSH cluster synchronizer and remote execution
   - `110-update-remote-installer.md` & `111-update-remote-probe.md`: Sibling repository discovery and dynamic remote installer

2. **Feature Modules & Sub-systems**:
   - [`01-vscode-project-manager-sync/`](./01-vscode-project-manager-sync/): VS Code Project Manager synchronization
   - [`03-commit-in/`](./03-commit-in/): Multi-commit batching and commit-in replay workflow
   - [`03-general/`](./03-general/): CLI design patterns, PowerShell build/deploy helpers
   - [`03-tasks/`](./03-tasks/): Task automation, shortcut standardizations
   - [`04-generic-cli/`](./04-generic-cli/): Subcommand architectures, progress tracking, batch execution
   - [`07-generic-release/`](./07-generic-release/) & [`16-generic-release/`](./16-generic-release/): Release metadata, asset packaging, and verification
   - [`08-generic-update/`](./08-generic-update/): Self-update handoff, deploy path resolution, cleanup
   - [`08-json-schemas/`](./08-json-schemas/): Formal JSON schemas and contract definitions for CLI outputs
   - [`09-pipeline/`](./09-pipeline/) & [`09-pipeline-extend-v2/`](./09-pipeline-extend-v2/): CI/CD and release pipelines
   - [`15-distribution-and-runner/`](./15-distribution-and-runner/): Install contracts and runner specifications
   - [`16-scan-clone-cluster-interactive-macros/`](./16-scan-clone-cluster-interactive-macros/): Interactive macro recording, cluster daemon, and terminal UI
   - [`17-mv-rm-resolver-replace/`](./17-mv-rm-resolver-replace/): Repository moving, removal, and replacement engines
   - [`18-install-cg-ssh/`](./18-install-cg-ssh/): Coding guideline SSH installer
   - [`19-ssh-executor/`](./19-ssh-executor/): Remote SSH execution specifications
   - [`23-app-db/`](./23-app-db/): App-specific data model and SQLite tables
   - [`24-app-ui-design-system/`](./24-app-ui-design-system/): Terminal UI and docs styling
   - [`25-app-spec-audit/`](./25-app-spec-audit/) & [`26-coding-guideline-audit/`](./26-coding-guideline-audit/): Quality audits and reports

---

## Placement Rule

All application-specific features and implementation details belong in `21-app/`.
All bug investigations, root cause analysis, and fixes belong in `22-app-issues/`.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Issues | [../22-app-issues/00-overview.md](../22-app-issues/00-overview.md) |
| Coding Guidelines | [../02-coding-guidelines/00-overview.md](../02-coding-guidelines/00-overview.md) |
| Error Management | [../03-error-manage/00-overview.md](../03-error-manage/00-overview.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/00-overview.md](../01-spec-authoring-guide/00-overview.md) |
