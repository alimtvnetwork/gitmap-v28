# App Issues

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 22 App Issues.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `22-app-issues/`.
- [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Overview

App-specific issue analysis, root-cause analysis, bug documentation, and solution guidance — for whatever project this repo ships. Whether the app is a web app, Chrome extension, browser plugin, CLI tool, mobile app, WordPress plugin, or desktop app, **its bug reports and post-mortems live here.**

This folder tracks problems encountered during application development, their diagnosis, and their resolution.

---

## Placement Rule

Any content that analyzes bugs, failures, root causes, or fixes for application-level work belongs here, regardless of the app's runtime. General coding-principle violations or cross-cutting concerns belong in the core fundamentals range (`01–20`).

---

## Contents

_No app issue analyses added yet. Add issue files as numbered entries within this folder._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Specs | [../21-app/01-index.md](../21-app/01-index.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/01-index.md](../01-spec-authoring-guide/01-index.md) |

---

## Verification

_Auto-generated section — see `02-spec/22-app-issues/97-acceptance-criteria.md` for the full criteria index._

### AC-AI-001: App issues triage conformance: Index

**Given** Audit issue write-ups for the required Reproduction / Cause / Fix / Prevention sections.
**When** Run the verification command shown below.
**Then** Every issue file contains all four sections and references at least one commit or PR.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
