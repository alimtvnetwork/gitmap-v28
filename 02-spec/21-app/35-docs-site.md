# Documentation Site — Interactive CLI Reference

## Overview

Replace the placeholder React frontend with a full documentation
site for gitmap. The site showcases all CLI commands, usage examples,
configuration guides, and architecture notes.

---

## Design Direction

- **Tone**: Developer-focused, clean, terminal-inspired aesthetic.
- **Typography**: Monospace display font for headings, clean sans for body.
- **Color**: Dark theme by default with a light toggle. Accent color
  derived from terminal green (`hsl(142, 71%, 45%)`).
- **Layout**: Sidebar navigation + main content area.

---

## Pages

### Home (`/`)

Hero section with:
- Project name and tagline
- Quick install command (copyable)
- Feature highlights (3–4 cards)
- Link to getting started

### Commands (`/commands`)

Searchable list of all CLI commands with:
- Command name and alias
- One-line description
- Expandable detail with flags and examples

### Getting Started (`/getting-started`)

Step-by-step guide:
1. Install gitmap
2. Run first scan
3. Clone repos
4. Set up shell alias for `gcd`

### Configuration (`/config`)

Documentation for:
- `config.json` structure
- `git-setup.json` options
- Profile management
- Environment variables

### Architecture (`/architecture`)

High-level overview:
- Project structure diagram
- Data flow (scan → store → clone)
- Database schema
- Output artifacts

---

## Technical Approach

- Built with React + Vite + Tailwind (existing stack).
- Content sourced from spec files, rendered as components.
- Responsive: mobile sidebar collapses to hamburger menu.
- Code blocks with syntax highlighting and copy button.
- Search across all commands using client-side filter.

---

## Component Structure

| Component | Purpose |
|-----------|---------|
| `Layout.tsx` | Sidebar + content shell |
| `Sidebar.tsx` | Navigation links |
| `CommandCard.tsx` | Expandable command reference |
| `CodeBlock.tsx` | Syntax-highlighted code with copy |
| `SearchBar.tsx` | Filter commands and content |
| `FeatureCard.tsx` | Home page feature highlight |
| `InstallBlock.tsx` | Copyable install command |

---

## Route Structure

| Route | Page |
|-------|------|
| `/` | Home |
| `/commands` | Command reference |
| `/getting-started` | Setup guide |
| `/config` | Configuration docs |
| `/architecture` | Architecture overview |

---

## SEO

- Title per page: `<Page> — gitmap docs`
- Meta descriptions per route
- Single H1 per page
- Semantic HTML (`nav`, `main`, `article`, `section`)
- JSON-LD for SoftwareApplication on home page

---

## Constraints

- All components under 200 lines.
- Use design tokens from `index.css` — no hardcoded colors.
- Dark/light mode via CSS variables.
- Mobile-first responsive design.
