# 02 — Shared implementation patterns

**Spec:** 275 · **Updated:** 2026-10-17

Every new command MUST follow these patterns. They exist to keep 30
commands consistent and reviewable.

---

## 1. Package layout

New packages: `cli/cmd<name>/` (lowercase, no underscores).
Existing packages are extended in place (e.g. `cli/cmdmerge/` for the mr
family — do NOT create `cli/cmdmr/`).

Per file: one concern, ≤300 lines. Typical split:
- `<cmd>.go` — flag parsing + flow orchestration
- `<cmd>_exec.go` — external calls (gh, git, ssh)
- `<cmd>_help.go` — help text + helpdoc example entry data

## 2. Dispatcher registration

Register in `cli/cmd/rootdata.go` following the existing table:
```go
{[]string{"mr list", "mrl"}, func() error { return cmdmerge.RunMrList(argsTail()) }},
```
Constants for command names go in `cli/constants/` (see
`constants_movemerge.go` pattern). Aliases are a kept feature — register
every alias from the proposal.

## 3. Help text

Every command implements `--help` printing:
- one-line description
- `Usage:` line with the exact invocation shape
- flags list
- one concrete example

## 4. Help-example registry (3-part format)

Add an entry to the `cli/helpdoc/examples*.go` registry for every new
command (see program 273: `MergeWriteHelpExamples` in
`cli/helpdoc/examples_merge.go`):
- WHAT IT DOES — one paragraph
- INVOCATION — concrete command line
- EXPECTED OUTPUT — realistic sample output

This wires the command into `gitmap help example <cmd>` automatically.

## 5. Dry-run discipline

Commands with side effects (mr merge, repo import, pe rerun, hooks
install, release train) MUST support `--dry-run` printing exactly what
would happen. Default-safe: `release train` runs in dry-run unless
`--yes` is passed.

## 6. Read-only commands

Purely informational commands (list, view, triage, health, verify,
audit, drift, diff, report, progress, version, topic, howto, preview)
make no mutations and need no dry-run flag.

## 7. Error handling

- Return `*appfault.AppError`, never bare `error`.
- Positive booleans only (`is*`, `has*` prefixes).
- Missing required flags: print usage + the missing flag name, exit non-zero.
- External tool missing (`gh` not installed): clear message naming the tool.

## 8. Output

- Human table on TTY, matching existing `pe`/`summary` style.
- No ANSI codes when piped (dumb-terminal safe, per program 273).
