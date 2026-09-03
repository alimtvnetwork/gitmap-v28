# Research

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 10 Research.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `10-research/`.
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

Dedicated folder for all exploratory and evaluative work that supports the spec system. This is the **single canonical location** for research content at the root spec level.

---

## What Belongs Here

| Content Type | Examples |
|-------------|----------|
| Comparative studies | Framework X vs Framework Y |
| Technology evaluations | Assessing a new library or tool |
| Exploratory technical notes | Proof-of-concept findings |
| Game development research | Engine comparisons, architecture patterns |
| Language evaluations | Assessing a new language for the stack |

## Placement Rule

All root-level research content MUST be placed in this folder (`02-spec/10-research/`) unless explicitly categorized elsewhere. Language-specific research within coding guidelines belongs in `02-spec/02-coding-guidelines/10-research/`.

---

## Contents

_No research documents added yet. Add research files as numbered entries (e.g., `01-framework-comparison.md`)._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding Guidelines Research | [../02-coding-guidelines/10-research/01-index.md](../02-coding-guidelines/10-research/01-index.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/01-index.md](../01-spec-authoring-guide/01-index.md) |
| Consolidated Guidelines | [../17-consolidated-guidelines/15-root-research.md](../17-consolidated-guidelines/15-root-research.md) |

---

## Verification

_Auto-generated section — see `02-spec/10-research/97-acceptance-criteria.md` for the full criteria index._

### AC-RES-001: Research-folder conformance: Index

**Given** Validate research note structure (front-matter, dated filenames, source links).
**When** Run the verification command shown below.
**Then** Every research note has a date prefix, a `Source:` line, and a `Decision:` or `Outcome:` section.

**Verification command:**

```bash
python3 linter-scripts/check-spec-folder-refs.py
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
