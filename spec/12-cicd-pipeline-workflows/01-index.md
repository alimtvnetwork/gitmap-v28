# CI/CD Pipeline Workflows

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for CI/CD Pipeline Workflows.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

1. [ ] `/goal` Read and understand all numbered specifications under `12-cicd-pipeline-workflows/`.
2. [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
3. [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
4. [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.

---

**Version:** 4.1.0
**Updated:** 2026-08-30
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose

Central location for all CI/CD pipeline specifications, deployment automation, and related infrastructure-as-code documentation. All pipeline-related content — build pipelines, deployment workflows, environment promotion strategies, and CI/CD tooling configurations — MUST be documented in this folder.

---

## Scope

This module covers two distinct pipeline archetypes, shared conventions, reusable quality guards, and cross-cutting concerns:

| Archetype | Subfolder | Description |
|-----------|-----------|-------------|
| Browser Extension Deploy | `01-browser-extension-deploy/` | Node.js/pnpm multi-component builds, zip packaging, Chrome Web Store |
| Go Binary Deploy | `02-go-binary-deploy/` | Cross-compiled Go binaries, tar.gz/zip, install scripts, code signing |
| Reusable CI Guards | `03-reusable-ci-guards/` | 13 language-agnostic baseline diff gating and quality guards |
| Shared Conventions | Root files | Common patterns used across all pipeline types |

---

## Feature Inventory

### Root (Shared Conventions & Workflows)

| # | File | Description | Status |
|---|------|-------------|--------|
| 02 | [02-ci-pipeline.md](./02-ci-pipeline.md) | Core CI pipeline execution matrix and stages | ✅ Active |
| 03 | [03-shared-conventions.md](./03-shared-conventions.md) | Platform, triggers, concurrency, version resolution, checksums | ✅ Active |
| 04 | [04-github-release-standard.md](./04-github-release-standard.md) | Release body assembly, pre-release detection, asset matrix | ✅ Active |
| 05 | [05-release-pipeline.md](./05-release-pipeline.md) | Release workflow, deployment triggers, and AI release sync protocol | ✅ Active |
| 06 | [06-vulnerability-scanning.md](./06-vulnerability-scanning.md) | Standalone and in-CI vulnerability scanning patterns | ✅ Active |
| 07 | [07-install-script-generation.md](./07-install-script-generation.md) | Reusable PS1+Bash installer pattern, placeholder strategy, checksum verification | ✅ Active |
| 08 | [08-installation-flow.md](./08-installation-flow.md) | End-to-end install: one-liners, terminal output, upgrade, uninstall | ✅ Active |
| 09 | [09-changelog-integration.md](./09-changelog-integration.md) | Changelog format, CI extraction, release body assembly, terminal display | ✅ Active |
| 10 | [10-code-signing.md](./10-code-signing.md) | SignPath integration, feature-flag gating, signature verification | ✅ Active |
| 11 | [11-self-update-mechanism.md](./11-self-update-mechanism.md) | Generic CLI self-update blueprint: deploy path, rename-first, handoff, cleanup | ✅ Active |
| 12 | [12-version-and-help.md](./12-version-and-help.md) | Version display, help system, command-level docs, CI verification | ✅ Active |
| 13 | [13-environment-variable-setup.md](./13-environment-variable-setup.md) | `env` command: persistent variables, PATH registration, auto-home | ✅ Active |
| 14 | [14-release-body-and-changelog.md](./14-release-body-and-changelog.md) | Changelog extraction, release body template, asset matrix assembly | ✅ Active |
| 15 | [15-terminal-output-standards.md](./15-terminal-output-standards.md) | Output formatting: icons, tables, progress, errors, CI summaries | ✅ Active |
| 16 | [16-binary-icon-branding.md](./16-binary-icon-branding.md) | Windows binary icon embedding via `go-winres`: icon, manifest, version info | ✅ Active |
| 17 | [17-release-pipeline-issues-rca.md](./17-release-pipeline-issues-rca.md) | 🔴 Unified Root-Cause Analysis ledger of 13 CI/CD failure post-mortems and standing rules | ✅ Active |
| 18 | [18-lint-gating-rules.md](./18-lint-gating-rules.md) | Strict lint gating strategies, baseline diff rules, and CI guards master index | ✅ Active |
| 19 | [19-blue-green-deployment.md](./19-blue-green-deployment.md) | Zero-downtime blue/green deployment strategy | ✅ Active |
| 20 | [20-flaky-test-quarantine.md](./20-flaky-test-quarantine.md) | Flaky test detection and automated quarantine pattern | ✅ Active |
| 21 | [21-contract-testing.md](./21-contract-testing.md) | Microservice and API consumer contract testing | ✅ Active |
| 22 | [22-e2e-testing-pattern.md](./22-e2e-testing-pattern.md) | End-to-end integration and smoke test runner patterns | ✅ Active |
| 99 | [99-consistency-report.md](./99-consistency-report.md) | Consistency validation report for CI/CD workflows | ✅ Active |

---

## Subfolders

1. [ ] `/learn` [01-browser-extension-deploy/01-index.md](./01-browser-extension-deploy/01-index.md) — Chrome extension automated build and packaging.
2. [ ] `/learn` [02-go-binary-deploy/01-index.md](./02-go-binary-deploy/01-index.md) — Cross-platform Go binary packaging, signing, and release.
3. [ ] `/learn` [03-reusable-ci-guards/01-index.md](./03-reusable-ci-guards/01-index.md) — 13 modular, language-agnostic CI guards and checkers.

---

## Verification

_Auto-generated section — see `02-spec/12-cicd-pipeline-workflows/97-acceptance-criteria.md` for the full criteria index._

### AC-CI-001: CI/CD pipeline conformance: Index

**Given** Validate `.github/workflows/*.yml` against the documented job matrix.
**When** Run the verification command shown below.
**Then** Required jobs (`lint`, `cross-links`, `sync-drift`) are present; concurrency groups follow the `<workflow>-<ref>` pattern; `permissions:` is least-privilege.

**Verification command:**

```bash
npm run sync && npm run lint && npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
