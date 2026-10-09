# Subtask 02 — `gitmap autofix` CLI wiring, dispatch, help, skill surface (task 246)

## File box
- `cli/cmdautofix/fix.go` — create (CLI entry; ≤ ~300 lines).
- `cli/cmd/rootutility.go` — TWO surgical edits only: one `cmdautofix` import
  line + one dispatch entry in `utilityDesktopEntries()`. Nothing else.
- Do NOT touch `cli/cmd/llm/llm_train.go`, `.agents/skills/gitmap/SKILL.md`,
  `helptext/*.md`, or any readme/index — those are lead-owned (D13/D14).
- NOTE: despite this subtask's filename, the Go package is `cli/cmdautofix/`
  (`cli/cmdfix/` is taken by the git-state `fix` command — spec 01-overview D2).

## Context
Subtask 01 builds the engine + 8 categories. This subtask wires the CLI:
flags, dispatch, `--help`, exit codes, and the `SkillSnippetMD()` surface the
lead embeds into `gitmap llm train` and the Antigravity skill.
Design: spec `02-spec/21-app/246-parallel-fix-command/02-fix-command.md`
(§1, §4, §6–§8). Rules: `gitmap aum search` / `gitmap find` / `gitmap cat`
only for reads (TOTAL BAN on grep/rg/git grep/Select-String); relative paths;
`go build ./...` from `cli/` to verify — NEVER `go test` without the owner's
explicit command.

## `fix.go` checklist
- [ ] `func RunAutofixCmd(args []string) error` — entry point.
- [ ] Flag parsing with `flag.NewFlagSet("autofix", flag.ContinueOnError)`
      (same shape as `cli/cmd/llm/llm_train.go:113`):
      `--apply`, `--workers`/`-w` (default `runtime.NumCPU()`),
      `--ext`, `--category`, `--uri-pattern`, `--json`.
- [ ] Validate: workers ≥ 1; every `--category` name in the 8-category set;
      scan root readable. Any failure → usage error, exit code 2 semantics.
- [ ] Orchestrate the engine; print human report (default) or `--json` report;
      map outcome → exit codes per spec §4: 0 clean · 1 violations found /
      unfixable findings remain · 2 tool error.
      (Follow how sibling commands signal exit codes — search for the
      established pattern; do not invent a new one.)
- [ ] `func SkillSnippetMD() string` returning the markdown block from spec
      §7 VERBATIM (the lead wires it into llm train + SKILL.md).
- [ ] `--help` handling: build the help content as Go structs (program 243
      DRY rule — structs are the source of truth; `helptext/*.md` is generated
      by the lead, never hand-edited here).

## Dispatch registration (`cli/cmd/rootutility.go`)
In `utilityDesktopEntries()` — the file-utility family next to `cat`/`touch`/
`mkfile` (pattern at `cli/cmd/rootutility.go:477`) — add:

```go
{[]string{"autofix", "afx"}, func() error { return cmdautofix.RunAutofixCmd(argsTail()) }},
```

plus the `cmdautofix` import. Do NOT register in `rootcore.go` (that table owns
the git-state `fix`). Verify `afx` has no dispatch collision (spec D2 recorded
it free — re-confirm with `gitmap aum search` before committing).

## `--help` output shape (must render exactly this structure)

```
Usage: gitmap autofix [path] [flags]  (alias: gitmap afx)

Fix file hygiene issues in parallel — native port of the 03-ai-scripts autofixers.

Flags:
  --apply            Write fixes (default: dry-run audit only)
  -w, --workers N    Worker threads (default: CPU count)
  --ext .go,.md      Only scan these extensions (default: all text files)
  --category a,b     Subset of categories (default: all)
  --uri-pattern P    REPO_FILE_URI override for the paths category
  --json             Machine-readable report

Categories:
  encoding    UTF-8 no-BOM + LF normalization
  newlines    CRLF→LF, trim trailing whitespace, one final newline
  naming      boolean-comparison style audit (report-only)
  paths       Windows absolute path sanitizer (docs)
  gofmt       gofmt -w over .go files
  misspell    British→American spelling, case-preserving
  markdown    collapse 3+ blank lines in .md
  guidelines  composite: newlines + naming

Exit codes: 0 clean · 1 violations found · 2 tool error
```

## Deliverable
Reply with: the exact added lines (unified diff) for `cli/cmd/rootutility.go`,
the `go build ./...` exit code from `cli/`, the rendered `gitmap autofix --help`
output, and confirmation that `SkillSnippetMD()` returns the spec §7 text
verbatim. Then stop.
