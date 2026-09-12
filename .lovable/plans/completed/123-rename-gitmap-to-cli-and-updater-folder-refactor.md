# Completed Plan 123: Rename gitmap to cli and gitmap-updater to cli-updater Folder Refactor

## Task Lifecycle & Step Audit
- **Main Task Origin**: User prompt requesting:
  1. Remove `gitmap/gitmap.json`.
  2. Rename `gitmap/` folder to `cli/` and replace all references across codebase, including test inventory, runner, linters, workflows, and docs.
  3. Rename `gitmap-updater/` to `cli-updater/` and replace all references.
  4. Fix all paths everywhere across the entire repository.
- **Continuous Loop Budget**: N = 200
- **Total Steps Executed**: 14 (Decomposition: 6 subtasks; Refactoring: 1,608 Go source files; Paths updated: 38 files; Tests indexed: 3,534; Linters verified: 7/7).
- **Completion Status**: 100% Complete & Verified.
- **Verification Summary**:
  - `go build ./...` inside `cli/`: PASS (code 0)
  - `go build ./...` inside `cli-updater/`: PASS (code 0)
  - `go vet ./...` inside `cli/`: PASS (code 0)
  - `go vet ./...` inside `cli-updater/`: PASS (code 0)
  - `python .github/scripts/tests/test_ci_scripts.py`: PASS (18/18 tests)
  - AST Linters: 0 violations across 2,800+ source files.
  - Test Inventory: 3,534 tests indexed under `cli/...` (0 empty targets, 0 absolute paths).

---

## Consolidated Subtasks Execution Record

### Subtask [01-remove-gitmap-json.md]

# Subtask 01: Remove gitmap.json and Stale Artifacts

## Objective
Remove obsolete `gitmap/gitmap.json` and stale scratch file `fix_root.go`.

## Context & Invariants
- `gitmap/gitmap.json` is a legacy artifact containing empty `[]`.
- Runtime scans produce `.gitmap/output/gitmap.json` and do not depend on `gitmap/gitmap.json`.
- `fix_root.go` was a temporary root scratch script.

## Execution Steps
1. Delete `gitmap/gitmap.json`.
2. Delete `fix_root.go`.
3. Verify git status tracks the deletion.


---

### Subtask [02-git-mv-folders.md]

# Subtask 02: Git Move Folders (gitmap -> cli, gitmap-updater -> cli-updater)

## Objective
Atomically rename directories while preserving 100% git commit lineage.

## Context & Invariants
- Moving `gitmap/` to `cli/` eliminates root naming confusion.
- Moving `gitmap-updater/` to `cli-updater/` standardizes updater naming.
- `git mv` preserves rename detection (R100) in git history.

## Execution Steps
1. Execute `git mv gitmap cli`.
2. Execute `git mv gitmap-updater cli-updater`.
3. Verify git status shows renames tracked properly.


---

### Subtask [03-go-modules-and-imports-refactor.md]

# Subtask 03: Go Modules, Constants, and Imports Refactoring

## Objective
Update module names in `go.mod` files and transmutate all Go import statements.

## Context & Invariants
- `cli/go.mod`: `module github.com/alimtvnetwork/gitmap-v28/cli`
- `cli-updater/go.mod`: `module github.com/alimtvnetwork/gitmap-v28/cli-updater`
- Over 1,600 Go source files contain imports starting with `github.com/alimtvnetwork/gitmap-v28/gitmap`.
- All must be systematically updated to `github.com/alimtvnetwork/gitmap-v28/cli`.
- `cli-updater/main.go` import must point to `github.com/alimtvnetwork/gitmap-v28/cli-updater/cmd`.
- Constants:
  - `cli/constants/constants_update.go`: `UpdaterBin = "cli-updater"`
  - `cli-updater/cmd/constants.go`: `UpdaterCopy = "cli-updater-tmp-%d.exe"`
  - `cli/constants/deploy-manifest.json`: `"sourceRepoSubdir": "cli"`
  - `cli/constants/deploy_manifest.go`: fallback `"cli"`
  - `cli/powershell.json`: `"goSource": "./cli"`
  - `cli/.goreleaser.yaml`: ldflags module paths updated to `cli/...`.

## Execution Steps
1. Edit `cli/go.mod` and `cli-updater/go.mod`.
2. Batch replace import paths in all `.go` files across `cli/`.
3. Update `cli-updater/main.go` and `cli-updater/cmd/constants.go`.
4. Update `cli/constants/constants_update.go`, `deploy-manifest.json`, `deploy_manifest.go`, `powershell.json`, and `.goreleaser.yaml`.
5. Run `go build ./...` in `cli/` and `cli-updater/` to ensure syntax validity.


