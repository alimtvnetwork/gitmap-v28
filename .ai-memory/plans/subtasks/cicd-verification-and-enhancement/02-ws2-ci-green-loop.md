# Subtask 02 — WS2: CI green loop

**Version:** 1.0.0
**Updated:** 2026-10-17
**Owner:** CI worker
**Status:** pending (starts after WS1 lead review)

---

## Goal

Drive CI from its current state to **100% green** by running the exact loop
in `02-spec/21-app/cicd-verification-and-enhancement/02-ci-green-loop-procedure.md`.

## Inputs

- `02-ci-green-loop-procedure.md` — the loop (normative).
- `01-overview.md` §4 — hard rules; §6 — ground truth.

## Steps

1. Run `gitmap pe` and record the current failure set.
2. For each failure, in priority order (blocking workflows first):
   a. Read the failing run's log via `gh` — quote the failing command and
      error line; never diagnose from memory or guesses.
   b. Classify: flaky infra (re-run first) vs deterministic (fix root cause).
   c. Apply the smallest fix addressing the root cause.
   d. Commit atomically: `gitmap cpf "<module> - <summary>"` (hyphen, no
      colons) and push; verify with `git ls-remote`.
   e. Re-run `gitmap pe`; repeat until the failure clears.
3. One failure per commit. Coordinate sequencing with the lead if other
   workers have uncommitted files in the tree.
4. If a failure survives two different fixes, STOP and escalate to the lead
   with: workflow/job/step, quoted log lines, what was tried, exact error.
5. When `gitmap pe` shows zero errors and the dashboard is green on HEAD,
   write the completion report (see acceptance criteria).

## Acceptance criteria

- [ ] `gitmap pe` reports zero pipeline errors (summary lines quoted in report).
- [ ] CI dashboard green on the latest HEAD commit — all workflows, all
      platforms, including the coverage-floor gate (no >2% drop vs baseline).
- [ ] Every fix is one atomic pushed commit; no binaries committed; no raw
      `git` commands used.
- [ ] No version bumps were made during the loop.
- [ ] Completion report delivered to the lead containing: pe summary, one
      line per fix commit (failure → change), the green HEAD sha, and the
      explicit line "WS2 COMPLETE — CI 100% green on <sha>; WS5 unblocked."

## Out of scope

Version bumps and releases (WS5); new features; touching the hard-exclusion
list; weakening any CI gate to force a pass.
