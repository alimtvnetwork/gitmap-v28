# Subtask 02 — Design-system theme registry integration (program 273, WS2)

## File box

- `cli/termout/themes.go` — NEW (≤300 lines): `Theme{Slug, Roles}` table,
  `map[string]Theme` registry, semantic Role constants.
- `cli/cmd/themeflag.go` — surgical: validate `--theme` against registry keys
  (keep `failInvalidThemeChoice` semantics, list new slugs).
- `cli/termout/theme_theme.go`, `filter.go` — rewire to the registry
  (keep ModeMono SGR-strip and pipe auto-mono behavior byte-identical).
- Touch nothing else.

## Context

Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/02-theme-design.md`.
Source of truth: coding-guidelines-v24
`02-spec/07-design-system/16-theme-catalogue-and-palettes.md` +
`17-theme-tokens.json`. New slugs: `navy-purple` (default), `tokyo-night`,
`vscode-dark-plus`, `monokai-pro`. Legacy aliases: `bright`→`navy-purple`,
`standard`→`vscode-dark-plus`, `monochrome`/`mono`→SGR strip (unchanged).
Roles per theme: `fg`, `fg-muted`, `fg-subtle`, `primary`, `accent`, `warn`,
`error`, `ok` — semantic only, never literal colors in code. Honor
Von Restorff single-accent + WCAG AA. Rules: gitmap search tools only;
relative paths; ≤300 lines; `go build ./...` from `cli/` — NEVER `go test`.

## Checklist

- [ ] Theme table with exact hex ground from the catalogue; SGR codes
  documented per role (nearest-code mapping recorded in code comments).
- [ ] Heatmap roles use ONLY catalogue scales (Plasma / RAG / density grid).
- [ ] `--theme` rejects unknown slugs with the valid list.
- [ ] `gitmap <cmd> | cat` byte-identical across all themes (pipe test).
- [ ] Docs UI help screens render from the same registry.
- [ ] Commit: `gitmap cpf "theme - design-system theme registry"` + push
  (verify `git ls-remote`).
