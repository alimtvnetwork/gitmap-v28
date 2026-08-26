# Research

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

All root-level research content MUST be placed in this folder (`spec/10-research/`) unless explicitly categorized elsewhere. Language-specific research within coding guidelines belongs in `spec/02-coding-guidelines/10-research/`.

---

## Contents

_No research documents added yet. Add research files as numbered entries (e.g., `01-framework-comparison.md`)._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding Guidelines Research | [../02-coding-guidelines/10-research/00-overview.md](../02-coding-guidelines/10-research/00-overview.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/00-overview.md](../01-spec-authoring-guide/00-overview.md) |
| Consolidated Guidelines | [../17-consolidated-guidelines/12-root-research.md](../17-consolidated-guidelines/12-root-research.md) |

---

## Verification

_Auto-generated section — see `spec/10-research/97-acceptance-criteria.md` for the full criteria index._

### AC-RES-000: Research-folder conformance: Overview

**Given** Validate research note structure (front-matter, dated filenames, source links).  
**When** Run the verification command shown below.  
**Then** Every research note has a date prefix, a `Source:` line, and a `Decision:` or `Outcome:` section.

**Verification command:**

```bash
python3 linter-scripts/check-spec-folder-refs.py
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
