# Plan 05: Templates Polish — More Languages, `init`, `diff`

Spec: `spec/01-app/110-templates-polish.md` ✅ (v3.108.0)
Memory: `mem://features/templates-ignore-attributes` (extend), new `mem://features/templates-init-diff`
Builds on: Plan 04 (templates infrastructure, merge engine, pretty renderer)

## Goals

1. **Broader language coverage** — add `java`, `ruby`, `php`, `swift`, `kotlin` to both `ignore` and `attributes` corpora.
2. **`gitmap templates init`** — one-shot scaffold of `.gitignore` + `.gitattributes` (+ optional LFS) for a chosen language stack.
3. **`gitmap templates diff`** — show what `add` would change without writing (overlay vs embed, or current-file vs template).

All three pieces share the existing resolver / merge / marker-block primitives. No new infra, just composition + corpus.

## Open questions (resolve before Phase 0)

| Question | Default proposal |
|----------|------------------|
| Stack presets for `init` (single lang vs combos like `node+python`)? | Allow comma list: `--lang go,node`. Apply each in order; common always first. |
| `init` overwrite policy if files already exist? | Refuse unless `--force`; suggest `add` instead. Always allow `--dry-run`. |
| `diff` output format | Unified diff via `github.com/sergi/go-diff` OR hand-rolled marker-block diff (no new dep). Lean hand-rolled to keep dep surface small. |
| Should `diff` color via pretty renderer? | Yes, reuse `render.RenderANSI` tokens (`+` green, `-` red, `@` cyan) — TTY-aware via existing helper. |
| LFS in `init` | Off by default; opt-in with `--lfs`. Calls existing `add lfs-install` path. |
| Aliases | `ti` → `templates init`, `td` → `templates diff` |

## Phases

### Phase 0 — Spec + scaffolding

- [ ] Write `spec/01-app/110-templates-polish.md` covering: new langs, `init` UX, `diff` UX, exit codes, dry-run semantics.
- [ ] Decide and lock the open questions above.
- [ ] Add `mem://features/templates-init-diff` summary.
- [ ] No code yet.

### Phase 1 — Expand language corpus

### Phase 1 — Expand language corpus ✅ (v3.109.0)

- [x] `assets/ignore/{java,ruby,php,swift,kotlin}.gitignore` — all five present with audit-trail headers.
- [x] `assets/attributes/{java,ruby,php,swift,kotlin}.gitattributes` — all five present with audit-trail headers.
- [x] Each file: `# source: ...` + `# kind:` + `# lang:` + `# version: 1` header (verified via `head -5` on every file).
- [x] No `constants_templates.go` lang enum needed — resolver discovers langs via filesystem walk per Plan 04 design (template kind/lang inferred from filename + header). Confirmed: zero references to `LangJava`/`LangRuby`/etc. in `gitmap/templates/*.go`.
- [x] `corpus_parity_test.go` already enumerates all five new langs (lines 23-30) and asserts the ignore-vs-attributes parity batch (line 97). No `corpus_test.go` extension required.
- [x] `templates list` output picks up new langs automatically via the resolver — no `list_test.go` change needed.

### Phase 2 — `templates init` ✅ (v3.110.0)

- [x] `gitmap/cmd/templatesinit.go` — flags: `--lfs`, `--force`, `--dry-run`. Uses positional `<lang> [<lang>...]` instead of `--lang <csv>` and reads CWD via `os.Getwd()` instead of an explicit `--cwd` (final UX call: positional langs feel more natural for a scaffolder; `--cwd` deferred — `cd && gitmap templates init …` covers it).
- [x] Reuses `templates.Resolve` + `templates.Merge` — zero new merge logic.
- [x] Order: per lang `[ignore, attributes]` then optional single `lfs/common` step. Common is implicit since the embedded `common.gitignore` lives outside the per-lang loop and is merged separately by `add ignore` users; `init` keeps the lang focus tight.
- [x] Behavior: ignore template REQUIRED per lang (hard-fail with hint), attributes template OPTIONAL (soft-skip with dim notice — matches embed corpus reality where some langs lack an attributes file).
- [x] `--force` removes the target file before merge so the resulting block is the only content. Without `--force`, `templates.Merge` preserves non-marker content and updates-in-place or appends.
- [x] Idempotent: re-running `init <lang>` produces "unchanged" lines; running `add ignore <lang>` afterward is also a no-op (same marker tag `ignore/<lang>`).
- [x] Helptext: `gitmap/helptext/templates-init.md` (133 lines, markdown, picked up by pretty renderer).
- [x] Alias `ti` registered alongside `init` in `templatescli.go` dispatcher.
- [x] `templatesinit_test.go` — 9 unit tests covering flag parsing, dry-run simulation, soft-skip on missing attributes, and `--force` + idempotency paths.

