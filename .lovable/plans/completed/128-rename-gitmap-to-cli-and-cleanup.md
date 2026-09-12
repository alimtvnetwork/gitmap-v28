# Completed Plan 128: Rename gitmap to cli, gitmap-updater to cli-updater, and Clean Up gitmap.json

## Status
- **State**: Completed
- **Created**: 2026-09-12
- **Completed**: 2026-09-12

## Objectives & High-Level Architecture
1. **Remove `gitmap.json`**:
   - Permanently removed the root `gitmap.json` artifact per user explicit command.
2. **Rename `gitmap` folder to `cli` and Sync All References**:
   - Renamed folder path references (`gitmap/...` -> `cli/...`) across all scripts, Go source code, tests, docs, specifications, and test inventory.
   - Updated critical path lookups:
     - `cli/store/erd_parity_test.go`: `constantsDirRel = "cli/constants"`.
     - `cli/tests/fixrepo_test/gofmt_e2e_test.go`: `findRepoRoot` uses `cli/go.mod` and `cmd.Dir = filepath.Join(repoRoot, "cli")`.
     - `cli/release/selfrelease_resolve.go`: `isGitmapSourceRepo` validates `cli/constants/constants.go`.
     - `src/test/version-sync.test.ts`: `goConstantsPath` points to `cli/constants/constants.go`.
     - `run.ps1`, `run.sh`, `install.sh`, `scripts/`, `hooks/pre-commit`, `.github/scripts/`.
3. **Rename `gitmap-updater` to `cli-updater`**:
   - Replaced all active references of `gitmap-updater` with `cli-updater` in:
     - `cli/cmd/update.go`, `cli/cmd/updateremoteinstall.go`, `cli/cmd/llmdocsgroups.go`, `cli/cmd/llmdocssections.go`
     - `cli-updater/cmd/check.go`, `cli-updater/cmd/github.go`, `cli-updater/cmd/usage.go`, `cli-updater/cmd/worker.go`
     - `cli/helptext/update.md`
     - `.github/workflows/release.yml`
     - `src/pages/InstallGitmap.tsx`
     - `readme.md`, `what-to-read.md`, `plan.md`, `spec/`, and `.lovable/`
4. **Clean up Test Inventory & Temporary Tracking**:
   - Regenerated `.lovable/test-inventory.json` with 3,534 tests across 134 packages (50 slow, 3484 fast) with 100% relative paths.
   - Cleaned legacy paths from `.lovable/temp/recent-file-changes.json` and recorded all modified files under atomic file lock.
5. **Quality Gates & Linters**:
   - Boolean Guidelines Linter: 0 violations.
   - Nested If Linter: 0 violations.
   - Relative Path Linter: 0 violations across 6,721 files.
   - golangci-lint (strict): 0 issues across all packages.
   - Full CI/CD Local Runner: 38/38 quality gates passed (100% green).

## Subtask Execution Summary
- **Subtask 1 (`01-remove-gitmap-json.md`)**: Removed root `gitmap.json`. Confirmed absent.
- **Subtask 2 (`02-replace-gitmap-updater-references.md`)**: Replaced all occurrences of `gitmap-updater` with `cli-updater` across code, docs, workflows, and specs.
- **Subtask 3 (`03-replace-gitmap-folder-paths-with-cli.md`)**: Replaced folder path occurrences of `gitmap/` with `cli/` across tests, scripts, docs, and Go resolution helpers. Preserved executable name `gitmap.exe` and repo slug `github.com/alimtvnetwork/gitmap-v28`.
- **Subtask 4 (`04-verify-test-inventory-and-scripts.md`)**: Regenerated test inventory (3,534 tests, 0 invalid paths), cleaned and recorded recent changes under lock.
- **Subtask 5 (`05-quality-gate-verification-and-consolidation.md`)**: Cleaned unused modularization stubs, fixed boolean naming (`isImportPathSwapNeeded`), passed all 38 CI/CD quality gates, and consolidated documentation.

## Verification
- `python 03-ai-scripts/06-cicd-local-runner.py --no-tests`: 38/38 gates passed (27.74s).
- `gofmt -l cli/ cli-updater/`: 0 dirty files.
- `go build -v .` in `cli/`: compiled cleanly.
- `go build -v .` in `cli-updater/`: compiled cleanly.
- Strict relative paths: 0 absolute paths across 6,721 files.
