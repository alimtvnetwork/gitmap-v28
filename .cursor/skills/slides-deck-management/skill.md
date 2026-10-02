---
name: slides-deck-management
description: Author, update, configure, and verify interactive presentation slides and deck registries in slides-app.
---

# Slides Deck Management & Presentation Authoring

This skill provides step-by-step guidance for managing, creating, and modifying interactive presentation slides within `slides-app/`.

## Architecture Overview

- **Deck Registry (`slides-app/src/deck/registry.ts`):** Canonical manifest mapping numeric sequence prefixes (`00-title`, `01-table-of-contents`, etc.) to React slide components and metadata (`SlideSection`, `severity`).
- **Slide Components (`slides-app/src/slides/`):** Individual slide components following React functional component conventions, Tailwind CSS utility classes, and Lucide icons.
- **Section Grouping (`groupBySection`):** Dynamically filters and groups slides by logical domain (e.g. Core Guidelines, Error Management, Database Conventions, React Guidelines).

## Creating or Modifying a Slide

1. **Component Creation:**
   - Create the slide component in `slides-app/src/slides/<NN>-<slug>.tsx`.
   - Keep the file bounded (<100 lines per component).
   - Use declarative UI, Lucide icons for visual anchors, and highlight good vs bad patterns.

2. **Registry Registration:**
   - Import the new slide component in `slides-app/src/deck/registry.ts`.
   - Register it in the `DECK` array with its title, section, and path.

3. **Verification:**
   - Run Vite build to ensure TypeScript types and exports resolve cleanly:
     ```bash
     npm --prefix slides-app run build
     ```