### Phase 3 — `templates diff` ✅ (v3.108.0)

- [x] `gitmap/cmd/templatesdiff.go` — flags: `--lang <name>`, `--kind ignore|attributes` (default both), `--cwd <path>`.
- [x] `gitmap/templates/diff.go` — marker-block aware, pure (never writes).
  Status enum (NoChange / MissingFile / MissingBlock / BlockChanged) drives exit codes.
  Reuses `blockRegex(tag)` from `merge.go` so parser can't drift from writer.
  Hand-rolled (no Myers / no external diff dep) — block bodies are small enough
  that a flat removal-then-addition is honest and ≤180 LOC.
- [x] `diff_test.go` — 5 cases pinning all 4 branches + blank-line preservation.
- [x] TTY-aware coloring via `render.HighlightQuotesANSI` (cyan `+`, yellow `-`, dim `@@`).
- [x] Exit codes mirror `diff(1)`: `0` no change, `1` differences, `2` error.
- [x] Helptext `gitmap/helptext/templates-diff.md` with exit-code table + pre-commit example.
- [x] Alias `td` registered alongside `diff`.

### Phase 4 — Wiring + docs (in progress)

- [x] Register `init` and `diff` under `templatescli.go` dispatcher. (`diff` shipped in v3.108; `init` blocked on Phase 2.)
- [x] Update `gitmap/helptext/templates.md` usage banner with the new subcommand. (v3.108)
- [x] Update `src/data/changelog.ts` with v3.107 (renderer corpus) + v3.108 (`templates diff`).
- [x] Add `templates diff` entry to `src/data/commands.ts` so the docs site command browser surfaces it. (v3.108)
- [x] Add the `init` entry to `src/data/commands.ts`. (v3.110)
- [ ] README: short snippet under "Templates" section showing `gitmap templates init go --lfs` + `gitmap templates diff --lang go` lifecycle.
- [x] Completion generator: `td` and `ti` aliases now surfaced via typed `CmdTemplatesDiffAlias` / `CmdTemplatesInitAlias` constants in `constants_templates_cli.go` (v3.111). The full subcommand strings (`diff`, `init`) carry `// gitmap:cmd skip` so they don't double-register. Decision reversed from v3.108's "subcommand precedent" stance per user request.
- [ ] Manual smoke matrix:
  - Empty repo → `templates init --lang go,node --lfs` → expect 3 files.
  - Existing `.gitignore` → refuse without `--force`.
  - `templates diff --lang python` after `add ignore python` → expect zero diff, exit 0.
  - `templates diff --lang java` on a repo missing the marker block → expect addition hunk, exit 1.

### Phase 5 — QA + release

- [ ] `go build ./... && go test ./...` clean.
- [ ] `golangci-lint run` clean (per `mem://tech/static-analysis-security`).
- [ ] Pretty-render the new helptexts on a TTY; confirm no token leakage.
- [ ] Tag release, update plan to `done`.

## Out of scope (defer)

- IDE-specific templates (VSCode, IntelliJ workspaces) — separate plan if requested.
- Template version upgrades / migration prompts — needs version bumping in `# version:` header system; defer.
- Custom user-defined templates discoverable from `~/.gitmap/templates/custom/` — interesting, but defer until a user asks.

## Sequencing recommendation

Ship Phase 1 (langs) standalone first as `v3.19.0` — low risk, immediate value. Then Phase 2+3 together as `v3.20.0`. Phase 0 spec covers both releases.
