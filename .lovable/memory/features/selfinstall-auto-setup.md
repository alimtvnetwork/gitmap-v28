---
name: self-install auto-runs setup
description: `gitmap self-install` invokes `gitmap setup` as a final non-fatal step so the shell wrapper (gcd / gitmap function) is installed without a separate command. v4.43.0+.
type: feature
---

# Feature: self-install auto-runs setup (v4.43.0+)

## Rule
`gitmap self-install` MUST run `gitmap setup` after the install
script completes successfully but before printing the reminder line.
Setup is idempotent (marker `# gitmap shell wrapper v2`), so re-runs
on every install are safe.

## Why
First-time users were hitting `! Shell wrapper not active — 'gitmap cd'
printed the path but cannot change your directory.` because `setup`
was a separate manual step. Auto-running it closes the gap so a fresh
install gives working `gcd` after one terminal restart.

## Behavior
- `runSelfInstall` calls `autoRunSetupAfterInstall()` between `MsgSelfInstallDone` and `MsgSelfInstallReminder`.
- Wrapped in `recover()` so a setup panic never breaks the install summary.
- Failures print `(setup auto-run skipped: <reason>)` to stderr and continue.

## Files
- `gitmap/cmd/selfinstall.go::autoRunSetupAfterInstall`.

## Limits
- Cannot re-source the parent shell from a child process — user STILL needs `. $PROFILE` or a new terminal. The `cd` wrapper-not-active warning remains as a final UX safety net.
