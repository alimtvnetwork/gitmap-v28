# Plan 04: Templates (ignore / attributes) + Pretty Renderer

Spec: `spec/01-app/109-templates-ignore-attributes-pretty.md`
Memory: `mem://features/templates-ignore-attributes`

## Resolved open questions

| Question                                  | Decision                                      |
|-------------------------------------------|-----------------------------------------------|
| Alias names                               | `ai` / `aa` / `alfs` / `tl` / `ts`            |
| Real `add` subcommand router              | Yes — `gitmap/cmd/rootadd.go` + dispatchAdd   |
| Audit trail for curated common template   | Per-file `# source:` header + version int     |
| Read-only install fallback                | `~/.gitmap/templates/` overlay materialized   |
| Pretty renderer fixture corpus            | `gitmap/render/testdata/pretty/*.{in.md,want.txt}` |

## Phases

### Phase 0 — Scaffolding (THIS PR)

- [x] Spec, memory, plan written
- [ ] `gitmap/templates/embed.go` — `//go:embed assets/**`
- [ ] `gitmap/templates/resolver.go` — overlay > embed lookup
- [ ] `gitmap/templates/materialize.go` — first-run extract to `~/.gitmap/templates/`
- [ ] `gitmap/templates/paths.go` — user-templates dir per OS
- [ ] `gitmap/templates/assets/.keep` placeholder (Phase 1 fills it)
- [ ] `gitmap/constants/constants_templates.go` — kind/lang/marker constants
- [ ] No CLI wiring yet (next phase)

### Phase 1 — Seed corpus

- common + go + node + python + rust + csharp for both `ignore` and `attributes`
- `lfs/common.gitattributes` (binary patterns, NO `*.svg`)
- Each file has `# source: ... # version: 1` header

### Phase 2 — `add` router + `add ignore` ✅ (v3.105.0)

- [x] `gitmap/cmd/rootadd.go` with `dispatchAdd`
- [x] `gitmap/cmd/addignoreattrs.go` (shared engine for ignore + attributes)
- [x] `gitmap/templates/merge.go` marker-block engine (already shipped)
- [x] Single-pass dedupe with blank-line preservation
- [x] Sorted-tag invariant: `go,node` and `node,go` share `ignore/go+node`
- [x] Idempotence pinned by tests + by `templates.Merge`'s
      `bytes.Equal(prior, next)` short-circuit

### Phase 3 — `add attributes` + `add lfs-install` ✅ (v3.105.0)

- [x] `add attributes` shares `executeAddTemplate` with `add ignore`,
      parameterized via `addTemplateSpec` (kind / target / banner)
- [x] `add lfs-install` already shipped (separate command — runs
      `git lfs install --local` then merges `lfs/common.gitattributes`)

### Phase 4 — discovery commands

- `templates list` (groups by kind, shows overlay vs embed source)
- `templates show <kind> <lang>` (prints to stdout)

### Phase 5 — Pretty renderer

- `gitmap/render/pretty.go` with the 4 rules
- `gitmap/render/testdata/pretty/case-001..NNN`
- Table-driven test that loops fixtures

### Phase 6 — Wire pretty into CLI

- [x] Help output (`gitmap help <cmd>`, `--help`/`-h`) — `helptext.Print`
  routes through `render.RenderANSI` when stdout is a TTY. Opt-out via
  `GITMAP_NO_PRETTY=1`. Pipes / redirects keep raw markdown so editors
  and `less` see the source bytes unchanged. New `helptext.PrintRaw`
  exposed for callers that want to bypass pretty unconditionally.
- [x] Changelog display (`gitmap changelog`) — left on the existing
  structured renderer (`renderChangelogEntry`). It is already richer
  than the generic markdown pipeline (depth-aware bullets, hanging
  indent, ordered-list markers, wrap width, inline `**bold**` /
  `` `code` ``) and operates on `release.ChangelogEntry` structs, not
  raw markdown. Routing changelog through `render.RenderANSI` would be
  a regression. Decision: pretty renderer stays scoped to free-form
  markdown sources (helptext today; future doc commands).

## Non-goals

- No download-on-demand. Templates are embedded, audit-trailed, versioned.
- No yaml/toml config of templates. Filesystem overlay only.
- No editor integration. Files on disk only.
