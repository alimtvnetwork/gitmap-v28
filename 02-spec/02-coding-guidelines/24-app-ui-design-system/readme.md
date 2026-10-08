# App UI — Design System (AI Execution Prompt)

> **/goal** Standardize application UI design system architecture, component hierarchies, responsive styling, and design token hierarchies.
> **/learn** Master 10-step color ramp scales, token-driven typography, headless component primitives, and strict component size bounds (<= 300 lines).

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Define design system tokens using CSS variables or Tailwind tokens without hardcoded hex literals in components.
- [ ] `/learn` Enforce component decomposition so that every UI component remains under the 300-line canonical size limit.
- [ ] `/goal` Implement responsive, accessible layouts following ARIA standards and keyboard navigation contracts.
- [ ] `/learn` Validate that all design specifications use relative paths and align with the master design system token catalog.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below. This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Draft
**Ambiguity:** None

---

## Purpose

App-specific UI and design-system specifications, theming rules, component patterns, and layout conventions. This folder captures UI decisions that are specific to the application (web app, Chrome extension, plugin, CLI, etc.) rather than cross-language or cross-project.

---

## Contents

_No content yet. Add design system documents as numbered files within this folder._

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding Guidelines Overview | [../readme.md](../readme.md) |
| Consolidated Summary | [../../17-consolidated-guidelines/19-app-design-system-and-ui.md](../../17-consolidated-guidelines/19-app-design-system-and-ui.md) |

> Note: The consolidated guideline filename retains the historic `app-design-system-and-ui` name for backward compatibility; the source folder uses the canonical slug `24-app-ui-design-system`.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-APP-004: Application UI/UX Design System Specifications

**Given** Application UI components, design tokens, color scales, and layout patterns in `02-spec/02-coding-guidelines/24-app-ui-design-system/`.
**When** Audited against this reference specification and coding guidelines.
**Then** All UI components adhere to token hierarchies, 300-line limits, and accessible design system contracts with zero violations and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/24-app-ui-design-system --check-only
```
**Expected:** exit 0. Zero violations.
