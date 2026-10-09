# 02 — Acceptance Criteria and Evidence (task 245: fix backup-branch last_error guard)

## Fix under test (F1)

`requireCleanTree()` in `cli/cmdspace/backup_branch_git.go` ignores porcelain lines
that are untracked (`??`) gitmap state — path `.gitmap` or anything under
`.gitmap/...`. Tracked modifications still count as dirty. The error-log writer
(`cli/cmd/root.go:404`) is UNCHANGED: `.gitmap/last_error.log` still lands in CWD
on failures.

## Acceptance table

| ID | Criterion | PASS condition | Evidence artifact |
|----|-----------|----------------|-------------------|
| F01 | A failed `space backup-branch` (e.g. target branch already exists) followed by the same command with `--force` SUCCEEDS. | Second invocation exits 0 and the branch points at HEAD, with `.gitmap/` present as untracked (the guard ignores gitmap's own untracked state). | `e2e-245-guard-fix.sh` section F01 in `run-245.log`: exit codes, `git rev-parse <branch>` == HEAD, `git status --porcelain` showing `?? .gitmap/`. |
| F02 | Guard still refuses on real tracked dirt. | With a modified tracked file, `space backup-branch --force` exits 1 with the exact `errDirtyTree` message. | `run-245.log` F02 section: exit code 1, verbatim refusal message, `git status --porcelain` showing the tracked modification. |
| F03 | Guard still refuses on OTHER untracked files. | With an untracked `notes.txt` and no `.gitmap/` present at invocation, `--force` exits 1 with the exact `errDirtyTree` message. | `run-245.log` F03 section: porcelain shows `?? notes.txt`, exit code 1, verbatim refusal message. |
| F04 | Error-log feature intact. | After a refusal, `.gitmap/last_error.log` exists in CWD and is valid JSON containing the command name. | `run-245.log` F04 section: the file's JSON content (validated with a JSON parser), recorded verbatim. |
| F05 | `gitmap error export` still works. | Exports the last error to a user-specified file, OR reports "no recent error" cleanly when no error is on record — either is a PASS; the script documents which path was taken. | `run-245.log` F05 section: exported file content, or the "no recent error" output, whichever occurred. |
| F06 | No regressions — full task-244 suite re-run against the fixed binary. | `RESULT: 28 passed, 0 failed` (A08 now passes; E01 + D04-preflight + all groups green). Task-244 scripts are used unchanged, run against the fixed binary. | `run-245.log` tail showing the RESULT line. |

## Evidence locations

- All e2e work lives OUTSIDE the repo at `~/workspace/gitmap-e2e-244/`: new script
  `e2e-245-guard-fix.sh`, log `run-245.log`.
- Nothing test-related is committed to the repo. Test fixtures, scratch repos, and
  logs are disposable and stay outside the tree.

## Out of scope

- Unit tests (`go test`) are NOT run per the standing rule (never run tests without
  the owner's explicit command). Verification is e2e-only via the built binary.

## Failure policy

- Any FAIL → the fix is rejected, root-caused, and re-verified. The task does not
  close with a failing acceptance ID.
