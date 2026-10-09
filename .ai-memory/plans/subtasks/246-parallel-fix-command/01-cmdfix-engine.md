# Subtask 01 — `cli/cmdautofix` engine + 8 category ports (task 246)

## File box
`cli/cmdautofix/` ONLY (new package; create it). Do not touch any other file.
NOTE: despite this subtask's filename, the Go package is `cli/cmdautofix/`
(`cli/cmdfix/` is taken by the git-state `fix` command — spec 01-overview D2).

## Context
Port the 8 autofix scripts in `03-ai-scripts/` (read-only reference; NEVER
modify them) to native Go, behind one parallel engine. Full design: spec
`02-spec/21-app/246-parallel-fix-command/02-fix-command.md` (§2–§6).
Rules: reads only via `gitmap aum search` / `gitmap find` / `gitmap cat`
(TOTAL BAN on grep/rg/git grep/Select-String); relative paths everywhere;
every file ≤ ~300 lines; `go build ./...` from `cli/` to verify — NEVER
`go test` without the owner's explicit command.

## Build order
1. `engine.go` — walker, null-byte probe, worker pool, types, category registry.
2. The 8 category files (any order; each is independent).
3. `guidelines.go` last (composites `newlines` + `naming`).

## Engine (`engine.go`) checklist
- [ ] `Options` struct: `Apply bool`, `Workers int`, `Exts []string`,
      `Categories []string`, `URIPattern string`, `JSON bool`.
- [ ] `Violation{Path, Category, Line, Detail}` and
      `Category{Name, Exts, Check, Fix}` types per spec §6 (Path always
      relative, slash-separated).
- [ ] Walker: `filepath.WalkDir` from the scan root; skip `.git/` dirs;
      apply `--ext` filter; emit candidate paths on `fileCh`.
- [ ] Binary probe: first 8 KiB, `0x00` byte → skip silently (not a violation).
- [ ] Worker pool: N = `Options.Workers` (validated ≥ 1 by `fix.go`); each
      worker applies ALL selected categories to one file (`Check`, then `Fix`
      when `Apply` and `Fix != nil`); results on `resultCh`.
- [ ] Aggregator sorts by path → deterministic output under any schedule.
- [ ] UTF-8 validity gate lives in the `encoding` category (see below), not
      the engine — other categories operate on bytes.

## Byte-safety rules (apply to EVERY category file)
- Never write when transformed bytes equal the input bytes.
- Never write lossy bytes: undecodable input is a violation (dry-run) or a
  skip (`--apply`), never a replacement-character write.
- LF discipline on all read/write buffers (mirror `02-shared-engine.py`
  `read_file_lf`/`write_file_lf` semantics).
- Preserve original file mode on write (`os.Stat` → write with same mode).
- Zero nesting beyond one level; positive boolean names (CODE RED rules).

## Per-category checklists

### `encoding.go` — from `10-encoding-normalizer.py`
- [ ] Detect BOM (`EF BB BF`, `FF FE`, `FE FF`) at byte level; strip it.
- [ ] Normalize CRLF→LF as part of the fix (encoding implies newline form).
- [ ] `Check`: violation if BOM present OR bytes not valid UTF-8.
- [ ] `Fix`: return stripped/normalized bytes; if still not `utf8.Valid`,
      return the ORIGINAL bytes plus a violation (skip-write rule).

### `newlines.go` — from `04-newline-fixer.py`
- [ ] `Check`: violation on any CRLF, any trailing-whitespace line, or file
      not ending in exactly one `\n`.
- [ ] `Fix`: CRLF→LF, trim trailing whitespace per line, exactly one final
      newline. Empty file → stays empty (no newline invented).

### `naming.go` — from `08-naming-autofixer.py` — REPORT-ONLY (`Fix` = nil)
- [ ] Flag `== true`, `=== true`, `== True` (CODE RED boolean style).
- [ ] Skip comment lines: `//`, `#`, `/*`…`*/` line comments per file type —
      keep the script's line-based heuristic, do not build a parser.
- [ ] `Fix` is nil: even under `--apply` this category only reports (D11).
      Document why in a one-line comment.

### `paths.go` — from `07-relative-path-fixer.py`
- [ ] `const defaultRepoURI = "coding-guidelines"`; effective pattern =
      `Options.URIPattern` if set, else the constant (D10).
- [ ] Mirror script 07's transform exactly (read the script; do not redesign):
      sanitize Windows absolute paths containing the pattern in docs.
- [ ] `Check`/`Fix` pair; default behavior byte-identical to the script.

### `gofmt.go` — from `26-go-code-formatter.py`
- [ ] Exec `gofmt -w <file>` via `os/exec` (D8 — matches the script; NOT
      `go/format` stdlib, so output matches the repo toolchain).
- [ ] `gofmt` missing from PATH → return a tool error (engine surfaces it;
      overall exit becomes 2 per spec §4).
- [ ] Scope: `.go` files only. `Check`: run `gofmt -d` and treat any diff
      output as a violation; `Fix`: `gofmt -w`.

### `misspell.go` — from `27-misspell-auditor.py`
- [ ] Port the 30-word British→American dictionary verbatim from the script.
- [ ] CASE-PRESERVING replacement (D9 — fixes the script's known bug):
      `ALL-UPPER`→upper replacement; `Title`→title replacement;
      `lower`→dictionary form; anything else→dictionary (lowercase) form.
- [ ] Word-boundary matching (regexp `\b` equivalent); never touch the
      script's unimplemented "delegates to misspell binary" docstring —
      do NOT port it, do not shell out.
- [ ] `Check` lists each hit (line, word); `Fix` rewrites preserving case.

### `markdown.go` — from `31-md-gap-fixer.py`
- [ ] Scope: `.md` files only.
- [ ] `Check`: violation on any run of 3+ consecutive newlines.
- [ ] `Fix`: collapse 3+ newlines → exactly 2.

### `guidelines.go` — from `05-guideline-autofixer.py` — COMPOSITE
- [ ] Runs the `newlines` fix, then the `naming` check (D12).
- [ ] Obeys the unified dry-run/`--apply` rule — do NOT copy script 05's
      inverted default-fix behavior.

## Deliverable
Reply with: the 10 created file paths (relative), `go build ./...` exit code
from `cli/`, and one line per category stating its Check/Fix behavior. Then stop.
Do NOT wire dispatch, help, or skill surfaces — that is subtask 02.
