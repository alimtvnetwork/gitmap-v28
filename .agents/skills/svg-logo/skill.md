---
name: svg-logo
description: "Executes the SVG Icon & Vector Graphic Creation - Design Workflow prompt. This prompt instructs Lovable (and similar web design AI environments) on the exact technical standards for generating scalable, responsive, and accessible SVG vector icons and graphics. Use when the user asks to run svg-logo, or the task is about logos, banners, thumbnails, or React UI design."
---

# SVG Icon & Vector Graphic Creation - Design Workflow

Source prompt: `01-prompts/10-ui-and-design/03-svg-logo.md`

## Instructions

1. Read `01-prompts/10-ui-and-design/03-svg-logo.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

This prompt instructs Lovable (and similar web design AI environments) on the exact technical standards for generating scalable, responsive, and accessible SVG vector icons and graphics.

## Core Architectural Invariants

- **100% Transparent Canvas (`alpha = 0`):** All SVGs MUST have a 100% transparent canvas. NEVER add a solid background, dark card backdrop, or background `<rect>`.
- **Dual-Mode Inverted Artwork Clarity:**
  - Dark Mode means inverted high-luminance strokes ($L > 0.70$, e.g. `#ffffff`, `#fcd34d`, `#38bdf8`) on a 100% transparent canvas—NEVER a dark box container (`<rect>`).
  - Light Mode uses rich, deep saturated hues ($L < 0.40$, e.g. `#0f172a`, `#1e1b4b`, `#1d4ed8`) on a 100% transparent canvas.
- **SVG Precision & Quiet Zone:**
  - Enforce an 8–10% inner quiet zone margin across all boundaries to prevent clipping at small favicon sizes (16px/32px) and app dock icons.
  - Round all coordinates, paths, and dimensions to a maximum of 2 decimal places.
  - Structure complex fills inside `<defs><linearGradient>` with unique lowercase IDs.
  - Enforce accessibility attributes (`<title>`, `<desc>`, `role="img"`, `aria-labelledby`).
  - Use `fill="currentColor"` or `stroke="currentColor"` for monochrome glyphs.
- **Design Token Alignment:** Map SVG variables and color styles directly to design system tokens (`--pres-accent`, `--pres-text`, `--primary`, `--foreground`).
