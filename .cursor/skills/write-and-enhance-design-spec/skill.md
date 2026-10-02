---
name: write-and-enhance-design-spec
description: Autonomously ingest source designs (codebases, specifications, presentation decks, UI components, screenshots) and author or enhance public design-system specifications with mathematical precision, structured tables, and zero hallucinations.
---

# Write and Enhance a Design Spec That a Blind AI Can Follow

> **[/goal](slashCommand;goal)** Autonomously ingest source designs (codebases, specifications, presentation decks, UI components, screenshots) and author or enhance public design-system specifications so comprehensive, mathematically exact, and rigorous that any "blind AI" or human engineer with zero external context can follow them blindly to design websites, presentations, blogs, animations, menus, buttons, builder modes, or images without guessing a single value or hallucinating.
> **[/learn](slashCommand;learn)** A specification is finished ONLY when every component, token, physics parameter, layout slot, and interaction timing in the source has its exact values written down in structured tables, every gap is audited, and the completion gate at the end passes 100%.

**Source prompt:** `01-prompts/10-ui-and-design/09-write-and-enhance-design-spec.md`

---

## 1. When to Use

Activate this skill when:
- Authoring new design-system specifications under `02-spec/07-design-system/`.
- Auditing and enhancing existing UI/UX, website, blog, or slide specifications.
- Reverse-engineering token systems, component dimensions, and physics curves from source implementations.
- Establishing closed sets, typography scales, spacing tokens, and contrast compliance tables.

---

## 2. Hard Architectural Rules

1. **Total Ban on Invented Values:** Every hex, HSL, RGB, OKLCH, pixel, millisecond, easing, and limit MUST originate from `SOURCE` or grounded mathematics. If source is silent, explicitly state: `Not specified in source. Do not invent.` Never write "about", "around", or "~".
2. **Strict Anonymization Mandate:** Never leak client company names, internal repositories, or personal names into public specifications, comments, sample text, or image paths. Always use generic enterprise designations (`{BRAND}`, `{COMPANY}`, `Executive Leadership`).
3. **Pure DOM Text Rendering Mandate (Zero Baked-In Text):** In websites, presentations, and UI designs, all headlines, body copy, bullets, metrics, credentials, and captions MUST be live DOM HTML elements (`<h1>`, `<p>`, `<span>`). Raster images are strictly reserved for photographic visual plates, author avatars, and partner logos.
4. **Strict Relative Git Paths Only:** All file paths, markdown links, subtask paths, and citations MUST be strictly relative paths starting from the repository root (e.g. `02-spec/07-design-system/11-button-system.md`). TOTAL BAN on absolute filesystem paths and `file:///` URIs.
5. **Lowercase File Naming Convention:** All generated specification files MUST use strictly lowercase kebab-case naming (e.g. `34-slide-layout-catalog.md`). No uppercase characters, spaces, or underscores.
6. **Strict Boolean Standards:** Positive booleans MUST ALWAYS be evaluated implicitly (`if isReady`). NEVER evaluate explicitly against true (`if isReady == true` is banned). NEVER combine positive and negative checks in the same condition.
7. **Closed Sets Rule:** Every selectable options list (button variants, sizes, slide layout types, theme IDs, pill preset colors, transition kinds, hotkeys, builder modes) MUST be specified as a closed set followed by: `Anything else does not exist.`
8. **Values Beside the Rule:** Never simply write "use the primary color". Document the semantic token, HEX value, RGB coordinates, and HSL values in the same row.
9. **Strict File Size Cap:** Specification files MUST remain bounded (at most 300 lines per file). Complex topics must be broken into discrete, single-responsibility files.

---

## 3. Systematic 5-Phase Workflow

1. **Phase 1: Comprehensive Source Inventory:** Scan source directories, token files, stylesheets, component primitives, hooks, route templates, and existing specs. Construct structured inventory tables for tokens, typography, components, sections, motion, slides, builders, and images.
2. **Phase 2: Gap & Fidelity Audit:** Compare every inventory item against target specs (`exact`, `wrong`, `missing`).
3. **Phase 3: Bounded Spec Authoring & Remediation:** Author or patch files in bounded micro-batches with complete tables and closed sets.
4. **Phase 4: Multi-Theme & Accessibility Verification:** Check all 10 themes and WCAG AA contrast compliance ($4.5:1$ body, $3:1$ large text).
5. **Phase 5: Blind-AI Completion Gate:** Verify zero vague phrases, pure DOM text mandate, relative paths, lowercase names, and 100% table coverage.
