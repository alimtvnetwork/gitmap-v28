# Subtask 265-stdout-ci-gates · 04 — Deprecation Policy

Scope: implementer of the command-deprecation policy + the `fix` unknown-subcommand redirect for task 265.
Spec home: `02-spec/21-app/265-stdout-ci-gates/04-deprecation-and-pe-parallel.md` Part A (read it first).
Rules: files ≤300 lines, lowercase filenames, positive-prefixed booleans (`hasX`, `isX`), `*apperror.AppError` for errors, relative paths only. Do NOT run `go build`, `go test`, or any git commands — the lead verifies the build and owns commits/releases.

## Steps

- [ ] Policy doc: add the deprecation policy to the canonical CLI-conventions docs home (follow wherever the flag-alias deprecation is documented; if there is no such home, create `docs/cli-deprecation-policy.md`, lowercase). Content: `Msg*` constant + stderr one-liner + two-release warning window + "aliases never break" + the test contract.
- [ ] Add the `Msg*` constant(s) in `cli/constants/constants_cli.go`, mirroring the existing flag-alias constants (e.g. the `--concurrency is deprecated; use --workers` constant).
- [ ] Emit the stderr one-liner at command dispatch time, mirroring `cli/cmdscan/flags.go:152`: one line on stderr, then continue the old code path unchanged.
- [ ] `fix` redirect warning: in `cli/cmdautofix/fix.go` (unknown-subcommand path at lines 43–47, routed from `cli/cmd/rootutility.go:477`), replace the bare usage error with the deprecation-style warning naming `stash` / `wip` / `discard` as the intended targets, then route to the suggested subcommand (or show its usage).
- [ ] Tests: assert (a) the stderr warning is emitted exactly once, and (b) the deprecated path still executes its old behavior.
- [ ] Keep every touched file ≤300 lines; split by concern if a dispatch file grows past the limit.

## Evidence

(policy doc path, constant name(s), test names + results as reported by the lead)
