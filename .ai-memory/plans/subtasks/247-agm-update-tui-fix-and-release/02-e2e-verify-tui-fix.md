# Subtask 02 — e2e-verify the agm update TUI fix (task 247)

## File box
`~/workspace/gitmap-e2e-247/` ONLY (outside the repo — nothing test-related is
ever committed). Repo reads read-only via gitmap tools.

## Phase A — author (do NOT run yet)
1. `~/workspace/gitmap-e2e-247/mock-install.sh`: staged fake installer — prints
   stage markers with sleeps (so progress is observable), supports `--fail` to exit
   1 after printing an error. Must be runnable via the implementation's test seam
   (env override documented in spec 01; if the worker implemented a different seam
   name, read `cli/cmdinstall/agm_update.go` to discover it — read-only).
2. `~/workspace/gitmap-e2e-247/e2e-247-agm-update.sh` (bash, `set -u -o pipefail`,
   `export TERM=xterm-256color`; `GITMAP_BIN` from environment): covers
   - U01: run update against mock installer, assert progress lines appear DURING
     the run (timestamped log), ends with success summary.
   - U02: `gitmap aum search` for `os.Stdout` assignments in
     `cli/cmdinstall/agm_update.go` + `installagmanager.go` update dispatch fns —
     expect zero hits.
   - U03: `--fail` mock → structured error panel, assert NO Go stack trace
     (`goroutine` string absent).
   - U04: assert no stray raw escape sequences outside renderer control in output.
   - U05: piped stdout (non-TTY) → plain lines, zero ANSI escapes.
   - U06: `gitmap agm update --help` exit 0.
   One `_record`-style PASS/FAIL line per ID + summary. Report the script paths
   and STOP — the lead provides the fixed binary for Phase B.

## Phase B — run (only when the lead sends the fixed binary path)
1. Run the script with the provided binary; capture to `run-247.log`.
2. Report per-ID PASS/FAIL with evidence lines.

## Rules
- Never write into the repo. Never run `go test`. No network calls in tests —
  the mock installer is local. The binary under test is the lead-provided one.
