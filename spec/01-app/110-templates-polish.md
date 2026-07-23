# Templates Polish: Expanded Languages, `init`, and `diff`

Status: Spec (v3.108.0 — Phases 1+2+3 shipped; Phase 4 docs wiring + Phase 5 QA outstanding)
Owner: gitmap CLI
Related memory: `mem://features/templates-ignore-attributes`
Related plan: `.lovable/memory/plans/05-templates-polish-plan.md`
Builds on: `spec/01-app/109-templates-ignore-attributes-pretty.md`

## 1. Goals

1. **Broader language coverage** — extend the curated `ignore` /
   `attributes` corpora to cover Java, Ruby, PHP, Swift, and Kotlin in
   addition to the v3.105.0 set (common, go, node, python, rust, csharp).
2. **`gitmap templates init <lang>...`** — one-shot scaffold of
   `.gitignore` + `.gitattributes` (+ optional LFS attrs) for a chosen
   language stack, so a brand-new repo is set up with a single command.
3. **`gitmap templates diff --lang <name>`** — show what
   `add ignore <lang>` / `add attributes <lang>` *would* change in the
   current repo, without writing to disk. Standard `diff(1)` exit codes
   (`0` / `1` / `2`) so it slots into pre-commit hooks and CI.

All three pieces compose the existing primitives (`templates.Resolve`,
`templates.Merge`, `blockRegex`) — no new infrastructure, just corpus +
CLI plumbing + a tiny block-aware diff helper.

## 2. CLI Surface (delta from spec 109)

| Command | Alias | Action |
|---------|-------|--------|
| `gitmap templates init <lang>...` | `ti` | Scaffold `.gitignore` + `.gitattributes` (+ `--lfs`) for one or more languages |
| `gitmap templates diff --lang <l> [--kind k] [--cwd p]` | `td` | Print unified-style hunks for what `add` would change; never writes |

`list` / `show` from spec 109 also gain filter flags (`--kind`,
`--lang`) and a `--raw` / `--pretty` / `--no-pretty` triad documented
in spec 109 §4 (kept here for cross-reference).

### 2.1 Resolved open questions (from plan 05)

| Question | Decision |
|----------|----------|
| Stack presets for `init` | Positional language list: `templates init go node python`. Common always merged first. |
| `init` overwrite policy | Refuse if either target file already has a non-empty gitmap block, unless `--force`. `--dry-run` always allowed. |
| `diff` output format | Hand-rolled, block-scoped — flat `-` then `+`, banner per kind. No external diff dep, no Myers alignment. Block bodies are small enough that LCS would be false sophistication. |
| `diff` color | Reuse `render.HighlightQuotesANSI` token replacer. `+` cyan, `-` yellow, `@@` dim. TTY-aware via `render.StdoutIsTerminal()`. |
| LFS in `init` | Off by default; opt-in with `--lfs`. Calls the existing `add lfs-install` merge path. |
| Aliases | `ti` → `templates init`, `td` → `templates diff` |
| `diff` exit codes | Mirror `diff(1)`: `0` no change, `1` differences, `2` error. |

## 3. Language Corpus (Phase 1 — shipped v3.105.0)

The `assets/ignore/` and `assets/attributes/` directories ship 11
languages:

```
common  go  node  python  rust  csharp  java  ruby  php  swift  kotlin
```

Each file carries a `# source:` audit header and a `# version: <int>`
line per the spec 109 §4 contract.

### 3.1 Adding a new language

1. Drop `<lang>.gitignore` into `gitmap/templates/assets/ignore/`.
2. Drop `<lang>.gitattributes` into `gitmap/templates/assets/attributes/`.
3. Both files MUST start with the audit header
   `# source: <origin>` and `# version: 1`.
4. `corpus_test.go` enforces presence and non-emptiness automatically —
   no test code changes needed.

## 4. `templates init` (Phase 2 — shipped v3.105.0)

### 4.1 Behavior

- Resolves `common` + each requested language in argument order via
  `templates.Resolve(kind, lang)`.
- Merges into `<cwd>/.gitignore` and `<cwd>/.gitattributes` using the
  same `templates.Merge` engine as `add ignore` / `add attributes`.
- Sorted-tag invariant from spec 109 §5 still applies — order of
  language args does not change the resulting block tag.
- With `--lfs`, also merges `lfs/common.gitattributes` (same path as
  the standalone `add lfs-install` command).

### 4.2 Flags

| Flag | Default | Purpose |
|------|---------|---------|
| `--lfs` | off | Also merge LFS attributes |
| `--dry-run` | off | Preview every block, never touch disk |
| `--force` | off | Replace existing target files outright |

### 4.3 Exit codes

