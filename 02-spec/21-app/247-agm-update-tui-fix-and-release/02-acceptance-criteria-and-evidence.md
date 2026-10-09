# Task 247 — Acceptance Criteria & Evidence

Fix under test: `gitmap agm update` (`cli/cmdinstall/agm_update.go`) replaces its
broken output handling with a coordinated pterm-based GitLab-style progress UI fed
by captured (not inherited) child output, plus a non-TTY plain fallback.
Release: version `6.517.0`, annotated tag `v6.517.0`.

## Acceptance table

| ID | Criterion | PASS condition | Evidence artifact |
|----|-----------|----------------|-------------------|
| U01 | Linux update shows live progress (spinner + staged output), never minutes of silence | Mock-installer e2e: progress lines appear DURING the run (timestamps increasing), ending with a success summary | `~/workspace/gitmap-e2e-247/run-247-u01.log` |
| U02 | No stdio inheritance in the agm update paths | `gitmap aum search` for `os.Stdout` assignments in `cli/cmdinstall/agm_update.go` + `cli/cmdinstall/installagmanager.go` returns zero hits in the update dispatch functions | Search output captured in `~/workspace/gitmap-e2e-247/run-247-u02.log` |
| U03 | Failure shows a structured error panel, not a raw dump + stack trace | Mock-installer failure e2e: a bounded error section; NO Go stack trace anywhere in the output | `~/workspace/gitmap-e2e-247/run-247-u03.log` |
| U04 | No screen corruption from concurrent writers | Code assertion: all child output in the agm update path goes through pipes/capture into the single renderer (verified via `gitmap aum search` of the update dispatch); mock e2e output contains no stray raw escape sequences outside the renderer's control | `~/workspace/gitmap-e2e-247/run-247-u04.log` |
| U05 | Non-TTY fallback | Run with stdout NOT a terminal (piped): plain line output, zero ANSI escape sequences in the captured output | `~/workspace/gitmap-e2e-247/run-247-u05.log` |
| U06 | Real-path smoke: dispatch wiring unchanged | `gitmap agm update --help` exits 0 and shows the command with all existing flags intact | `~/workspace/gitmap-e2e-247/run-247-u06.log` |
| U07 | Release | `version.json` = `6.517.0`; annotated tag `v6.517.0` with message `Release v6.517.0` exists on origin; changelog entry present; `gitmap version` reports `6.517.0` | `~/workspace/gitmap-e2e-247/run-247-u07.log` |

## Evidence locations

All e2e scripts, mock installer, and logs live OUTSIDE the repo at
`~/workspace/gitmap-e2e-247/` (e.g. `e2e-247-u01.sh`, `mock-installer.sh`,
`run-247-u0x.log`). Nothing test-related is committed to the repo. The mock
installer script simulates a staged installer emitting progress lines and
exercising both success and failure paths.

## Out of scope

- `go test` never runs (standing rule) — verification is e2e-only via the built binary.
- `gitmap update` self-updater paths are a separate follow-up, not covered here.

## Failure policy

Any FAIL → fix is rejected, root-caused, and re-verified; the task never closes
with a failing acceptance ID. The release tag `v6.517.0` is cut ONLY after all
of U01–U06 pass.
