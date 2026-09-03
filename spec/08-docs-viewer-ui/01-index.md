# Docs Viewer UI — Overview

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 08 Docs Viewer Ui.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `08-docs-viewer-ui/`.
- [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** 90%
**Ambiguity:** 10%

---

## Overview

UI enhancement specification for the interactive documentation viewer at `/docs`. Covers typography, syntax highlighting, navigation, and usability improvements.

---

## Keywords

docs-viewer, typography, ubuntu, poppins, syntax-highlighting, highlight-js, fullscreen, keyboard-navigation, copy-markdown, shortcuts

---

## Scoring

| Criterion | Score | Notes |
|-----------|-------|-------|
| Completeness | 90% | All 6 features specified |
| Testability | 85% | UI interactions are manually testable |
| Ambiguity | 10% | Clear requirements |
| Cross-Refs | 80% | Links to spec authoring guide |

**Health Score:** 90%

---

## Document Inventory

| # | File | Description |
|---|------|-------------|
| 00 | `01-index.md` | This file — master index |
| 01 | `02-fundamentals.md` | Architecture, typography system, navigation model |
| 02 | `02-features/01-index.md` | Feature index |
| 02.01 | `02-features/01-typography.md` | Ubuntu + Poppins font integration |
| 02.02 | `02-features/02-syntax-highlighting.md` | highlight.js code block rendering |
| 02.03 | `02-features/03-fullscreen-mode.md` | Fullscreen toggle for doc content |
| 02.04 | `02-features/04-keyboard-navigation.md` | Arrow key file/folder navigation |
| 02.05 | `02-features/05-copy-markdown.md` | Copy raw markdown to clipboard |
| 02.06 | `02-features/06-shortcuts-overlay.md` | Keyboard shortcuts help modal |
| 02.07 | `02-features/06-ui-theme-animations.md` | Theme colors and animation system |
| 02.08 | `02-features/07-visual-rendering-guide.md` | Complete visual rendering & component guide |
| 99 | `99-consistency-report.md` | Consistency report |

---

## Cross-References

- [Spec Authoring Guide](../01-spec-authoring-guide/01-index.md) — Conventions
- [App Project Template](../01-spec-authoring-guide/07-app-project-template.md) — Template used

---

*Overview — updated: 2026-04-03*

---

## Verification

_Auto-generated section — see `02-spec/08-docs-viewer-ui/97-acceptance-criteria.md` for the full criteria index._

### AC-UI-001: Docs viewer UI conformance: Index

**Given** Render the docs viewer against the spec tree fixture.
**When** Run the verification command shown below.
**Then** Keyboard navigation, syntax highlighting, fullscreen toggle, and copy-markdown all function without console errors.

**Verification command:**

```bash
npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
