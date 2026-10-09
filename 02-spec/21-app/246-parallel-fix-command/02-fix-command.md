# Spec 246.2 — `gitmap fix` command design (owner refinement 2026-10-09)

## §1 Command shape — parent with per-category subcommands

```
gitmap fix                                    → parent help (lists subcommands)
gitmap fix <category> [path] [--workers N] [-y|--yes] [--ext .go,.md] [--json] [--no-cache]
gitmap fix all        [path] [--workers N] [-y|--yes] [--ext .go,.md] [--json] [--no-cache]
```

Categories (one subcommand each):

| Subcommand | Source script | Scope |
|---|---|---|
| `encoding` | `10-encoding-normalizer.py` | all text files — normalize to UTF-8 no-BOM + LF |
| `newlines` | `04-newline-fixer.py` | all text files — CRLF→LF, trim trailing whitespace, one final newline |
| `naming` | `08-naming-autofixer.py` | code files — audit `== true` / `=== true` / `== True` (CODE RED). REPORT-ONLY always |
| `paths` | `07-relative-path-fixer.py` | docs — sanitize Windows absolute paths; `--uri-pattern` (default `coding-guidelines`, D10) |
| `gofmt` | `26-go-code-formatter.py` | `.go` — exec `gofmt -w` (D8); missing binary → tool error |
| `misspell` | `27-misspell-auditor.py` | text files — 30-word British→American dict, CASE-PRESERVING (D9) |
| `markdown` | `31-md-gap-fixer.py` | `.md` — collapse 3+ newlines → 2 |
| `guidelines` | `05-guideline-autofixer.py` | composite — runs `newlines` then `naming` (D12) |
| `all` | — | runs every category, one combined summary, one prompt |

Why `fix` (not `autofix`): OWNER DECISION 2026-10-09 — `fix` is the parent, full
stop. This displaces the existing git-state `fix` (`cli/cmd/rootcore.go:79` →
`cmdfix.RunFix`). Displacement plan (D2-revised):
- REMOVE `{[]string{"fix"}, ...}` from `coreBasicMaintenanceEntries()` in
  `cli/cmd/rootcore.go`. The `stash` / `wip` / `discard` entries stay UNTOUCHED —
  git-state remediation remains fully accessible through them.
- ADD `{[]string{"fix"}, func() error { return cmdautofix.RunFixCmd(argsTail()) }}`
  in `utilityDesktopEntries()` (`cli/cmd/rootutility.go`), next to the other file
  utilities. Package stays `cli/cmdautofix/` (internal name; command is `fix`).
- Unknown subcommand (e.g. old `gitmap fix ls`) → error listing valid
  subcommands PLUS the redirect line: `git-state remediation moved to
  'gitmap stash' / 'gitmap wip' / 'gitmap discard'`.

## §2 Check-driven behavior (owner: "fixes are check-driven, not blind")

Each subcommand runs its category's CHECK (the "gate") and fixes EXACTLY what the
check surfaces — no blind filesystem rewriting.

Default flow (no flags):
1. **Scan** — walk `[path]` (default `.`), run the check on every in-scope file
   (parallel workers; scan cache consulted).
2. **Summary** (MANDATORY) — printed always, even when zero findings:
   ```
   Scan complete in 1.24s — 1,203 files scanned (312 from cache).
   CATEGORY    FILES FLAGGED   FIXES
   encoding    12              12
   newlines    48              52
   TOTAL       60              64
   ```
   Fields: files scanned, files modified (0 before apply), per-category fix
   counts, time taken. `naming` findings appear under FIXES as `0 (report-only)`.
3. **Prompt** — `Apply these fixes? [y/N]` (default N; empty = N). `y`/`Y`/`yes`
   applies; anything else aborts with exit 0 and "No changes written."
4. **Apply** — rewrite only flagged files (compare-before-write, §5), print
   `Modified N files in M.MM s.`

`-y` / `--yes`: skip the prompt, apply immediately after the summary.
`--json`: machine-readable report (summary + violations) on stdout; prompt is
skipped under `--json` unless `-y` is also given (CI: `gitmap fix all --json`
exits 1 on findings, applies nothing).

`fix all`: runs all 8 checks (categories may run concurrently; files are walked
once), prints ONE combined per-category summary, ONE prompt, then applies.

`fix naming`: check-only category — prints findings, no prompt (nothing to apply).

Exit codes (CI contract, D6): `0` clean / nothing to do / applied cleanly;
`1` findings in check-only flow, or unfixable findings remain (e.g. `naming`
hits); `2` tool error (bad flag, unknown subcommand, unreadable path, `gofmt`
missing).

## §3 Parallelism + scan cache (owner: "parallelly ... changeable from the CLI")

Worker pool over FILES (D5), mirroring script 26's thread-pool:

```
walk(path, exts) → fileCh → N workers → resultCh → aggregator (sorts by path)
```

- `--workers N` / `-w N` sets N; default `runtime.NumCPU()`; N < 1 → usage error.
- Aggregator sorts by path → deterministic output regardless of scheduling.

