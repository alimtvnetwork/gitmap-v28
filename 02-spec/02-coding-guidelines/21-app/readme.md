# App (AI Execution Prompt)

> **/goal** Establish application-specific specification standards, feature workflows, and architecture boundaries distinct from foundational cross-cutting guidelines.
> **/learn** Understand the strict placement boundary dividing general core fundamentals (01–20) from app-specific implementations and feature models (21+).

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Place all application feature workflows, subsystem specs, and business domains in `21-app/`.
- [ ] `/learn` Keep foundational, reusable language principles in core directories (01–20); do not mix general and app-specific rules.
- [ ] `/goal` Enforce consistent numbered file naming and structured markdown architecture within `21-app/`.
- [ ] `/learn` Verify zero broken links and 100% relative paths across all app specification files.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

App-specific specification content. This folder contains implementation specs, feature definitions, workflows, and architecture decisions that are specific to application-level work — as opposed to foundational, cross-cutting guidelines.

---

## Placement Rule

Any content that defines a specific application feature, workflow, or implementation detail belongs here. Foundational, reusable principles belong in the core fundamentals range (01–20).

---

## Contents

_No content yet. Add app-specific specs as numbered files within this folder._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Issues | [../22-app-issues/readme.md](../22-app-issues/readme.md) |
| Coding Guidelines Spec | [../readme.md](../readme.md) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-APP-001: Application Specifications Index & Navigation

**Given** Application-specific specifications and feature architecture workflows in `02-spec/02-coding-guidelines/21-app/`.
**When** Audited against this reference specification and coding guidelines.
**Then** All app specifications are indexed, zero inconsistencies or invalid cross-references are detected, and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/21-app --check-only
```
**Expected:** exit 0. Zero violations.
