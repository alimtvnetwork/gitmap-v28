# 245 — Fix `space backup-branch` last-error guard (fix scope and plan)

## 1. Defect statement

Running `space backup-branch` twice for the same task string leaves the user
stuck: the first run refuses because the backup branch already exists (correct),
but that refusal itself writes gitmap's error log into the repo, so the second
run — even with `--force` — is refused again with "uncommitted changes"
(incorrect). The user must manually delete the `.gitmap/` directory before
`--force` will work.

## 2. Root cause

Every dispatch error goes through `cli/cmd/root.go:209` (`handleGlobalError`)
→ `cli/cmd/root.go:390` (`persistLastError`) → `cli/cmd/root.go:404`
(`writeLastErrorFile`), which does `os.MkdirAll(".gitmap", 0755)` and writes
`.gitmap/last_error.log` relative to the current working directory — i.e. inside
the user's repo.

`cli/cmdspace/backup_branch_git.go:44` (`requireCleanTree`) runs a raw
`git status --porcelain` (`:45`) and refuses on ANY non-empty output (`:50-51`).
In `cli/cmdspace/backup_branch.go` the guard (`:30`) runs BEFORE
`ensureBranchAbsence(branch, force)` (`:34`).

Chain: run N refuses (branch exists) → `writeLastErrorFile` creates untracked
`.gitmap/` → run N+1's `requireCleanTree` sees the porcelain output → refuses
with `errDirtyTree` before `--force` is ever consulted. The guard that protects
the HEAD snapshot is defeated by gitmap's own diagnostic litter.

Readers of the log that must keep working (untouched by this fix):
`cli/cmderrors/error_cmd.go:36`, `cli/cmdpipeline/pipeline_logs_target.go:277`,
`cli/cmdstorage/storage_reset.go:86`. The DB write (`persistToErrorsDB`) targets
the global dir, not CWD — unaffected.

## 3. Fix design (chosen: F1)

Change `requireCleanTree()` in `cli/cmdspace/backup_branch_git.go` to ignore
porcelain lines that refer ONLY to untracked gitmap state: status code `??`
AND path `.gitmap`, `.gitmap/`, or anything under `.gitmap/...`. Every other
porcelain line (including a *tracked* `.gitmap/...` modification, however
unusual) still counts as dirty.

Rationale: untracked `.gitmap/` was never part of the HEAD snapshot the backup
takes — `git branch` snapshots tracked state only — so gitmap's own diagnostic
state must not defeat gitmap's own safety check. Requiring the `??` status keeps
the rule conservative: no real user change (tracked modifications, staged
changes) is ever masked.

Explicitly rejected:
- F2 (skip the error-file write for validation errors): changes `gitmap error`
  behavior for all commands; bigger blast radius, alters the feature contract.
- F3 (move the log out of the repo): changes the storage contract and breaks
  all three readers above; a redesign, not a fix.

## 4. Known limitation

The `.gitmap/last_error.log` litter still lands in the repo on failures —
unchanged behavior. This fix only stops the guard from tripping on it. Removing
or relocating the litter is out of scope.

## 5. File box

`cli/cmdspace/backup_branch_git.go` ONLY. No other source file may be touched.

Test-file note: `cli/cmdspace/` currently contains NO test files (checked), so
per the constraint no new test file may be added — verification is via build
plus e2e (section 6).

## 6. Test plan

1. `go build ./...` from `cli/` — build only. NO `go test` per the standing
   rule (tests run only on the owner's explicit command).
2. E2E A08 chain in a disposable scratch repo (artifacts live OUTSIDE the repo,
   under `~/workspace/gitmap-e2e-244/`, never committed):
   - pre-create branch `backup/<slug>` → run `space backup-branch --no-push
     "<task>"` → expect exit 1, "already exists";
   - `git status --porcelain` → `.gitmap/` present but untracked;
   - run `space backup-branch --no-push --force "<task>"` → expect exit 0,
     branch recreated at HEAD (no manual cleanup);
   - sanity: a real dirty tree (tracked file modified) still refuses with
     `errDirtyTree`, and `--force` does not override it.
3. Regression: re-run the full task-244 e2e suite (`run-all.sh`) — expect the
   same result as baseline (only A08 now passing).

## 7. Non-goals

- No refactoring of the error-log feature; `cli/cmd/root.go` is not touched.
- No change to what gets written to `.gitmap/last_error.log` or when.
- No change to the refusal message or to any other `space` subcommand.
- No changes to the three log readers or the errors DB.
