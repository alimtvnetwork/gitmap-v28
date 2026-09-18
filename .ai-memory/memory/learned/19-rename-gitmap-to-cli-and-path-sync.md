# 19-rename-gitmap-to-cli-and-path-sync.md — Learned Conventions: CLI and Updater Path Refactoring

## 1. Directory Structure vs Product Identity Distinction
- **Folder Path**: The Go CLI package directory is `cli/` (NOT `gitmap/`). All filesystem paths, scripts, imports, help text references, and test finders must point to `cli/...`.
- **Updater Directory & Name**: The updater binary and sub-package is `cli-updater/` (NOT `gitmap-updater`).
- **Product Identity**: The executable binary name remains `gitmap.exe` (Windows) / `gitmap` (Unix), the primary command is `gitmap`, and the GitHub repository module path remains `github.com/alimtvnetwork/gitmap-v28`.
- **Runtime Directory**: `.gitmap/` (with a leading dot) remains the local database and output directory (`.gitmap/repo.db`, `.gitmap/output/`).

## 2. Critical Path Invariants in Go Source
- `cli/store/erd_parity_test.go`: `constantsDirRel = "cli/constants"`.
- `cli/tests/fixrepo_test/gofmt_e2e_test.go`: `findRepoRoot` checks `cli/go.mod`, and build directory is `filepath.Join(repoRoot, "cli")`.
- `cli/release/selfrelease_resolve.go`: `isGitmapSourceRepo` validates `cli/constants/constants.go`.
- `src/test/version-sync.test.ts`: `goConstantsPath` points to `../../cli/constants/constants.go`.

## 3. Test Inventory & Change Tracking
- Test inventory generator (`03-ai-scripts/33-test-inventory-generator.py`) maps all tests to `cli/...`.
- `recent-file-changes.json` records modified files with repository-relative paths under atomic file lock.
- Clean legacy paths immediately to prevent stale tracking during release ceremonies.

## 4. Coding Guidelines & Linters
- Zero nested `if` statements.
- Boolean functions must strictly use positive prefixes (`is`, `has`, `can`, `allow`), avoiding banned `should...` prefixes (e.g. `isImportPathSwapNeeded`).
- Total ban on absolute paths or `file:///` URIs.
- Unused stubs from modularization passes must be pruned to maintain `golangci-lint (strict)` cleanliness.
