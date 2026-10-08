# App Issues (AI Execution Prompt)

> **/goal** Standardize application-level issue tracking, root cause analysis (RCA), bug diagnoses, and resolution protocols.
> **/learn** Master the 4-part RCA methodology (Root Cause, Impact, Detection Gap, Prevention Rule) and isolate app-specific bugs from foundational cross-cutting violations.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Document all application-specific bug analyses and root-cause reports with structured 4-part RCA formatting.
- [ ] `/learn` Separate application-level runtime bugs from core coding guideline principle violations.
- [ ] `/goal` Maintain strict traceability between reported app issues and verified code fixes.
- [ ] `/learn` Verify that all cross-references to issue logs and architecture specs use strictly relative paths.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

App-specific issue analysis, root cause analysis, bug documentation, and solution guidance. This folder tracks problems encountered during application development, their diagnosis, and their resolution.

---

## Placement Rule

Any content that analyzes bugs, failures, root causes, or fixes for application-level work belongs here. General coding principle violations or cross-cutting concerns belong in the core fundamentals range (01–20).

---

## Contents

_No content yet. Add app issue analyses as numbered files within this folder._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Specs | [../21-app/readme.md](../21-app/readme.md) |
| Coding Guidelines Spec | [../readme.md](../readme.md) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-APP-002: Non-CI/CD Application Issues & Bug Catalog Index

**Given** Non-CI/CD application issues, bug analyses, and 4-part RCA documentation in `02-spec/02-coding-guidelines/22-app-issues/`.
**When** Audited against this reference specification and coding guidelines.
**Then** All app issues are documented with 4-part RCA formatting, zero invalid cross-references exist, and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/22-app-issues --check-only
```
**Expected:** exit 0. Zero violations.
