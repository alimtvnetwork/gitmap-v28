# 135-repo-scoped-temp-storage-and-prebuild-clean.md: Repository-Scoped Temp Storage & Mandatory Pre-Build Cleanup

**Status: completed**
**Execution Loop Count: 8 steps completed in continuous N-step self-loop (N=150 budget)**
**Task Initiation: Initiated via user request to confirm build artifact locations (distinguishing workspace bin from OS/user temp), strictly enforce repository-scoped subdirectories under OS temp (<temp>/gitmap/<category>), mandate clean-before-build execution to prevent disk waste, and update all prompts, skills, and coding guidelines.**

---

## 1. Executive Summary & Build Location Confirmation

### Build Locations Analysis
1. **Workspace Build Artifacts**: The CI/CD local runner compiles the main application executable to `bin/gitmap.exe` directly within the repository workspace.
2. **OS/User Temp Artifacts**: Smoke tests (`smoke-installer.py`, `smoke-history-pin.py`, `smoke-history-purge.py`), heavy test suites (`cliexit_helpers_test.go`), and purge engines previously wrote test binaries and workdirs directly into the root of the OS user temp directory (`os.TempDir()` / `tempfile.gettempdir()`).

### Core Solution
1. **Repository Namespacing in OS/User Temp**:
   All operations targeting OS/user temp must strictly write inside a repository-named directory: `<temp_dir>/gitmap/<subfolder>/` (e.g. `build/`, `test/`, `purge/`, `downloads/`). No loose files or arbitrary folders in the root of OS temp.
2. **Mandatory Pre-Build Cleanup (Storage Respect & Reuse)**:
   Before any build command executes, stale binaries and compilation artifacts in the destination build directory must be purged. Storage must be reused efficiently without accumulating orphaned binaries.
3. **Prompts, Skills, and Coding Guidelines Updates**:
   Updated `.lovable/strictly-avoid.md`, `.lovable/coding-guidelines.md`, `spec/02-coding-guidelines/01-cross-language/29-no-generated-artifacts.md`, `.agents/skills/coding-guidelines/skill.md`, `.agents/skills/ci-cd-create/skill.md`, created `.agents/skills/temp-storage-and-build-hygiene/SKILL.md`, and updated prompts in `01-prompts/`.

---

## 2. Mandatory Rules & Invariants Followed

1. **Rule 1 (Strict Relative Git Paths)**: All markdown paths, documentation, and references are strictly relative to repo root. No absolute paths or `file:///` URIs.
2. **Rule 2 (Strict Sizing & Style)**: Every Go function <= 15 lines. Mandatory blank line before every return statement.
3. **Rule 3 (Affirmative Booleans)**: All boolean variables and struct fields use affirmative prefixes (`is*`, `has*`).
4. **Rule 4 (Zero Swallowed Errors)**: AppError / proper error wrapping with context.
5. **Rule 5 (Cache Tracking)**: All modified files recorded in `.lovable/temp/recent-file-changes.json` under atomic lock via `python 03-ai-scripts/33-test-inventory-generator.py --record`.
6. **Rule 6 (Consolidated Commit & Push)**: Committed all touched files atomically into a single commit and pushed to `origin main`.
7. **Rule 7 (Constraint Enforcement)**: No auto-releases, no routine unit test execution, no `06-cicd-local-runner.py`.

---

## 3. Architecture & Implementation Details

### 3.1 Go tempdir Package & Pre-Build Clean
- Created `cli/tempdir/tempdir.go`:
  - `RepoTempDir(subdirs ...string) string`: returns `<os.TempDir()>/gitmap/<subdirs>` with automatic `os.MkdirAll`.
  - `BuildTempDir() string`: returns `RepoTempDir("build")`.
  - `TestTempDir() string`: returns `RepoTempDir("test")`.
  - `ClearRepoBuildTempDir() error`: sweeps and purges all existing artifacts from `BuildTempDir()` before any build executes.
- Updated `cli/tests/heavy_test/cliexit_helpers_test.go`:
  - Calls `tempdir.ClearRepoBuildTempDir()` before compiling.
  - Builds test executable to `tempdir.BuildTempDir()`.
- Updated `cli/cmdpurge/purge_engine.go` to place purge backup archives in `tempdir.RepoTempDir("purge", ...)`.
- Updated `cli/cmdclone/clonevscode.go` to use `tempdir.RepoTempDir("vscode", ...)`.
- Created unit tests in `cli/tempdir/tempdir_test.go` (100% passing).

### 3.2 Python Smoke Tests & Purge Utilities Scoping
- `.github/scripts/smoke-installer.py`:
  - Added `get_repo_temp_dir(*subdirs)` scoping all workdirs under `<tempfile.gettempdir()>/gitmap/<subdirs>`.
  - Added `clear_repo_build_temp()` wiping `<temp>/gitmap/build` before source build runs.
  - Scoped mock release directory and smoke test workdir to `<temp>/gitmap/test`.
- `.github/scripts/smoke-history-pin.py` and `.github/scripts/smoke-history-purge.py`:
  - Scoped workdirs to `<tempfile.gettempdir()>/gitmap/test`.
- `03-ai-scripts/30-purge-history.py` and `03-ai-scripts/33-git-history-tracer-and-purger.py`:
  - Scoped purge backups to `<tempfile.gettempdir()>/gitmap/purge`.

### 3.3 Specs, Prompts & Skills Updates
- `.lovable/strictly-avoid.md`: Added total ban on un-namespaced OS temp writes and un-cleaned builds.
- `.lovable/coding-guidelines.md`: Added Temp Directory & Storage Hygiene section.
- `spec/02-coding-guidelines/01-cross-language/29-no-generated-artifacts.md`: Added Storage Hygiene, Temp Directory Isolation & Pre-Build Cleanup section.
- `.agents/skills/coding-guidelines/skill.md` & `.agents/skills/ci-cd-create/skill.md`: Added rule R17 to review checklists.
- Created `.agents/skills/temp-storage-and-build-hygiene/SKILL.md`.
- `01-prompts/16-ci-cd/02-cicd-run-ps1.md` & `01-prompts/04-coding-standards/01-coding-guidelines.md`: Updated with repository temp scoping and pre-build cleanup constraints.

---

## 4. Verification & Testing

- Unit tests: `go test -v ./tempdir` passed 100%.
- Linters:
  - `python linter-scripts/check-nested-ifs.py`: 0 violations across 2,858 files.
  - `python linter-scripts/check-boolean-guidelines.py`: 0 violations across 2,858 files.
  - `python linter-scripts/check-relative-paths.py`: 0 violations across 6,809 files.
  - `gofmt -w cli/`: Clean.
  - `python -m py_compile`: Clean on all touched Python scripts.
  - `go vet ./tempdir ./cmdclone ./cmdpurge ./tests/heavy_test`: Clean (exit 0).
