# Spec 246.2 — `gitmap autofix` command design

## §1 Command shape

```
gitmap autofix [path] [--apply] [--workers N | -w N] [--ext .go,.md]
               [--category encoding,newlines] [--uri-pattern P] [--json]
```

Alias: `gitmap afx` (identical behavior; aliases are a kept feature).

Why `autofix` and not `fix`: `gitmap fix` already exists — `cli/cmd/rootcore.go:79`
dispatches it to `cmdfix.RunFix`, the git-state remediation command
(stash/wip/discard flows), implemented in `cli/cmdfix/`. Reusing the name would
collide on dispatch, on the package name, and in user muscle memory. `autofix` /
`afx` were verified free of dispatch entries and package collisions (D2).

| Flag | Meaning |
|---|---|
| `[path]` | Positional scan root; default `.` |
| `--apply` | Write fixes. Absent = dry-run audit only |
| `--workers N`, `-w N` | Worker threads; default `runtime.NumCPU()` |
| `--ext .go,.md` | Comma-separated extension filter; default: all text files |
| `--category a,b` | Comma-separated subset of §2; default: all. Unknown name → exit 2 |
| `--uri-pattern P` | Overrides the `paths` category's repo-URI pattern (D10) |
| `--json` | Machine-readable report on stdout |

## §2 Categories (the owner's grouping)

| Category | Source script | Scope | Behavior |
|---|---|---|---|
| `encoding` | `10-encoding-normalizer.py` | all text files | Normalize to UTF-8 no-BOM + LF. Detects BOM/CRLF at byte level |
| `newlines` | `04-newline-fixer.py` | all text files | CRLF→LF, trim trailing whitespace per line, exactly one final newline |
| `naming` | `08-naming-autofixer.py` | code files (`.go .py .js .ts .java` …) | Audit `== true` / `=== true` / `== True` (CODE RED style), skipping comment lines. REPORT-ONLY, even under `--apply` (D11) |
| `paths` | `07-relative-path-fixer.py` | docs (`.md .txt` …) | Sanitize Windows absolute paths; pattern = `--uri-pattern` (default `coding-guidelines`, D10) |
| `gofmt` | `26-go-code-formatter.py` | `.go` | Exec `gofmt -w` (D8); missing binary → tool error |
| `misspell` | `27-misspell-auditor.py` | text files | 30-word British→American dictionary, CASE-PRESERVING replacement (D9) |
| `markdown` | `31-md-gap-fixer.py` | `.md` | Collapse 3+ consecutive newlines → 2 |
| `guidelines` | `05-guideline-autofixer.py` | composite | Runs `newlines` fix then `naming` check (D12) |

## §3 Parallelism

Worker pool over FILES (D5), mirroring script 26's thread-pool design:

```
walk(path, exts) → fileCh → N workers → resultCh → aggregator
```

- Each worker takes one file, applies every selected category to it
  (`Check`; then `Fix` when `--apply`), and emits its violations/fix records.
- The aggregator sorts results by path so terminal and `--json` output are
  deterministic regardless of worker scheduling.
- `--workers` / `-w` sets N; default `runtime.NumCPU()`; N < 1 → usage error.

## §4 Exit codes (CI contract, D6)

- `0` — clean: no violations (dry-run) or everything fixed (`--apply`).
- `1` — violations found in dry-run, or unfixable findings remain after
  `--apply` (e.g. `naming` report-only hits).
- `2` — tool error: bad flag, unknown category, unreadable path, `gofmt`
  binary missing.

Dry-run with violations prints the violation list and exits 1 — this is the
CI gate (`gitmap autofix` with no flags fails the build on dirty hygiene).

## §5 Byte-safety rules (D7 — non-negotiable)

1. **Binary probe**: read the first 8 KiB; if it contains a `0x00` byte, skip
   the file silently. Never a violation, never written.
2. **UTF-8 validity**: `encoding` category only. If bytes are not valid UTF-8
   after BOM/CRLF handling, report a violation in dry-run and SKIP the file
   under `--apply` — never write lossy/replacement bytes.
3. **Compare before write**: compute the fixed bytes; write only if they differ
   from the original. No-op rewrites must not touch mtime.
