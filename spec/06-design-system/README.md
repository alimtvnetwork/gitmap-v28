# Design System Specifications

Complete design system reference for the gitmap-v27 documentation site.
All visual decisions are traceable to a single source of truth.

## File Index

| File | Purpose |
|------|---------|
| `01-colors-and-themes.md` | **Root reference** — all CSS variables, color tokens, light/dark values |
| `02-typography.md` | Font families, sizing, heading styles |
| `03-code-blocks.md` | Syntax highlighting, line features, code block layout |
| `04-component-patterns.md` | Reusable UI patterns, interactive states, spacing |
| `05-syntax-token-colors.md` | Highlight.js token color map |

## Architecture Rule

> `01-colors-and-themes.md` is the **root file**. Every other file in this
> folder references its tokens by CSS variable name. If a color changes,
> it changes in `01` and propagates everywhere automatically.
