# 02 — Theme Design (from coding-guidelines-v24)

## 1. Source of truth

All palettes below are taken from the coding-guidelines-v24 design system —
`02-spec/07-design-system/16-theme-catalogue-and-palettes.md` and the
machine-readable `17-theme-tokens.json`. **No color value may be invented.**
Where the terminal cannot render a hex value, the spec mandates the nearest
256-color/16-color SGR code, documented per theme. Semantic tokens are bound
first; literal color classes in component code are forbidden (design-system
§7 checklist).

## 2. New named themes (terminal registry)

The `--theme` registry gains named themes. Legacy modes stay as compatibility
aliases: `bright` → `navy-purple`, `standard` → `vscode-dark-plus`,
`monochrome`/`mono` → SGR strip (unchanged behavior).

| Theme slug | Design-system source | Hex ground (bg / primary / accent) | Why chosen |
|------------|----------------------|------------------------------------|------------|
| `navy-purple` | Family 1: Navy Blue & Purple ("Midnight Cyber") | `#0A0F1D` / `#A855F7` / `#38BDF8` | Recommended for **dev tools**; default |
| `tokyo-night` | Family 2: Tokyo Night (Storm) | `#1F2335` / `#7AA2F7` / `#7DCFFF` | Recommended for **modern CLI tools** |
| `vscode-dark-plus` | Family 2: VS Code Dark+ | `#1E1E1E` / `#007ACC` / `#4EC9B0` | Instant familiarity for engineers |
| `monokai-pro` | Family 2: Monokai Pro | `#2D2A2E` / `#FF6188` / `#FFD866` | High-impact telemetry/status |

Terminal-role mapping per theme (from the catalogue's semantic tokens):

| Terminal role | Maps from | Used for |
|---------------|-----------|----------|
| `fg` | `--fg` | Body text, headings |
| `fg-muted` | `--fg-muted` | Metadata, secondary lines |
| `fg-subtle` | `--fg-subtle` | Labels, captions |
| `primary` | `--primary` | Success CTA, active indicator |
| `accent` | `--accent` | Data highlights, links |
| `warn` | heatmap scale Warning `#F59E0B` | Warnings |
| `error` | heatmap Critical `#DC2626` | Errors |
| `ok` | heatmap Optimal `#10B981` | Success states |

Heatmap work (telemetry, activity grids) uses ONLY the catalogue scales:
Sequential Plasma for intensity, diverging Red-Amber-Green for quality/risk,
GitHub-style density grid for activity (§4 of the catalogue). Random
decorative colors are banned.

## 3. Design-system rules honored

- **Single-accent Von Restorff:** exactly one prominent accent element per
  screen region; no double-shouting CTAs.
- **60/30/10 (terminal interpretation):** 60% base text, 30% structural
  surfaces/type, 10% accent.
- **Semantic decoupling:** help and UI code reference `primary`/`accent`/
  `fg-muted`, never raw hex or literal color classes.
- **WCAG 2.2 AA:** body text on background ≥ 4.5:1 (all four themes meet this
  per the catalogue's contrast ratings).

## 4. Data-driven registry

One Go registry (`cli/termout/themes.go` replacing per-mode branches):
`map[string]Theme{Slug, Roles: map[Role]SGR, ...}` — themes, role mappings,
and SGR codes live in ONE table. `--theme <slug>` is validated against the
registry keys; unknown values fail with the valid list (existing
`failInvalidThemeChoice` behavior preserved). The docs UI help screens render
from the same registry — never hand-maintained twice.

## 5. Dumb-terminal safety

Piped output (non-TTY) must not emit ANSI: keep the existing auto-mono
behavior — pipe detection forces SGR-strip regardless of theme. New themes
MUST pass the pipe test: `gitmap <cmd> | cat` output byte-identical with and
without `--theme <slug>`.