---

### Subtask [04-ci-scripts-and-linters-path-updates.md]

# Subtask 04: CI/CD Runner, Linters, Workflows, and Bootstrap Paths

## Objective
Update all path references across Python runners, AST linters, GitHub Actions workflows, and bootstrap scripts.

## Context & Invariants
- `03-ai-scripts/06-cicd-local-runner.py`:
  - `CLUSTER_GO_*` and `CLUSTER_REPO_TEXT` patterns to `cli/` and `cli-updater/`.
  - `GATE_SPECS` configs, tool scripts, and relevant patterns.
  - `build_or_update_test_inventory` and `run_package_tests_worker` directory roots and `cwd`.
- `03-ai-scripts/33-test-inventory-generator.py`: heavy test candidate path `repo_root / "cli" / prefix`.
- `03-ai-scripts/14-version-sync-checker.py`: `cli/constants/constants.go`.
- `03-ai-scripts/29-release-bumper.py` and `29-release-orchestrator.py`: `cli/constants/constants.go`.
- `03-ai-scripts/30-db-struct-enum-generator.py`: `cli/pipelinedb`.
- `linter-scripts/`:
  - `check-enum-and-boolean.py`: `"cli" in filepath.parts and "constants" in filepath.parts`.
  - `check-enum-guidelines.py`, `check-function-formatting.py`, `check-function-signatures.py`, `check-schema-guidelines.py`: `ROOT_DIR / 'cli'`.
  - `check-error-management.py`: `"cli-updater"`.
- `.github/scripts/`:
  - `check-constants-collisions.py`, `check-constants-naming.py`, `check-deploy-layout.py`, `check-no-golden-allow-leak.py`, `fix-repo-gofmt-audit.py`, `misspell-changed.py`, `test_ci_scripts.py`.
- `.github/workflows/`:
  - 10 workflow YAMLs: `working-directory: cli`, `cli/go.mod`, `cli/go.sum`, `cli/dist/*`, `cli/helptext/*`, `cli-updater`.
- `run.ps1` and `run.sh`:
  - `$CliDir = Join-Path $RepoRoot "cli"`, `CLI_DIR="$REPO_ROOT/cli"`, and build ldflags.

## Execution Steps
1. Update `03-ai-scripts/06-cicd-local-runner.py` and `33-test-inventory-generator.py`.
2. Update remaining `03-ai-scripts/` utilities.
3. Update all `linter-scripts/*.py`.
4. Update `.github/scripts/*.py`.
5. Update all `.github/workflows/*.yml`.
6. Update `run.ps1` and `run.sh`.


---

### Subtask [05-test-inventory-regeneration.md]

# Subtask 05: Test Inventory Regeneration

## Objective
Regenerate `.lovable/test-inventory.json` with 100% relative `cli/...` test targets.

## Context & Invariants
- Test inventory must map every test function to its corresponding target file.
- All packages and files must begin with `cli/...` (or `cli-updater/...`).
- 0 empty target_file entries, 0 absolute paths.
- Test execution remains strictly disabled.

## Execution Steps
1. Execute `python 03-ai-scripts/33-test-inventory-generator.py --force-run-all --slow-threshold 4.0`.
2. Verify output JSON has 0 references to `gitmap/` in package or target fields.
3. Verify test counts match expectations (~3,534 tests).


---

### Subtask [06-verification-and-consolidation.md]

# Subtask 06: Verification, Linter Checks, and Plan Consolidation

## Objective
Verify compilation, execute AST linters, and consolidate all subtasks into completed plan.

## Context & Invariants
- `go build ./...` and `go vet ./...` in both `cli/` and `cli-updater/`.
- Linters must pass with 0 violations:
  - `python linter-scripts/check-nested-ifs.py`
  - `python linter-scripts/check-enum-and-boolean.py`
  - `python linter-scripts/check-schema-guidelines.py`
  - `python linter-scripts/check-function-formatting.py`
  - `python linter-scripts/check-function-signatures.py`
  - `python linter-scripts/check-error-management.py`
  - `python 03-ai-scripts/14-version-sync-checker.py`
- Consolidation:
  - Combine subtasks into `.lovable/plans/completed/123-rename-gitmap-to-cli-and-updater-folder-refactor.md`.
  - Remove subtasks directory and pending plan.
  - Update `.lovable/plans/01-index.md`.
  - Update `walkthrough.md`.

## Execution Steps
1. Run `go build ./...` and `go vet ./...` in `cli/` and `cli-updater/`.
2. Run AST linters.
3. Consolidate plan and clean up subtasks.
4. Update indices and memory.


---
