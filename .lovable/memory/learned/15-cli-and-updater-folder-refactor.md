# Learned Conventions: CLI & CLI-Updater Folder Architecture Refactoring

## 1. Context & Architectural Motivation
To eliminate ambiguity between the repository root (`gitmap`) and inner Go packages, the core CLI package folder was refactored from `gitmap/` to `cli/`, and the standalone updater was refactored from `gitmap-updater/` to `cli-updater/`.

## 2. Core Conventions & Invariants
- **Go Module Paths**:
  - `cli/go.mod`: `module github.com/alimtvnetwork/gitmap-v28/cli`
  - `cli-updater/go.mod`: `module github.com/alimtvnetwork/gitmap-v28/cli-updater`
- **Import Statements**:
  - All internal imports across `cli/` packages and heavy tests use `"github.com/alimtvnetwork/gitmap-v28/cli/<package>"`.
  - The obsolete `"github.com/alimtvnetwork/gitmap-v28/gitmap"` prefix is completely eliminated.
- **Deploy Manifest Single Source of Truth**:
  - `cli/constants/deploy-manifest.json` specifies `"sourceRepoSubdir": "cli"`.
  - Scripts (`run.ps1`, `run.sh`, `install.sh`, `install.ps1`) consume this manifest to locate source packages dynamically.
- **Local Runner & Checker Anchoring**:
  - `03-ai-scripts/06-cicd-local-runner.py` scopes `CLUSTER_GO_*` patterns to `cli/**/*.go`.
  - CI runners execute within `cli/` working directory.
  - Test inventory paths are strictly relative and reference `cli/...` for all unit and integration tests.
- **Git Commit Lineage**:
  - Folder moves performed via `git mv` to maintain unbroken git blame and file history across ~1,700 source files.
