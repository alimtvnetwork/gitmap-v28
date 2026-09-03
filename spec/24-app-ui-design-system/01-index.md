# App UI — Design System

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 24 App Ui Design System.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `24-app-ui-design-system/`.
- [ ] `/learn` Adhere strictly to `.lovable/folder-structure.md` and `.lovable/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Draft
**Ambiguity:** None

---

## Keywords

`app-ui` · `app-design-system` · `theming` · `components` · `layout`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `01-index.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Application-specific UI and design-system specifications for whatever project this repo ships — web app, Chrome extension, CLI, plugin, mobile app, etc. Covers component patterns, theming decisions, layout conventions, and visual standards specific to this application.

---

## Document Inventory

| # | File | Purpose |
|---|------|---------|
| — | *(empty — awaiting content)* | — |

---

## Cross-References

- [Design System (Core)](../07-design-system/01-index.md) — Foundational design system spec
- [App](../21-app/01-index.md) — App-specific features and workflows
- [Consolidated Design System](../17-consolidated-guidelines/10-design-system.md) — Consolidated summary

---

*App UI — Design System — created 2026-04-10, renumbered 23→24 on 2026-04-16, slug renamed `24-app-design-system-and-ui` → `24-app-ui-design-system` on 2026-04-26*

---

## Verification

_Auto-generated section — see `02-spec/24-app-ui-design-system/97-acceptance-criteria.md` for the full criteria index._

### AC-ADS-001: App design-system conformance: Index

**Given** Scan app UI for raw colors and untokenized spacing; render Storybook (or equivalent) snapshot suite.
**When** Run the verification command shown below.
**Then** All components consume semantic tokens; snapshot diff is empty in light and dark themes.

**Verification command:**

```bash
npm run lint && npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
