# Spec 170: Repo Create Coding Guideline Auto-Sync, No-Sync Flag & Race Detector Guard

## Status
`active`

## Context & Motivation
1. **Coding Guideline Auto-Sync Failure in `create-repo`:**
   When running `gitmap create-repo <name> --cg`, previous versions failed to synchronize the repository coding guidelines (`02-spec/02-coding-guidelines`). This was caused by `resolveCGBaseDir()` attempting to inspect `filepath.Dir(filepath.Dir(os.Executable()))` when running from the compiled binary installation path (`%USERPROFILE%\AppData\Local\gitmap-cli\gitmap.exe`), where no source specs reside. As a result, it fell back to generating only an empty `00-overview.md` file rather than synchronizing the complete 30-document coding guidelines suite.
2. **Selective Workspace Synchronization (`--no-sync`):**
   Users working with lightweight or automated repository creation requested a `--no-sync` / `--skip-sync` flag to bypass launching GitHub Desktop and opening VS Code, while still registering with Google Antigravity (`[vsc: skipped] [desktop: skipped] [agy: ok]`).
3. **CI/CD Race Detector & Linter Guard:**
   Concurrent execution in `commitin/workspace.CloneInputs` triggered a data race warning in `workspace_test.go:TestCloneInputsStagesAllThreeKinds` when appending to a shared slice without synchronization. Furthermore, unformatted test files caused the GitHub Actions `Lint Script Unit Tests` job (`test_ci_scripts.py`) to fail.

## Architectural Changes

### 1. Robust Coding Guideline Resolution (`cli/cmd/repo_create_init.go`)
- Multi-tier resolution strategy:
  1. `constants.RepoPath` (when set at build time).
  2. `GITMAP_CODING_GUIDELINES_DIR` environment variable.
  3. Current working directory and parent directory walk (`resolveCGBaseDirFromCwd()`).
  4. Local development drives (`./coding-guidelines`, `C:\work\coding-guidelines`).
  5. Dynamic remote installer fallback via `cmdcg.RunCgScriptInRepo(absDir)` if no local source directory exists.
- In `applyCGFiles`: recursively copy the full `02-spec/02-coding-guidelines` directory into the destination repository, stage, and commit with `docs(spec): synchronize latest coding guidelines to 02-spec/02-coding-guidelines/`.

### 2. Workspace Synchronization Control (`cli/cmd/create_ops.go` & `cli/workspacesync/workspacesync.go`)
- Added `IsNoSync bool` to `createRepoParams`, capturing `--no-sync`, `--skip-sync`, or `GITMAP_NO_SYNC=1`.
- Added `workspacesync.SyncAgyOnly(repoPath, repoName)`:
  - Skips VS Code Project Manager.
  - Skips GitHub Desktop integration.
  - Synchronizes Antigravity configuration cleanly.
- Routed sync execution in `create_ops.go` through `applyWorkspaceSync(params, absDir)`.

### 3. Concurrency Protection & Go Format Cleanliness
- Added `sync.Mutex` in `cli/cmd/commitin/workspace/workspace_test.go` to guard `cloneTargets` slice appends across concurrent worker goroutines.
- Formatted `cli/cmdpull/pull_flags_test.go` using `gofmt -w` to satisfy `.github/scripts/go-format-check.py`.

## Verification & Outcomes
- `gitmap create-repo <dir> --local --cg --no-sync` verified:
  - 30 coding guidelines specifications copied into `02-spec/02-coding-guidelines`.
  - Terminal output confirmed: `-> sync: [vsc: skipped] [desktop: skipped] [agy: ok]`.
- All 18 Python CI unit tests in `.github/scripts/tests/test_ci_scripts.py` pass cleanly.
- `check-relative-paths.py` passed across 8,102 files with 0 violations.
- `check-nested-ifs.py` passed with 0 violations.