- `0` — every requested merge succeeded (or `--dry-run` ran clean).
- `1` — refused due to existing non-empty target without `--force`,
  or any merge step returned an error.

## 5. `templates diff` (Phase 3 — shipped v3.108.0)

### 5.1 Behavior

`templates diff` is the read-only counterpart to `add`. Given a
language and a kind, it loads the on-disk target (`./.gitignore` or
`./.gitattributes`) and reports the delta between its current
`# >>> gitmap:<kind>/<lang> >>>` block body and what `templates.Merge`
*would* write.

The diff is **block-scoped**: only the gitmap-managed marker block
participates. Hand edits OUTSIDE the block are invisible — same
contract as `add` itself, so the diff's silence is meaningful.

### 5.2 Status enum

`gitmap/templates/diff.go` exposes a `DiffStatus` enum that drives
both the CLI exit code and the printed hunks:

| Status | Trigger | Hunks |
|--------|---------|-------|
| `DiffNoChange` | block exists, body matches template byte-for-byte | none |
| `DiffMissingFile` | target file absent | every template line as `+` |
| `DiffMissingBlock` | file exists, no gitmap block for this tag | every template line as `+` |
| `DiffBlockChanged` | block exists, body differs | flat `-`-then-`+` paired hunk |

### 5.3 Flags

| Flag | Required | Default | Purpose |
|------|----------|---------|---------|
| `--lang <name>` | yes | — | Language to diff |
| `--kind <ignore\|attributes>` | no | both | Restrict to one kind |
| `--cwd <path>` | no | `.` | Repo directory |

### 5.4 Exit codes (mirror `diff(1)`)

| Code | Meaning |
|------|---------|
| `0` | No changes — script can short-circuit |
| `1` | Differences found — `add` would change something |
| `2` | Error — bad flag value, unknown language, or I/O failure |

### 5.5 Output format

```
@@ gitmap:ignore/node @@
-*.log
+*.log
+node_modules/
```

Banner per kind, then `-` lines for current content, then `+` lines
for would-be content. TTY-aware coloring via
`render.HighlightQuotesANSI`: cyan `+`, yellow `-`, dim `@@`. Pipes
and redirects stay raw so downstream tools (`grep`, `wc`, custom
linters) can re-parse cleanly.

### 5.6 Reuses, never re-implements

- `extractBlockBody` reuses `merge.go`'s `blockRegex(tag)` — the
  parser CANNOT drift from the writer. Adding a new marker shape
  forces a single-file change.
- `splitDiffLines` deliberately preserves intra-body blank lines so
  visual separators in templates show up as `+` / `-` rather than
  silently collapsing.
- `Diff` is pure — it never writes — so a future `add --dry-run` can
  reuse the same computation without touching disk.

## 6. Non-goals (deferred)

- IDE-specific templates (VSCode workspace, IntelliJ `.idea/`) — out
  of scope; revisit if a user asks.
- Template version-upgrade prompts when `# version:` headers bump —
  needs a dedicated migration story; deferred.
- Custom user-defined templates discoverable from
  `~/.gitmap/templates/custom/<lang>.gitignore` — interesting, but
  the overlay already supports byte-for-byte replacement; defer until
  someone asks for additive composition.
- Myers / LCS alignment in `diff` — block bodies are typically <30
  lines; a flat `-`-then-`+` is honest and avoids a dependency.

## 7. Test matrix

| Surface | Test | File |
|---------|------|------|
| Corpus headers + non-empty | `TestCorpusHasHeader`, `TestCorpusNonEmpty` | `gitmap/templates/corpus_test.go` |
| Sorted-tag invariant | `TestSortedTag*` | `gitmap/cmd/addignoreattrs_test.go` |
| `list` filters | `TestFilterTemplates*` | `gitmap/cmd/templatescli_filter_test.go` |
| Pretty renderer | 9 fixtures + ANSI swap + unterminated quote | `gitmap/render/pretty_test.go` |
| `Diff` status branches | `TestDiffMissingFile/MissingBlock/NoChange/BlockChanged` | `gitmap/templates/diff_test.go` |
| Blank-line preservation | `TestDiffPreservesBlankLines` | `gitmap/templates/diff_test.go` |

## 8. Sequencing recap

| Release | Scope |
|---------|-------|
| v3.105.0 | Phase 1 (11-lang corpus) + Phase 2 (`templates init`) |
| v3.106.0 | `templates list --kind/--lang` filters |
| v3.107.0 | Pretty-renderer corpus rounded to 9 fixtures |
| v3.108.0 | Phase 3 (`templates diff` + alias `td`) |

Phase 4 (changelog.ts entry, docs site sidebar, completion-generator
visibility for `td`) and Phase 5 (full local `go test` + lint sweep)
remain — see plan 05 for the open boxes.
