# App UI — Design System

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
| `00-overview.md` present | ✅ |
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

- [Design System (Core)](../07-design-system/00-overview.md) — Foundational design system spec
- [App](../21-app/00-overview.md) — App-specific features and workflows
- [Consolidated Design System](../17-consolidated-guidelines/07-design-system.md) — Consolidated summary

---

*App UI — Design System — created 2026-04-10, renumbered 23→24 on 2026-04-16, slug renamed `24-app-design-system-and-ui` → `24-app-ui-design-system` on 2026-04-26*

---

## Verification

_Auto-generated section — see `spec/24-app-ui-design-system/97-acceptance-criteria.md` for the full criteria index._

### AC-ADS-000: App UI / design-system conformance: Overview

**Given** Scan app UI for raw colors and untokenized spacing; render Storybook (or equivalent) snapshot suite.  
**When** Run the verification command shown below.  
**Then** All components consume semantic tokens; snapshot diff is empty in light and dark themes.

**Verification command:**

```bash
npm run lint && npm run test
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
