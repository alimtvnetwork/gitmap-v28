# App

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

- `22-app-issues/` — bug reports and root-cause analyses for this app
- `23-app-db/` — database schema and queries for this app
- `24-app-ui-design-system/` — UI components and design tokens for this app

---

## Contents

_No app-specific specs added yet. Add specs as numbered files within this folder._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Issues | [../22-app-issues/00-overview.md](../22-app-issues/00-overview.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/00-overview.md](../01-spec-authoring-guide/00-overview.md) |

---

## Verification

_Auto-generated section — see `spec/21-app/97-acceptance-criteria.md` for the full criteria index._

### AC-APP-000: App-level conformance: Overview

**Given** Run the application's integration smoke suite.  
**When** Run the verification command shown below.  
**Then** Boot sequence completes; health endpoint returns 200; no unhandled promise rejections appear in the log.

**Verification command:**

```bash
npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
