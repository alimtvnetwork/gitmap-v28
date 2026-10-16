# 02 — CI Green Loop Procedure

**Version:** 1.0.0
**Updated:** 2026-10-17

---

## 1. Purpose

The exact, repeatable loop WS2 follows to drive CI from its current state to
**100% green**. No improvisation: every iteration runs the same steps in the
same order.

## 2. The loop

```
1. CHECK    gitmap pe                      (pipeline errors / CI status)
2. READ     gh logs of the failing run     (the ACTUAL log text — never guess)
3. DIAGNOSE identify the single failing step + root cause from the log
4. FIX      smallest change that fixes the root cause
5. COMMIT   gitmap cpf "<module> - <summary>" (hyphen, NO colons) + push
6. RE-CHECK gitmap pe again
7. REPEAT   until 100% green — no partial credit, no "green except"
```

One failure fixed per commit. Never batch unrelated fixes into one commit;
never fix what the log does not show.

## 3. Reading failure logs (no guesses)

- Pull the failing workflow run's log with `gh` and read the **actual
  output text**. The diagnosis must quote the failing command and the error
  line — a diagnosis without a quoted log line is not a diagnosis.
- Extract, in order: (a) which workflow + job + step failed; (b) the exact
  command that failed; (c) the error text verbatim; (d) whether it is
  flaky (infra/timeout) or deterministic (code/config).
- Flaky infra failures (runner timeout, network blip) → re-run the workflow
  first; only fix code when the failure reproduces.
- Deterministic failures → fix the root cause in the repo, not the symptom
  in the workflow file. Prefer fixing the code over loosening the gate.
- Known gates that must stay green (do not weaken them to "pass"):
  coverage floor (fail on >2% drop vs baseline, program 272), all platform
  builds, lint gates.

## 4. What counts as green

- `gitmap pe` reports zero pipeline errors.
- The CI dashboard shows all workflows green on the latest commit of the
  working branch — including scheduled/nightly workflows if they ran.
- Green means the **latest** commit. A green older commit with a red HEAD is
  not green.

## 5. Commit discipline during the loop

- Every fix: atomic `gitmap cpf "<module> - <summary>"`, hyphen, no colons,
  and every commit MUST push (verify with `git ls-remote`).
- Never run raw `git` commands; never amend or force-push.
- Never commit binaries or test artifacts.
- While other workers have uncommitted files in the tree, coordinate with the
  lead before committing — `gitmap cpf` stages broadly, so sequencing is the
  lead's call.
- No version bumps inside the loop. Bumping is WS5's job and is BLOCKED until
  this loop reports green.

## 6. Stop conditions

- STOP and report to the lead when: the same failure survives two different
  fixes (escalate — do not keep guessing); a failure is clearly external
  (GitHub incident, third-party service down); a fix would require an
  owner decision (scope change, rejected-list territory).
- Report format: failing workflow/job/step, quoted log lines, what was tried,
  exact error text, what is needed (owner call / later wave / infra wait).

## 7. Handoff to WS5

WS2's completion report to the lead must contain:

1. `gitmap pe` output showing zero errors (copy the summary lines).
2. List of every fix commit (one line each: what failed → what was changed).
3. The exact HEAD sha that is green.
4. Explicit statement: "WS2 COMPLETE — CI 100% green on <sha>; WS5 unblocked."

WS5 (release) is BLOCKED until this report lands. Starting a bump or a
release on a red tree is a program failure, not a shortcut.
