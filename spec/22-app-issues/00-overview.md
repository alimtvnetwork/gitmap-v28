# App Issues

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
| App Specs | [../21-app/00-overview.md](../21-app/00-overview.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/00-overview.md](../01-spec-authoring-guide/00-overview.md) |

---

## Verification

_Auto-generated section — see `spec/22-app-issues/97-acceptance-criteria.md` for the full criteria index._

### AC-AI-000: App issues triage conformance: Overview

**Given** Audit issue write-ups for the required Reproduction / Cause / Fix / Prevention sections.  
**When** Run the verification command shown below.  
**Then** Every issue file contains all four sections and references at least one commit or PR.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
