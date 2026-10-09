# Subtask 01 — implement the agm update TUI fix (task 247)

## File box
`cli/cmdinstall/agm_update.go`, `cli/cmdinstall/installagmanager.go`, + ONE new
helper file/package under `cli/cmdinstall/`. Do not touch any other file — in
particular NOT `cli/cmdagent/agent_subtask.go` (another workstream), NOT
`cli/cmdupdate/`, `cli/cmddownload/`, `cli/glyphs/`, NOT version files.

## What to build (per `02-spec/21-app/247-agm-update-tui-fix-and-release/01-rca-scope-and-plan.md`)
1. Replace `runAgmUpdateLinuxQuiet`'s `CombinedOutput()` silence with a live
   pterm-based progress UI: spinner + streaming installer-output tail + stage
   checkmarks (GitLab-style). pterm v0.12.83 is already in go.mod; follow the
   existing `cli/cloner/batchprogress.go` / `cli/cluster/pool.go` patterns.
2. Replace raw stdio inheritance in `dispatchAgManagerWindowsWithVersion`
   (`installagmanager.go:163`) and `dispatchAgManagerUnixWithVersion` (`:173-175`)
   with piped capture feeding the same single renderer — one mutex-guarded writer
   owns the screen; nothing writes to os.Stdout directly mid-run.
3. Failure path: structured error panel (bounded, last N lines of installer output),
   NO Go stack trace dump.
4. Non-TTY fallback: when stdout is not a terminal, plain line output, no escape
   sequences.
5. Test seam: env override for the installer source (e.g. `GITMAP_AGM_INSTALL_SCRIPT`
   pointing at a local script, or URL override) — default behavior unchanged when
   unset. The e2e worker drives the update against a mock installer through this.
6. `cli/glyphs/` stays untouched; render through the existing pipeline.

## Rules
- GitMap tools only for reads; total ban on rg/grep/git grep.
- Small files, zero nesting, positive booleans (CODE RED). Keep new helper ≤300 lines.
- `go build ./...` from `cli/` to verify. NEVER `go test`.

## Deliverable
The unified diff of your changes, `go build` exit code, built binary path. Then stop.
