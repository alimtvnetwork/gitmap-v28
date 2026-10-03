# Completed Plan: 65-commit-non-git-rca

## Overview
- **Task ID:** 65-commit-non-git-rca
- **Status:** COMPLETED
- **Description:** Fix `gitmap commit all` execution crash in non-git directory (`D:\work`), add child git repository batch commit delegation, beautify terminal UI for `gitmap cpf` and `commit-push` (CRLF warning noise filtering, 4-space padded git output, completion summary card), remediate CI/CD failure gates (legacy refs check and gofmt unformatted files), author 4-part RCA, and release via minor version bump.

## User Request (Verbatim)
```text
PS D:\work> gitmap commit all                                                                                  
  INFO Staging all changes...                                                                                  
fatal: not a git repository (or any of the parent directories): .git                                           
gitmap commit: execute failed: [E9000:EXECUTION] git add failed:: exit status 128 (at=cmd/commit_cmd.go:52)    
Stack Trace:                                                                                                   
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.executeCommit (cmd/commit_cmd.go:52)                        
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runCommit (cmd/commit_cmd.go:24)                            
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.coreBasicOpEntries.func21 (cmd/rootcore.go:113)             
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatchTable (cmd/rootdispatch.go:24)                   
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatchCore (cmd/rootcore.go:22)                           
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:481)                                  
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:145)                               
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:112)                                       
    at main.main (cli/main.go:7)                                                                               
PS D:\work>                                                                                                    

please fix this with RCA and release withminor bump pelase
```
Also previous user request:
"please make this output better and release agian and check gitmap pe please until ghreen"
with screenshot of `gitmap cpf "chrome profile for gitmap"`.

## Deliverables & Verified Outcomes
1. **CLI Non-Git Repository Preflight Guard**:
   - `cli/cmd/commit_cmd.go`: Added `!isGitRepoCWD()` guard in `dispatchCommit`. When invoked outside a git repository, it diverts to `handleNonGitRepoCommit`.
   - `cli/cmd/commit_push.go`: Added `!isGitRepoCWD()` guard in `executeCommitPush` returning `handleNonGitRepoCommitPush()` cleanly.
   - Stack traces suppressed by returning an `ErrorTypeAbort` `AppError` with `reported: true`.

2. **Workspace Child Repository Batch Delegation**:
   - `cli/cmd/commit_batch.go`: Added `isAllCommitRequested`, `stripAllFlags`, and `executeChildReposBatch`.
   - When running `gitmap commit all` from a workspace root containing child Git repositories, it discovers child repos via `fsutil.DiscoverChildGitRepos(cwd)` and commits dirty repositories.

3. **Terminal UI Sanitization & Formatting**:
   - `cli/cmd/commit_ui.go`: Added `isGitNoiseLine` filtering out CRLF/LF line ending translation warnings.
   - `execGitPaddedFiltered`: Indents git output by 4 spaces and drops noise lines.
   - `renderCommitPushSummaryCard`: Renders boxed terminal card displaying branch, commit SHA, message, and remote status.

4. **CI/CD Failure Gate Resolution**:
   - `02-spec/01-spec-authoring-guide/13-root-readme-conventions.md`: Updated outdated `gitmap-v6` references to `gitmap-v28`. Verified `check-legacy-refs.py` passes with 0 findings.
   - `cli/cmdchromeprofile/chromeprofile_paths.go`, `cli/cmdssh/exports.go`, `cli/cmdports/ports.go`, `cli/cmdports/ports_test.go`, `cli/cmdports/ports_types.go`: Formatted with `gofmt`. Verified `go-format-check.py` passes with all 3898 files clean.
   - Unit tests: `test_ci_scripts.py` passes 18/18 tests.

5. **RCA Documentation**:
   - `02-spec/22-app-issues/65-commit-non-git-repo-crash-and-commit-all-delegation-rca.md`: Authored full 4-part RCA and indexed in `readme.md` and `01-index.md`.

## Verification Evidence
- `python .github/scripts/check-legacy-refs.py .` -> Exit 0 (0 legacy references).
- `python .github/scripts/go-format-check.py --check-only` -> Exit 0 (All 3898 files gofmt-clean).
- `python .github/scripts/tests/test_ci_scripts.py` -> Exit 0 (Ran 18 tests, OK).
- `python linter-scripts/check-relative-paths.py` -> Exit 0 (0 absolute paths).
- `python linter-scripts/check-forbidden-strings.py` -> Exit 0 (All forbidden-string rules passed).