4. **LF discipline**: all reads/writes go through LF-normalized buffers
   (same semantics as `02-shared-engine.py`'s `read_file_lf`/`write_file_lf`).
5. **Permissions preserved**: `os.Stat` the original mode; write back with it.

## §6 Package layout — new package `cli/cmdautofix/`

One file per concern, each ≤ ~300 lines (standing rule):

```
cli/cmdautofix/
  fix.go        RunAutofixCmd(args []string) — flag parsing, orchestration,
                exit codes, --json report, SkillSnippetMD()
  engine.go     walker, null-byte probe, worker pool, Violation type,
                Category registry (Check/Fix function pairs)
  encoding.go   category: encoding   (from script 10)
  newlines.go   category: newlines   (from script 04)
  naming.go     category: naming     (from script 08, report-only)
  paths.go      category: paths      (from script 07, --uri-pattern)
  gofmt.go      category: gofmt      (from script 26, exec gofmt)
  misspell.go   category: misspell   (from script 27, case-preserving)
  markdown.go   category: markdown   (from script 31)
  guidelines.go category: guidelines (from script 05, composite)
```

Engine core types (implementer may refine names, not semantics):

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
    // Check returns violations without touching disk.
    Check func(path string, src []byte, opts Options) []Violation
    // Fix returns the rewritten bytes; nil = report-only category.
    Fix func(path string, src []byte, opts Options) ([]byte, []Violation)
}
```

Dispatch registration (subtask 02 owns the exact edit): one entry in
`utilityDesktopEntries()` in `cli/cmd/rootutility.go`, next to the other
file utilities (`cat`/`touch`/`mkfile`), following the existing pattern
(`cli/cmd/rootutility.go:477`):

```go
{[]string{"autofix", "afx"}, func() error { return cmdautofix.RunAutofixCmd(argsTail()) }},
```

plus the `cmdautofix` import. NOT in `rootcore.go` — that table owns the
git-state `fix`.

## §7 LLM train + skill integration (D13)

`gitmap llm train` must teach AI agents the autofix workflow, and the
Antigravity skill must document it. Boundary: the 246 subtasks expose the
content; the lead wires it into the shared surfaces.

- `cmdautofix.SkillSnippetMD() string` returns the static markdown block below.
- Lead wiring (NOT in 246 file boxes): embed the snippet in `gitmap llm train`
  output (`cli/cmd/llm/llm_train.go`, as a new curriculum stage/section) and
  append it to `.agents/skills/gitmap/SKILL.md`.

Snippet text (source of truth — lead copies verbatim):

```markdown
### gitmap autofix — parallel file-hygiene fixer
Dry-run audit (CI gate): `gitmap autofix` — exits 1 on violations.
Apply fixes: `gitmap autofix --apply`.
Worker threads: `gitmap autofix --apply -w 8` (default: CPU count).
Categories (grouping): encoding, newlines, naming, paths, gofmt,
misspell, markdown, guidelines (composite). Select a subset:
`gitmap autofix --category encoding,newlines`.
Scope: `gitmap autofix [path] [--ext .go,.md] [--json]`.
Exit codes: 0 clean · 1 violations found · 2 tool error.
Byte-safe: binaries skipped (null-byte probe); files rewritten only
when bytes differ; invalid-UTF-8 files are reported, never lossy-written.
`naming` (boolean-comparison style) is report-only by design.
AI rule: run `gitmap autofix` (dry-run) before committing hygiene-sensitive
work; never hand-roll sed/regex loops for encoding, newlines, or spelling.
```

Recursive directive for the train (lead's wording, same intent): an AI that has
run `gitmap llm train` must prefer `gitmap autofix` over ad-hoc scripts for
file hygiene, and must teach the next model the same.

## §8 Help (program 243 DRY rule, D14)

Help content lives in Go structs inside `cli/cmdautofix/` (source of truth);
`helptext/*.md` topics are GENERATED from them — never hand-edited.
`gitmap autofix --help` renders the usage, flags, category table, and exit codes
(exact shape defined in subtask 02). `gitmap help autofix` resolves to the same
topic.
