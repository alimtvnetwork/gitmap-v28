# UI and Design Prompts Library

> **/goal** Master, execute, and author canonical AI prompts for UI/UX engineering, presentation slide decks, visual builders, precision mega menus, and design specification generation.
> **/learn** All prompts in this directory use strictly relative repository paths, enforce pure DOM typography without baked-in text, require zero private/company names, and follow grounded specifications in `02-spec/07-design-system/`.

**Version:** 4.2.0
**Status:** Active

---

## Catalog of UI and Design Prompts

| File Name | Purpose & Workflow | Key Reading Specifications |
|:---|:---|:---|
| [`01-logo-create.md`](./01-logo-create.md) | Logo design, brand marks, and identity generation. | `02-spec/07-design-system/01-svg/` |
| [`02-react-ui-fixes-update.md`](./02-react-ui-fixes-update.md) | Modern React component fixes, Tailwind refactoring, and CSS variables. | `02-spec/07-design-system/11-button-system.md` |
| [`03-svg-logo.md`](./03-svg-logo.md) | Scalable, responsive SVG vector icon and brand logo authoring. | `02-spec/07-design-system/01-svg/` |
| [`04-youtube-thumbnail-create.md`](./04-youtube-thumbnail-create.md) | 16:9 high-CTR YouTube thumbnail design, typography, and safe zones. | `02-spec/07-design-system/37-image-specifications.md` |
| [`05-linkedin-profile-banner.md`](./05-linkedin-profile-banner.md) | High-authority personal LinkedIn profile banner authoring (1584×396). | `02-spec/07-design-system/37-image-specifications.md` |
| [`06-linkedin-company-banner.md`](./06-linkedin-company-banner.md) | Enterprise organization LinkedIn banner authoring (1128×191). | `02-spec/07-design-system/37-image-specifications.md` |
| [`07-follow-ui-ux-design-system.md`](./07-follow-ui-ux-design-system.md) | Core prompt for AI to construct marketing websites, precision mega menus, buttons, and sections. | `02-spec/07-design-system/10-header-navigation.md`, `33-mega-menu-components.md`, `11-button-system.md` |
| [`08-create-slide-deck.md`](./08-create-slide-deck.md) | Core prompt for AI to construct 16:9 presentation slide decks, HUD controller, step sound engine, webcam PIP, handouts, and Slide Builder Mode. | `02-spec/07-design-system/31-slide-controller-buttons.md`, `32-slide-color-options.md`, `34-slide-layout-catalog.md`, `35-slide-builder-canvas-inspector.md`, `42-slide-step-and-sound-system.md`, `43-slide-webcam-overlay.md`, `44-slide-presenter-inspector-and-handouts.md` |
| [`09-write-and-enhance-design-spec.md`](./09-write-and-enhance-design-spec.md) | Master prompt for blind AI to ingest source codebases/decks and author or enhance design specifications. | `02-spec/07-design-system/readme.md`, `32-slide-color-options.md`, `34-slide-layout-catalog.md`, `35-slide-builder-canvas-inspector.md`, `36-website-content-builder-mode.md` |
| [`10-blind-spec-creator-and-enhancer.md`](./10-blind-spec-creator-and-enhancer.md) | Exhaustive meta-instruction protocol for blind AI to create and enhance design specifications with exact checklists and zero hallucination. | `02-spec/07-design-system/readme.md`, `10-header-navigation.md`, `11-button-system.md`, `31-slide-controller-buttons.md`, `34-slide-layout-catalog.md`, `42-slide-step-and-sound-system.md`, `43-slide-webcam-overlay.md`, `44-slide-presenter-inspector-and-handouts.md` |

---

## Core Mandates Across All Prompts

1. **Strict Relative Git Paths:** All citations and markdown links MUST be strictly relative paths starting from the repository root (e.g., `02-spec/07-design-system/11-button-system.md`). Absolute paths and file-scheme URIs are totally banned.
2. **Pure DOM Typography (Zero Baked-In Text):** In websites, slides, and UI components, all text must be live DOM elements (`<h1>`, `<p>`, `<span>`) styled with CSS tokens.
3. **Strict Anonymization:** Never include private, client, or internal company names in public design systems or prompts. Always use generic placeholders.
4. **Closed Sets Rule:** Selectable options (variants, sizes, layouts, themes, hotkeys) must be declared as closed sets.
