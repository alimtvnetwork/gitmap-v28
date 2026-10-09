# Plan 245 — fix backup-branch clean-tree guard vs last_error.log side effect

## Defect (A08, from task 244)
`writeLastErrorFile` (`cli/cmd/root.go:404`) creates `.gitmap/` + `.gitmap/last_error.log`
CWD-relative on EVERY dispatch error (`handleGlobalError` :209 → `persistLastError` :390).
`requireCleanTree()` (`cli/cmdspace/backup_branch_git.go:44`) runs raw
`git status --porcelain` and refuses on any non-empty output. In
`cli/cmdspace/backup_branch.go` the guard (:30) runs BEFORE
`ensureBranchAbsence(branch, force)` (:34). Result: a failed `space backup-branch`
(e.g. "branch exists") dirties the repo, so the next invocation's guard refuses
before `--force` is ever consulted. Blast radius: EVERY failing command run inside
a repo dirties that repo's CWD.

## Fix (F1, per spec 01)
`requireCleanTree()` ignores porcelain lines that refer ONLY to untracked gitmap
state: status `??` AND path `.gitmap` or `.gitmap/...`. Tracked modifications
(including a tracked `.gitmap/...`) still count as dirty. Rationale: untracked
`.gitmap/` was never part of the HEAD snapshot `git branch` takes; gitmap's own
diagnostic state must not defeat gitmap's own safety check. The error-log writer
is UNCHANGED — `gitmap error export` and pipeline error-logs keep working.

Rejected: F2 (skip write for validation errors — changes `gitmap error` behavior),
F3 (move log out of repo — changes storage contract + all readers).
Known limitation: the `.gitmap/last_error.log` litter still lands on failures;
only the guard is fixed.

## File boxes (disjoint)
- W1 (implement): `cli/cmdspace/backup_branch_git.go` ONLY.
- W2 (verify): `~/workspace/gitmap-e2e-244/` (new `e2e-245-guard-fix.sh`, `run-245.log`) ONLY.
- Lead: spec/plan/readme/push.

## Execution
1. W1 implements the filter, runs `go build ./...` (build ONLY — no `go test`
   per standing rule), reports the diff + binary path.
2. Lead personally verifies the build (process rule: never trust self-reported builds).
3. W2 authors `e2e-245-guard-fix.sh` (F01–F05), then runs it against the fixed
   binary + the full 244-suite regression (F06: expect 28 passed, 0 failed).
4. Secrets gate + relative-path linter on changed files; atomic commit
   (`Feature: ...`) with ONLY this task's files; push; verify.

## Acceptance
F01 failed run → `--force` succeeds (exit 0, branch at HEAD, `.gitmap/` untracked).
F02 tracked dirt → `--force` exits 1, exact `errDirtyTree` message.
F03 other untracked file → exits 1, exact message.
F04 `.gitmap/last_error.log` exists, valid JSON, contains command name.
F05 `gitmap error export` works (or cleanly reports none).
F06 244-suite re-run: 28 passed, 0 failed.
Failure policy: any FAIL → rejected, root-caused, re-verified; never close dirty.

## Specs
- `02-spec/21-app/245-fix-backup-branch-lasterror-guard/01-fix-scope-and-plan.md`
- `02-spec/21-app/245-fix-backup-branch-lasterror-guard/02-acceptance-criteria-and-evidence.md`

## Result (2026-10-09)
- Implementation: `requireCleanTree()` now ignores untracked `.gitmap/` porcelain lines
  (`treeHasUserDirt` / `isGitmapUntrackedDirt` in `cli/cmdspace/backup_branch_git.go`).
  `go build ./...` exit 0 (lead-verified); binary `/tmp/gitmap-245-lead`.
- E2E (`run-245.log`): F01–F05 all PASS — the A08 chain is fixed (refusal leaves
  `.gitmap/` untracked, `--force` recreates the branch at HEAD); guard still refuses
  on tracked dirt (F02) and other untracked files (F03); error-log feature intact
  (F04/F05).
- Regression (`run-245-f06-groups.log`): A 10/10, B 7/7, C 4/4, D 5/5 — 26 passed,
  0 failed (A08 now passes). `run-all.sh`'s E01 gate pins the 244 binary's sha256
  by design, so the orchestrator was not re-run against the 245 binary; groups
  were run directly instead.
- Verdict: all acceptance IDs F01–F06 PASS. Task complete.
