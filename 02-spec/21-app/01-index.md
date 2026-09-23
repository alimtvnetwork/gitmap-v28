# App

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 21 App.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `21-app/`.
- [ ] `/learn` Adhere strictly to `.ai-memory/folder-structure.md` and `.ai-memory/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Overview

App-specific specification content at the root spec level. This folder contains implementation specs, feature definitions, workflows, and architecture decisions for whatever project this repo ships — web app, Chrome extension, browser plugin, CLI tool, mobile app, WordPress plugin, desktop app, or any other deliverable.

Whatever the app is, **its product-level documentation lives here.** Foundational, cross-cutting guidelines (naming, error handling, design tokens, CI/CD, etc.) belong in the core fundamentals range (`01–20`).

---

## Placement Rule

Any content that defines a specific application feature, workflow, screen, command, or implementation detail belongs here, regardless of the app's runtime (browser, Node, PHP, Go, native, extension manifest, etc.). Foundational, reusable principles belong in the core fundamentals range (`01–20`).

Sibling folders for app-scoped concerns:

- `22-app-issues/` — bug reports and root-cause analyzes for this app
- `23-app-db/` — database schema and queries for this app
- `24-app-ui-design-system/` — UI components and design tokens for this app

---

## Contents

See [00-overview.md](./00-overview.md) for the complete application specification directory and architectural roadmap.

### Recent Specifications
- [Spec 129: PR Commit Engines & SQLite Split-DB](./129-pr-commit-engines-and-sqlite-split-db.md)
- [Spec 130: Pipeline-AI Live Error Streaming & Remediation](./130-pipeline-ai-live-error-streaming-and-remediation.md)
- [Spec 131: Developer Tools Cache Cleanup](./131-developer-tools-cache-cleanup.md)
- [Spec 132: SSH Multi-Node Execution, Copy/Move & Env](./132-ssh-multinode-exec-copy-mv-and-env.md)
- [Spec 133: SSH Interactive Join, Password Vault & Cluster](./133-ssh-interactive-join-password-vault-and-cluster.md)
- [Spec 134: Antigravity IDE-First Integration & Queue Protocol](./134-antigravity-ide-first-integration-and-queue-protocol.md)
- [Spec 135: Precompiled Test Warmup & Quad-Process Runner](./135-precompiled-test-warmup-and-dual-queue-runner.md)
- [Spec 136: Git Pull Efficient & Split-DB](./136-git-pull-efficient-and-split-database.md)
- [Spec 137: VM Cluster E2E SSH & OS Integration](./137-vm-cluster-e2e-ssh-and-os-integration.md)
- [Spec 138: VM Cluster Lifecycle, E2E Verification & OS Utility Suite](./138-vm-cluster-lifecycle-and-osutil-suite.md)
- [Spec 139: PE Multi-Line Extraction Intelligence & Declarative JSON Format Profiles](./139-pe-custom-format-and-error-capture.md)
- [Spec 140: SSH Fleet Update, Node Inspection, Smart Command Routing, and Macro Engine](./140-ssh-fleet-update-nodes-macros.md)

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Issues | [../22-app-issues/01-index.md](../22-app-issues/01-index.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/01-index.md](../01-spec-authoring-guide/01-index.md) |

---

## Verification

_Auto-generated section — see `02-spec/21-app/97-acceptance-criteria.md` for the full criteria index._

### AC-APP-001: App-level conformance: Index

**Given** Run the application's integration smoke suite.
**When** Run the verification command shown below.
**Then** Boot sequence completes; health endpoint returns 200; no unhandled promise rejections appear in the log.

**Verification command:**

```bash
npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