Scan-result cache (owner: "caching ... so repeated runs are fast"):
- Location: `os.UserCacheDir()/gitmap/fix-scan-cache.json` (fallback: os.TempDir).
- Entry per absolute path: `{modtime_unix, size, categories: [...], violations: [...]}`.
- Hit rule: stat matches modtime+size AND requested categories ⊆ cached categories
  → reuse violations, skip re-check. Counted as "from cache" in the summary.
- Written back after every run. `--no-cache` bypasses read and write.

## §4 Flags

| Flag | Meaning |
|---|---|
| `[path]` | Positional scan root; default `.` |
| `-y`, `--yes` | Apply without prompting |
| `--workers N`, `-w N` | Worker threads; default `runtime.NumCPU()` |
| `--ext .go,.md` | Comma-separated extension filter; default: category scope |
| `--uri-pattern P` | `paths` category repo-URI pattern (default `coding-guidelines`) |
| `--json` | Machine-readable report; implies no prompt unless `-y` |
| `--no-cache` | Bypass the scan-result cache |

## §5 Byte-safety rules (D7 — non-negotiable)

1. **Binary probe**: first 8 KiB contains `0x00` → skip silently.
2. **UTF-8 validity** (`encoding` only): not valid UTF-8 after BOM/CRLF handling →
   report violation, SKIP under apply — never write lossy bytes.
3. **Compare before write**: write only if fixed bytes differ (no mtime touch).
4. **LF discipline**: LF-normalized buffers throughout.
5. **Permissions preserved**: write back with original mode.

## §6 Package layout — `cli/cmdautofix/` (command: `fix`)

One file per concern, each ≤ ~300 lines:

```
cli/cmdautofix/
  fix.go        RunFixCmd(args []string) — subcommand dispatch, parent help,
                prompt, summary rendering, exit codes, --json, SkillSnippetMD()
  engine.go     walker, null-byte probe, worker pool, scan cache,
                Violation type, Category registry (Check/Fix pairs)
  encoding.go   category: encoding   (from script 10)
  newlines.go   category: newlines   (from script 04)
  naming.go     category: naming     (from script 08, report-only)
  paths.go      category: paths      (from script 07, --uri-pattern)
  gofmt.go      category: gofmt      (from script 26, exec gofmt)
  misspell.go   category: misspell   (from script 27, case-preserving)
  markdown.go   category: markdown   (from script 31)
  guidelines.go category: guidelines (from script 05, composite)
```

```go
type Violation struct {
    Path     string // relative path, always slash-separated
    Category string
    Line     int    // 1-based; 0 = file-level
    Detail   string
}

type Category struct {
    Name string
    Exts []string // nil = all text files
    Check func(path string, src []byte, opts Options) []Violation
    Fix   func(path string, src []byte, opts Options) ([]byte, []Violation) // nil = report-only
}
```

Public entries (exact signatures — the LLM train calls these in-process):
- `func RunFixCmd(args []string) error` — full CLI (dispatch used by rootutility).
- `func SkillSnippetMD() string` — static markdown for train/skill wiring.

Dispatch changes (subtask 02 owns the exact edits):
1. `cli/cmd/rootcore.go`: DELETE `{[]string{"fix"}, func() error { return cmdfix.RunFix(argsTail(), "") }},`.
   `stash`/`wip`/`discard` lines stay byte-identical.
2. `cli/cmd/rootutility.go` `utilityDesktopEntries()`: ADD
   `{[]string{"fix"}, func() error { return cmdautofix.RunFixCmd(argsTail()) }},`
   plus the `cmdautofix` import.

## §7 LLM train + skill integration (D13)

- `cmdautofix.SkillSnippetMD() string` returns the static markdown block below.
- The train's heal phase (spec 03) calls `RunFixCmd` in-process for the content
  sub-step (dry-run; never `-y` from the train).

Snippet text (source of truth):

```markdown
### gitmap fix — parallel file-hygiene fixer (check-driven)
Check: `gitmap fix <category> [path]` — scans, prints summary, prompts
`Apply these fixes? [y/N]`. Categories: encoding, newlines, naming, paths,
gofmt, misspell, markdown, guidelines, all.
Apply without prompting: `gitmap fix all -y`.
Worker threads: `gitmap fix all -w 8` (default: CPU count).
Summary (always printed): files scanned, files modified, per-category fix
counts, time taken. Exit codes: 0 clean · 1 findings remain · 2 tool error.
Byte-safe: binaries skipped (null-byte probe); files rewritten only when
bytes differ; invalid-UTF-8 files are reported, never lossy-written.
`naming` (boolean-comparison style) is report-only by design.
Note: `gitmap fix` (git-state: stash/wip/discard) moved to
`gitmap stash` / `gitmap wip` / `gitmap discard`.
AI rule: run `gitmap fix all` (check) before committing hygiene-sensitive
work; never hand-roll sed/regex loops for encoding, newlines, or spelling.
```

## §8 Help (program 243 DRY rule, D14)

Help content lives in Go structs inside `cli/cmdautofix/` (source of truth);
`helptext/*.md` topics are GENERATED — never hand-edited.
`gitmap fix --help` renders parent help (subcommand table); 
`gitmap fix <category> --help` renders category help (exact shapes in subtask 02).
